package main

import (
	"time"
)

type progressFn func(time.Duration, time.Duration, time.Time, float64, time.Time)

type flight struct {
	start    time.Time
	end      time.Time
	duration time.Duration
	tzDiff   int

	onProgress progressFn
}

func NewFlight(start, end time.Time, onProgress progressFn) flight {
	_, startOffset := start.Zone()
	_, endOffset := end.Zone()
	return flight{
		start,
		end,
		end.Sub(start),
		endOffset - startOffset,
		onProgress,
	}
}

func (f *flight) Start() {
	for {
		now := time.Now()
		remain := f.end.Sub(now)
		// Want to check this, but time.Now()'s location comes back as "Local"
		// if time.Now().Location() != f.start.Location() {
		// 	panic("Start time zone not the same as machine time zone; can't calcualte wall-clock arrival time")
		// }
		wallclockEnd := f.end.In(time.Now().Location())
		prog := now.Sub(f.start)
		pct := prog.Seconds() / f.duration.Seconds()

		relTime := f.start.Add(prog + time.Duration(pct*float64(f.tzDiff))*time.Second)

		f.onProgress(prog, remain, wallclockEnd, pct, relTime)

		time.Sleep(1 * time.Second)
	}
}
