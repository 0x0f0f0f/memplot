# memplot

A small utility written in golang to quickly plot memory usage of processes.
`memplot` constantly samples memory usage of a process, for a given 
amount of time and then outputs an image file. Painless and straightforward.
Supported image extensions are: .eps, .jpg, .jpeg, .pdf, .png, .svg, .tex, .tif and .tiff

## Installation
Binaries for Linux and MacOS are available in the
[releases](https://github.com/0x0f0f0f/memplot/releases) page
```
go get -u -v github.com/0x0f0f0f/memplot/cmd
```

## Usage

```
Usage of memplot:
Arguments following options will be interpreted as the command to spawn and sample
  -dur duration
    	total profiling time. a value of 0 means that the program
    	will be sampled until it is no longer running
  -height string
    	plot image height (can be cm or in) (default "12cm")
  -o string
    	output image file name. supported extensions are:
    	.eps, .jpg, .jpeg, .pdf, .png, .svg, .tex, .tif and .tiff
    	(default "output-plot.png")
  -pid int
    	pid of the process to analyze (default -1)
  -sd duration
    	sample size in time (default 5ms)
  -vsz
    	plot virtual size
  -width string
    	plot image width (can be cm or in) (default "16cm")
```

## Example Plot 
```
memplot -pid 25273 -width 8in -height 8in -dur 60s -sd 50ms -o plot.png
```
or 
```
memplot -width 8in -height 8in -dur 60s -sd 50ms -o plot.png firefox
```
![](https://raw.githubusercontent.com/0x0f0f0f/memplot/master/plot.png)

## Roadmap 
- [ ] Plot other useful values
- [ ] Colors
- [ ] Option to change the font
- [ ] Memory autoscaling
- [X] Spawn child processes to sample
