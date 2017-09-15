package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var options struct {
	width        int
	height       int
	target       int
	reward       int
	penalty      int
	max          int
	stepTime     int
	seed         int
	parallel     bool
	windowWidth  int
	windowHeight int
}

const (
	WindowTitle = "C328 Subset Sum Automata"
)

func main() {
	runtime.LockOSThread()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&options.width, "width", 30, "width of world in cells")
	flag.IntVar(&options.height, "height", 30, "height of world in cells")
	flag.IntVar(&options.target, "target", 5, "target value")
	flag.IntVar(&options.reward, "reward", 3, "reward value")
	flag.IntVar(&options.penalty, "penalty", 1, "penalty value")
	flag.IntVar(&options.max, "max", 15, "max value")
	flag.IntVar(&options.stepTime, "steptime", 250, "time per step in milliseconds")
	flag.IntVar(&options.seed, "seed", -1, "seed")
	flag.BoolVar(&options.parallel, "parallel", false, "run in parallel")
	flag.IntVar(&options.windowWidth, "windowwidth", 800, "width of window in pixels")
	flag.IntVar(&options.windowHeight, "windowheight", 600, "height of window in pixels")
	flag.Parse()

	if options.seed < 0 {
		options.seed = time.Now().Nanosecond()
	}

	w := NewWorld(options.width, options.height)
	wScratch := NewWorld(options.width, options.height)

	stepFn := StepSerial
	if options.parallel {
		stepFn = StepParallel
	}

	rand.Seed(int64(options.seed))
	for y := 0; y < w.Height(); y++ {
		for x := 0; x < w.Width(); x++ {
			w.Set(x, y, Cell(rand.Intn(options.target+options.reward)))
		}
	}

	var window *sdl.Window
	var renderer *sdl.Renderer
	var err error

	cellWidth := options.windowWidth / options.width
	cellHeight := options.windowHeight / options.height

	window, err = sdl.CreateWindow(WindowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, options.windowWidth, options.windowHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(1)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(2)
	}
	defer renderer.Destroy()

	renderer.Clear()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		sdl.Delay(uint32(options.stepTime))

		renderer.Clear()
		renderer.SetDrawColor(0, 0, 0, 0x20)
		renderer.FillRect(&sdl.Rect{0, 0, int32(options.windowWidth), int32(options.windowHeight)})

		for y := 0; y < w.Height(); y++ {
			for x := 0; x < w.Width(); x++ {
				cell := w.Get(x, y)
				renderer.SetDrawColor(uint8(cell*16), uint8(cell*16), uint8(cell*16), 0x20)
				renderer.FillRect(&sdl.Rect{int32(x * cellWidth), int32(y * cellHeight), int32((x + 1) * cellWidth), int32((y + 1) * cellHeight)})
			}
		}

		renderer.Present()

		stepFn(w, wScratch, Cell(options.target), Cell(options.reward), Cell(options.penalty), Cell(options.max))
		w, wScratch = wScratch, w

	}
}
