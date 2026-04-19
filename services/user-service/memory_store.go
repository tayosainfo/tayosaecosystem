package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// MemoryStore is the dev / no-DATABASE_URL backend.
type MemoryStore struct {
	mu              sync.RWMutex
	usersByPhone    map[string]User
	usersByEmail    map[string]User
	usersByAuthMail map[string]User
	usersByID       map[string]User
	onboarding      map[string]OnboardingProfile
	geoRows         []geoRow
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		usersByPhone:    map[string]User{},
		usersByEmail:    map[string]User{},
		usersByAuthMail: map[string]User{},
		usersByID:       map[string]User{},
		onboarding:      map[string]OnboardingProfile{},
	}
}

func (m *MemoryStore) Close() {}

func (m *MemoryStore) Ping() error { return nil }

func (m *MemoryStore) indexUser(u User) {
	m.usersByPhone[u.PhoneE164] = u
	m.usersByAuthMail[u.AuthEmail] = u
	m.usersByID[u.ID] = u
	if u.ContactEmail != "" {
		m.usersByEmail[u.ContactEmail] = u
	}
	if u.InsforgeEmail != "" {
		m.usersByEmail[u.InsforgeEmail] = u
	}
}

func (m *MemoryStore) FindByPhone(phoneE164 string) (User, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.usersByPhone[phoneE164]
	return u, ok
}

func (m *MemoryStore) FindByEmailKey(normalizedEmail string) (User, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if u, ok := m.usersByEmail[normalizedEmail]; ok {
		return u, true
	}
	if u, ok := m.usersByAuthMail[normalizedEmail]; ok {
		return u, true
	}
	return User{}, false
}

func (m *MemoryStore) FindByUserID(userID string) (User, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.usersByID[userID]
	return u, ok
}

func (m *MemoryStore) CreateIdentityWithOnboarding(u User, ob OnboardingProfile, passwordHash *string) error {
	if passwordHash != nil {
		u.PasswordHash = *passwordHash
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.indexUser(u)
	m.onboarding[ob.UserID] = ob
	return nil
}

func (m *MemoryStore) UpdateIdentity(u User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.indexUser(u)
	return nil
}

func (m *MemoryStore) SetContactEmailVerified(phoneE164 string, verified bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.usersByPhone[phoneE164]
	if !ok {
		return nil
	}
	u.ContactEmailChecked = verified
	m.indexUser(u)
	return nil
}

func (m *MemoryStore) UpsertOnboarding(p OnboardingProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p.LastUpdatedAt = time.Now()
	m.onboarding[p.UserID] = p
	return nil
}

func (m *MemoryStore) GetOnboarding(userID string) (OnboardingProfile, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.onboarding[userID]
	return v, ok
}

func (m *MemoryStore) GeoDistinct(level, parent string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return geoDistinctFromRows(m.geoRows, level, parent), nil
}

func (m *MemoryStore) GeoRecordExists(district, county, subCounty, parish, village string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, row := range m.geoRows {
		if strings.EqualFold(strings.TrimSpace(row.District), strings.TrimSpace(district)) &&
			strings.EqualFold(strings.TrimSpace(row.County), strings.TrimSpace(county)) &&
			strings.EqualFold(strings.TrimSpace(row.SubCounty), strings.TrimSpace(subCounty)) &&
			strings.EqualFold(strings.TrimSpace(row.Parish), strings.TrimSpace(parish)) &&
			strings.EqualFold(strings.TrimSpace(row.Village), strings.TrimSpace(village)) {
			return true, nil
		}
	}
	return false, nil
}

func (m *MemoryStore) GroupPolicyStats() (GroupPolicyStats, error) {
	return GroupPolicyStats{}, nil
}

func (m *MemoryStore) EnsureGeoSeeded() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.geoRows) > 0 {
		return nil
	}
	csvPath := filepath.Join("..", "..", "uganda_geo_data_2025-11-26.csv")
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil || len(rows) <= 1 {
		return err
	}
	for _, row := range rows[1:] {
		if len(row) < 5 {
			continue
		}
		m.geoRows = append(m.geoRows, geoRow{
			District:  strings.TrimSpace(row[0]),
			County:    strings.TrimSpace(row[1]),
			SubCounty: strings.TrimSpace(row[2]),
			Parish:    strings.TrimSpace(row[3]),
			Village:   strings.TrimSpace(row[4]),
		})
	}
	return nil
}

func geoDistinctFromRows(geoRows []geoRow, level, parent string) []string {
	parent = strings.ToLower(strings.TrimSpace(parent))
	unique := map[string]bool{}
	var values []string
	for _, row := range geoRows {
		var candidate string
		switch strings.ToLower(strings.TrimSpace(level)) {
		case "district":
			candidate = row.District
		case "county":
			if parent != "" && strings.ToLower(row.District) != parent {
				continue
			}
			candidate = row.County
		case "subcounty":
			if parent != "" && strings.ToLower(row.County) != parent {
				continue
			}
			candidate = row.SubCounty
		case "parish":
			if parent != "" && strings.ToLower(row.SubCounty) != parent {
				continue
			}
			candidate = row.Parish
		case "village":
			if parent != "" && strings.ToLower(row.Parish) != parent {
				continue
			}
			candidate = row.Village
		default:
			return nil
		}
		if candidate != "" && !unique[candidate] {
			unique[candidate] = true
			values = append(values, candidate)
		}
	}
	return values
}
