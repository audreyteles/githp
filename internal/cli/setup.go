package cli

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Define colors
const (
	blue     = lipgloss.Color("#4287f5")
	darkGray = lipgloss.Color("#767676")
)

const (
	addFiles int = iota
	commitMessage
)

// Store labels colors variables
var (
	inputStyle    = lipgloss.NewStyle().Foreground(blue)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

// Form Model
type Form struct {
	state    int
	cursor   int
	err      error
	choices  []string
	textarea textarea.Model
	selected map[int]struct{}
}

// Define errMsg of type error
type errMsg error

// InitialForm Define the Form (inputs, tables, pages...)
func InitialForm() Form {
	var files []string

	files = ListFilesChanged()

	commitMessage := textarea.New()
	commitMessage.Placeholder = "Your commit message here..."
	commitMessage.Focus()

	return Form{
		err:      nil,
		choices:  files,
		state:    addFiles,
		textarea: commitMessage,
		selected: make(map[int]struct{}),
	}
}

func (f Form) Init() tea.Cmd {
	return nil
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		// Quit the app
		case "ctrl+c", "q", "esc":
			return f, tea.Quit
		// Go down
		case "down", "j":
			if f.cursor < len(f.choices)-1 {
				f.cursor++
			}
		// Go up
		case "up":
			if f.cursor > 0 {
				f.cursor--
			}
		case "tab":
			if f.state == commitMessage {
				directory := os.Args[1]

				//Opens an already existing repository.
				r, _ := git.PlainOpen(directory)
				w, _ := r.Worktree()

				for index := range f.selected {
					_, err := w.Add(f.choices[index])

					if err != nil {
						log.Println(err)
					}
				}

				_, err := w.Commit(f.textarea.Value(), &git.CommitOptions{})

				if err != nil {
					fmt.Println(err)
				}

				println("Changes were committed successfully!")

				return f, tea.Quit

			} else {
				f.state = commitMessage
			}

		case "shift+tab":
			f.state = addFiles
		// Select an option
		case " ":
			if f.state == addFiles {
				// Toggle selection
				_, ok := f.selected[f.cursor]
				if ok {
					delete(f.selected, f.cursor)
				} else {
					f.selected[f.cursor] = struct{}{}
				}
			}
		// Confirm selections
		case "ctrl+shift":

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
	var s string

	switch f.state {
	case addFiles:
		s = "Select the files to add:\n\n"

		for i, file := range f.choices {
			cursor := " "
			if f.cursor == i {
				cursor = ">"
			}
			checked := " "
			if _, ok := f.selected[i]; ok {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, file)
		}
	case commitMessage:
		s += fmt.Sprintf("\nEnter your commit message:\n\n")
		s += f.textarea.View()

		s += "\nQUIT NOW!"

	}
	return s
}

func ListFilesChanged() []string {
	var out bytes.Buffer
	var files []string

	// Get modified files
	cmd := exec.Command("git", "ls-files", "--modified")
	cmd.Stdout = &out

	// Run command
	err := cmd.Run()

	// Check if it has an error
	if err != nil {
		fmt.Printf("An error has occurred %v\n", err)
		return nil
	}

	// Split output by lines
	lines := strings.Split(out.String(), "\n")

	// Iterate lines and get all file names
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}

	// Return files modified
	return files
}
