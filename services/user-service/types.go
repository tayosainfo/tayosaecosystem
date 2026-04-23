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
	SupabaseUserID      string     `json:"supabaseUserId,omitempty"`
	SupabaseLoginEmail  string     `json:"-"`
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

type UserConsents struct {
	UserID             string     `json:"userId"`
	TermsAcceptedAt    *time.Time `json:"termsAcceptedAt,omitempty"`
	PrivacyAcceptedAt  *time.Time `json:"privacyAcceptedAt,omitempty"`
	TermsVersion       string     `json:"termsVersion,omitempty"`
	PrivacyVersion     string     `json:"privacyVersion,omitempty"`
	LastUpdatedAt      time.Time  `json:"lastUpdatedAt"`
}

type KYCDocument struct {
	DocType    string    `json:"docType"`
	DocSide    string    `json:"docSide,omitempty"`
	StorageKey string    `json:"storageKey"`
	CreatedAt  time.Time `json:"createdAt"`
}

type KYCProfile struct {
	UserID                     string     `json:"userId"`
	Status                     string     `json:"status"`
	DateOfBirth                *time.Time `json:"dateOfBirth,omitempty"`
	Gender                     string     `json:"gender,omitempty"`
	Nationality                string     `json:"nationality,omitempty"`
	OccupationStatus           string     `json:"occupationStatus,omitempty"`
	IDType                     string     `json:"idType,omitempty"`
	IDNumber                   string     `json:"idNumber,omitempty"`
	NOKFullName                string     `json:"nokFullName,omitempty"`
	NOKRelationship            string     `json:"nokRelationship,omitempty"`
	NOKPhone                   string     `json:"nokPhone,omitempty"`
	NOKEmail                   string     `json:"nokEmail,omitempty"`
	SourceOfFunds              string     `json:"sourceOfFunds,omitempty"`
	PEPStatus                  *bool      `json:"pepStatus,omitempty"`
	SACCOMembershipDisclosures string     `json:"saccoMembershipDisclosures,omitempty"`
	SubmittedAt                *time.Time `json:"submittedAt,omitempty"`
	ReviewedAt                 *time.Time `json:"reviewedAt,omitempty"`
	ReviewNote                 string     `json:"reviewNote,omitempty"`
	ReviewedBy                 string     `json:"reviewedBy,omitempty"`
	LastUpdatedAt              time.Time  `json:"lastUpdatedAt"`
}

type SaccoMembership struct {
	UserID                   string    `json:"userId"`
	Status                   string    `json:"status"`
	District                 string    `json:"district,omitempty"`
	County                   string    `json:"county,omitempty"`
	SubCounty                string    `json:"subCounty,omitempty"`
	Parish                   string    `json:"parish,omitempty"`
	Village                  string    `json:"village,omitempty"`
	StreetPlot               string    `json:"streetPlot,omitempty"`
	MobileMoneyProvider      string    `json:"mobileMoneyProvider,omitempty"`
	MobileMoneyNumber        string    `json:"mobileMoneyNumber,omitempty"`
	SecondaryMoMoNumber      string    `json:"secondaryMoMoNumber,omitempty"`
	ContributionFrequency    string    `json:"contributionFrequency,omitempty"`
	SavingsGoalAmount        float64   `json:"savingsGoalAmount,omitempty"`
	SavingsGoalPurpose       string    `json:"savingsGoalPurpose,omitempty"`
	SharesToPurchase         int       `json:"sharesToPurchase,omitempty"`
	EntranceFeePaymentMethod string    `json:"entranceFeePaymentMethod,omitempty"`
	LastUpdatedAt            time.Time `json:"lastUpdatedAt"`
}

type KibiinaPreference struct {
	UserID                 string    `json:"userId"`
	Action                 string    `json:"action"`
	InviteCode             string    `json:"inviteCode,omitempty"`
	GroupName              string    `json:"groupName,omitempty"`
	ContributionAmount     float64   `json:"contributionAmount,omitempty"`
	CycleFrequency         string    `json:"cycleFrequency,omitempty"`
	MaxGroupSize           int       `json:"maxGroupSize,omitempty"`
	PayoutOrderPreference  string    `json:"payoutOrderPreference,omitempty"`
	NotificationPreference string    `json:"notificationPreference,omitempty"`
	LanguagePreference     string    `json:"languagePreference,omitempty"`
	LastUpdatedAt          time.Time `json:"lastUpdatedAt"`
}

type AdminKYCQueueItem struct {
	UserID       string     `json:"userId"`
	FullName     string     `json:"fullName,omitempty"`
	PhoneE164    string     `json:"phoneE164,omitempty"`
	ContactEmail string     `json:"contactEmail,omitempty"`
	Status       string     `json:"status"`
	IDType       string     `json:"idType,omitempty"`
	IDNumber     string     `json:"idNumber,omitempty"`
	SubmittedAt  *time.Time `json:"submittedAt,omitempty"`
	ReviewedAt   *time.Time `json:"reviewedAt,omitempty"`
}

type geoRow struct {
	District  string
	County    string
	SubCounty string
	Parish    string
	Village   string
}
