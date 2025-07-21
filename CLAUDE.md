# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the echo-refactor tool within the FLOAT Workspace ecosystem. It's a terminal-based note-taking application built with Go and the Bubble Tea framework, designed for rapid thought capture using the FLOAT methodology.

## Key Requirements

Based on app.md specifications:
- **No edit mode** - only log mode with multi-line support
- **Quick-log mode**: Single-line input, Enter to send
- **Multi-line mode**: Using TextArea component, Ctrl+S to send
- **Markdown support**: In multi-line mode via glamour
- **Persistence**: Save entries to file (similar to float-go-echorefactor)

## Commands

```bash
# Build the application
go build -o float-echo

# Run the application
go run main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code (if golangci-lint is installed)
golangci-lint run
```

## Architecture Guidelines

### Component Structure
- Use Bubble Tea's Model-Update-View pattern
- Leverage existing bubbles components (especially textarea)
- Keep the UI minimal and keyboard-driven

### Persistence Strategy
- Save to `.float/echo.md` or similar location
- Append entries with timestamps
- Use FLOAT metadata format (YAML frontmatter for sessions, inline :: markers)
- Consider the existing float-go-echorefactor persistence pattern

### FLOAT Integration
- Support `::` notation for metadata (e.g., `ctx::`, `mode::`, `type::`)
- Entries should follow FLOAT patterns but don't need complex parsing
- Focus on capture speed over feature complexity
- Integrate with FLOAT dropzone patterns if needed

### Key Dependencies
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - Pre-built components
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/glamour` - Markdown rendering

## Development Workflow

1. Reference the existing float-go-echorefactor implementation for patterns
2. Keep the interface minimal - no edit mode, just capture
3. Use TextArea from bubbles for multi-line input
4. Implement mode switching between quick-log and multi-line
5. Ensure Ctrl+S works consistently in multi-line mode
6. Implement file persistence with proper FLOAT formatting

## Important Context

- This is part of the larger FLOAT ecosystem for consciousness technology
- The tool prioritizes rapid capture over editing capabilities
- Terminal UI should be clean and distraction-free
- Persistence should follow FLOAT conventions for integration with other tools

## Testing Approach

- Test the Model's Update function for state transitions
- Verify keyboard shortcuts work as expected
- Ensure mode switching maintains user input
- Test persistence layer (file writes, formatting)
- Test edge cases like empty submissions

## File Format

Entries should be saved in FLOAT-compatible markdown format:
```markdown
---
created: 2025-07-21T10:30:00.000000Z
tool: float-echo
---

## 2025-07-21 10:30:00

Entry content here...
- ctx:: working on echo refactor
- mode:: development

## 2025-07-21 10:31:00

Another entry...
```

When implementing, focus on creating a frictionless capture experience that aligns with FLOAT's philosophy of preserving authentic thought streams while ensuring data persistence for later processing.