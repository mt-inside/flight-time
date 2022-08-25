package main

import (
	"time"
)

type progressFn func(time.Duration, time.Duration, float64, time.Time)

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
		prog := now.Sub(f.start)
		pct := prog.Seconds() / f.duration.Seconds()

		relTime := f.start.Add(prog + time.Duration(pct*float64(f.tzDiff))*time.Second)

		f.onProgress(prog, remain, pct, relTime)

		time.Sleep(1 * time.Second)
	}
}
