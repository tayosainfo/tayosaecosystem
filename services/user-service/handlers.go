package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func respond(w http.ResponseWriter, code int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func decodeJSON(r *http.Request, target any) error {
	if r.Body == nil {
		return errors.New("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(target)
}

// userPublicProfile is the stable JSON shape for "user" on register, login, and extensions like /users/me.
func userPublicProfile(u User) map[string]any {
	out := map[string]any{
		"id":                   u.ID,
		"fullName":             u.FullName,
		"phoneE164":            u.PhoneE164,
		"contactEmail":         u.ContactEmail,
		"contactEmailVerified": u.ContactEmailChecked,
		"createdAt":            u.CreatedAt,
	}
	if strings.TrimSpace(u.InsforgeUserID) != "" {
		out["insforgeUserId"] = u.InsforgeUserID
	}
	if u.DateOfBirth != nil {
		out["dateOfBirth"] = u.DateOfBirth.Format("2006-01-02")
	}
	if u.Nationality != "" {
		out["nationality"] = u.Nationality
	}
	return out
}

// insforgeUpstreamHTTP maps errors from insforgePost / InsforgeRequestError to an HTTP status and message.
func insforgeUpstreamHTTP(err error) (status int, msg string) {
	var ifErr *InsforgeRequestError
	if errors.As(err, &ifErr) {
		s := ifErr.Status
		if s >= 400 && s <= 599 {
			return s, ifErr.Message
		}
		return http.StatusBadGateway, ifErr.Message
	}
	return http.StatusBadGateway, err.Error()
}

// clientTypeQuery forwards InsForge ?client_type=web|mobile|desktop|server from the incoming request URL.
func clientTypeQuery(r *http.Request) url.Values {
	q := url.Values{}
	if ct := strings.TrimSpace(r.URL.Query().Get("client_type")); ct != "" {
		q.Set("client_type", ct)
	}
	return q
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// pgx surfaces SQLSTATE 23505 for unique_violation
	return strings.Contains(err.Error(), "23505")
}

func usesPostgres() bool {
	return databaseDSN() != ""
}

func lookupUserByEmailKey(email string) (User, bool) {
	return activeStore.FindByEmailKey(email)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FullName    string `json:"fullName"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		DateOfBirth string `json:"dateOfBirth"`
		Nationality string `json:"nationality"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	dobPtr, err := parseOptionalDateOfBirth(req.DateOfBirth)
	if err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	nationality, err := normalizeNationality(req.Nationality)
	if err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	phoneE164, err := normalizePhone(req.Phone)
	if err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	contactEmail := normalizeEmail(req.Email)
	if req.FullName == "" || req.Password == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "fullName and password are required"})
		return
	}

	if insforgeConfigured() && contactEmail == "" {
		respond(w, http.StatusBadRequest, map[string]any{
			"error": "email is required so InsForge can deliver verification and password-reset codes",
		})
		return
	}

	if _, exists := activeStore.FindByPhone(phoneE164); exists {
		respond(w, http.StatusConflict, map[string]any{"error": "phone already registered"})
		return
	}
	if contactEmail != "" {
		if _, exists := activeStore.FindByEmailKey(contactEmail); exists {
			respond(w, http.StatusConflict, map[string]any{"error": "email already registered"})
			return
		}
	}

	authEmail := makeAuthEmail(phoneE164)
	insforgeMail := contactEmail
	if insforgeMail == "" {
		insforgeMail = authEmail
	}

	if insforgeConfigured() {
		signupResp, _, err := insforgePostWithQuery("/api/auth/users", clientTypeQuery(r), map[string]any{
			"email":    insforgeMail,
			"password": req.Password,
			"name":     req.FullName,
		})
		if err != nil {
			var ifErr *InsforgeRequestError
			if errors.As(err, &ifErr) && ifErr.Status == http.StatusConflict {
				respond(w, http.StatusConflict, map[string]any{"error": ifErr.Message})
				return
			}
			code, msg := insforgeUpstreamHTTP(err)
			respond(w, code, map[string]any{"error": msg})
			return
		}

		ifUserID := extractInsForgeSignupUserID(signupResp)
		if ifUserID == "" && insforgeAdminConfigured() {
			if id, aerr := insforgeAdminFindUserIDByEmail(insforgeMail); aerr == nil && id != "" {
				ifUserID = id
			}
		}
		if ifUserID == "" {
			probe, _, _ := insforgePostAnonAnyHTTPStatus("/api/auth/sessions", clientTypeQuery(r), map[string]any{
				"email":    insforgeMail,
				"password": req.Password,
			})
			if id := extractInsForgeSignupUserID(probe); id != "" {
				ifUserID = id
			}
		}
		if ifUserID == "" {
			needsVerify := mapGetBool(signupResp, "requireEmailVerification") || !insforgeSignupHasAccessToken(signupResp)
			if needsVerify {
				respond(w, http.StatusCreated, map[string]any{
					"requireEmailVerification": true,
					"pendingLocalProfile":      true,
					"email":                    insforgeMail,
					"message":                  "InsForge created your account. Check your email for the verification code, then return here or use Verify Email. After verifying, submit the same phone and password once on the verify page to finish your Tayosa profile.",
				})
				return
			}
			payload := map[string]any{
				"error":                "InsForge sign-up response missing user id",
				"hint":                 "Some InsForge responses omit user.id. Set INSFORGE_ADMIN_API_KEY (project admin) for lookup, or inspect insforgeResponseKeys. If email verification is on, you should still have received a 201 pending response — contact support with these keys.",
				"insforgeResponseKeys": topLevelJSONKeys(signupResp),
			}
			respond(w, http.StatusBadGateway, payload)
			return
		}

		u := User{
			ID:             ifUserID,
			FullName:       req.FullName,
			PhoneE164:      phoneE164,
			AuthEmail:      authEmail,
			ContactEmail:   contactEmail,
			InsforgeEmail:  insforgeMail,
			InsforgeUserID: ifUserID,
			DateOfBirth:    dobPtr,
			Nationality:    nationality,
			CreatedAt:      time.Now(),
		}
		ob := OnboardingProfile{
			UserID:         u.ID,
			Phase:          1,
			TrustScoreSeed: 10,
			LastUpdatedAt:  time.Now(),
		}
		if err := activeStore.CreateIdentityWithOnboarding(u, ob, nil); err != nil {
			if isUniqueViolation(err) {
				respond(w, http.StatusConflict, map[string]any{"error": "identity already exists"})
				return
			}
			respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}

		respond(w, http.StatusCreated, map[string]any{
			"user":                     userPublicProfile(u),
			"session":                  sessionPayloadFromInsForge(signupResp, u.ID),
			"requireEmailVerification": signupResp["requireEmailVerification"],
		})
		return
	}

	u := User{
		ID:            "usr_" + strings.ReplaceAll(strings.TrimPrefix(phoneE164, "+"), " ", ""),
		FullName:      req.FullName,
		PhoneE164:     phoneE164,
		AuthEmail:     authEmail,
		ContactEmail:  contactEmail,
		InsforgeEmail: insforgeMail,
		DateOfBirth:   dobPtr,
		Nationality:   nationality,
		Password:      req.Password,
		CreatedAt:     time.Now(),
	}
	ob := OnboardingProfile{
		UserID:         u.ID,
		Phase:          1,
		TrustScoreSeed: 10,
		LastUpdatedAt:  time.Now(),
	}

	var passHash *string
	if usesPostgres() {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			respond(w, http.StatusInternalServerError, map[string]any{"error": "could not hash password"})
			return
		}
		s := string(hash)
		passHash = &s
		u.Password = ""
	}

	if err := activeStore.CreateIdentityWithOnboarding(u, ob, passHash); err != nil {
		if isUniqueViolation(err) {
			respond(w, http.StatusConflict, map[string]any{"error": "identity already exists"})
			return
		}
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	if !usesPostgres() {
		devMapsMu.Lock()
		devVerifyCodes[u.InsforgeEmail] = "123456"
		devMapsMu.Unlock()
	}

	token := "dev-token-" + u.ID
	respond(w, http.StatusCreated, map[string]any{
		"user":    userPublicProfile(u),
		"session": map[string]any{"accessToken": token, "userId": u.ID},
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	identifier := strings.TrimSpace(req.Identifier)
	var user User
	var found bool

	if strings.Contains(identifier, "@") {
		email := normalizeEmail(identifier)
		user, found = activeStore.FindByEmailKey(email)
	} else {
		if phone, err := normalizePhone(identifier); err == nil {
			user, found = activeStore.FindByPhone(phone)
		}
	}

	if !found {
		if insforgeConfigured() && strings.Contains(identifier, "@") {
			emailTry := normalizeEmail(identifier)
			if emailTry != "" {
				_, _, err := insforgePostWithQuery("/api/auth/sessions", clientTypeQuery(r), map[string]any{
					"email":    emailTry,
					"password": req.Password,
				})
				if err != nil {
					var ifErr *InsforgeRequestError
					if errors.As(err, &ifErr) {
						if ifErr.Status == http.StatusTooManyRequests {
							respond(w, http.StatusTooManyRequests, map[string]any{"error": ifErr.Message})
							return
						}
						if ifErr.Status >= 500 {
							respond(w, ifErr.Status, map[string]any{"error": ifErr.Message})
							return
						}
						if ifErr.Status >= 400 {
							payload := map[string]any{"error": ifErr.Message}
							if insforgeIndicatesEmailNotVerified(ifErr.Status, ifErr.Message) {
								payload["requireEmailVerification"] = true
								payload["email"] = emailTry
								st := ifErr.Status
								if st == http.StatusUnauthorized {
									st = http.StatusForbidden
								}
								respond(w, st, payload)
								return
							}
							respond(w, ifErr.Status, payload)
							return
						}
					}
				}
			}
		}
		respond(w, http.StatusUnauthorized, map[string]any{"error": "Invalid credentials. Use your registered phone number or email. If you have not verified your email yet, sign in with the same email you used at registration."})
		return
	}

	if insforgeConfigured() {
		sessionResp, _, err := insforgePostWithQuery("/api/auth/sessions", clientTypeQuery(r), map[string]any{
			"email":    insforgeLoginEmail(user),
			"password": req.Password,
		})
		if err != nil {
			var ifErr *InsforgeRequestError
			if errors.As(err, &ifErr) {
				if ifErr.Status == http.StatusTooManyRequests {
					respond(w, http.StatusTooManyRequests, map[string]any{"error": ifErr.Message})
					return
				}
				if ifErr.Status >= 500 {
					respond(w, ifErr.Status, map[string]any{"error": ifErr.Message})
					return
				}
				if ifErr.Status >= 400 {
					payload := map[string]any{"error": ifErr.Message}
					if insforgeIndicatesEmailNotVerified(ifErr.Status, ifErr.Message) {
						payload["requireEmailVerification"] = true
						payload["email"] = insforgeLoginEmail(user)
						st := ifErr.Status
						if st == http.StatusUnauthorized {
							st = http.StatusForbidden
						}
						respond(w, st, payload)
						return
					}
					respond(w, ifErr.Status, payload)
					return
				}
			}
			respond(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
			return
		}
		userObj, _ := sessionResp["user"].(map[string]any)
		if id := mapGetString(userObj, "id"); id != "" {
			user.InsforgeUserID = id
			user.ID = id
		}
		_ = activeStore.UpdateIdentity(user)
		if refreshed, ok := activeStore.FindByUserID(user.ID); ok {
			user = refreshed
		}
		respond(w, http.StatusOK, map[string]any{
			"session": sessionPayloadFromInsForge(sessionResp, user.ID),
			"user":    userPublicProfile(user),
		})
		return
	}

	if user.PasswordHash != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			respond(w, http.StatusUnauthorized, map[string]any{"error": "Invalid credentials. Use your registered phone number or email."})
			return
		}
	} else if user.Password != req.Password {
		respond(w, http.StatusUnauthorized, map[string]any{"error": "Invalid credentials. Use your registered phone number or email."})
		return
	}

	respond(w, http.StatusOK, map[string]any{
		"session": map[string]any{
			"accessToken": "dev-token-" + user.ID,
			"userId":      user.ID,
		},
		"user": userPublicProfile(user),
	})
}

func resendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	email := normalizeEmail(req.Email)
	if email == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "email is required"})
		return
	}

	target := email
	if u, ok := lookupUserByEmailKey(email); ok {
		target = insforgeLoginEmail(u)
	}

	if insforgeConfigured() {
		// send-verification is a server-side operation — use the admin/API key.
		// The anon key does not have permission to trigger verification emails in InsForge.
		var (
			out map[string]any
			err error
		)
		if insforgeAdminConfigured() {
			out, _, err = insforgeAdminPost("/api/auth/email/send-verification", clientTypeQuery(r), map[string]any{"email": target})
		} else {
			// Fallback to anon key if no admin key is set (may still fail depending on InsForge project settings)
			out, _, err = insforgePostWithQuery("/api/auth/email/send-verification", clientTypeQuery(r), map[string]any{"email": target})
		}
		if err != nil {
			code := http.StatusBadGateway
			var ifErr *InsforgeRequestError
			if errors.As(err, &ifErr) && ifErr.Status >= 400 {
				code = ifErr.Status
			}
			payload := map[string]any{"error": err.Error()}
			if code >= 500 || strings.Contains(strings.ToLower(err.Error()), "verification token") {
				payload["hint"] = "InsForge rejected creating or emailing a verification token. Ensure INSFORGE_ADMIN_API_KEY is set in your .env, outbound email (SMTP or your provider) is configured in the InsForge project dashboard, and the project is active."
			}
			respond(w, code, payload)
			return
		}
		respond(w, http.StatusOK, out)
		return
	}

	devMapsMu.Lock()
	devVerifyCodes[target] = "123456"
	devMapsMu.Unlock()
	respond(w, http.StatusOK, map[string]any{"success": true, "message": "Verification code sent"})
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		FullName    string `json:"fullName"`
		Phone       string `json:"phone"`
		DateOfBirth string `json:"dateOfBirth"`
		Nationality string `json:"nationality"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	email := normalizeEmail(req.Email)
	if email == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "email is required"})
		return
	}

	target := email
	if u, ok := lookupUserByEmailKey(email); ok {
		target = insforgeLoginEmail(u)
	}

	if insforgeConfigured() {
		out, _, err := insforgePostWithQuery("/api/auth/email/verify", clientTypeQuery(r), map[string]any{
			"email": target,
			"otp":   req.OTP,
		})
		if err != nil {
			code, msg := insforgeUpstreamHTTP(err)
			if code < 400 {
				code = http.StatusUnauthorized
			}
			respond(w, code, map[string]any{"error": msg})
			return
		}
		userObj, _ := out["user"].(map[string]any)
		ev := false
		if v, ok := userObj["emailVerified"].(bool); ok {
			ev = v
		}
		ifUserID := extractInsForgeSignupUserID(out)
		if ifUserID == "" {
			ifUserID = extractInsForgeSignupUserID(map[string]any{"user": userObj})
		}
		if ifUserID == "" {
			respond(w, http.StatusBadGateway, map[string]any{
				"error":                "InsForge verify response missing user id",
				"insforgeResponseKeys": topLevelJSONKeys(out),
			})
			return
		}

		insforgeMail := firstNonEmpty(mapGetString(userObj, "email"), target)
		u0, userOk := lookupUserByEmailKey(email)
		if userOk {
			_ = activeStore.SetContactEmailVerified(u0.PhoneE164, ev)
			var user User
			if refreshed, ok := activeStore.FindByUserID(u0.ID); ok {
				user = refreshed
			} else {
				user = u0
			}
			respond(w, http.StatusOK, map[string]any{
				"session": sessionPayloadFromInsForge(out, user.ID),
				"user":    userPublicProfile(user),
			})
			return
		}

		if strings.TrimSpace(req.FullName) == "" || strings.TrimSpace(req.Phone) == "" {
			respond(w, http.StatusBadRequest, map[string]any{
				"error":               "fullName and phone are required to finish your Tayosa profile after email verification",
				"pendingLocalProfile": true,
			})
			return
		}
		phoneE164, err := normalizePhone(req.Phone)
		if err != nil {
			respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		dobPtr, err := parseOptionalDateOfBirth(req.DateOfBirth)
		if err != nil {
			respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		nationality, err := normalizeNationality(req.Nationality)
		if err != nil {
			respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		if _, exists := activeStore.FindByPhone(phoneE164); exists {
			respond(w, http.StatusConflict, map[string]any{"error": "phone already registered"})
			return
		}
		contactEmail := email
		if contactEmail != "" {
			if _, exists := activeStore.FindByEmailKey(contactEmail); exists {
				respond(w, http.StatusConflict, map[string]any{"error": "email already registered"})
				return
			}
		}
		authEmail := makeAuthEmail(phoneE164)
		u := User{
			ID:             ifUserID,
			FullName:       strings.TrimSpace(req.FullName),
			PhoneE164:      phoneE164,
			AuthEmail:      authEmail,
			ContactEmail:   contactEmail,
			InsforgeEmail:  insforgeMail,
			InsforgeUserID: ifUserID,
			DateOfBirth:    dobPtr,
			Nationality:    nationality,
			CreatedAt:      time.Now(),
		}
		ob := OnboardingProfile{
			UserID:         u.ID,
			Phase:          1,
			TrustScoreSeed: 10,
			LastUpdatedAt:  time.Now(),
		}
		if err := activeStore.CreateIdentityWithOnboarding(u, ob, nil); err != nil {
			if isUniqueViolation(err) {
				respond(w, http.StatusConflict, map[string]any{"error": "identity already exists"})
				return
			}
			respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		_ = activeStore.SetContactEmailVerified(u.PhoneE164, ev)
		var user User
		if refreshed, ok := activeStore.FindByUserID(u.ID); ok {
			user = refreshed
		} else {
			user = u
		}
		respond(w, http.StatusOK, map[string]any{
			"session":                  sessionPayloadFromInsForge(out, user.ID),
			"user":                     userPublicProfile(user),
			"requireEmailVerification": false,
		})
		return
	}

	devMapsMu.Lock()
	expected := devVerifyCodes[target]
	if expected == "" || req.OTP != expected {
		devMapsMu.Unlock()
		respond(w, http.StatusUnauthorized, map[string]any{"error": "Invalid verification code"})
		return
	}
	delete(devVerifyCodes, target)
	devMapsMu.Unlock()
	respond(w, http.StatusOK, map[string]any{
		"accessToken": "dev-token-verified-" + strings.ReplaceAll(email, "@", "_"),
		"user":        map[string]any{"email": email, "emailVerified": true},
	})
}

func sendResetPasswordEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	email := normalizeEmail(req.Email)
	if email == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "email is required"})
		return
	}
	target := email
	if u, ok := lookupUserByEmailKey(email); ok {
		target = insforgeLoginEmail(u)
	}

	if insforgeConfigured() {
		out, _, err := insforgePost("/api/auth/email/send-reset-password", map[string]any{"email": target})
		if err != nil {
			code := http.StatusBadGateway
			var ifErr *InsforgeRequestError
			if errors.As(err, &ifErr) && ifErr.Status >= 400 {
				code = ifErr.Status
			}
			respond(w, code, map[string]any{"error": err.Error()})
			return
		}
		respond(w, http.StatusOK, out)
		return
	}

	devMapsMu.Lock()
	devResetCodes[target] = "123456"
	devMapsMu.Unlock()
	respond(w, http.StatusOK, map[string]any{"success": true, "message": "Password reset email sent"})
}

func exchangeResetPasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	email := normalizeEmail(req.Email)
	if email == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "email is required"})
		return
	}
	target := email
	if u, ok := lookupUserByEmailKey(email); ok {
		target = insforgeLoginEmail(u)
	}

	if insforgeConfigured() {
		out, _, err := insforgePost("/api/auth/email/exchange-reset-password-token", map[string]any{
			"email": target,
			"code":  req.Code,
		})
		if err != nil {
			respond(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
			return
		}
		respond(w, http.StatusOK, out)
		return
	}

	devMapsMu.Lock()
	defer devMapsMu.Unlock()
	if devResetCodes[target] == "" || devResetCodes[target] != req.Code {
		respond(w, http.StatusUnauthorized, map[string]any{"error": "Invalid or expired reset code"})
		return
	}
	token := "reset-token-" + strings.ReplaceAll(email, "@", "_")
	devResetTokens[token] = email
	delete(devResetCodes, target)
	respond(w, http.StatusOK, map[string]any{
		"token":     token,
		"expiresAt": time.Now().Add(15 * time.Minute).Format(time.RFC3339),
	})
}

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NewPassword string `json:"newPassword"`
		OTP         string `json:"otp"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if len(req.NewPassword) < 6 {
		respond(w, http.StatusBadRequest, map[string]any{"error": "newPassword must be at least 6 characters"})
		return
	}

	if insforgeConfigured() {
		out, _, err := insforgePost("/api/auth/email/reset-password", map[string]any{
			"newPassword": req.NewPassword,
			"otp":         req.OTP,
		})
		if err != nil {
			code := http.StatusUnauthorized
			var ifErr *InsforgeRequestError
			if errors.As(err, &ifErr) && ifErr.Status >= 400 {
				code = ifErr.Status
			}
			respond(w, code, map[string]any{"error": err.Error()})
			return
		}
		respond(w, http.StatusOK, out)
		return
	}

	devMapsMu.Lock()
	email := devResetTokens[req.OTP]
	if email == "" {
		devMapsMu.Unlock()
		respond(w, http.StatusUnauthorized, map[string]any{"error": "Invalid reset token"})
		return
	}
	delete(devResetTokens, req.OTP)
	devMapsMu.Unlock()

	u, found := activeStore.FindByEmailKey(email)
	if !found {
		u, found = activeStore.FindByEmailKey(normalizeEmail(email))
	}
	if !found {
		// try auth email map via second lookup — FindByEmailKey already covers auth
		respond(w, http.StatusOK, map[string]any{"message": "Password reset successfully. Please login with your new password."})
		return
	}

	if usesPostgres() {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			respond(w, http.StatusInternalServerError, map[string]any{"error": "could not hash password"})
			return
		}
		u.PasswordHash = string(hash)
		u.Password = ""
	} else {
		u.Password = req.NewPassword
	}
	_ = activeStore.UpdateIdentity(u)

	respond(w, http.StatusOK, map[string]any{
		"message": "Password reset successfully. Please login with your new password.",
	})
}

func onboardingHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID       string            `json:"userId"`
		Phase        int               `json:"phase"`
		KYC          map[string]any    `json:"kyc"`
		Membership   map[string]any    `json:"membership"`
		Kibiina      map[string]any    `json:"kibiina"`
		ReferralCode string            `json:"referralCode"`
		Geo          map[string]string `json:"geo"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if req.UserID == "" || req.Phase < 2 || req.Phase > 4 {
		respond(w, http.StatusBadRequest, map[string]any{"error": "valid userId and phase 2-4 required"})
		return
	}
	if req.Geo == nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": "geo is required"})
		return
	}
	for _, key := range []string{"district", "county", "sub_county", "parish", "village"} {
		if strings.TrimSpace(req.Geo[key]) == "" {
			respond(w, http.StatusBadRequest, map[string]any{
				"error": "geo must include district, county, sub_county, parish, and village (matching the geo selector cascade)",
			})
			return
		}
	}
	ok, err := activeStore.GeoRecordExists(
		req.Geo["district"], req.Geo["county"], req.Geo["sub_county"], req.Geo["parish"], req.Geo["village"])
	if err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	if !ok {
		respond(w, http.StatusBadRequest, map[string]any{
			"error": "geo does not match any row in the official Uganda hierarchy (see GET /api/v1/geo)",
		})
		return
	}

	if authedUserID(r) != req.UserID {
		respond(w, http.StatusForbidden, map[string]any{"error": "cannot modify another user's onboarding"})
		return
	}

	p := OnboardingProfile{
		UserID:         req.UserID,
		Phase:          req.Phase,
		KYC:            req.KYC,
		Membership:     req.Membership,
		Kibiina:        req.Kibiina,
		ReferralCode:   req.ReferralCode,
		Geo:            req.Geo,
		TrustScoreSeed: 10,
		LastUpdatedAt:  time.Now(),
	}
	if err := activeStore.UpsertOnboarding(p); err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	respond(w, http.StatusOK, map[string]any{"profile": p})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	uid := authedUserID(r)
	u, ok := activeStore.FindByUserID(uid)
	if !ok {
		respond(w, http.StatusNotFound, map[string]any{"error": "user not found in platform store"})
		return
	}
	ob, _ := activeStore.GetOnboarding(uid)
	respond(w, http.StatusOK, map[string]any{
		"user":       userPublicProfile(u),
		"onboarding": ob,
	})
}

func geoLookupHandler(w http.ResponseWriter, r *http.Request) {
	level := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("level")))
	parent := strings.TrimSpace(r.URL.Query().Get("parent"))
	switch level {
	case "district", "county", "subcounty", "parish", "village":
	default:
		respond(w, http.StatusBadRequest, map[string]any{"error": "level must be district|county|subcounty|parish|village"})
		return
	}
	values, err := activeStore.GeoDistinct(level, parent)
	if err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	respond(w, http.StatusOK, map[string]any{"items": values, "count": len(values)})
}

func parishGroupPolicyHandler(w http.ResponseWriter, r *http.Request) {
	st, err := activeStore.GroupPolicyStats()
	if err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	respond(w, http.StatusOK, map[string]any{
		"parishSaccoPolicy":              "one parish SACCO boundary",
		"villageKibiinaPolicy":           "unlimited village-level merry-go-round groups per parish SACCO",
		"source":                         "uganda_geo_data_2025-11-26.csv",
		"parishSaccosRegistered":         st.ParishSaccos,
		"villageKibiinaGroupsRegistered": st.VillageKibiinaGroups,
	})
}
