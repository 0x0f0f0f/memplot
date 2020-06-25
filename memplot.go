package memplot

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"time"
)

// Memory data for a given instant
type MemoryInstant struct {
	MemoryInfo *process.MemoryInfoStat
	Instant    time.Duration
}

type MemoryCollection struct {
	Pid            int32
	StartTime      time.Time
	SampleDuration time.Duration // Time between samples
	Samples        []MemoryInstant
}

// Gather a process resident size in memory in batch
func NewMemoryCollection(pid int32, sd, duration time.Duration) (*MemoryCollection, error) {
	numsamples := duration / sd
	if numsamples < 2 {
		return nil, errors.New("There must be at least two samples. Sample Duration too short")
	}

	proc, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	elapsed := time.Since(start)
	var mem *process.MemoryInfoStat
	coll := &MemoryCollection{
		Pid:            pid,
		StartTime:      start,
		SampleDuration: sd,
		Samples:        make([]MemoryInstant, 0),
	}

	for elapsed <= duration {
		elapsed = time.Since(start)
		mem, err = proc.MemoryInfo()
		if err != nil {
			return nil, err
		}

		instant := MemoryInstant{
			MemoryInfo: mem,
			Instant:    elapsed,
		}

		// fmt.Printf("instant=%+v\n", instant)

		coll.Samples = append(coll.Samples, instant)
		time.Sleep(sd)
	}

	return coll, nil
}

// Gather RSS points from a memory collection
func (m *MemoryCollection) GatherRSSXYs() plotter.XYs {
	pts := make(plotter.XYs, len(m.Samples))
	for i, s := range m.Samples {
		pts[i].X = s.Instant.Seconds()
		pts[i].Y = float64(m.Samples[i].MemoryInfo.RSS) / 1024
	}

	return pts
}

// Gather VSZ points from a memory collection
func (m *MemoryCollection) GatherVSZXYs() plotter.XYs {
	pts := make(plotter.XYs, len(m.Samples))
	for i, s := range m.Samples {
		pts[i].X = s.Instant.Seconds()
		pts[i].Y = float64(m.Samples[i].MemoryInfo.VMS) / 1024
	}

	return pts
}

type PlotOptions struct {
	PlotRss bool
	PlotVsz bool
}

// Plot a memory collection
func (m *MemoryCollection) Plot(opt PlotOptions) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}

	p.Title.Text = fmt.Sprintf("Memory Plot of PID %d", m.Pid)
	p.X.Label.Text = "Time (Seconds)"
	p.Y.Label.Text = "KiloBytes"
	// Draw a grid behind the area
	p.Add(plotter.NewGrid())

	if opt.PlotRss {
		// RSS line plotter and style
		rssData := m.GatherRSSXYs()
		rssLine, err := plotter.NewLine(rssData)
		if err != nil {
			return nil, err
		}
		rssLine.LineStyle.Width = vg.Points(1)
		rssLine.LineStyle.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}

		// Add the plotters to the plot, with legend entries
		p.Add(rssLine)
		p.Legend.Add("RSS", rssLine)
	}

	// TODO add another Y axis for vsz
	if opt.PlotVsz {
		// RSS line plotter and style
		vszData := m.GatherVSZXYs()
		vszLine, err := plotter.NewLine(vszData)
		if err != nil {
			return nil, err
		}
		vszLine.LineStyle.Width = vg.Points(1)
		vszLine.LineStyle.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}

		// Add the plotters to the plot, with legend entries
		p.Add(vszLine)
		p.Legend.Add("VSZ", vszLine)
	}

	return p, nil
}

func SavePlot(p *plot.Plot, width, height vg.Length, filename string) error {
	return p.Save(width, height, filename)
}
