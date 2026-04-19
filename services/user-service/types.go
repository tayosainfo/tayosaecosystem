package main

import "time"

type User struct {
	ID                  string     `json:"id"`
	FullName            string     `json:"fullName"`
	PhoneE164           string     `json:"phoneE164"`
	AuthEmail           string     `json:"-"`
	ContactEmail        string     `json:"contactEmail,omitempty"`
	Password            string     `json:"-"` // in-memory dev plain text (never persist)
	PasswordHash        string     `json:"-"` // bcrypt from Postgres
	InsforgeUserID      string     `json:"insforgeUserId,omitempty"`
	InsforgeEmail       string     `json:"-"`
	DateOfBirth         *time.Time `json:"-"` // surfaced in /users/me as dateOfBirth (YYYY-MM-DD)
	Nationality         string     `json:"nationality,omitempty"`
	PhoneVerifiedAt     time.Time  `json:"phoneVerifiedAt,omitempty"`
	ContactEmailChecked bool       `json:"contactEmailVerified"`
	CreatedAt           time.Time  `json:"createdAt"`
}

// GroupPolicyStats counts rows backing parish SACCO / village Kibiina policy.
type GroupPolicyStats struct {
	ParishSaccos         int `json:"parishSaccosRegistered"`
	VillageKibiinaGroups int `json:"villageKibiinaGroupsRegistered"`
}

type OnboardingProfile struct {
	UserID            string                 `json:"userId"`
	Phase             int                    `json:"phase"`
	KYC               map[string]any         `json:"kyc,omitempty"`
	Membership        map[string]any         `json:"membership,omitempty"`
	Kibiina           map[string]any         `json:"kibiina,omitempty"`
	ReferralCode      string                 `json:"referralCode,omitempty"`
	Geo               map[string]string      `json:"geo,omitempty"`
	TrustScoreSeed    int                    `json:"trustScoreSeed"`
	LastUpdatedAt     time.Time              `json:"lastUpdatedAt"`
	AdditionalContext map[string]interface{} `json:"additionalContext,omitempty"`
}

type geoRow struct {
	District  string
	County    string
	SubCounty string
	Parish    string
	Village   string
}
