package main

import (
	"flag"
	"fmt"
	"github.com/0x0f0f0f/memplot"
	"gonum.org/v1/plot/vg"
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

// Simply exit on error
func checke(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	opts := memplot.PlotOptions{
		PlotRss: true,
		PlotVsz: false,
		// PlotNumThreads: false,
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Arguments following options will be"+
			" interpreted as the command to spawn and sample\n")
		flag.PrintDefaults()
	}
	// Default sample duration time
	defaultSd, err := time.ParseDuration("5ms")
	check(err)

	// Total sampling time
	defaultDur, err := time.ParseDuration("0s")
	check(err)

	defaultFilename := "output-plot.png"
	pidPtr := flag.Int("pid", -1, "pid of the process to analyze")
	filenamePtr := flag.String("o", defaultFilename, "output image file name. "+
		"supported extensions are:\n.eps, .jpg, .jpeg, .pdf, "+
		".png, .svg, .tex, .tif and .tiff\n")
	sdPtr := flag.Duration("sd", defaultSd, "sample size in time")
	durPtr := flag.Duration("dur", defaultDur, "total profiling time. a value of 0 means"+
		" that the program\nwill be sampled until it is no longer running")
	// To plot or not VSZ
	flag.BoolVar(&opts.PlotVsz, "vsz", false, "plot virtual size")
	// flag.BoolVar(&opts.PlotNumThreads, "nthreads", false, "plot the number of threads")
	widthStr := flag.String("width", "16cm", "plot image width (can be cm or in)")
	heightStr := flag.String("height", "12cm", "plot image height (can be cm or in)")

	flag.Parse()

	// Parse the image size
	widthImage, err := vg.ParseLength(*widthStr)
	checke(err)
	heightImage, err := vg.ParseLength(*heightStr)
	checke(err)

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
		pidChan := make(chan int, 1)
		go func() {
			cmd := exec.Command(args[0], args[1:]...)
			err := cmd.Start()
			check(err)
			pidChan <- cmd.Process.Pid
			cmd.Wait()
		}()
		*pidPtr = <-pidChan
	} else {
		if len(flag.Args()) > 0 {
			fmt.Println("A pid was specified. Ignoring additional arguments")
		}
	}

	// Create and sample
	fmt.Fprintln(os.Stderr, "Collecting data from PID", *pidPtr, "...")
	if *durPtr == 0 {
		fmt.Fprintln(os.Stderr, "Warning: sampling will continue "+
			"until program is no longer running")
	}
	coll, err := memplot.NewCollection(int32(*pidPtr), *sdPtr, *durPtr)
	check(err)

	fmt.Fprintln(os.Stderr, "Generating plot...")
	plot, err := coll.Plot(opts)
	check(err)

	fmt.Fprintln(os.Stderr, "Saving plot..")
	memplot.SavePlot(plot, widthImage, heightImage, *filenamePtr)

}
