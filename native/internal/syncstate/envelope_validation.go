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

// validBookmarkEnvelope enforces Bookmarks' negotiated Sync envelope boundary.
// Tombstones deliberately carry no application payload; live records must.
func validBookmarkEnvelope(envelope Envelope) bool {
	capability, ok := bookmarksItemsCapability()
	if !ok || envelope.Dataset != capability.Dataset ||
		envelope.SchemaVersion != capability.SchemaVersion ||
		envelope.RecordID == "" || len(envelope.RecordID) > maxSyncRecordIDBytes ||
		envelope.Revision == 0 || envelope.UpdatedAt.IsZero() ||
		strings.TrimSpace(envelope.OriginDevice) == "" {
		return false
	}
	if envelope.Deleted {
		return capability.Delete && envelope.Payload == nil
	}
	return capability.Write && envelope.Payload != nil
}
