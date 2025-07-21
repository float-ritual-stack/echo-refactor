package main

import (
	"fmt"
	"log"
	"os"
	
	tea "github.com/charmbracelet/bubbletea"
	
	"github.com/evan/float-echo/ui"
)

func main() {
	// Create the model
	model, err := ui.NewModel()
	if err != nil {
		log.Fatalf("Error creating model: %v", err)
	}
	
	// Create the program
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}