package syncstate

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestFetchBookmarksRejectsNonconformingEnvelopeShape(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "unnegotiated schema",
			body: `{"dataset":"bookmarks.items","count":1,"records":[{"dataset":"bookmarks.items","schemaVersion":2,"recordId":"bookmark-1","revision":1,"updatedAt":"2026-08-28T20:00:00Z","originDevice":"device-1","deleted":false,"payload":{"url":"https://one.invalid"}}]}`,
		},
		{
			name: "tombstone retains payload",
			body: `{"dataset":"bookmarks.items","count":1,"records":[{"dataset":"bookmarks.items","schemaVersion":1,"recordId":"bookmark-1","revision":1,"updatedAt":"2026-08-28T20:00:00Z","originDevice":"device-1","deleted":true,"payload":{"url":"https://must-not-survive.invalid"}}]}`,
		},
		{
			name: "live record missing payload",
			body: `{"dataset":"bookmarks.items","count":1,"records":[{"dataset":"bookmarks.items","schemaVersion":1,"recordId":"bookmark-1","revision":1,"updatedAt":"2026-08-28T20:00:00Z","originDevice":"device-1","deleted":false}]}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := RetrievalClient{
				BaseURL: "https://sync.invalid", BearerToken: "session-token",
				Client: retrievalDoerFunc(func(*http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(test.body)), Header: make(http.Header)}, nil
				}),
			}
			if _, err := client.FetchBookmarks(context.Background()); err == nil {
				t.Fatal("nonconforming retrieval envelope must fail closed")
			}
		})
	}
}

func TestFetchBookmarksAcceptsPayloadFreeTombstone(t *testing.T) {
	body := `{"dataset":"bookmarks.items","count":1,"records":[{"dataset":"bookmarks.items","schemaVersion":1,"recordId":"bookmark-1","revision":2,"updatedAt":"2026-08-28T20:00:00Z","originDevice":"device-1","deleted":true}]}`
	client := RetrievalClient{
		BaseURL: "https://sync.invalid", BearerToken: "session-token",
		Client: retrievalDoerFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
		}),
	}
	records, err := client.FetchBookmarks(context.Background())
	if err != nil {
		t.Fatalf("FetchBookmarks: %v", err)
	}
	if len(records) != 1 || !records[0].Deleted || records[0].Payload != nil {
		t.Fatalf("unexpected tombstone: %+v", records)
	}
}
