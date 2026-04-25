package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// adminCheckStatusHandler handles POST /api/v1/admin/check-status
// Checks if the current user has admin role
func adminCheckStatusHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	
	if strings.TrimSpace(req.Email) == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "email is required"})
		return
	}
	
	// Find user by email
	user, ok := activeStore.FindByEmail(req.Email)
	if !ok {
		// User not found - return non-admin status
		respond(w, http.StatusOK, map[string]any{
			"isAdmin": false,
			"role":    "user",
			"email":   req.Email,
		})
		return
	}
	
	// Return user's role
	respond(w, http.StatusOK, map[string]any{
		"isAdmin": user.Role == "admin",
		"role":    user.Role,
		"email":   user.AuthEmail,
	})
}

// adminUsersListHandler handles GET /api/v1/admin/users
// Returns paginated list of users with search and filtering
func adminUsersListHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	statusFilter := strings.TrimSpace(r.URL.Query().Get("status"))
	kycFilter := strings.TrimSpace(r.URL.Query().Get("kyc_status"))
	
	// Calculate offset
	offset := (page - 1) * limit
	
	// Get users from store with filters
	users, total, err := activeStore.ListUsersWithFilters(search, statusFilter, kycFilter, limit, offset)
	if err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	
	// Calculate total pages
	totalPages := (total + limit - 1) / limit
	
	// Format response
	userList := make([]map[string]any, 0, len(users))
	for _, u := range users {
		kyc, _ := activeStore.GetKYCProfile(u.ID)
		
		userList = append(userList, map[string]any{
			"user_id":     u.ID,
			"full_name":   u.FullName,
			"auth_email":  u.AuthEmail,
			"phone_e164":  u.PhoneE164,
			"role":        u.Role,
			"status":      u.Status,
			"kyc_status":  kyc.Status,
			"created_at":  u.CreatedAt.Format(time.RFC3339),
			"last_login":  formatOptionalTime(u.LastLogin),
		})
	}
	
	respond(w, http.StatusOK, map[string]any{
		"users":      userList,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	})
}

// adminUserDetailHandler handles GET /api/v1/admin/users/{userId}
// Returns comprehensive user information
func adminUserDetailHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "userId is required"})
		return
	}
	
	// Get user from store
	user, ok := activeStore.FindByUserID(userID)
	if !ok {
		respond(w, http.StatusNotFound, map[string]any{"error": "user not found"})
		return
	}
	
	// Get related data
	kyc, _ := activeStore.GetKYCProfile(userID)
	onboarding, _ := activeStore.GetOnboarding(userID)
	sacco, _ := activeStore.GetSaccoMembership(userID)
	kibiina, _ := activeStore.GetKibiinaPreference(userID)
	
	respond(w, http.StatusOK, map[string]any{
		"user": map[string]any{
			"user_id":           user.ID,
			"full_name":         user.FullName,
			"auth_email":        user.AuthEmail,
			"phone_e164":        user.PhoneE164,
			"contact_email":     user.ContactEmail,
			"role":              user.Role,
			"status":            user.Status,
			"created_at":        user.CreatedAt.Format(time.RFC3339),
			"last_login":        formatOptionalTime(user.LastLogin),
			"role_assigned_at":  formatOptionalTime(user.RoleAssignedAt),
			"role_assigned_by":  user.RoleAssignedBy,
			"date_of_birth":     formatOptionalDate(user.DateOfBirth),
			"nationality":       user.Nationality,
		},
		"kyc": map[string]any{
			"status":       kyc.Status,
			"submitted_at": formatOptionalTime(kyc.SubmittedAt),
			"reviewed_at":  formatOptionalTime(kyc.ReviewedAt),
			"reviewed_by":  kyc.ReviewedBy,
		},
		"onboarding": map[string]any{
			"phase":           onboarding.Phase,
			"last_updated_at": onboarding.LastUpdatedAt.Format(time.RFC3339),
		},
		"sacco": map[string]any{
			"status": sacco.Status,
		},
		"kibiina": map[string]any{
			"action": kibiina.Action,
		},
	})
}

// adminUserStatusHandler handles PATCH /api/v1/admin/users/{userId}/status
// Updates user account status
func adminUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "userId is required"})
		return
	}
	
	var req struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	
	// Validate status
	validStatuses := map[string]bool{"active": true, "suspended": true, "deactivated": true}
	if !validStatuses[req.Status] {
		respond(w, http.StatusBadRequest, map[string]any{"error": "status must be active, suspended, or deactivated"})
		return
	}
	
	if strings.TrimSpace(req.Reason) == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "reason is required"})
		return
	}
	
	// Get admin user ID from context (set by middleware)
	adminID := authedUserID(r)
	
	// Update user status
	if err := activeStore.UpdateUserStatus(userID, req.Status, adminID, req.Reason); err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	
	// Log audit event
	emitAudit(adminID, "user_status_change", "user:"+userID, fmt.Sprintf("status changed to %s: %s", req.Status, req.Reason))
	
	// TODO: If suspended, invalidate user sessions via Supabase Admin API
	
	respond(w, http.StatusOK, map[string]any{
		"userId": userID,
		"status": req.Status,
		"message": "User status updated successfully",
	})
}

// adminUserRoleHandler handles PATCH /api/v1/admin/users/{userId}/role
// Updates user role (admin/user)
func adminUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "userId is required"})
		return
	}
	
	var req struct {
		Role string `json:"role"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respond(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	
	// Validate role
	if req.Role != "admin" && req.Role != "user" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "role must be admin or user"})
		return
	}
	
	// Get admin user ID from context
	adminID := authedUserID(r)
	
	// Update user role
	if err := activeStore.UpdateUserRole(userID, req.Role, adminID); err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	
	// Log audit event
	action := "role_granted"
	if req.Role == "user" {
		action = "role_revoked"
	}
	emitAudit(adminID, action, "user:"+userID, fmt.Sprintf("role changed to %s", req.Role))
	
	respond(w, http.StatusOK, map[string]any{
		"userId": userID,
		"role":   req.Role,
		"message": "User role updated successfully",
	})
}

// adminUserResetPasswordHandler handles POST /api/v1/admin/users/{userId}/reset-password
// Triggers password reset email for user
func adminUserResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "userId is required"})
		return
	}
	
	// Get user
	user, ok := activeStore.FindByUserID(userID)
	if !ok {
		respond(w, http.StatusNotFound, map[string]any{"error": "user not found"})
		return
	}
	
	// Get admin user ID
	adminID := authedUserID(r)
	
	// Trigger password reset via Supabase
	if supabaseConfigured() {
		email := supabaseLoginEmail(user)
		_, _, err := supabasePost("/auth/v1/recover", map[string]any{"email": email})
		if err != nil {
			respond(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
			return
		}
	}
	
	// Log audit event
	emitAudit(adminID, "admin_password_reset", "user:"+userID, "admin triggered password reset")
	
	respond(w, http.StatusOK, map[string]any{
		"message": "Password reset email sent to user",
	})
}

// adminUserActivityHandler handles GET /api/v1/admin/users/{userId}/activity
// Returns user activity log
func adminUserActivityHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		respond(w, http.StatusBadRequest, map[string]any{"error": "userId is required"})
		return
	}
	
	// Get activity logs from store (last 90 days)
	since := time.Now().AddDate(0, 0, -90)
	activities, err := activeStore.GetUserActivity(userID, since)
	if err != nil {
		respond(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	
	respond(w, http.StatusOK, map[string]any{
		"activity": activities,
		"count":    len(activities),
	})
}

// Helper functions

func formatOptionalTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func formatOptionalDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}
