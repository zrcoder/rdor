package dialog

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/style/color"
)

const (
	DefaultWidth   = 80
	DefaultHeight  = 9
	DefaultPadding = 15
	starCh         = "★"
	starOutlineCh  = "☆"
)

var (
	DefaultContentStyle = lipgloss.NewStyle()
	DefaultBorderStyle  = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#874BFD")).
				Padding(1, 1)
	DefaultWhiteSpaceForground = color.Faint
	starStyle                  = lipgloss.NewStyle().Foreground(color.Orange)
)

func Error(message string) *Dialog {
	return newDialog(message, "Error!", kindError)
}

func Success(message string) *Dialog {
	return newDialog(message, "Success!", kindSuccess)
}

func (d *Dialog) Height(h int) *Dialog {
	d.height = h
	return d
}

func (d *Dialog) Width(w int) *Dialog {
	d.width = w
	return d
}

func (d *Dialog) Stars(total, erned int) *Dialog {
	d.totalStars = total
	d.ernedStars = erned
	return d
}

func (d *Dialog) WhiteSpaceChars(chs string) *Dialog {
	d.whiteSpaceChars = chs
	return d
}

func (d *Dialog) String() string {
	fix(d)
	cw, ch := d.width-d.padding*2-2, d.height-d.padding*2
	d.message = d.contentStyle.Width(cw).Height(ch).Align(lipgloss.Center).Render(d.message)
	if d.totalStars > 0 {
		stars := starStyle.Render(strings.Repeat(starCh, d.ernedStars)) + strings.Repeat(starOutlineCh, d.totalStars-d.ernedStars)
		d.message = lipgloss.JoinVertical(lipgloss.Center, stars, d.message)
	}
	return lipgloss.Place(d.width, d.height,
		lipgloss.Center, lipgloss.Center,
		d.borderStyle.Render(d.message),
		lipgloss.WithWhitespaceChars(d.whiteSpaceChars),
		lipgloss.WithWhitespaceForeground(d.whiteSpaceForground),
	)
}

func newDialog(message, placeholder string, kind kind) *Dialog {
	if message == "" {
		message = placeholder
	}
	return &Dialog{message: message, kind: kind}
}

func fix(d *Dialog) {
	if d.width <= 0 {
		d.width = DefaultWidth
	}
	if d.height <= 0 {
		d.height = DefaultHeight
	}
	if d.padding <= 0 {
		d.padding = DefaultPadding
	}
	if d.contentStyle == nil {
		d.contentStyle = &DefaultContentStyle
	}
	if d.borderStyle == nil {
		d.borderStyle = &DefaultBorderStyle
	}
	if d.whiteSpaceForground == nil {
		d.whiteSpaceForground = DefaultWhiteSpaceForground
	}
	contentStyle := *d.contentStyle
	switch d.kind {
	case kindError:
		contentStyle = contentStyle.Foreground(color.Red)
	case kindSuccess:
		contentStyle = contentStyle.Foreground(color.Green)
	}
	d.contentStyle = &contentStyle
}

type kind = int

const (
	kindPlain kind = iota
	kindError
	kindSuccess
	kindFailure
)

type Dialog struct {
	whiteSpaceForground lipgloss.TerminalColor
	borderStyle         *lipgloss.Style
	contentStyle        *lipgloss.Style
	message             string
	whiteSpaceChars     string
	width               int
	height              int
	padding             int
	kind                kind
	totalStars          int
	ernedStars          int
}
