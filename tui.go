package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const timeZoneFmt = "15:04 -0700 MST"
const timeRelFmt = "15:04"

type model struct {
	flight  flight
	elapsed time.Duration
	remain  time.Duration
	ratio   float64
	relTime time.Time

	progress progress.Model
}

type progressMsg struct {
	elapsed time.Duration
	remain  time.Duration
	ratio   float64
	relTime time.Time
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4
		return m, nil

	case progressMsg:
		m.elapsed = msg.elapsed
		m.remain = msg.remain
		m.ratio = msg.ratio
		m.relTime = msg.relTime
		return m, nil

	default:
		return m, nil
	}
}

var labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render
var infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render
var relTimeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true).Render
var remainStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true).Render

func (m model) View() string {
	return "\n" +
		fmt.Sprintf("%v ✈️ %v    ⏱ %v",
			m.flight.start.Format(timeZoneFmt),
			m.flight.end.Format(timeZoneFmt),
			m.flight.duration,
		) +
		"\n" +
		fmt.Sprintf("(UTC 🛫%v 🕰%v 🛬%v)",
			m.flight.start.In(time.UTC).Format(timeRelFmt),
			time.Now().In(time.UTC).Format(timeRelFmt),
			m.flight.end.In(time.UTC).Format(timeRelFmt),
		) +
		"\n" +
		fmt.Sprintf(labelStyle("TZ diff %sh  Rel time %s"),
			infoStyle(strconv.Itoa(m.flight.tzDiff/60/60)),
			relTimeStyle(m.relTime.Format(timeRelFmt)),
		) +
		"\n" +
		"\n" +
		m.progress.ViewAs(m.ratio) +
		"\n" +
		"\n" +
		fmt.Sprintf(labelStyle("Elapsed %v  Remain %v"),
			infoStyle(m.elapsed.String()),
			remainStyle(m.remain.String()),
		) +
		"\n" +
		"\n"
}
