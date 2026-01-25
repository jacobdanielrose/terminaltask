package editmenu

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	datepicker "github.com/ethanefung/bubble-datepicker"
)

type EscapeEditMsg struct{}
type SaveTaskMsg struct {
	Title string
	Desc  string
	Date  time.Time
}

type Model struct {
	showTitle  bool
	showHelp   bool
	Title      string
	TaskTitle  textinput.Model
	Desc       textinput.Model
	focusIdx   int // 0=title, 1=desc, 2=duedate
	keymap     *EditTaskKeyMap
	isNew      bool
	styles     Styles
	help       help.Model
	width      int
	height     int
	DatePicker datepicker.Model
}

func New(width, height int) Model {
	title := textinput.New()
	title.Prompt = "Title: "
	title.PromptStyle.Underline(true)
	title.Placeholder = "Title"
	title.Focus()

	desc := textinput.New()
	desc.Prompt = "Description: "
	desc.PromptStyle.Underline(true)
	desc.Placeholder = "Description"

	titleStr := "Editing..."

	return Model{
		showTitle:  true,
		showHelp:   true,
		Title:      titleStr,
		TaskTitle:  title,
		Desc:       desc,
		focusIdx:   0,
		keymap:     newEditTaskKeyMap(),
		styles:     DefaultStyles(),
		help:       help.New(),
		width:      width,
		height:     height,
		DatePicker: datepicker.New(time.Now()),
	}
}

func (m *Model) SetShowTitle(v bool) {
	m.showTitle = v
}
