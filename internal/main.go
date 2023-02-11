package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
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
	count    int
	err      error
	overDisk *disk
}

func New() *main {
	return &main{
		setting:  true,
		keys:     keysSetting,
		keysHelp: help.New(),
		buf:      &strings.Builder{},
	}
}

func (m *main) setted(n int) {
	m.setting = false
	m.disks = n
	m.count = 0
	m.overDisk = nil
	m.err = nil
	m.showHelp = false
	m.keys = keysSetted
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

func (m *main) Init() tea.Cmd { return nil }

func (m *main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	set := func(key string) (tea.Model, tea.Cmd) {
		n, _ := strconv.Atoi(key)
		m.setted(n)
		return m, nil
	}
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
	case tea.KeyMsg:
		m.err = nil
		key := msg.String()
		switch key {
		case "q":
			return m, tea.Quit
		case "h":
			m.showHelp = !m.showHelp
			return m, nil
		case "r":
			if m.setting {
				m.err = errDiskNum
			} else {
				m.setting = true
			}
		case "1", "2", "3":
			if m.setting {
				return set(key)
			}
			return m, m.pick(key)
		case "4", "5", "6", "7":
			if m.setting {
				return set(key)
			}
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
		if m.setting {
			m.writeSettingView()
		} else {
			m.writeHead()
			m.writePoles()
			m.writeGround()
			m.writeLabels()
			m.writeState()
		}
	}
	m.writeKeysHelp()
	return m.buf.String()
}

func (m *main) pick(id string) tea.Cmd {
	if m.success() {
		return nil
	}
	return func() tea.Msg {
		i, _ := strconv.Atoi(id)
		curPile := m.piles[i-1]
		if m.overDisk == nil && curPile.empty() {
			return nil
		}
		m.count++
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
				curPile.push(p.pop())
				p.overOne = false
				m.overDisk = nil
			}
		}
		return nil
	}
}

func (m *main) writeHead() {
	if m.success() {
		minSteps := 1<<m.disks - 1
		steps := m.count / 2
		stars := 5
		if steps == minSteps {
			m.buf.WriteString(helpStyle.Render("Fantastic! you earned all the stars! "))
			m.buf.WriteString(starStyle.Render(strings.Repeat(successCh, stars)))
		} else {
			s := fmt.Sprintf("Done! can you complete it in %d step(s)? ", minSteps)
			m.buf.WriteString(helpStyle.Render(s))
			if steps-minSteps <= minSteps/2 {
				stars = 3
			} else {
				stars = 1
			}
			m.buf.WriteString(starStyle.Render(strings.Repeat(successCh, stars)))
		}
	}
	m.writeBlankLine()
}

func (m *main) success() bool {
	last := m.piles[len(m.piles)-1]
	return len(last.disks) == m.disks
}

func (m *main) writeSettingView() {
	m.writeLine("how many disks do you like? (1-7)")
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
	if m.err != nil {
		m.writeError(m.err)
		m.writeBlankLine()
	} else {
		m.writeLine(fmt.Sprintf("step: %d", m.count/2))
	}
	m.writeBlankLine()
}

func (m *main) writeKeysHelp() {
	m.buf.WriteString(m.keysHelp.FullHelpView(m.keys.FullHelp()))
	m.writeBlankLine()
}

func (m *main) writeHelpInfo() {
	m.buf.WriteString(helpStyle.Render(helpInfo))
	m.writeBlankLine()
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
