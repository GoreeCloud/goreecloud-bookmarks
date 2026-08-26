package syncstate

// Capability describes a stable Bookmarks dataset exposed to GoreeCloud Sync.
// Internal persistence and owner-scoping remain implementation details of
// Bookmarks and are not coupled to the Sync transport.
type Capability struct {
	Dataset       string `json:"dataset"`
	Application   string `json:"application"`
	SchemaVersion int    `json:"schemaVersion"`
	Read          bool   `json:"read"`
	Write         bool   `json:"write"`
	Delete        bool   `json:"delete"`
}

var capabilities = []Capability{
	{Dataset: "bookmarks.items", Application: "bookmarks", SchemaVersion: 1, Read: true, Write: true, Delete: true},
	{Dataset: "bookmarks.collections", Application: "bookmarks", SchemaVersion: 1, Read: true, Write: true, Delete: true},
	{Dataset: "bookmarks.assignments", Application: "bookmarks", SchemaVersion: 1, Read: true, Write: true, Delete: true},
}

// Capabilities returns a defensive copy suitable for protocol advertisement.
func Capabilities() []Capability {
	return append([]Capability(nil), capabilities...)
}
