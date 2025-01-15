package cli

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v5"
	"os"
)

// Define colors
const (
	blue     = lipgloss.Color("#4287f5")
	darkGray = lipgloss.Color("#767676")
)

// Store labels colors variables
var (
	inputStyle    = lipgloss.NewStyle().Foreground(blue)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

// Form Model
type Form struct {
	inputs   []textinput.Model
	textarea textarea.Model
	cursor   int
	err      error
}

// Define errMsg of type error
type errMsg error

// InitialForm InitialModel Define the Form (inputs, tables, pages...)
func InitialForm() Form {
	commitMessage := textarea.New()
	commitMessage.Placeholder = "Your commit message here..."
	commitMessage.Focus()

	return Form{
		textarea: commitMessage,
		err:      nil,
	}
}

func (f Form) Init() tea.Cmd {
	return textarea.Blink
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnd:
			directory := os.Args[1]

			// Opens an already existing repository.
			r, _ := git.PlainOpen(directory)
			w, _ := r.Worktree()

			_, err := w.Commit(f.textarea.Value(), &git.CommitOptions{})

			if err != nil {
				fmt.Println(err)
			}

			println("Changes were committed successfully!")
			
			return f, tea.Quit

		case tea.KeyEsc:
			if f.textarea.Focused() {
				f.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return f, tea.Quit
		default:
			if !f.textarea.Focused() {
				cmd = f.textarea.Focus()
				cmds = append(cmds, cmd)
			}

		}

	// We handle errors just like any other message
	case errMsg:
		f.err = msg
		return f, nil
	}

	f.textarea, cmd = f.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return f, tea.Batch(cmds...)
}

func (f Form) View() string {
	return fmt.Sprintf(
		"Welcome to githp!\n\n%s\n\n%s",
		f.textarea.View(),
		"(ctrl+c to quit)",
	) + "\n\n"
}
