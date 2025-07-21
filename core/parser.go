package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	
	"github.com/evan/float-echo/float"
)

// ParsedEntry represents the result of parsing an entry
type ParsedEntry struct {
	Type     string
	Content  string
	Metadata map[string]string
}

// Parser handles parsing of FLOAT-style entries
type Parser struct {
	typeRegex *regexp.Regexp
	metaRegex *regexp.Regexp
}

// NewParser creates a new parser instance
func NewParser() *Parser {
	// Create regex for type detection (word followed by ::)
	typeRegex := regexp.MustCompile(`^(\w+)::(.*)`)
	
	// Create regex for inline metadata extraction
	metaRegex := regexp.MustCompile(`\b(\w+)::\s*([^\s,]+)`)
	
	return &Parser{
		typeRegex: typeRegex,
		metaRegex: metaRegex,
	}
}

// ParseEntry parses a raw input string into a structured entry
func (p *Parser) ParseEntry(input string) ParsedEntry {
	input = strings.TrimSpace(input)
	
	result := ParsedEntry{
		Type:     "log",
		Content:  input,
		Metadata: make(map[string]string),
	}
	
	// Remove leading bullet point if present
	cleanInput := input
	if strings.HasPrefix(input, "- ") {
		cleanInput = strings.TrimPrefix(input, "- ")
	}
	
	// Check for type prefix
	if matches := p.typeRegex.FindStringSubmatch(cleanInput); len(matches) == 3 {
		typePrefix := matches[1]
		content := strings.TrimSpace(matches[2])
		
		// Look up the type in our type map
		for _, entryType := range float.EntryTypes {
			if entryType.Name == typePrefix {
				result.Type = typePrefix
				result.Content = content
				break
			}
		}
	}
	
	// Extract inline metadata from content
	p.extractMetadata(&result)
	
	return result
}

// extractMetadata finds and extracts inline metadata markers
func (p *Parser) extractMetadata(entry *ParsedEntry) {
	// Common metadata patterns to extract
	metaPatterns := map[string]*regexp.Regexp{
		"mode":    regexp.MustCompile(`\bmode::\s*([^\s,]+)`),
		"project": regexp.MustCompile(`\bproject::\s*([^\s,]+)`),
		"bridge":  regexp.MustCompile(`\bbridge::\s*(CB-\d{8}-\d{4}-\w{4})`),
		"uid":     regexp.MustCompile(`\buid::\s*([^\s,]+)`),
		"priority": regexp.MustCompile(`\bpriority::\s*(high|medium|low)`),
	}
	
	for key, regex := range metaPatterns {
		if matches := regex.FindStringSubmatch(entry.Content); len(matches) > 1 {
			entry.Metadata[key] = matches[1]
		}
	}
	
	// Extract mentions (@username)
	mentionRegex := regexp.MustCompile(`@(\w+)`)
	mentions := mentionRegex.FindAllStringSubmatch(entry.Content, -1)
	if len(mentions) > 0 {
		var mentionList []string
		for _, m := range mentions {
			mentionList = append(mentionList, m[1])
		}
		entry.Metadata["mentions"] = strings.Join(mentionList, ",")
	}
	
	// Extract tags (#tag)
	tagRegex := regexp.MustCompile(`#(\w+)`)
	tags := tagRegex.FindAllStringSubmatch(entry.Content, -1)
	if len(tags) > 0 {
		var tagList []string
		for _, t := range tags {
			tagList = append(tagList, t[1])
		}
		entry.Metadata["tags"] = strings.Join(tagList, ",")
	}
	
	// Extract dates (YYYY-MM-DD)
	dateRegex := regexp.MustCompile(`\b(\d{4}-\d{2}-\d{2})\b`)
	if matches := dateRegex.FindStringSubmatch(entry.Content); len(matches) > 1 {
		entry.Metadata["date"] = matches[1]
	}
	
	// Extract URLs
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	if matches := urlRegex.FindStringSubmatch(entry.Content); len(matches) > 0 {
		entry.Metadata["url"] = matches[0]
	}
}

// ExtractBridgeCommand checks if the input is a bridge command
func (p *Parser) ExtractBridgeCommand(input string) (bool, string) {
	bridgeCreateRegex := regexp.MustCompile(`^bridge::create\s+(.+)`)
	bridgeRestoreRegex := regexp.MustCompile(`^bridge::restore\s+(.+)`)
	
	if matches := bridgeCreateRegex.FindStringSubmatch(input); len(matches) > 1 {
		return true, "create:" + matches[1]
	}
	
	if matches := bridgeRestoreRegex.FindStringSubmatch(input); len(matches) > 1 {
		return true, "restore:" + matches[1]
	}
	
	return false, ""
}

// GenerateBridgeID creates a new bridge ID in FLOAT format
func GenerateBridgeID() string {
	now := time.Now()
	return fmt.Sprintf("CB-%s-%s-%s",
		now.Format("20060102"),
		now.Format("1504"),
		randString(4))
}