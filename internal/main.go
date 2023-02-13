package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Run(args []string) {
	var err error
	_, err = tea.NewProgram(New()).Run()
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}

type errMsg error

type main struct {
	disks    int
	piles    []*pile
	keys     keyMap
	keysHelp help.Model

	setting  bool
	showHelp bool
	buf      *strings.Builder
	steps    int
	err      error
	overDisk *disk
}

func New() *main {
	return &main{
		setting:  true,
		keys:     keys,
		keysHelp: help.New(),
		buf:      &strings.Builder{},
	}
}

func (m *main) Init() tea.Cmd {
	m.keys.Piles.SetEnabled(false)
	m.keys.Disks.SetEnabled(true)
	m.keys.Reset.SetEnabled(false)
	m.keysHelp.ShowAll = true
	return nil
}

func (m *main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
	case tea.KeyMsg:
		m.err = nil
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			m.keysHelp.ShowAll = !m.showHelp
			m.keys.Reset.SetEnabled(!m.setting && !m.showHelp)
		case key.Matches(msg, m.keys.Reset):
			if m.setting {
				m.err = errDiskNum
			} else {
				m.setting = true
			}
			return m, m.Init()
		case key.Matches(msg, m.keys.Disks):
			n, _ := strconv.Atoi(msg.String())
			m.setted(n)
		case key.Matches(msg, m.keys.Piles):
			return m, m.pick(msg.String())
		default:
			if m.setting {
				m.err = errDiskNum
			}
		}
	}
	return m, nil
}

func (m *main) View() string {
	m.buf.Reset()
	if m.showHelp {
		m.writeHelpInfo()
	} else {
		m.writeHead()
		if m.setting {
			m.writeSettingView()
		} else {
			m.writePoles()
			m.writeGround()
			m.writeLabels()
			m.writeState()
		}
	}
	m.writeKeysHelp()
	m.writeBlankLine()
	return m.buf.String()
}

func (m *main) setted(n int) {
	m.setting = false
	m.disks = n
	m.steps = 0
	m.overDisk = nil
	m.err = nil
	m.showHelp = false
	m.keys.Disks.SetEnabled(false)
	m.keys.Piles.SetEnabled(true)
	m.keys.Reset.SetEnabled(true)
	m.piles = make([]*pile, 3)
	for i := range m.piles {
		m.piles[i] = &pile{}
	}
	disks := make([]*disk, n)
	for i := 1; i <= n; i++ {
		disks[n-i] = &disk{
			id:   i,
			view: diskStyles[i-1].Render(strings.Repeat(diskCh, i*diskWidthUnit)),
		}
	}
	m.piles[0].disks = disks
}

func (m *main) pick(key string) tea.Cmd {
	if m.success() {
		return nil
	}
	idx := map[string]int{
		"1": 0,
		"j": 0,
		"2": 1,
		"k": 1,
		"3": 2,
		"l": 2,
	}
	i := idx[key]
	curPile := m.piles[i]
	return func() tea.Msg {
		if m.overDisk == nil && curPile.empty() {
			return nil
		}
		if m.overDisk == nil {
			curPile.overOne = true
			m.overDisk = curPile.top()
			return nil
		}
		if !curPile.empty() && m.overDisk.id > curPile.top().id {
			return errMsg(errCantMove)
		}
		if !curPile.empty() && m.overDisk == curPile.top() {
			curPile.overOne = false
			m.overDisk = nil
			return nil
		}
		for _, p := range m.piles {
			if p.overOne {
				m.steps++
				curPile.push(p.pop())
				p.overOne = false
				m.overDisk = nil
			}
		}
		return nil
	}
}

func (m *main) writeHead() {
	m.buf.WriteString(head)
}
func (m *main) success() bool {
	last := m.piles[len(m.piles)-1]
	return len(last.disks) == m.disks
}

func (m *main) writeSettingView() {
	m.writeLine(settingHint)
	if m.err != nil {
		m.writeError(m.err)
	}
	m.writeBlankLine()
}

func (m *main) writePoles() {
	views := make([]string, len(m.piles))
	for i, p := range m.piles {
		views[i] = p.view()
	}
	poles := lipgloss.JoinHorizontal(
		lipgloss.Top,
		views...,
	)
	m.buf.WriteString(poles)
	m.writeBlankLine()
}

func (m *main) writeGround() {
	m.buf.WriteString(strings.Repeat(groundCh, (pileWidth*3 + horizontalSepBlanks*4)))
	m.writeBlankLine()
}

func (m *main) writeLabels() {
	n := horizontalSepBlanks + (pileWidth-len(pole1Label))/2
	m.buf.WriteString(blanks(n))
	m.buf.WriteString(pole1Label)
	n = (pileWidth-len(pole1Label))/2 + horizontalSepBlanks + (pileWidth-len(pole2Label))/2
	m.buf.WriteString(blanks(n))
	m.buf.WriteString(pole2Label)
	n = (pileWidth-len(pole2Label))/2 + horizontalSepBlanks + (pileWidth-len(pole3Label))/2
	m.buf.WriteString(blanks(n))
	m.buf.WriteString(pole3Label)
	m.writeBlankLine()
	m.writeBlankLine()
}

func (m *main) writeState() {
	if m.success() {
		minSteps := 1<<m.disks - 1
		totalStart := 5
		if m.steps == minSteps {
			m.buf.WriteString(infoStyle.Render("Fantastic! you earned all the stars! "))
			m.buf.WriteString(starStyle.Render(strings.Repeat(starCh, totalStart)))
		} else {
			s := fmt.Sprintf("Done! Taken %d steps, can you complete it in %d step(s)? ", m.steps, minSteps)
			m.buf.WriteString(infoStyle.Render(s))
			stars := 3
			if m.steps-minSteps > minSteps/2 {
				stars = 1
			}
			s = strings.Repeat(starCh, stars) + strings.Repeat(starOutlineCh, totalStart-stars)
			m.buf.WriteString(starStyle.Render(s))
		}
		m.writeBlankLine()
	} else if m.err != nil {
		m.writeError(m.err)
		m.writeBlankLine()
	} else {
		m.writeLine(fmt.Sprintf("steps: %d", m.steps))
	}
	m.writeBlankLine()
}

func (m *main) writeKeysHelp() {
	m.buf.WriteString(m.keysHelp.View(m.keys))
	m.writeBlankLine()
}

func (m *main) writeHelpInfo() {
	m.buf.WriteString(helpInfo)
}

func (m *main) writeBlankLine() {
	m.buf.WriteByte('\n')
}

func (m *main) writeError(err error) {
	m.buf.WriteString(errorStyle.Render(err.Error()))
}

func (m *main) writeLine(s string) {
	m.buf.WriteString(s)
	m.writeBlankLine()
}
