package syncstate

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidDeviceIdentity = errors.New("invalid Bookmarks sync device identity")
	ErrSyncSubmissionFailed  = errors.New("bookmark sync submission failed")
)

type Envelope struct {
	Dataset       string         `json:"dataset"`
	SchemaVersion int            `json:"schemaVersion"`
	RecordID      string         `json:"recordId"`
	Revision      uint64         `json:"revision"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	OriginDevice  string         `json:"originDevice"`
	Deleted       bool           `json:"deleted"`
	Payload       map[string]any `json:"payload,omitempty"`
}

type RecordProof struct {
	DeviceID  string `json:"deviceId"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
}

type DeviceIdentity struct {
	DeviceID   string
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

func SignedBookmarkEnvelope(item ItemRecord, revision uint64, identity DeviceIdentity) (Envelope, RecordProof, error) {
	capability, ok := bookmarksItemsCapability()
	if !ok || !capability.Write || strings.TrimSpace(item.ID) == "" || len(item.ID) > maxSyncRecordIDBytes || strings.TrimSpace(item.URL) == "" || item.UpdatedAt.IsZero() || revision == 0 {
		return Envelope{}, RecordProof{}, ErrInvalidBookmarkItem
	}
	if strings.TrimSpace(identity.DeviceID) == "" || len(identity.PublicKey) != ed25519.PublicKeySize || len(identity.PrivateKey) != ed25519.PrivateKeySize {
		return Envelope{}, RecordProof{}, ErrInvalidDeviceIdentity
	}
	payload := map[string]any{
		"id": item.ID, "url": item.URL, "title": item.Title, "note": item.Note,
		"tags": append([]string(nil), item.Tags...),
		"createdAt": item.CreatedAt.UTC().Format(time.RFC3339Nano),
		"updatedAt": item.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	envelope := Envelope{
		Dataset: capability.Dataset, SchemaVersion: capability.SchemaVersion, RecordID: item.ID,
		Revision: revision, UpdatedAt: item.UpdatedAt.UTC(), OriginDevice: identity.DeviceID,
		Payload: payload,
	}
	if !validBookmarkEnvelope(envelope) {
		return Envelope{}, RecordProof{}, ErrInvalidBookmarkItem
	}
	message, err := proofMessage(envelope)
	if err != nil {
		return Envelope{}, RecordProof{}, err
	}
	proof := RecordProof{
		DeviceID: identity.DeviceID,
		PublicKey: base64.RawURLEncoding.EncodeToString(identity.PublicKey),
		Signature: base64.RawURLEncoding.EncodeToString(ed25519.Sign(identity.PrivateKey, message)),
	}
	return envelope, proof, nil
}

type HTTPDoer interface{ Do(*http.Request) (*http.Response, error) }

type SubmissionClient struct {
	BaseURL     string
	BearerToken string
	Client      HTTPDoer
}

func (c SubmissionClient) SubmitBookmark(ctx context.Context, envelope Envelope, proof RecordProof) error {
	token := strings.TrimSpace(c.BearerToken)
	if strings.TrimSpace(c.BaseURL) == "" || token == "" || c.Client == nil || !validBookmarkEnvelope(envelope) {
		return ErrSyncSubmissionFailed
	}
	body, err := json.Marshal(struct {
		Record Envelope    `json:"record"`
		Proof  RecordProof `json:"proof"`
	}{Record: envelope, Proof: proof})
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.BaseURL, "/")+"/api/v1/sync/bookmarks/items", bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := c.Client.Do(request)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSyncSubmissionFailed, err)
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 64<<10))
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("%w: status %d", ErrSyncSubmissionFailed, response.StatusCode)
	}
	return nil
}

func proofMessage(record Envelope) ([]byte, error) {
	payload, err := json.Marshal(record.Payload)
	if err != nil {
		return nil, err
	}
	digest := sha256.Sum256(payload)
	return []byte(fmt.Sprintf(
		"GC-SYNC-RECORD/1\n%s\n%d\n%s\n%d\n%s\n%s\n%t\n%s",
		record.Dataset, record.SchemaVersion, record.RecordID, record.Revision,
		record.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"),
		record.OriginDevice, record.Deleted, hex.EncodeToString(digest[:]),
	)), nil
}
