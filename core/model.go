package core

import (
	"time"
)

// Entry represents a single log entry in the system
type Entry struct {
	ID        string            `json:"id"`
	Content   string            `json:"content"`
	Timestamp time.Time         `json:"timestamp"`
	UpdatedAt time.Time         `json:"updated_at"`
	Type      string            `json:"type"`
	Indent    int               `json:"indent"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Session represents a FLOAT echo session
type Session struct {
	Created    time.Time `json:"created"`
	LastActive time.Time `json:"last_active"`
	Tool       string    `json:"tool"`
}

// Store represents the persistent storage structure
type Store struct {
	Entries []Entry  `json:"entries"`
	Session *Session `json:"session"`
}

// NewEntry creates a new entry with the given content
func NewEntry(content string) *Entry {
	now := time.Now()
	return &Entry{
		ID:        generateID(),
		Content:   content,
		Timestamp: now,
		UpdatedAt: now,
		Type:      "log",
		Indent:    0,
		Metadata:  make(map[string]string),
	}
}

// NewSession creates a new session
func NewSession() *Session {
	now := time.Now()
	return &Session{
		Created:    now,
		LastActive: now,
		Tool:       "float-echo",
	}
}

// generateID creates a unique identifier for entries
func generateID() string {
	return time.Now().Format("20060102-150405") + "-" + randString(4)
}

// Simple random string generator for IDs
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}