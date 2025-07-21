package float

// EntryType represents the different types of FLOAT entries
type EntryType struct {
	Name     string
	Prefix   string
	Priority int
	Color    string
}

// Common FLOAT entry types with their prefixes and priorities
var EntryTypes = []EntryType{
	// Core consciousness markers
	{Name: "ctx", Prefix: "ctx::", Priority: 100, Color: "#00ff00"},
	{Name: "bridge", Prefix: "bridge::", Priority: 95, Color: "#ff00ff"},
	{Name: "mode", Prefix: "mode::", Priority: 90, Color: "#00ffff"},
	
	// Project and organization
	{Name: "project", Prefix: "project::", Priority: 85, Color: "#ffff00"},
	{Name: "dispatch", Prefix: "dispatch::", Priority: 80, Color: "#ff8800"},
	{Name: "todo", Prefix: "todo::", Priority: 75, Color: "#ff0088"},
	
	// Communication and meetings
	{Name: "meeting", Prefix: "meeting::", Priority: 70, Color: "#8800ff"},
	{Name: "claude", Prefix: "claude::", Priority: 65, Color: "#0088ff"},
	{Name: "talk", Prefix: "talk::", Priority: 60, Color: "#88ff00"},
	
	// Technical markers
	{Name: "error", Prefix: "error::", Priority: 55, Color: "#ff0000"},
	{Name: "fix", Prefix: "fix::", Priority: 50, Color: "#00ff88"},
	{Name: "idea", Prefix: "idea::", Priority: 45, Color: "#ff88ff"},
	
	// Metadata markers
	{Name: "highlight", Prefix: "highlight::", Priority: 40, Color: "#ffff88"},
	{Name: "note", Prefix: "note::", Priority: 35, Color: "#88ffff"},
	{Name: "log", Prefix: "log::", Priority: 30, Color: "#aaaaaa"},
	
	// System markers
	{Name: "sys", Prefix: "sys::", Priority: 25, Color: "#888888"},
	{Name: "meta", Prefix: "meta::", Priority: 20, Color: "#666666"},
	
	// Default
	{Name: "text", Prefix: "", Priority: 0, Color: "#ffffff"},
}

// TypeMap provides quick lookup by prefix
var TypeMap = make(map[string]EntryType)

func init() {
	for _, t := range EntryTypes {
		if t.Prefix != "" {
			TypeMap[t.Prefix] = t
		}
	}
}

// Metadata keys commonly used in FLOAT
const (
	MetaMode      = "mode"
	MetaProject   = "project"
	MetaBridge    = "bridge"
	MetaPriority  = "priority"
	MetaTimestamp = "timestamp"
	MetaUID       = "uid"
	MetaSource    = "source"
)