# memplot

A small utility written in golang to quickly plot memory usage of processes.
Still in a very early stage.
`memplot` constantly samples memory usage of a process, for a given 
duration of time and then outputs a .png file. Painless and straightforward.

## Installation
```
go get -u -v github.com/0x0f0f0f/memplot/cmd
```

## Usage

```
Usage of memplot:
Any argument following the options will be interpreted as the command to spawn and sample
  -dur duration
    	total profiling time (default 10s)
  -height string
    	plot image height (can be cm or in) (default "12cm")
  -o string
    	output image file name (default "output-plot.png")
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
![](https://raw.githubusercontent.com/0x0f0f0f/memplot/master/plot.png)
