package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/0x0f0f0f/memplot"
	"gonum.org/v1/plot/vg"
	"os"
	"time"
)

// Simply panic on error
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	opts := memplot.PlotOptions{
		PlotRss: true,
		PlotVsz: false,
	}

	pidPtr := flag.Int("pid", -1, "pid of the process to analyze")

	// Default sample duration time
	defaultSd, err := time.ParseDuration("5ms")
	check(err)

	// Total sampling time
	defaultDur, err := time.ParseDuration("10s")
	check(err)

	sdPtr := flag.Duration("sd", defaultSd, "sample size in time")
	durPtr := flag.Duration("dur", defaultDur, "total profiling time")

	flag.BoolVar(&opts.PlotVsz, "vsz", false, "plot virtual size")

	flag.Parse()
	if *pidPtr <= 0 {
		panic(errors.New("Invalid PID. Please specify a PID using -pid flag"))
	}

	// Create and sample
	fmt.Fprintln(os.Stderr, "Analyzing PID", *pidPtr, "...")
	coll, err := memplot.NewMemoryCollection(int32(*pidPtr), *sdPtr, *durPtr)
	check(err)

	fmt.Fprintln(os.Stderr, "Generating plot...")
	plot, err := coll.Plot(opts)
	check(err)

	fmt.Fprintln(os.Stderr, "Saving plot..")
	memplot.SavePlot(plot, 8*vg.Inch, 8*vg.Inch, "plot.png")

}
