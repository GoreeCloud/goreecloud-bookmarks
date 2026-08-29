package syncstate

import "strings"

const (
	bookmarksItemsDataset = "bookmarks.items"
	maxSyncRecordIDBytes  = 512
)

func bookmarksItemsCapability() (Capability, bool) {
	for _, capability := range capabilities {
		if capability.Dataset == bookmarksItemsDataset {
			return capability, true
		}
	}
	return Capability{}, false
}

// validBookmarkEnvelope validates the direction-neutral negotiated envelope
// shape. Read/write/delete permissions are enforced at the operation boundary.
func validBookmarkEnvelope(envelope Envelope) bool {
	capability, ok := bookmarksItemsCapability()
	if !ok || envelope.Dataset != capability.Dataset ||
		envelope.SchemaVersion != capability.SchemaVersion ||
		envelope.RecordID == "" || len(envelope.RecordID) > maxSyncRecordIDBytes ||
		envelope.Revision == 0 || envelope.UpdatedAt.IsZero() ||
		strings.TrimSpace(envelope.OriginDevice) == "" {
		return false
	}
	// Privacy Shield data minimization: tombstones carry no application payload;
	// live records carry application state.
	return envelope.Deleted ? envelope.Payload == nil : envelope.Payload != nil
}

func canSubmitBookmarkEnvelope(envelope Envelope) bool {
	capability, ok := bookmarksItemsCapability()
	if !ok || !validBookmarkEnvelope(envelope) {
		return false
	}
	if envelope.Deleted {
		return capability.Delete
	}
	return capability.Write
}
