# FLOAT Echo

A terminal-based note-taking application built with Go and Bubble Tea, designed for rapid thought capture using the FLOAT methodology.

## Features

- **Fast capture**: Optimized for quick entry of thoughts and ideas
- **FLOAT type system**: Automatic recognition of entry types (ctx::, mode::, project::, etc.)
- **Multi-line support**: Toggle between single-line and multi-line input modes
- **Entry management**: Navigate, indent, and reorder entries
- **Persistence**: Automatically saves entries to `~/.float/echo.json`
- **Metadata extraction**: Automatically extracts mentions, tags, dates, and URLs

## Installation

```bash
go build -o float-echo
./float-echo
```

## Keyboard Shortcuts

- **Enter**: Submit entry (single-line mode) or new line (multi-line mode)
- **Ctrl+S**: Submit entry in any mode
- **Ctrl+E**: Toggle between single-line and multi-line modes
- **Alt+↑/↓**: Navigate through entries
- **Tab/Shift+Tab**: Indent/outdent selected entry
- **[/]**: Move selected entry up/down
- **Escape**: Clear selection
- **Ctrl+C**: Save and quit

## Entry Types

Prefix your entries with type markers for automatic categorization:

- `ctx::` - Context markers
- `bridge::` - Bridge references
- `mode::` - Cognitive state
- `project::` - Project association
- `todo::` - Task items
- `meeting::` - Meeting notes
- `error::` - Error tracking
- `idea::` - Ideas and inspiration

## File Format

Entries are saved in JSON format at `~/.float/echo.json`:

```json
{
  "entries": [{
    "id": "20250721-103000-a1b2",
    "content": "Working on new feature",
    "timestamp": "2025-07-21T10:30:00Z",
    "type": "ctx",
    "indent": 0,
    "metadata": {
      "project": "float-echo"
    }
  }],
  "session": {
    "created": "2025-07-21T10:00:00Z",
    "last_active": "2025-07-21T10:30:00Z",
    "tool": "float-echo"
  }
}
```

## Export

The tool can export entries to FLOAT-compatible markdown format for integration with other tools in the ecosystem.

## Development

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling

See [CLAUDE.md](CLAUDE.md) for development guidelines.