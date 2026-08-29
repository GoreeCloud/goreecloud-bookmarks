package syncstate

import (
	"crypto/ed25519"
	"crypto/rand"
	"strings"
	"testing"
	"time"
)

func TestSignedBookmarkTombstoneIsPayloadFree(t *testing.T) {
	capability, ok := bookmarksItemsCapability()
	if !ok {
		t.Fatal("bookmarks.items capability missing")
	}
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	envelope, proof, err := SignedBookmarkTombstone("bookmark-1", 4, time.Unix(200, 0).UTC(), DeviceIdentity{
		DeviceID: "device-b", PublicKey: publicKey, PrivateKey: privateKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !envelope.Deleted || envelope.Payload != nil || envelope.Dataset != capability.Dataset || envelope.SchemaVersion != capability.SchemaVersion {
		t.Fatalf("unexpected tombstone: %+v", envelope)
	}
	if proof.DeviceID != "device-b" || proof.Signature == "" {
		t.Fatalf("unexpected proof: %+v", proof)
	}
}

func TestSignedBookmarkTombstoneRejectsOversizedRecordID(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	if _, _, err := SignedBookmarkTombstone(strings.Repeat("b", maxSyncRecordIDBytes+1), 1, time.Unix(200, 0).UTC(), DeviceIdentity{
		DeviceID: "device-b", PublicKey: publicKey, PrivateKey: privateKey,
	}); err == nil {
		t.Fatal("oversized tombstone record ID must fail before signing")
	}
}
