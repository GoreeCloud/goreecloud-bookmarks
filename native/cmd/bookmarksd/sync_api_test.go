package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSyncCapabilitiesExposeMetadataWithoutOwnerData(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/sync/capabilities", nil)
	server{}.syncCapabilities(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	body := response.Body.String()
	for _, required := range []string{
		`"application":"bookmarks"`,
		`"bookmarks.items"`,
		`"bookmarks.collections"`,
		`"bookmarks.assignments"`,
		`"records_exposed":false`,
		`"owner_ids_exposed":false`,
		`"credentials_exposed":false`,
		`"production_approved":false`,
	} {
		if !strings.Contains(body, required) {
			t.Fatalf("sync capability response missing %s: %s", required, body)
		}
	}
}
