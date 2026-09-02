package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBookmarkCollectionReadsAndDeletesFailClosedWithoutAssignmentStore(t *testing.T) {
	app, bookmark, _ := newBookmarkCollectionFixture(t)
	app.assignments = nil

	for _, test := range []struct {
		name   string
		method string
		invoke func(http.ResponseWriter, *http.Request)
	}{
		{name: "get", method: http.MethodGet, invoke: app.getBookmarkCollection},
		{name: "delete", method: http.MethodDelete, invoke: app.removeBookmarkCollection},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/api/v1/bookmarks/"+bookmark.ID+"/collection", nil)
			request.SetPathValue("id", bookmark.ID)
			response := httptest.NewRecorder()
			test.invoke(response, request)
			if response.Code != http.StatusServiceUnavailable {
				t.Fatalf("assignment store unavailable status=%d body=%s", response.Code, response.Body.String())
			}
		})
	}
}
