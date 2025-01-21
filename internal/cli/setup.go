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
	blue     = lipgloss.Color("#0380fc")
	darkGray = lipgloss.Color("#767676")
	white    = lipgloss.Color("#ffffff")
	red      = lipgloss.Color("#FE5F86")
	green    = lipgloss.Color("#02BA84")
)

const (
	addFiles int = iota
	commitMessage
)

// Store labels colors variables
var (
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
	blueStyle     = lipgloss.NewStyle().Foreground(white).BorderForeground(blue).Padding(0, 1).Border(lipgloss.RoundedBorder())
	whiteStyle    = lipgloss.NewStyle().Foreground(white)
	successStyle  = lipgloss.NewStyle().Foreground(green)
	errorStyle    = lipgloss.NewStyle().Foreground(red)
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
						os.Exit(1)
					}
				}

				_, err := w.Commit(f.textarea.Value(), &git.CommitOptions{})

				if err != nil {
					println(errorStyle.Bold(true).Render("Cannot commit changes..."))
				} else {
					println(successStyle.Bold(true).Render("Changes were committed successfully!"))
				}
				return f, tea.Quit
			} else {
				f.state = commitMessage
			}
		case "shift+tab":
			f.state = addFiles
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
		s += blueStyle.Width(42).Align(lipgloss.Center).Render(" Welcome to GITHP ")

		var addFilesForm string

		addFilesForm += whiteStyle.Width(40).Render("\nSelect the files to add:") + "\n\n"

		for i, file := range f.choices {
			cursor := " "
			if f.cursor == i {
				cursor = whiteStyle.Render(">")
			}
			checked := " "
			if _, ok := f.selected[i]; ok {
				checked = successStyle.Render("✔")
			}
			addFilesForm +=
				fmt.Sprintf(
					"%s [%s] %s\n",
					cursor, checked, whiteStyle.Render(file))
		}
		addFilesForm += fmt.Sprintf("\nPress %s to select", successStyle.Bold(true).Render("SPACE"))
		addFilesForm += fmt.Sprintf("\nPress %s to go next", continueStyle.Bold(true).Render("TAB"))
		s += "\n" + blueStyle.Render(addFilesForm) + "\n"
		s += "\n"

	case commitMessage:
		var commitMessageForm string

		commitMessageForm += whiteStyle.Width(30).Render("\nEnter your commit message:") + "\n\n"
		commitMessageForm += f.textarea.View()
		commitMessageForm += fmt.Sprintf("\n\nPress %s to go next", continueStyle.Bold(true).Render("TAB"))
		commitMessageForm += fmt.Sprintf("\nPress %s to go back", continueStyle.Bold(true).Render("SHIFT + TAB"))

		s += blueStyle.Render(commitMessageForm) + "\n"
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
