package syncstate

import "strings"

const (
	bookmarksItemsDataset       = "bookmarks.items"
	bookmarksItemsSchemaVersion = 1
	maxSyncRecordIDBytes        = 512
)

// validBookmarkEnvelope enforces Bookmarks' negotiated Sync envelope boundary.
// Tombstones deliberately carry no application payload; live records must.
func validBookmarkEnvelope(envelope Envelope) bool {
	if envelope.Dataset != bookmarksItemsDataset ||
		envelope.SchemaVersion != bookmarksItemsSchemaVersion ||
		envelope.RecordID == "" || len(envelope.RecordID) > maxSyncRecordIDBytes ||
		envelope.Revision == 0 || envelope.UpdatedAt.IsZero() ||
		strings.TrimSpace(envelope.OriginDevice) == "" {
		return false
	}
	if envelope.Deleted {
		return envelope.Payload == nil
	}
	return envelope.Payload != nil
}
