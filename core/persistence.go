package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Persistence handles saving and loading data
type Persistence struct {
	filepath string
}

// NewPersistence creates a new persistence handler
func NewPersistence() (*Persistence, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home directory: %w", err)
	}
	
	floatDir := filepath.Join(home, ".float")
	if err := os.MkdirAll(floatDir, 0755); err != nil {
		return nil, fmt.Errorf("creating .float directory: %w", err)
	}
	
	return &Persistence{
		filepath: filepath.Join(floatDir, "echo.json"),
	}, nil
}

// Load reads the store from disk
func (p *Persistence) Load() (*Store, error) {
	file, err := os.Open(p.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty store if file doesn't exist
			return &Store{
				Entries: []Entry{},
				Session: NewSession(),
			}, nil
		}
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()
	
	var store Store
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&store); err != nil {
		return nil, fmt.Errorf("decoding JSON: %w", err)
	}
	
	// Update session last active time
	if store.Session != nil {
		store.Session.LastActive = time.Now()
	} else {
		store.Session = NewSession()
	}
	
	return &store, nil
}

// Save writes the store to disk atomically
func (p *Persistence) Save(store *Store) error {
	// Update session last active time
	if store.Session != nil {
		store.Session.LastActive = time.Now()
	}
	
	// Create temporary file
	tempFile := p.filepath + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	
	// Write JSON with indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(store); err != nil {
		file.Close()
		os.Remove(tempFile)
		return fmt.Errorf("encoding JSON: %w", err)
	}
	
	if err := file.Close(); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("closing temp file: %w", err)
	}
	
	// Atomically rename temp file to actual file
	if err := os.Rename(tempFile, p.filepath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("renaming file: %w", err)
	}
	
	return nil
}

// ExportMarkdown exports entries to a markdown file
func (p *Persistence) ExportMarkdown(store *Store, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating export file: %w", err)
	}
	defer file.Close()
	
	// Write YAML frontmatter
	fmt.Fprintf(file, "---\n")
	fmt.Fprintf(file, "created: %s\n", store.Session.Created.Format(time.RFC3339))
	fmt.Fprintf(file, "updated: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "tool: float-echo\n")
	fmt.Fprintf(file, "---\n\n")
	
	// Write entries grouped by date
	var currentDate string
	for _, entry := range store.Entries {
		date := entry.Timestamp.Format("2006-01-02")
		if date != currentDate {
			fmt.Fprintf(file, "\n## %s\n\n", date)
			currentDate = date
		}
		
		// Write timestamp and entry
		timestamp := entry.Timestamp.Format("15:04:05")
		indent := ""
		for i := 0; i < entry.Indent; i++ {
			indent += "  "
		}
		
		fmt.Fprintf(file, "%s**%s** ", indent, timestamp)
		if entry.Type != "log" && entry.Type != "text" {
			fmt.Fprintf(file, "`%s::` ", entry.Type)
		}
		fmt.Fprintf(file, "%s\n", entry.Content)
		
		// Write metadata if present
		if len(entry.Metadata) > 0 {
			for key, value := range entry.Metadata {
				fmt.Fprintf(file, "%s  - %s:: %s\n", indent, key, value)
			}
		}
	}
	
	return nil
}

// Backup creates a timestamped backup of the current store
func (p *Persistence) Backup() error {
	// Read current file
	source, err := os.Open(p.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Nothing to backup
		}
		return fmt.Errorf("opening source file: %w", err)
	}
	defer source.Close()
	
	// Create backup filename
	backupDir := filepath.Dir(p.filepath)
	backupFile := filepath.Join(backupDir, fmt.Sprintf("echo-backup-%s.json", time.Now().Format("20060102-150405")))
	
	// Create backup file
	dest, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("creating backup file: %w", err)
	}
	defer dest.Close()
	
	// Copy content
	if _, err := io.Copy(dest, source); err != nil {
		return fmt.Errorf("copying to backup: %w", err)
	}
	
	return nil
}