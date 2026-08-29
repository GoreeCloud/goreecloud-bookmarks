package syncstate

import (
	"crypto/ed25519"
	"encoding/base64"
	"strings"
	"time"
)

// SignedBookmarkTombstone creates a payload-free deletion envelope. Owner
// identity remains in authenticated Sync/account context rather than payload.
func SignedBookmarkTombstone(recordID string, revision uint64, updatedAt time.Time, identity DeviceIdentity) (Envelope, RecordProof, error) {
	recordID = strings.TrimSpace(recordID)
	capability, ok := bookmarksItemsCapability()
	if !ok || !capability.Delete || recordID == "" || len(recordID) > maxSyncRecordIDBytes || revision == 0 || updatedAt.IsZero() {
		return Envelope{}, RecordProof{}, ErrInvalidBookmarkItem
	}
	if strings.TrimSpace(identity.DeviceID) == "" || len(identity.PublicKey) != ed25519.PublicKeySize || len(identity.PrivateKey) != ed25519.PrivateKeySize {
		return Envelope{}, RecordProof{}, ErrInvalidDeviceIdentity
	}
	envelope := Envelope{
		Dataset: capability.Dataset, SchemaVersion: capability.SchemaVersion, RecordID: recordID,
		Revision: revision, UpdatedAt: updatedAt.UTC(), OriginDevice: identity.DeviceID,
		Deleted: true,
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
