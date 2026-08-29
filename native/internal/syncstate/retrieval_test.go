package syncstate

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type retrievalDoerFunc func(*http.Request) (*http.Response, error)

func (f retrievalDoerFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

func TestFetchBookmarksUsesAuthenticatedSessionAndValidatesDataset(t *testing.T) {
	client := RetrievalClient{
		BaseURL:     "https://sync.invalid/",
		BearerToken: "fixture",
		Client: retrievalDoerFunc(func(request *http.Request) (*http.Response, error) {
			if request.Method != http.MethodGet || request.URL.Path != "/api/v1/sync/bookmarks/items" {
				t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
			}
			if request.URL.Query().Get("limit") != "256" || request.URL.Query().Get("after") != "" {
				t.Fatalf("unexpected pagination query: %s", request.URL.RawQuery)
			}
			if got := request.Header.Get("Authorization"); got != "Bearer fixture" {
				t.Fatalf("unexpected authorization header")
			}
			body := `{"dataset":"bookmarks.items","count":1,"records":[{"dataset":"bookmarks.items","schemaVersion":1,"recordId":"bookmark-1","revision":3,"updatedAt":"2026-08-26T23:35:00Z","originDevice":"device-1","deleted":false,"payload":{"url":"https://example.invalid"}}]}`
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
		}),
	}

	records, err := client.FetchBookmarks(context.Background())
	if err != nil {
		t.Fatalf("FetchBookmarks: %v", err)
	}
	if len(records) != 1 || records[0].RecordID != "bookmark-1" || records[0].Revision != 3 {
		t.Fatalf("unexpected records: %+v", records)
	}
}

func TestFetchBookmarksFollowsStrictlyAdvancingPages(t *testing.T) {
	calls := 0
	client := RetrievalClient{
		BaseURL:     "https://sync.invalid",
		BearerToken: "fixture",
		Client: retrievalDoerFunc(func(request *http.Request) (*http.Response, error) {
			calls++
			after := request.URL.Query().Get("after")
			var body string
			switch calls {
			case 1:
				if after != "" {
					t.Fatalf("first after = %q", after)
				}
				body = `{"dataset":"bookmarks.items","count":1,"records":[{"dataset":"bookmarks.items","schemaVersion":1,"recordId":"bookmark-1","revision":1,"updatedAt":"2026-08-26T23:35:00Z","originDevice":"device-1","deleted":false,"payload":{"url":"https://one.invalid"}}],"nextAfter":"bookmark-1"}`
			case 2:
				if after != "bookmark-1" {
					t.Fatalf("second after = %q", after)
				}
				body = `{"dataset":"bookmarks.items","count":1,"records":[{"dataset":"bookmarks.items","schemaVersion":1,"recordId":"bookmark-2","revision":1,"updatedAt":"2026-08-26T23:36:00Z","originDevice":"device-1","deleted":false,"payload":{"url":"https://two.invalid"}}]}`
			default:
				t.Fatalf("unexpected extra request")
			}
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
		}),
	}

	records, err := client.FetchBookmarks(context.Background())
	if err != nil {
		t.Fatalf("FetchBookmarks: %v", err)
	}
	if calls != 2 || len(records) != 2 || records[0].RecordID != "bookmark-1" || records[1].RecordID != "bookmark-2" {
		t.Fatalf("unexpected pagination result: calls=%d records=%+v", calls, records)
	}
}

func TestFetchBookmarksRejectsRepeatingContinuation(t *testing.T) {
	calls := 0
	client := RetrievalClient{
		BaseURL:     "https://sync.invalid",
		BearerToken: "fixture",
		Client: retrievalDoerFunc(func(*http.Request) (*http.Response, error) {
			calls++
			body := `{"dataset":"bookmarks.items","count":1,"records":[{"dataset":"bookmarks.items","schemaVersion":1,"recordId":"bookmark-1","revision":1,"updatedAt":"2026-08-26T23:35:00Z","originDevice":"device-1","deleted":false,"payload":{"url":"https://one.invalid"}}],"nextAfter":"bookmark-1"}`
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
		}),
	}

	if _, err := client.FetchBookmarks(context.Background()); err == nil {
		t.Fatal("repeating continuation must fail closed")
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
}

func TestFetchBookmarksRequiresBearerBeforeTransport(t *testing.T) {
	called := false
	client := RetrievalClient{
		BaseURL: "https://sync.invalid",
		Client: retrievalDoerFunc(func(*http.Request) (*http.Response, error) {
			called = true
			return nil, nil
		}),
	}

	if _, err := client.FetchBookmarks(context.Background()); err == nil {
		t.Fatal("missing bearer session must fail closed")
	}
	if called {
		t.Fatal("transport must not be called without an authenticated Sync session")
	}
}

func TestFetchBookmarksRejectsCrossDatasetResponse(t *testing.T) {
	client := RetrievalClient{
		BaseURL:     "https://sync.invalid",
		BearerToken: "fixture",
		Client: retrievalDoerFunc(func(*http.Request) (*http.Response, error) {
			body := `{"dataset":"bookmarks.items","count":1,"records":[{"dataset":"search.history","schemaVersion":1,"recordId":"query-1","revision":1,"updatedAt":"2026-08-26T23:35:00Z","originDevice":"device-1","deleted":false,"payload":{}}]}`
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
		}),
	}

	if _, err := client.FetchBookmarks(context.Background()); err == nil {
		t.Fatal("cross-dataset retrieval response must fail closed")
	}
}

func TestFetchBookmarksRejectsOversizedRecordID(t *testing.T) {
	oversized := strings.Repeat("b", maxSyncRecordIDBytes+1)
	client := RetrievalClient{
		BaseURL:     "https://sync.invalid",
		BearerToken: "fixture",
		Client: retrievalDoerFunc(func(*http.Request) (*http.Response, error) {
			body := `{"dataset":"bookmarks.items","count":1,"records":[{"dataset":"bookmarks.items","schemaVersion":1,"recordId":"` + oversized + `","revision":1,"updatedAt":"2026-08-26T23:35:00Z","originDevice":"device-1","deleted":false,"payload":{"url":"https://one.invalid"}}],"nextAfter":"` + oversized + `"}`
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
		}),
	}

	if _, err := client.FetchBookmarks(context.Background()); err == nil {
		t.Fatal("oversized record ID and continuation must fail closed")
	}
}

func TestFetchBookmarksRejectsTrailingJSON(t *testing.T) {
	client := RetrievalClient{
		BaseURL:     "https://sync.invalid",
		BearerToken: "fixture",
		Client: retrievalDoerFunc(func(*http.Request) (*http.Response, error) {
			body := `{"dataset":"bookmarks.items","count":0,"records":[]} {}`
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
		}),
	}

	if _, err := client.FetchBookmarks(context.Background()); err == nil {
		t.Fatal("trailing JSON document must fail closed")
	}
}

func TestFetchBookmarksRejectsOversizedResponse(t *testing.T) {
	client := RetrievalClient{
		BaseURL:     "https://sync.invalid",
		BearerToken: "fixture",
		Client: retrievalDoerFunc(func(*http.Request) (*http.Response, error) {
			body := strings.Repeat(" ", maxRetrievalBodyBytes+1)
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
		}),
	}

	if _, err := client.FetchBookmarks(context.Background()); err == nil {
		t.Fatal("oversized retrieval response must fail closed")
	}
}
