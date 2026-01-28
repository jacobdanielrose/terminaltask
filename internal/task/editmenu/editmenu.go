package editmenu

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	datepicker "github.com/ethanefung/bubble-datepicker"
	"github.com/google/uuid"
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

	calendarPadding = 10
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
		Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})

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
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

func (e EditTaskKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		e.SaveField,
		e.EscapeEditMode,
		e.SaveTask,
		e.Help,
	}
}

func (e EditTaskKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			e.SaveField,
			e.EscapeEditMode,
			e.SaveTask,
			e.Help,
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

	// User-editable fields
	TaskTitle  textinput.Model
	Desc       textinput.Model
	DatePicker datepicker.Model

	// Layout / dimensions
	width  int
	height int

	// UI state
	focusIdx  int
	showTitle bool
	showHelp  bool
	statusMsg string

	// UI components and styles
	styles Styles
	help   help.Model

	// Input bindings
	keymap *EditTaskKeyMap
}

// New constructs a new edit menu model.
func New(
	width, height int,
	initialTitle, initialDesc string,
	initialTime time.Time,
) Model {
	titleInput := newTitleInput(initialTitle)
	descInput := newDescInput(initialDesc)

	windowTitle := defaultWindowTitle
	if initialTitle != "" {
		windowTitle = initialTitle
	}

	if initialTime.IsZero() {
		initialTime = time.Now()
	}

	return Model{
		// Identity / basic metadata
		Title: windowTitle,
		IsNew: false,

		// User-editable fields
		TaskTitle:  titleInput,
		Desc:       descInput,
		DatePicker: datepicker.New(initialTime),

		// Layout / dimensions
		width:  width,
		height: height,

		// UI state
		focusIdx:  0,
		showTitle: true,
		showHelp:  true,
		statusMsg: "",

		// UI components and styles
		styles: DefaultStyles(),
		help:   help.New(),

		// Input bindings
		keymap: newEditTaskKeyMap(),
	}
}

func newTitleInput(initial string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = defaultTitlePrompt
	ti.PromptStyle.Underline(true)
	ti.Placeholder = defaultTitlePlaceholder
	ti.SetValue(initial)
	ti.Focus()
	return ti
}

func newDescInput(initial string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = defaultDescPrompt
	ti.PromptStyle.Underline(true)
	ti.Placeholder = defaultDescPlaceholder
	ti.SetValue(initial)
	return ti
}

//
// Update / commands
//

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) showStatus(msg string) tea.Cmd {
	m.statusMsg = msg
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case clearStatusMsg:
		m.statusMsg = ""
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.SaveField):
			m.focusIdx = (m.focusIdx + 1) % 3
			m.setFocus()

		case key.Matches(msg, m.keymap.SaveTask):
			if m.DatePicker.Time.Before(time.Now().Truncate(24 * time.Hour)) {
				return m, m.showStatus("You cannot pick a date in the past!")
			}
			m.focusIdx = 0
			m.setFocus()
			return m, func() tea.Msg {
				return SaveTaskMsg{
					TaskID: m.TaskID,
					Title:  m.TaskTitle.Value(),
					Desc:   m.Desc.Value(),
					Date:   m.DatePicker.Time,
					Done:   false,
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

	var cmd tea.Cmd

	m.TaskTitle, cmd = m.TaskTitle.Update(msg)
	cmds = append(cmds, cmd)

	m.Desc, cmd = m.Desc.Update(msg)
	cmds = append(cmds, cmd)

	m.DatePicker, cmd = m.DatePicker.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) setFocus() {
	m.TaskTitle.Blur()
	m.Desc.Blur()
	m.DatePicker.Blur()

	switch m.focusIdx {
	case 0:
		m.TaskTitle.Focus()
	case 1:
		m.Desc.Focus()
	case 2:
		m.DatePicker.SelectDate()
		m.DatePicker.SetFocus(datepicker.FocusCalendar)
	}
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.help.Width = width
	m.height = height
}

func (m *Model) SetShowTitle(v bool) {
	m.showTitle = v
}

func (m *Model) SetShowHelp(v bool) {
	m.showHelp = v
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
	m.TaskTitle.TextStyle = m.styles.Normal
	m.TaskTitle.PromptStyle = m.styles.Normal
	m.Desc.TextStyle = m.styles.Normal
	m.Desc.PromptStyle = m.styles.Normal

	switch m.focusIdx {
	case 0:
		m.TaskTitle.TextStyle = m.styles.Normal
		m.TaskTitle.PromptStyle = m.styles.Focused
	case 1:
		m.Desc.TextStyle = m.styles.Normal
		m.Desc.PromptStyle = m.styles.Focused
	}

	calendar := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		PaddingLeft(calendarPadding).
		Render(m.DatePicker.View())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.TaskTitle.View(),
		m.Desc.View(),
		calendar,
	)
}

// ShowHelp returns whether or not the help is set to be rendered.
func (m Model) ShowHelp() bool {
	return m.showHelp
}

func (m Model) helpView() string {
	return m.styles.HelpStyle.Render(m.help.View(m.keymap))
}
