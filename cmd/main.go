package main

// A simple program demonstrating the textarea component from the Bubbles
// component library.

import (
	"githp/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	// Initialize cli app
	p := tea.NewProgram(cli.InitialForm())

	// If it has an error
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
