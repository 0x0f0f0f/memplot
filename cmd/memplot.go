package main

import (
	"flag"
	"fmt"
	"github.com/0x0f0f0f/memplot"
	"gonum.org/v1/plot/vg"
	"log"
	"os"
	"os/exec"
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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Any argument following the options will be"+
			" interpreted as the command to spawn and sample\n")
		flag.PrintDefaults()
	}
	// Default sample duration time
	defaultSd, err := time.ParseDuration("5ms")
	check(err)

	// Total sampling time
	defaultDur, err := time.ParseDuration("10s")
	check(err)

	defaultFilename := "output-plot.png"
	pidPtr := flag.Int("pid", -1, "pid of the process to analyze")
	filenamePtr := flag.String("o", defaultFilename, "output image file name")
	sdPtr := flag.Duration("sd", defaultSd, "sample size in time")
	durPtr := flag.Duration("dur", defaultDur, "total profiling time")

	// To plot or not VSZ
	flag.BoolVar(&opts.PlotVsz, "vsz", false, "plot virtual size")

	widthStr := flag.String("width", "16cm", "plot image width (can be cm or in)")
	heightStr := flag.String("height", "12cm", "plot image height (can be cm or in)")

	flag.Parse()

	// Run the PID passed with the -pid flag or the arguments following options
	if *pidPtr <= 0 {
		args := flag.Args()
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr,
				"Invalid PID. Please specify PID using -pid flag"+
					" or specify a command to exec and sample\n")
			flag.Usage()
			os.Exit(1)
		}
		cmd := exec.Command(args[0], args[1:]...)
		err := cmd.Start()
		check(err)
		*pidPtr = cmd.Process.Pid
	} else {
		if len(flag.Args()) > 0 {
			log.Println("A pid was specified. Ignoring arguments")
		}
	}

	widthImage, err := vg.ParseLength(*widthStr)
	check(err)
	heightImage, err := vg.ParseLength(*heightStr)
	check(err)

	// Create and sample
	fmt.Fprintln(os.Stderr, "Analyzing PID", *pidPtr, "...")
	coll, err := memplot.NewMemoryCollection(int32(*pidPtr), *sdPtr, *durPtr)
	check(err)

	fmt.Fprintln(os.Stderr, "Generating plot...")
	plot, err := coll.Plot(opts)
	check(err)

	fmt.Fprintln(os.Stderr, "Saving plot..")
	memplot.SavePlot(plot, widthImage, heightImage, *filenamePtr)

}
