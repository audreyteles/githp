package main

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
