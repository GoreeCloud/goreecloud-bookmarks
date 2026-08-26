package syncstate

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"
)

func TestSignedBookmarkTombstoneIsPayloadFree(t *testing.T) {
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
	if !envelope.Deleted || envelope.Payload != nil || envelope.Dataset != "bookmarks.items" {
		t.Fatalf("unexpected tombstone: %+v", envelope)
	}
	if proof.DeviceID != "device-b" || proof.Signature == "" {
		t.Fatalf("unexpected proof: %+v", proof)
	}
}
