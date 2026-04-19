package main

import (
	"strings"
	"testing"
)

func TestParseOptionalDateOfBirth(t *testing.T) {
	p, err := parseOptionalDateOfBirth("")
	if err != nil || p != nil {
		t.Fatalf("empty: got %v err=%v", p, err)
	}
	if _, err := parseOptionalDateOfBirth("99-99-99"); err == nil {
		t.Fatal("expected error for bad date")
	}
	p, err = parseOptionalDateOfBirth("2001-06-15")
	if err != nil || p == nil || p.Format("2006-01-02") != "2001-06-15" {
		t.Fatalf("valid: got %v err=%v", p, err)
	}
}

func TestNormalizeNationality(t *testing.T) {
	s, err := normalizeNationality("")
	if err != nil || s != "" {
		t.Fatalf("empty: %q err=%v", s, err)
	}
	s, err = normalizeNationality("  Uganda  ")
	if err != nil || s != "Uganda" {
		t.Fatalf("trim: %q err=%v", s, err)
	}
	if _, err := normalizeNationality(strings.Repeat("x", 65)); err == nil {
		t.Fatal("expected error for long nationality")
	}
}

func TestMemoryGeoRecordExistsFromCSV(t *testing.T) {
	m := NewMemoryStore()
	if err := m.EnsureGeoSeeded(); err != nil {
		t.Fatalf("seed: %v", err)
	}
	ok, err := m.GeoRecordExists("Abim", "Labwor County", "Abim", "Abongepach", "VillageAbongepach")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected known CSV row to match")
	}
	ok, err = m.GeoRecordExists("Abim", "Labwor County", "Abim", "Abongepach", "NonexistentVillage999")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected bogus village to miss")
	}
}

func TestMemoryGroupPolicyStats(t *testing.T) {
	m := NewMemoryStore()
	st, err := m.GroupPolicyStats()
	if err != nil {
		t.Fatal(err)
	}
	if st.ParishSaccos != 0 || st.VillageKibiinaGroups != 0 {
		t.Fatalf("memory store has no parish tables: %+v", st)
	}
}
