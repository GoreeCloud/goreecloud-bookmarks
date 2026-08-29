package syncstate

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
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
	capability, ok := bookmarksItemsCapability()
	if !ok {
		t.Fatal("bookmarks.items capability missing")
	}
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
	if envelope.Dataset != capability.Dataset || envelope.SchemaVersion != capability.SchemaVersion || envelope.RecordID != "bookmark-1" || proof.Signature == "" {
		t.Fatalf("unexpected signed bookmark: envelope=%+v proof=%+v", envelope, proof)
	}

	client := SubmissionClient{
		BaseURL:     "https://sync.internal",
		BearerToken: "session-token",
		Client: doerFunc(func(request *http.Request) (*http.Response, error) {
			if request.URL.Path != "/api/v1/sync/bookmarks/items" {
				t.Fatalf("unexpected path: %s", request.URL.Path)
			}
			if got := request.Header.Get("Authorization"); got != "Bearer session-token" {
				t.Fatalf("Authorization = %q", got)
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

func TestSignedBookmarkEnvelopeRejectsOversizedRecordID(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	item := ItemRecord{
		ID: strings.Repeat("b", maxSyncRecordIDBytes+1), URL: "https://example.com",
		CreatedAt: time.Unix(100, 0).UTC(), UpdatedAt: time.Unix(101, 0).UTC(),
	}
	if _, _, err := SignedBookmarkEnvelope(item, 1, DeviceIdentity{
		DeviceID: "device-a", PublicKey: publicKey, PrivateKey: privateKey,
	}); err == nil {
		t.Fatal("oversized record ID must fail before signing")
	}
}

func TestSubmitBookmarkRequiresBearerBeforeTransport(t *testing.T) {
	called := false
	client := SubmissionClient{
		BaseURL: "https://sync.internal",
		Client: doerFunc(func(*http.Request) (*http.Response, error) {
			called = true
			return nil, nil
		}),
	}
	if err := client.SubmitBookmark(context.Background(), Envelope{RecordID: "bookmark-1"}, RecordProof{}); err == nil {
		t.Fatal("missing bearer session must fail closed")
	}
	if called {
		t.Fatal("transport must not be called without an authenticated Sync session")
	}
}

func TestSubmitBookmarkRejectsNonconformingEnvelopeBeforeTransport(t *testing.T) {
	capability, ok := bookmarksItemsCapability()
	if !ok {
		t.Fatal("bookmarks.items capability missing")
	}
	base := Envelope{
		Dataset: capability.Dataset, SchemaVersion: capability.SchemaVersion,
		RecordID: "bookmark-1", Revision: 1, UpdatedAt: time.Unix(101, 0).UTC(),
		OriginDevice: "device-a", Payload: map[string]any{"url": "https://example.com"},
	}
	tests := []struct {
		name   string
		mutate func(*Envelope)
	}{
		{name: "unnegotiated schema", mutate: func(envelope *Envelope) { envelope.SchemaVersion++ }},
		{name: "tombstone payload", mutate: func(envelope *Envelope) { envelope.Deleted = true }},
		{name: "live record without payload", mutate: func(envelope *Envelope) { envelope.Payload = nil }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			envelope := base
			test.mutate(&envelope)
			called := false
			client := SubmissionClient{
				BaseURL: "https://sync.internal", BearerToken: "session-token",
				Client: doerFunc(func(*http.Request) (*http.Response, error) {
					called = true
					return nil, nil
				}),
			}
			if err := client.SubmitBookmark(context.Background(), envelope, RecordProof{}); err == nil {
				t.Fatal("nonconforming envelope must fail closed")
			}
			if called {
				t.Fatal("transport must not receive a nonconforming envelope")
			}
		})
	}
}

func TestSubmitBookmarkRejectsInvalidProofBeforeTransport(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	item := ItemRecord{
		ID: "bookmark-proof", URL: "https://example.com/proof", Title: "Proof",
		CreatedAt: time.Unix(100, 0).UTC(), UpdatedAt: time.Unix(101, 0).UTC(),
	}
	baseEnvelope, baseProof, err := SignedBookmarkEnvelope(item, 1, DeviceIdentity{
		DeviceID: "device-a", PublicKey: publicKey, PrivateKey: privateKey,
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		mutate func(*Envelope, *RecordProof)
	}{
		{name: "device mismatch", mutate: func(_ *Envelope, proof *RecordProof) { proof.DeviceID = "device-b" }},
		{name: "malformed public key", mutate: func(_ *Envelope, proof *RecordProof) { proof.PublicKey = "!" }},
		{name: "malformed signature", mutate: func(_ *Envelope, proof *RecordProof) { proof.Signature = "!" }},
		{name: "invalid signature", mutate: func(_ *Envelope, proof *RecordProof) {
			proof.Signature = base64.RawURLEncoding.EncodeToString(make([]byte, ed25519.SignatureSize))
		}},
		{name: "record changed after signing", mutate: func(envelope *Envelope, _ *RecordProof) { envelope.Revision++ }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			envelope := baseEnvelope
			proof := baseProof
			test.mutate(&envelope, &proof)
			called := false
			client := SubmissionClient{
				BaseURL: "https://sync.internal", BearerToken: "session-token",
				Client: doerFunc(func(*http.Request) (*http.Response, error) {
					called = true
					return nil, nil
				}),
			}
			if err := client.SubmitBookmark(context.Background(), envelope, proof); err == nil {
				t.Fatal("invalid record proof must fail closed")
			}
			if called {
				t.Fatal("transport must not receive an invalid record proof")
			}
		})
	}
}
