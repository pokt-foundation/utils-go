package domain

import (
	"testing"
)

func testExtractDomain(t *testing.T) {
	// Valid URLs
	validTests := []struct {
		inputURL       string
		expectedDomain string
		expectedErr    error
	}{
		{"https://www.google.com", "google.com", nil},
		{"https://github.com/openai", "github.com", nil},
		{"http://localhost:8080", "localhost", nil},
		{"https://node123.node.com:443", "node.com", nil},
		{"https://node92.whiteholestake.io:443", "whiteholestake.io", nil},
		{"https://node123.nodegroup.node.com:443", "node.com", nil},
		{"https://node123.node-provider.com:443", "node-provider.com", nil},
		{"not a valid url", "", errInvalidURL},
		{"http://", "", errInvalidURL},
		{"", "", errInvalidURL},
	}
	for _, tc := range validTests {
		domain, err := ExtractDomain(tc.inputURL)
		if err != tc.expectedErr {
			t.Errorf("expected error: %v, got: %v", tc.expectedErr, err)
		}
		if domain != tc.expectedDomain {
			t.Errorf("extractDomain(%q) = %q, expected %q", tc.inputURL, domain, tc.expectedDomain)
		}
	}

	// Invalid URLs
	invalidTests := []string{}
	for _, test := range invalidTests {
		_, err := extractDomain(test)
		if err == nil {
			t.Errorf("extractDomain(%q) did not return an error", test)
		}
	}
}
