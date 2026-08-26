package syncstate

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type doerFunc func(*http.Request) (*http.Response, error)

func (f doerFunc) Do(request *http.Request) (*http.Response, error) { return f(request) }

func TestSignedBookmarkEnvelopeAndSubmission(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	item := ItemRecord{
		ID: "bookmark-1", URL: "https://example.com", Title: "Example", Tags: []string{"test"},
		CreatedAt: time.Unix(100, 0).UTC(), UpdatedAt: time.Unix(101, 0).UTC(),
	}
	envelope, proof, err := SignedBookmarkEnvelope(item, 2, DeviceIdentity{DeviceID: "device-a", PublicKey: publicKey, PrivateKey: privateKey})
	if err != nil {
		t.Fatal(err)
	}
	if envelope.Dataset != "bookmarks.items" || envelope.RecordID != "bookmark-1" || proof.Signature == "" {
		t.Fatalf("unexpected signed bookmark: envelope=%+v proof=%+v", envelope, proof)
	}

	client := SubmissionClient{
		BaseURL: "https://sync.internal",
		Client: doerFunc(func(request *http.Request) (*http.Response, error) {
			if request.URL.Path != "/api/v1/sync/bookmarks/items" {
				t.Fatalf("unexpected path: %s", request.URL.Path)
			}
			var payload struct {
				Record Envelope    `json:"record"`
				Proof  RecordProof `json:"proof"`
			}
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			encoded, _ := json.Marshal(payload)
			for _, forbidden := range []string{"ownerId", "consentGranted", "purposeAllowed", "trusted"} {
				if strings.Contains(string(encoded), forbidden) {
					t.Fatalf("submission contains forbidden field %q: %s", forbidden, encoded)
				}
			}
			return &http.Response{StatusCode: http.StatusAccepted, Body: io.NopCloser(strings.NewReader("{}")), Header: make(http.Header)}, nil
		}),
	}
	if err := client.SubmitBookmark(context.Background(), envelope, proof); err != nil {
		t.Fatal(err)
	}
}
