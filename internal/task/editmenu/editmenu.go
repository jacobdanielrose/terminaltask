package editmenu

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/jacobdanielrose/terminaltask/internal/task"
)

//
// Constants
//

const (
	defaultTitlePrompt      = "Title: "
	defaultTitlePlaceholder = "Title"

	defaultDescPrompt      = "Description: "
	defaultDescPlaceholder = "Description"

	defaultWindowTitle = "Editing..."

	statusMsgDatePastError   = "Error: Date cannot be in the past"
	statusMsgTitleEmptyError = "Error: Title cannot be empty"
	statusMsgDescEmptyError  = "Error: Description cannot be empty"
)

//
// Messages
//

type EscapeEditMsg struct{}

type SaveTaskMsg struct {
	TaskID uuid.UUID
	Title  string
	Desc   string
	Date   time.Time
	Done   bool
	IsNew  bool
}

type ErrorMsg struct {
	ErrorStr string
}

type clearStatusMsg struct{}

//
// Styles
//

type Styles struct {
	TitleBar lipgloss.Style
	Title    lipgloss.Style

	HelpStyle lipgloss.Style

	Focused lipgloss.Style
	Blurred lipgloss.Style
	Normal  lipgloss.Style

	StatusMessage lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.TitleBar = lipgloss.NewStyle().Padding(0, 0, 1, 2) //nolint:mnd

	s.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2) //nolint:mnd

	s.Focused = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Padding(0, 0, 0, 2).MarginRight(1)
	s.Normal = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)

	s.StatusMessage = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF5555"})

	return s
}

//
// Keymap
//

type EditTaskKeyMap struct {
	SaveField      key.Binding
	EscapeEditMode key.Binding
	SaveTask       key.Binding
	Help           key.Binding
	Quit           key.Binding
}

func newEditTaskKeyMap() *EditTaskKeyMap {
	return &EditTaskKeyMap{
		SaveField: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "next field"),
		),
		EscapeEditMode: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit edit mode"),
		),
		SaveTask: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save task"),
		),
		Help: key.NewBinding(
			key.WithKeys("ctrl+o"),
			key.WithHelp("ctrl+o", "help"),
		),
	}
}

func (e EditTaskKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		e.SaveField,
		e.EscapeEditMode,
		e.SaveTask,
		// e.Help,
	}
}

func (e EditTaskKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			e.SaveField,
			e.EscapeEditMode,
			e.SaveTask,
			// e.Help,
			e.Quit,
		},
	}
}

//
// Model & constructor
//

// Model for the edit menu.
type Model struct {
	// Identity / basic metadata
	Title  string
	TaskID uuid.UUID
	IsNew  bool

	form Form

	// Layout / dimensions
	width  int
	height int

	// UI state
	showTitle bool
	showHelp  bool
	statusMsg string

	// UI components and styles
	styles Styles
	help   help.Model

	// Input bindings
	keymap *EditTaskKeyMap
}

// New constructs a new edit menu model with default styles.
func New(task task.Task) Model {
	return NewWithSize(0, 0, task)
}

// NewWithStyles constructs a new edit menu model with explicit styles
// for the outer edit menu and the inner form.
func NewWithStyles(task task.Task, menuStyles Styles, formStyles Styles) Model {
	return NewWithSizeAndStyles(0, 0, task, menuStyles, formStyles)
}

func NewWithSize(
	width, height int,
	task task.Task,
) Model {
	styles := DefaultStyles()
	return NewWithSizeAndStyles(width, height, task, styles, styles)
}

// NewWithSizeAndStyles is the core constructor used by all others.
// It allows callers (like the top-level app model) to inject styling
// for both the edit menu container and the inner form.
func NewWithSizeAndStyles(
	width, height int,
	task task.Task,
	menuStyles Styles,
	formStyles Styles,
) Model {
	var (
		title       = task.Title()
		description = task.Description()
		duedate     = task.DueDate
		done        = task.Done
		keymap      = newEditTaskKeyMap()
		isNew       = false
	)

	windowTitle := defaultWindowTitle
	if title != "" {
		windowTitle = title
	}

	if duedate.IsZero() {
		duedate = time.Now()
	}

	if task.IsEmpty() {
		isNew = true
	}

	return Model{
		// Identity / basic metadata
		Title: windowTitle,
		IsNew: isNew,

		// User-editable fields
		form: NewForm(title, description, duedate, done, keymap, formStyles),

		// Layout / dimensions
		width:  width,
		height: height,

		// UI state
		showTitle: true,
		showHelp:  true,
		statusMsg: "",

		// UI components and styles
		styles: menuStyles,
		help:   help.New(),

		// Input bindings
		keymap: keymap,
	}
}

//
// Update / commands
//

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) showStatus(msg string) tea.Cmd {
	m.statusMsg = msg
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.form, cmd = m.form.Update(msg)

	switch msg := msg.(type) {
	case clearStatusMsg:
		m.statusMsg = ""
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.SaveTask):
			if m.form.Date.Time.Before(time.Now().Truncate(24 * time.Hour)) {
				return m, m.showStatus(statusMsgDatePastError)
			}
			if m.form.Title.Value() == "" {
				return m, m.showStatus(statusMsgTitleEmptyError)
			}
			if m.form.Desc.Value() == "" {
				return m, m.showStatus(statusMsgDescEmptyError)
			}

			m.form = m.form.setFocus()
			return m, func() tea.Msg {
				return SaveTaskMsg{
					TaskID: m.TaskID,
					Title:  m.form.Title.Value(),
					Desc:   m.form.Desc.Value(),
					Date:   m.form.Date.Time,
					Done:   m.form.Done,
					IsNew:  m.IsNew,
				}
			}

		case key.Matches(msg, m.keymap.EscapeEditMode):
			return m, func() tea.Msg {
				return EscapeEditMsg{}
			}

		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m Model) SetSize(width int, height int) Model {
	m.width = width
	m.help.Width = width
	m.height = height
	return m
}

func (m Model) SetShowTitle(v bool) Model {
	m.showTitle = v
	return m
}

func (m Model) SetShowHelp(v bool) Model {
	m.showHelp = v
	return m
}

//
// View
//

func (m Model) View() string {
	var (
		sections    []string
		availHeight = m.height
	)

	if m.showTitle {
		v := m.titleView()
		sections = append(sections, v)
		availHeight -= lipgloss.Height(v)
	}

	var helpView string
	if m.showHelp {
		helpView = m.helpView()
		availHeight -= lipgloss.Height(helpView)
	}

	if m.statusMsg != "" {
		availHeight -= lipgloss.Height(m.statusMsg)
	}

	editContent := lipgloss.NewStyle().Height(availHeight).Render(m.editView())
	sections = append(sections, editContent)

	if m.statusMsg != "" {
		statusView := m.styles.StatusMessage.Align(lipgloss.Center).Render(m.statusMsg)
		sections = append(sections, statusView)
	}

	if m.showHelp {
		sections = append(sections, helpView)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) titleView() string {
	var (
		view          string
		titleBarStyle = m.styles.TitleBar
	)

	view += m.styles.Title.Render(m.Title)

	if len(view) > 0 {
		return titleBarStyle.Render(view)
	}

	return view
}

func (m Model) ShowTitle() bool {
	return m.showTitle
}

func (m Model) editView() string {
	// Delegate rendering of the editable fields to the Form submodel.
	return m.form.View()
}

// ShowHelp returns whether or not the help is set to be rendered.
func (m Model) ShowHelp() bool {
	return m.showHelp
}

func (m Model) helpView() string {
	return m.styles.HelpStyle.Render(m.help.View(m.keymap))
}
