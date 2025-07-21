package ui

import (
	"fmt"
	"strings"
	"time"
	
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	
	"github.com/evan/float-echo/core"
)

// InputMode represents the current input mode
type InputMode int

const (
	SingleLineMode InputMode = iota
	MultiLineMode
)

// Model represents the application state
type Model struct {
	// Core data
	store       *core.Store
	persistence *core.Persistence
	parser      *core.Parser
	
	// UI components
	viewport viewport.Model
	input    textarea.Model
	
	// UI state
	width         int
	height        int
	selectedIndex int
	inputMode     InputMode
	err           error
	message       string
	messageTimer  *time.Timer
}

// NewModel creates a new application model
func NewModel() (*Model, error) {
	// Initialize persistence
	persistence, err := core.NewPersistence()
	if err != nil {
		return nil, fmt.Errorf("creating persistence: %w", err)
	}
	
	// Load existing data
	store, err := persistence.Load()
	if err != nil {
		return nil, fmt.Errorf("loading store: %w", err)
	}
	
	// Create parser
	parser := core.NewParser()
	
	// Create input textarea
	ta := textarea.New()
	ta.Placeholder = "Type your entry..."
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.SetHeight(1)
	ta.Focus()
	
	// Create viewport for entries
	vp := viewport.New(80, 20)
	
	m := &Model{
		store:         store,
		persistence:   persistence,
		parser:        parser,
		viewport:      vp,
		input:         ta,
		selectedIndex: -1,
		inputMode:     SingleLineMode,
	}
	
	return m, nil
}

// Init implements tea.Model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.updateViewport(),
	)
}

// Update implements tea.Model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update viewport size
		headerHeight := 3
		inputHeight := 3
		statusHeight := 1
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - headerHeight - inputHeight - statusHeight
		
		// Update input width
		m.input.SetWidth(msg.Width - 4)
		
		cmds = append(cmds, m.updateViewport())
		
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			// Save before quitting
			m.persistence.Save(m.store)
			return m, tea.Quit
			
		case tea.KeyCtrlS:
			// Submit entry in any mode
			cmds = append(cmds, m.submitEntry())
			
		case tea.KeyEnter:
			if m.inputMode == SingleLineMode {
				// Submit in single-line mode
				cmds = append(cmds, m.submitEntry())
			} else {
				// Just add newline in multi-line mode
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				cmds = append(cmds, cmd)
			}
			
		case tea.KeyCtrlE:
			// Toggle input mode (Ctrl+E for "Edit mode")
			m.toggleInputMode()
			
		case tea.KeyTab:
			if m.selectedIndex >= 0 && m.selectedIndex < len(m.store.Entries) {
				// Indent selected entry
				m.store.Entries[m.selectedIndex].Indent++
				m.save()
				cmds = append(cmds, m.updateViewport())
			}
			
		case tea.KeyShiftTab:
			if m.selectedIndex >= 0 && m.selectedIndex < len(m.store.Entries) {
				// Outdent selected entry
				if m.store.Entries[m.selectedIndex].Indent > 0 {
					m.store.Entries[m.selectedIndex].Indent--
					m.save()
					cmds = append(cmds, m.updateViewport())
				}
			}
			
		case tea.KeyUp:
			if msg.Alt {
				// Navigate entries
				m.navigateUp()
				cmds = append(cmds, m.updateViewport())
			} else {
				// Pass to input
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				cmds = append(cmds, cmd)
			}
			
		case tea.KeyDown:
			if msg.Alt {
				// Navigate entries
				m.navigateDown()
				cmds = append(cmds, m.updateViewport())
			} else {
				// Pass to input
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				cmds = append(cmds, cmd)
			}
			
		case tea.KeyEscape:
			// Clear selection
			m.selectedIndex = -1
			cmds = append(cmds, m.updateViewport())
			
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "[":
				// Move entry up
				if m.selectedIndex > 0 && m.selectedIndex < len(m.store.Entries) {
					m.swapEntries(m.selectedIndex, m.selectedIndex-1)
					m.selectedIndex--
					m.save()
					cmds = append(cmds, m.updateViewport())
				}
			case "]":
				// Move entry down
				if m.selectedIndex >= 0 && m.selectedIndex < len(m.store.Entries)-1 {
					m.swapEntries(m.selectedIndex, m.selectedIndex+1)
					m.selectedIndex++
					m.save()
					cmds = append(cmds, m.updateViewport())
				}
			default:
				// Pass to input
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				cmds = append(cmds, cmd)
			}
			
		default:
			// Pass other keys to input
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
		
	default:
		// Handle other messages
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
		
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	
	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m *Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}
	
	var b strings.Builder
	
	// Header
	header := HeaderStyle.Render("FLOAT Echo") + " " +
		HelpStyle.Render("(Ctrl+C quit | Alt+↑/↓ navigate | Ctrl+E toggle mode)")
	b.WriteString(header + "\n")
	b.WriteString(strings.Repeat("─", m.width) + "\n")
	
	// Viewport with entries
	b.WriteString(m.viewport.View() + "\n")
	
	// Input area
	inputBorder := "─"
	if m.input.Focused() {
		inputBorder = "═"
	}
	b.WriteString(strings.Repeat(inputBorder, m.width) + "\n")
	b.WriteString(m.input.View() + "\n")
	
	// Status bar
	mode := "single-line"
	if m.inputMode == MultiLineMode {
		mode = "multi-line"
	}
	
	status := fmt.Sprintf(" Mode: %s | Entries: %d", mode, len(m.store.Entries))
	if m.message != "" {
		status += " | " + m.message
	}
	if m.err != nil {
		status += " | " + ErrorStyle.Render(m.err.Error())
	}
	
	b.WriteString(StatusBarStyle.Render(status))
	
	return b.String()
}

// Helper methods

func (m *Model) submitEntry() tea.Cmd {
	content := strings.TrimSpace(m.input.Value())
	if content == "" {
		return nil
	}
	
	// Check for bridge commands
	if isBridge, bridgeCmd := m.parser.ExtractBridgeCommand(content); isBridge {
		// Handle bridge command
		if strings.HasPrefix(bridgeCmd, "create:") {
			bridgeID := core.GenerateBridgeID()
			m.setMessage(fmt.Sprintf("Bridge created: %s", bridgeID))
		}
	}
	
	// Parse the entry
	parsed := m.parser.ParseEntry(content)
	
	// Create new entry
	entry := core.NewEntry(parsed.Content)
	entry.Type = parsed.Type
	entry.Metadata = parsed.Metadata
	
	// Insert after selected entry or at end
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.store.Entries) {
		// Insert after selected with same indent
		entry.Indent = m.store.Entries[m.selectedIndex].Indent
		m.store.Entries = append(
			m.store.Entries[:m.selectedIndex+1],
			append([]core.Entry{*entry}, m.store.Entries[m.selectedIndex+1:]...)...,
		)
		m.selectedIndex++
	} else {
		// Append at end
		m.store.Entries = append(m.store.Entries, *entry)
		m.selectedIndex = len(m.store.Entries) - 1
	}
	
	// Clear input
	m.input.SetValue("")
	
	// Save
	m.save()
	
	return m.updateViewport()
}

func (m *Model) toggleInputMode() {
	if m.inputMode == SingleLineMode {
		m.inputMode = MultiLineMode
		m.input.SetHeight(3)
	} else {
		m.inputMode = SingleLineMode
		m.input.SetHeight(1)
	}
}

func (m *Model) navigateUp() {
	if len(m.store.Entries) == 0 {
		return
	}
	
	if m.selectedIndex < 0 {
		m.selectedIndex = len(m.store.Entries) - 1
	} else {
		m.selectedIndex--
		if m.selectedIndex < 0 {
			m.selectedIndex = len(m.store.Entries) - 1
		}
	}
}

func (m *Model) navigateDown() {
	if len(m.store.Entries) == 0 {
		return
	}
	
	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	} else {
		m.selectedIndex++
		if m.selectedIndex >= len(m.store.Entries) {
			m.selectedIndex = 0
		}
	}
}

func (m *Model) swapEntries(i, j int) {
	m.store.Entries[i], m.store.Entries[j] = m.store.Entries[j], m.store.Entries[i]
}

func (m *Model) save() {
	if err := m.persistence.Save(m.store); err != nil {
		m.err = err
	}
}

func (m *Model) setMessage(msg string) {
	m.message = msg
	// Clear message after 3 seconds
	if m.messageTimer != nil {
		m.messageTimer.Stop()
	}
	m.messageTimer = time.AfterFunc(3*time.Second, func() {
		m.message = ""
	})
}

func (m *Model) updateViewport() tea.Cmd {
	// Render all entries
	var content strings.Builder
	
	for i, entry := range m.store.Entries {
		// Indentation
		indent := RenderIndent(entry.Indent)
		
		// Selection indicator
		if i == m.selectedIndex {
			content.WriteString(SelectedEntryStyle.Render(">"))
			content.WriteString(" ")
		} else {
			content.WriteString("  ")
		}
		
		// Timestamp
		timestamp := entry.Timestamp.Format("15:04:05")
		content.WriteString(TimestampStyle.Render(timestamp))
		content.WriteString(" ")
		
		// Type prefix
		content.WriteString(indent)
		content.WriteString(RenderTypePrefix(entry.Type))
		
		// Content
		content.WriteString(entry.Content)
		content.WriteString("\n")
	}
	
	m.viewport.SetContent(content.String())
	
	// Auto-scroll to bottom if near bottom
	if m.viewport.AtBottom() || m.selectedIndex == len(m.store.Entries)-1 {
		m.viewport.GotoBottom()
	}
	
	return nil
}