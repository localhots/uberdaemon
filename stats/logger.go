package stats

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger struct {
	base

	out      io.Writer
	interval time.Duration
}

func NewLogger(out io.Writer, interval time.Duration) *Logger {
	l := &Logger{
		out:      out,
		interval: interval,
	}
	l.init()
	go l.printWithInterval()

	return l
}

func NewStdoutLogger(interval time.Duration) *Logger {
	return NewLogger(os.Stdout, interval)
}

func (l *Logger) Print() {
	for name, s := range l.stats {
		fmt.Fprintf(l.out, "%s statistics:\n"+
			"Processed: %d\n"+
			"Errors:    %d\n"+
			"Min:       %.8fms\n"+
			"Max:       %.8fms\n"+
			"95%%:       %.8fms\n"+
			"Mean:      %.8fms\n"+
			"StdDev:    %.8fms\n",
			name,
			s.time.Count(),
			s.errors.Count(),
			float64(s.time.Min())/1000000,
			float64(s.time.Max())/1000000,
			s.time.Percentile(0.95)/1000000,
			s.time.Mean()/1000000,
			s.time.StdDev()/1000000,
		)
		s.time.Clear()
		s.errors.Clear()
	}
}

func (l *Logger) printWithInterval() {
	if l.interval == 0 {
		return
	}

	for range time.NewTicker(l.interval).C {
		l.Print()
	}
}
