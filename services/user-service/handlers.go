package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
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
	if strings.TrimSpace(u.SupabaseUserID) != "" {
		out["supabaseUserId"] = u.SupabaseUserID
	}
	if u.DateOfBirth != nil {
		out["dateOfBirth"] = u.DateOfBirth.Format("2006-01-02")
	}
	if u.Nationality != "" {
		out["nationality"] = u.Nationality
	}
	return out
}

// kept in supabase_client.go as supabaseUpstreamHTTP

// clientTypeQuery forwards ?client_type=web|mobile|desktop|server from the incoming request URL.
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
		FullName        string `json:"fullName"`
		Phone           string `json:"phone"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		DateOfBirth     string `json:"dateOfBirth"`
		Nationality     string `json:"nationality"`
		ReferralCode    string `json:"referralCode"`
		TermsAccepted   bool   `json:"termsAccepted"`
		PrivacyAccepted bool   `json:"privacyAccepted"`
		TermsVersion    string `json:"termsVersion"`
		PrivacyVersion  string `json:"privacyVersion"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	contactEmail := normalizeEmail(req.Email)
	if contactEmail == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "email is required"})
		return
	}
	if req.FullName == "" || req.Password == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "fullName and password are required"})
		return
	}
	if !req.TermsAccepted || !req.PrivacyAccepted {
		respond(w, http.StatusBadRequest, map[string]any{"error": "termsAccepted and privacyAccepted are required"})
		return
	}

	phoneE164, err := normalizePhone(req.Phone)
	if err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": "phone must be in Uganda format, e.g. 0700123456 or +256700123456"})
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
		respond(w, http.StatusConflict, map[string]any{"error": "phone number already registered"})
		return
	}
	if _, exists := activeStore.FindByEmailKey(contactEmail); exists {
		respond(w, http.StatusConflict, map[string]any{"error": "email already registered"})
		return
	}

	if !supabaseConfigured() {
		respond(w, http.StatusServiceUnavailable, map[string]any{"error": "auth backend not configured"})
		return
	}

	// Register with Supabase using OTP verification method
	signupResp, _, err := supabasePostWithQuery("/auth/v1/signup", clientTypeQuery(r), map[string]any{
		"email":    contactEmail,
		"password": req.Password,
		"phone":    phoneE164,
		"data": map[string]any{
			"name":  req.FullName,
			"phone": phoneE164,
		},
	})
	if err != nil {
		var sbErr *SupabaseRequestError
		if errors.As(err, &sbErr) && sbErr.Status == http.StatusConflict {
			respond(w, http.StatusConflict, map[string]any{"error": sbErr.Message})
			return
		}
		code, msg := supabaseUpstreamHTTP(err)
		respond(w, code, map[string]any{"error": msg})
		return
	}

	ifUserID := extractSupabaseSignupUserID(signupResp)
	if ifUserID == "" && supabaseServiceRoleConfigured() {
		if id, aerr := supabaseAdminFindUserIDByEmail(contactEmail); aerr == nil && id != "" {
			ifUserID = id
		}
	}

	// Send OTP for email verification
	needsVerify := !supabaseSignupHasAccessToken(signupResp)
	if needsVerify && supabaseConfigured() {
		// Send OTP using the dedicated OTP endpoint
		_, _, otpErr := supabasePostWithQuery("/auth/v1/otp", clientTypeQuery(r), map[string]any{
			"email": contactEmail,
			"type":  "signup",
		})
		if otpErr != nil {
			log.Printf("user-service: OTP send failed for %s: %v", contactEmail, otpErr)
		} else {
			log.Printf("user-service: OTP sent successfully for %s", contactEmail)
		}
	}

	if ifUserID == "" {
		// Supabase requires email verification before revealing the user ID.
		if needsVerify {
			// Store pending profile data so /verify-email can finish creating the local record.
			respond(w, http.StatusCreated, map[string]any{
				"requireEmailVerification": true,
				"pendingLocalProfile":      true,
				"email":                    contactEmail,
				"message":                  "Check your email for a 6-digit verification code.",
			})
			return
		}
		respond(w, http.StatusBadGateway, map[string]any{
			"error": "Supabase signup response missing user id",
		})
		return
	}

	u := User{
		ID:             ifUserID,
		FullName:       req.FullName,
		PhoneE164:      phoneE164,
		AuthEmail:      contactEmail,
		ContactEmail:   contactEmail,
		SupabaseLoginEmail:  contactEmail,
		SupabaseUserID: ifUserID,
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
	now := time.Now()
	_ = activeStore.UpsertUserConsents(UserConsents{
		UserID:            u.ID,
		TermsAcceptedAt:   &now,
		PrivacyAcceptedAt: &now,
		TermsVersion:      strings.TrimSpace(req.TermsVersion),
		PrivacyVersion:    strings.TrimSpace(req.PrivacyVersion),
	})
	refCode := referralCodeForUserID(u.ID)
	_ = activeStore.SetUserReferralCode(u.ID, refCode)
	if rc := strings.TrimSpace(req.ReferralCode); rc != "" {
		if referrerID, ok := activeStore.FindUserIDByReferralCode(rc); ok && referrerID != u.ID {
			createAffiliateReferral(rc, referrerID, u.ID)
		}
		ob.ReferralCode = rc
		_ = activeStore.UpsertOnboarding(ob)
	}
	emitAudit(u.ID, "register", "user", "account created")
	emitNotification("email", u.ContactEmail, "welcome_account_created")

	respond(w, http.StatusCreated, map[string]any{
		"user":                     userPublicProfile(u),
		"session":                  sessionPayloadFromSupabase(signupResp, u.ID),
		"requireEmailVerification": needsVerify,
		"referralCode":             refCode,
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
	if identifier == "" || req.Password == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "email and password are required"})
		return
	}

	// Resolve identifier to a Supabase email.
	// Email is always the primary credential; phone is accepted as a lookup convenience.
	loginEmail := ""
	if strings.Contains(identifier, "@") {
		loginEmail = normalizeEmail(identifier)
	} else {
		if phone, err := normalizePhone(identifier); err == nil {
			if u, ok := activeStore.FindByPhone(phone); ok {
				loginEmail = supabaseLoginEmail(u)
			}
		}
	}

	if loginEmail == "" {
		respond(w, http.StatusUnauthorized, map[string]any{"error": "Invalid email or password"})
		return
	}

	if !supabaseConfigured() {
		respond(w, http.StatusServiceUnavailable, map[string]any{"error": "auth backend not configured"})
		return
	}

	// Supabase login: POST /auth/v1/token?grant_type=password
	loginQuery := clientTypeQuery(r)
	loginQuery.Set("grant_type", "password")
	sessionResp, _, err := supabasePostWithQuery("/auth/v1/token", loginQuery, map[string]any{
		"email":    loginEmail,
		"password": req.Password,
	})
	if err != nil {
		var sbErr *SupabaseRequestError
		if errors.As(err, &sbErr) {
			if sbErr.Status == http.StatusTooManyRequests {
				respond(w, http.StatusTooManyRequests, map[string]any{"error": sbErr.Message})
				return
			}
			payload := map[string]any{"error": sbErr.Message}
			if supabaseIndicatesEmailNotVerified(sbErr.Status, sbErr.Message) {
				payload["requireEmailVerification"] = true
				payload["email"] = loginEmail
				st := sbErr.Status
				if st == http.StatusUnauthorized {
					st = http.StatusForbidden
				}
				respond(w, st, payload)
				return
			}
			respond(w, sbErr.Status, payload)
			return
		}
		respond(w, http.StatusUnauthorized, map[string]any{"error": "Invalid email or password"})
		return
	}

	// Sync Supabase user ID back to the local TAYOSA profile.
	sbUserObj, _ := sessionResp["user"].(map[string]any)
	sbUserID := mapGetString(sbUserObj, "id")

	user, found := activeStore.FindByEmailKey(loginEmail)
	if found && sbUserID != "" {
		user.SupabaseUserID = sbUserID
		user.ID = sbUserID
		_ = activeStore.UpdateIdentity(user)
		if refreshed, ok := activeStore.FindByUserID(sbUserID); ok {
			user = refreshed
		}
	}

	if !found {
		// Supabase login succeeded but no local TAYOSA profile yet.
		respond(w, http.StatusOK, map[string]any{
			"session": sessionPayloadFromSupabase(sessionResp, sbUserID),
			"user": map[string]any{
				"id":           sbUserID,
				"fullName":     mapGetString(sbUserObj, "name"),
				"contactEmail": loginEmail,
			},
		})
		return
	}

	respond(w, http.StatusOK, map[string]any{
		"session": sessionPayloadFromSupabase(sessionResp, user.ID),
		"user":    userPublicProfile(user),
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
		target = supabaseLoginEmail(u)
	}

	if supabaseConfigured() {
		out, _, err := supabaseResendVerification(clientTypeQuery(r), map[string]any{"email": target})
		if err != nil {
			code := http.StatusBadGateway
			var sbErr *SupabaseRequestError
			if errors.As(err, &sbErr) && sbErr.Status >= 400 {
				code = sbErr.Status
			}
			respond(w, code, map[string]any{"error": err.Error()})
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
		target = supabaseLoginEmail(u)
	}

	if supabaseConfigured() {
		// Supabase verify: POST /auth/v1/verify
		out, _, err := supabasePostWithQuery("/auth/v1/verify", clientTypeQuery(r), map[string]any{
			"type":  "signup",
			"token": req.OTP,
			"email": target,
		})
		if err != nil {
			code, msg := supabaseUpstreamHTTP(err)
			if code < 400 {
				code = http.StatusUnauthorized
			}
			respond(w, code, map[string]any{"error": msg})
			return
		}
		// Supabase verify response: { "access_token": "...", "user": { "id": "...", "email_confirmed_at": "..." } }
		userObj, _ := out["user"].(map[string]any)
		ev := mapGetString(userObj, "email_confirmed_at") != ""
		ifUserID := extractSupabaseSignupUserID(out)
		if ifUserID == "" {
			ifUserID = extractSupabaseSignupUserID(map[string]any{"user": userObj})
		}
		if ifUserID == "" {
			respond(w, http.StatusBadGateway, map[string]any{
				"error":              "Supabase verify response missing user id",
				"supabaseResponseKeys": topLevelJSONKeys(out),
			})
			return
		}

		supabaseMail := firstNonEmpty(mapGetString(userObj, "email"), target)
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
				"session": sessionPayloadFromSupabase(out, user.ID),
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
		authEmail := email
		u := User{
			ID:             ifUserID,
			FullName:       strings.TrimSpace(req.FullName),
			PhoneE164:      phoneE164,
			AuthEmail:      authEmail,
			ContactEmail:   contactEmail,
			SupabaseLoginEmail:  supabaseMail,
			SupabaseUserID: ifUserID,
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
			"session":                  sessionPayloadFromSupabase(out, user.ID),
			"user":                     userPublicProfile(user),
			"requireEmailVerification": false,
		})
		return
	}

	respond(w, http.StatusUnauthorized, map[string]any{"error": "Invalid verification code"})
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
		target = supabaseLoginEmail(u)
	}

	if supabaseConfigured() {
		// Supabase recover: POST /auth/v1/recover
		out, _, err := supabasePost("/auth/v1/recover", map[string]any{"email": target})
		if err != nil {
			code := http.StatusBadGateway
			var sbErr *SupabaseRequestError
			if errors.As(err, &sbErr) && sbErr.Status >= 400 {
				code = sbErr.Status
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
		target = supabaseLoginEmail(u)
	}

	if supabaseConfigured() {
		// Supabase: verify recovery OTP via /auth/v1/verify with type=recovery
		out, _, err := supabasePost("/auth/v1/verify", map[string]any{
			"type":  "recovery",
			"token": req.Code,
			"email": target,
		})
		if err != nil {
			respond(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
			return
		}
		// Supabase returns access_token on successful recovery verify
		accessToken := mapGetString(out, "access_token")
		respond(w, http.StatusOK, map[string]any{
			"token":     accessToken,
			"expiresAt": time.Now().Add(15 * time.Minute).Format(time.RFC3339),
		})
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

	if supabaseConfigured() {
		// Use the OTP as a Bearer token (from exchangeResetPasswordToken) to update password
		out, _, err := supabaseUserPut("/auth/v1/user", req.OTP, map[string]any{
			"password": req.NewPassword,
		})
		if err != nil {
			code := http.StatusUnauthorized
			var sbErr *SupabaseRequestError
			if errors.As(err, &sbErr) && sbErr.Status >= 400 {
				code = sbErr.Status
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
	cons, _ := activeStore.GetUserConsents(uid)
	kyc, _ := activeStore.GetKYCProfile(uid)
	sacco, _ := activeStore.GetSaccoMembership(uid)
	kibiina, _ := activeStore.GetKibiinaPreference(uid)
	shares, _, _ := activeStore.GetSharesUnits(uid)
	refCode, _ := activeStore.GetUserReferralCode(uid)
	respond(w, http.StatusOK, map[string]any{
		"user":       userPublicProfile(u),
		"onboarding": ob,
		"consents": map[string]any{
			"termsAcceptedAt":   cons.TermsAcceptedAt,
			"privacyAcceptedAt": cons.PrivacyAcceptedAt,
			"termsVersion":      cons.TermsVersion,
			"privacyVersion":    cons.PrivacyVersion,
		},
		"kyc":          kyc,
		"sacco":        sacco,
		"kibiina":      kibiina,
		"shares":       map[string]any{"balanceUnits": shares},
		"referralCode": refCode,
		"featureAccess": map[string]any{
			"canTransact":  kyc.Status == "approved" && sacco.Status == "enrolled",
			"canJoinKibiina": sacco.Status == "enrolled",
		},
	})
}

func onboardingKYCHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DateOfBirth                string `json:"dateOfBirth"`
		Gender                     string `json:"gender"`
		Nationality                string `json:"nationality"`
		OccupationStatus           string `json:"occupationStatus"`
		IDType                     string `json:"idType"`
		IDNumber                   string `json:"idNumber"`
		IDDocumentFrontKey         string `json:"idDocumentFrontKey"`
		IDDocumentBackKey          string `json:"idDocumentBackKey"`
		SelfieKey                  string `json:"selfieKey"`
		NOKFullName                string `json:"nokFullName"`
		NOKRelationship            string `json:"nokRelationship"`
		NOKPhone                   string `json:"nokPhone"`
		NOKEmail                   string `json:"nokEmail"`
		SourceOfFunds              string `json:"sourceOfFunds"`
		PEPStatus                  *bool  `json:"pepStatus"`
		SACCOMembershipDisclosures string `json:"saccoMembershipDisclosures"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.DateOfBirth) == "" || strings.TrimSpace(req.Gender) == "" || strings.TrimSpace(req.IDType) == "" ||
		strings.TrimSpace(req.IDNumber) == "" || strings.TrimSpace(req.NOKFullName) == "" || strings.TrimSpace(req.NOKRelationship) == "" ||
		strings.TrimSpace(req.NOKPhone) == "" || strings.TrimSpace(req.SourceOfFunds) == "" || req.PEPStatus == nil ||
		strings.TrimSpace(req.SACCOMembershipDisclosures) == "" || strings.TrimSpace(req.IDDocumentFrontKey) == "" ||
		strings.TrimSpace(req.IDDocumentBackKey) == "" || strings.TrimSpace(req.SelfieKey) == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "missing required KYC fields"})
		return
	}
	dob, err := parseOptionalDateOfBirth(req.DateOfBirth)
	if err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	nat, err := normalizeNationality(req.Nationality)
	if err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	uid := authedUserID(r)
	now := time.Now()
	kyc := KYCProfile{
		UserID:                     uid,
		Status:                     "pending",
		DateOfBirth:                dob,
		Gender:                     strings.TrimSpace(req.Gender),
		Nationality:                nat,
		OccupationStatus:           strings.TrimSpace(req.OccupationStatus),
		IDType:                     strings.TrimSpace(req.IDType),
		IDNumber:                   strings.TrimSpace(req.IDNumber),
		NOKFullName:                strings.TrimSpace(req.NOKFullName),
		NOKRelationship:            strings.TrimSpace(req.NOKRelationship),
		NOKPhone:                   strings.TrimSpace(req.NOKPhone),
		NOKEmail:                   strings.TrimSpace(req.NOKEmail),
		SourceOfFunds:              strings.TrimSpace(req.SourceOfFunds),
		PEPStatus:                  req.PEPStatus,
		SACCOMembershipDisclosures: strings.TrimSpace(req.SACCOMembershipDisclosures),
		SubmittedAt:                &now,
	}
	if err := activeStore.UpsertKYCProfile(kyc); err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	docs := []KYCDocument{
		{DocType: "id_document", DocSide: "front", StorageKey: strings.TrimSpace(req.IDDocumentFrontKey)},
		{DocType: "id_document", DocSide: "back", StorageKey: strings.TrimSpace(req.IDDocumentBackKey)},
		{DocType: "selfie", StorageKey: strings.TrimSpace(req.SelfieKey)},
	}
	_ = activeStore.ReplaceKYCDocuments(uid, docs)
	emitAudit(uid, "kyc_submit", "kyc", "kyc submitted")
	if u, ok := activeStore.FindByUserID(uid); ok && u.ContactEmail != "" {
		emitNotification("email", u.ContactEmail, "kyc_submitted")
	}
	respond(w, http.StatusAccepted, map[string]any{"kyc": kyc, "documents": docs})
}

func onboardingSaccoHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		District                 string  `json:"district"`
		County                   string  `json:"county"`
		SubCounty                string  `json:"subCounty"`
		Parish                   string  `json:"parish"`
		Village                  string  `json:"village"`
		StreetPlot               string  `json:"streetPlot"`
		MobileMoneyProvider      string  `json:"mobileMoneyProvider"`
		MobileMoneyNumber        string  `json:"mobileMoneyNumber"`
		SecondaryMoMoNumber      string  `json:"secondaryMoMoNumber"`
		ContributionFrequency    string  `json:"contributionFrequency"`
		SavingsGoalAmount        float64 `json:"savingsGoalAmount"`
		SavingsGoalPurpose       string  `json:"savingsGoalPurpose"`
		SharesToPurchase         int     `json:"sharesToPurchase"`
		EntranceFeePaymentMethod string  `json:"entranceFeePaymentMethod"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.District) == "" || strings.TrimSpace(req.County) == "" || strings.TrimSpace(req.SubCounty) == "" ||
		strings.TrimSpace(req.Parish) == "" || strings.TrimSpace(req.Village) == "" || strings.TrimSpace(req.MobileMoneyProvider) == "" ||
		strings.TrimSpace(req.MobileMoneyNumber) == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "missing required SACCO fields"})
		return
	}
	if req.SharesToPurchase < 1 {
		req.SharesToPurchase = 1
	}
	ok, err := activeStore.GeoRecordExists(req.District, req.County, req.SubCounty, req.Parish, req.Village)
	if err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	if !ok {
		respond(w, http.StatusBadRequest, map[string]any{"error": "geo does not match official hierarchy"})
		return
	}
	uid := authedUserID(r)
	if k, found := activeStore.GetKYCProfile(uid); !found || k.Status != "approved" {
		respond(w, http.StatusForbidden, map[string]any{"error": "KYC must be approved before SACCO enrollment"})
		return
	}
	mem := SaccoMembership{
		UserID:                   uid,
		Status:                   "enrolled",
		District:                 strings.TrimSpace(req.District),
		County:                   strings.TrimSpace(req.County),
		SubCounty:                strings.TrimSpace(req.SubCounty),
		Parish:                   strings.TrimSpace(req.Parish),
		Village:                  strings.TrimSpace(req.Village),
		StreetPlot:               strings.TrimSpace(req.StreetPlot),
		MobileMoneyProvider:      strings.TrimSpace(req.MobileMoneyProvider),
		MobileMoneyNumber:        strings.TrimSpace(req.MobileMoneyNumber),
		SecondaryMoMoNumber:      strings.TrimSpace(req.SecondaryMoMoNumber),
		ContributionFrequency:    strings.TrimSpace(req.ContributionFrequency),
		SavingsGoalAmount:        req.SavingsGoalAmount,
		SavingsGoalPurpose:       strings.TrimSpace(req.SavingsGoalPurpose),
		SharesToPurchase:         req.SharesToPurchase,
		EntranceFeePaymentMethod: strings.TrimSpace(req.EntranceFeePaymentMethod),
	}
	if err := activeStore.UpsertSaccoMembership(mem); err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	_ = activeStore.EnsureSharesLedger(uid, req.SharesToPurchase)
	emitAudit(uid, "sacco_enroll", "sacco_membership", "sacco enrollment submitted")
	respond(w, http.StatusOK, map[string]any{"membership": mem})
}

func onboardingKibiinaHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action                 string  `json:"action"`
		InviteCode             string  `json:"inviteCode"`
		GroupName              string  `json:"groupName"`
		ContributionAmount     float64 `json:"contributionAmount"`
		CycleFrequency         string  `json:"cycleFrequency"`
		MaxGroupSize           int     `json:"maxGroupSize"`
		PayoutOrderPreference  string  `json:"payoutOrderPreference"`
		NotificationPreference string  `json:"notificationPreference"`
		LanguagePreference     string  `json:"languagePreference"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.Action) == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "action is required"})
		return
	}
	uid := authedUserID(r)
	if s, ok := activeStore.GetSaccoMembership(uid); !ok || s.Status != "enrolled" {
		respond(w, http.StatusForbidden, map[string]any{"error": "SACCO membership is required before Kibiina setup"})
		return
	}
	p := KibiinaPreference{
		UserID:                 uid,
		Action:                 strings.TrimSpace(req.Action),
		InviteCode:             strings.TrimSpace(req.InviteCode),
		GroupName:              strings.TrimSpace(req.GroupName),
		ContributionAmount:     req.ContributionAmount,
		CycleFrequency:         strings.TrimSpace(req.CycleFrequency),
		MaxGroupSize:           req.MaxGroupSize,
		PayoutOrderPreference:  strings.TrimSpace(req.PayoutOrderPreference),
		NotificationPreference: strings.TrimSpace(req.NotificationPreference),
		LanguagePreference:     strings.TrimSpace(req.LanguagePreference),
	}
	if err := activeStore.UpsertKibiinaPreference(p); err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	emitAudit(uid, "kibiina_setup", "kibiina", "kibiina preference captured")
	respond(w, http.StatusOK, map[string]any{"kibiina": p})
}

func adminKYCDecisionHandler(w http.ResponseWriter, r *http.Request) {
	adminSecret := strings.TrimSpace(os.Getenv("ADMIN_API_KEY"))
	if adminSecret == "" || strings.TrimSpace(r.Header.Get("X-Admin-Secret")) != adminSecret {
		respond(w, http.StatusUnauthorized, map[string]any{"error": "admin authentication failed"})
		return
	}
	if r.Method == http.MethodGet {
		status := strings.TrimSpace(r.URL.Query().Get("status"))
		limit := 50
		items, err := activeStore.ListAdminKYCQueue(status, limit)
		if err != nil {
			respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		respond(w, http.StatusOK, map[string]any{"items": items, "count": len(items)})
		return
	}
	var req struct {
		Status     string `json:"status"`
		ReviewNote string `json:"reviewNote"`
		ReviewedBy string `json:"reviewedBy"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if req.Status != "approved" && req.Status != "rejected" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "status must be approved or rejected"})
		return
	}
	userID := strings.TrimSpace(r.URL.Query().Get("userId"))
	if userID == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "userId query parameter is required"})
		return
	}
	if err := activeStore.SetKYCDecision(userID, req.Status, strings.TrimSpace(req.ReviewedBy), strings.TrimSpace(req.ReviewNote)); err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	emitAudit(firstNonEmpty(req.ReviewedBy, "admin"), "kyc_"+req.Status, "kyc", "admin decision")
	if u, ok := activeStore.FindByUserID(userID); ok && u.ContactEmail != "" {
		emitNotification("email", u.ContactEmail, "kyc_"+req.Status)
	}
	respond(w, http.StatusOK, map[string]any{"status": req.Status, "userId": userID})
}

func adminSettingsHandler(w http.ResponseWriter, r *http.Request) {
	adminSecret := strings.TrimSpace(os.Getenv("ADMIN_API_KEY"))
	if adminSecret == "" || strings.TrimSpace(r.Header.Get("X-Admin-Secret")) != adminSecret {
		respond(w, http.StatusUnauthorized, map[string]any{"error": "admin authentication failed"})
		return
	}
	key := strings.TrimSpace(r.URL.Query().Get("key"))
	if key == "" {
		key = "fees"
	}
	switch r.Method {
	case http.MethodGet:
		v, ok, err := activeStore.GetAdminSetting(key)
		if err != nil {
			respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		if !ok {
			v = map[string]any{}
		}
		respond(w, http.StatusOK, map[string]any{"key": key, "value": v})
	case http.MethodPatch, http.MethodPost:
		var body map[string]any
		if err := decodeJSON(r, &body); err != nil {
			respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		if err := activeStore.SetAdminSetting(key, body); err != nil {
			respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		emitAudit("admin", "admin_setting_update", "admin_settings:"+key, "updated")
		respond(w, http.StatusOK, map[string]any{"ok": true})
	default:
		respond(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
	}
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
