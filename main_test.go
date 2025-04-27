package main

import (
	"net"
	"testing"
)

// MockGeoIPDB mocks the GeoIPLookup for tests
type MockGeoIPDB struct {
	FakeCountryCode string
	Err             error
}

func (m *MockGeoIPDB) LookupCountry(ip net.IP) (string, error) {
	return m.FakeCountryCode, m.Err
}

func TestValidateCountry_BlankInput(t *testing.T) {
	mock := &MockGeoIPDB{}
	info := ValidationInfo{}
	result := validate_country(info, mock)

	if !result.IsError {
		t.Errorf("Expected error for blank input")
	}
}

func TestValidateCountry_InvalidIp(t *testing.T) {
	mock := &MockGeoIPDB{}
	info := ValidationInfo{
		Ip:             "not-an-ip",
		CountryIsoCode: []string{"US"},
	}
	result := validate_country(info, mock)

	if !result.IsError {
		t.Errorf("Expected error for invalid IP")
	}
}

func TestValidateCountry_MatchingCountry(t *testing.T) {
	mock := &MockGeoIPDB{FakeCountryCode: "US"}
	info := ValidationInfo{
		Ip:             "8.8.8.8",
		CountryIsoCode: []string{"US"},
	}
	result := validate_country(info, mock)

	if result.IsError {
		t.Errorf("Did not expect error, got: %s", result.ErrorMessage)
	}
	if !result.IsValidCountry {
		t.Errorf("Expected valid country match")
	}
}

func TestValidateCountry_NoMatch(t *testing.T) {
	mock := &MockGeoIPDB{FakeCountryCode: "US"}
	info := ValidationInfo{
		Ip:             "8.8.8.8",
		CountryIsoCode: []string{"FR"},
	}
	result := validate_country(info, mock)

	if result.IsValidCountry {
		t.Errorf("Expected no valid country match")
	}
}
