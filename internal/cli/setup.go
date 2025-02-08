package cli

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"os/exec"
	"strings"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Cursor key.Binding
	Diff   key.Binding
	Next   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Space  key.Binding
	Commit key.Binding
	Back   key.Binding
}

var keys = keyMap{
	Cursor: key.NewBinding(
		key.WithKeys("up", "k", "down", "j"),
		key.WithHelp("↑/↓:", "move cursor"),
	),
	Diff: key.NewBinding(
		key.WithKeys("right", "l", "down", "j"),
		key.WithHelp("←/→:", "file diff"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab:", "next"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc:", "quit"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space: ", "toggle file"),
	),
	Commit: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab:", "commit"),
	),
	Back: key.NewBinding(
		key.WithKeys("ctrl+left"),
		key.WithHelp("ctrl+left:", "go back"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Commit, k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Cursor, k.Diff},
		{k.Space, k.Next, k.Help},
		{k.Quit},
	}
}

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
	displayCode
)

// Store labels colors variables
var (
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
	blueStyle     = lipgloss.NewStyle().Foreground(white).BorderForeground(blue).Padding(0, 2).Border(lipgloss.RoundedBorder())
	whiteStyle    = lipgloss.NewStyle().Foreground(white)
	successStyle  = lipgloss.NewStyle().Foreground(green)
	errorStyle    = lipgloss.NewStyle().Foreground(red)
)

// Form Model
type Form struct {
	state    int
	cursor   int
	err      error
	keys     keyMap
	choices  []string
	help     help.Model
	textarea textarea.Model
	viewport viewport.Model
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
	commitMessage.CharLimit = 1000
	commitMessage.SetWidth(64)

	var helpKeys help.Model

	helpKeys = help.New()

	helpKeys.ShowAll = true

	return Form{
		err:      nil,
		keys:     keys,
		choices:  files,
		help:     helpKeys,
		state:    addFiles,
		textarea: commitMessage,
		viewport: viewport.New(100, 10),
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
		case "esc", "ctrl+c":
			return f, tea.Quit
		// Go down
		case "down":
			if f.state == displayCode {
				f.viewport.LineDown(1)
			} else {
				if f.cursor < len(f.choices)-1 {
					f.cursor++
					f.viewport.SetContent(ViewFileDiff(f.choices[f.cursor]))
				}
			}

		// Go up
		case "up":
			if f.state == displayCode {
				f.viewport.LineUp(1)
			} else {
				if f.cursor > 0 {
					f.cursor--
					f.viewport.SetContent(ViewFileDiff(f.choices[f.cursor]))
				}
			}
		case "tab":
			if f.state == commitMessage {
				// Get current directory
				directory, err := os.Getwd()

				if err != nil {
					log.Println("Error getting current working directory")
				}

				// Opens an already existing repository
				r, _ := git.PlainOpen(directory)
				w, _ := r.Worktree()

				// Make git add to each file selected
				for index := range f.selected {
					_, err := w.Add(f.choices[index])

					if err != nil {
						log.Println(err)
						os.Exit(1)
					}
				}

				// Make a commit
				_, err = w.Commit(f.textarea.Value(), &git.CommitOptions{})

				if err != nil {
					println(errorStyle.Bold(true).Render("Cannot commit changes..."))
				} else {
					println(successStyle.Bold(true).Render("Changes were committed successfully!"))
				}
				return f, tea.Quit
			} else {
				f.help.ShowAll = false
				f.state = commitMessage
			}
		case "ctrl+left":
			f.state = addFiles
			f.help.ShowAll = true
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
		case "left":
			if f.state == displayCode {
				f.state = addFiles
			}
		case "right":
			if f.state == addFiles {
				f.state = displayCode
				f.viewport.SetContent(ViewFileDiff(f.choices[f.cursor]))
			}
		}

	// We handle errors just like any other message
	case errMsg:
		f.err = msg
		return f, nil
	}

	if f.state != addFiles {
		f.textarea, cmd = f.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}
	return f, tea.Batch(cmds...)

}

func (f Form) View() string {
	var s string

	switch f.state {
	case addFiles, displayCode:
		var addFilesForm string

		// TUI title
		s += blueStyle.Width(68).Align(lipgloss.Center).Render(" Welcome to GITHP ")

		// First screen title
		addFilesForm += whiteStyle.Width(64).Render("\nSelect the files to add:") + "\n\n"

		// Iterate all files to show them
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
		s += "\n"

		// Get the height from the main view and apply it to the viewport
		f.viewport.Height = lipgloss.Height(addFilesForm)

		// If the current state is displayCode (meaning the focus is on the code), display the viewport
		if f.state == displayCode {
			// Join horizontally the form to add files and
			s += lipgloss.JoinHorizontal(lipgloss.Top, blueStyle.Render(addFilesForm), blueStyle.Render(f.viewport.View()))
		} else {
			s += blueStyle.Render(addFilesForm)
		}

	case commitMessage:
		var commitMessageForm string

		commitMessageForm += whiteStyle.Width(64).Render("Enter your commit message:") + "\n\n"
		commitMessageForm += f.textarea.View()
		s += blueStyle.Padding(1).Render(commitMessageForm)
	}

	s += "\n"
	s += continueStyle.Padding(0, 0, 0, 1).Render(f.help.View(f.keys))
	s += "\n"

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

func ViewFileDiff(fileName string) string {
	var out bytes.Buffer
	var file string

	// Get modified files
	cmd := exec.Command("git", "diff", fileName)
	cmd.Stdout = &out

	// Run command
	err := cmd.Run()

	// Check if it has an error
	if err != nil {
		fmt.Printf("An error has occurred %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(out.String(), "\n")

	for _, line := range lines {
		if line != "" {
			if strings.Contains(line[0:1], "+") {
				file += successStyle.Render(line) + "\n"
			} else if strings.Contains(line[0:1], "-") {
				file += errorStyle.Render(line) + "\n"
			} else {
				file += line + "\n"
			}
		}
	}
	return file
}
