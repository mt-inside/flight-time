package main

import (
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const timeParseFmt = "2006-01-02 15:04"

// This is the only way to do this sadly
// - https://stackoverflow.com/questions/49084316/why-doesnt-gos-time-parse-parse-the-timezone-identifier
// - https://stackoverflow.com/questions/25368415/how-to-properly-parse-timezone-codes
func parseTimeWithZone(layout, str, loc string) time.Time {
	l, err := time.LoadLocation(loc)
	if err != nil {
		panic(err)
	}
	t, err := time.ParseInLocation(timeParseFmt, str, l)
	if err != nil {
		panic(err)
	}

	return t
}

func main() {
	start := parseTimeWithZone(timeParseFmt, os.Args[1], os.Args[2])
	end := parseTimeWithZone(timeParseFmt, os.Args[3], os.Args[4])

	var p *tea.Program

	flight := NewFlight(start, end, func(elapsed, remain time.Duration, wallclockEnd time.Time, pct float64, relTime time.Time) {
		p.Send(progressMsg{elapsed, remain, wallclockEnd, pct, relTime})
	})

	model := model{
		flight:   flight,
		progress: progress.New(progress.WithDefaultGradient()),
	}

	p = tea.NewProgram(model)

	go flight.Start()

	if err := p.Start(); err != nil {
		panic(err)
	}
}
