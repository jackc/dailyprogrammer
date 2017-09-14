package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var options struct {
	width    int
	height   int
	target   int
	reward   int
	penalty  int
	max      int
	stepTime int
	seed     int
	parallel bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&options.width, "width", 30, "width of world")
	flag.IntVar(&options.height, "height", 30, "height of world")
	flag.IntVar(&options.target, "target", 5, "target value")
	flag.IntVar(&options.reward, "reward", 3, "reward value")
	flag.IntVar(&options.penalty, "penalty", 1, "penalty value")
	flag.IntVar(&options.max, "max", 15, "max value")
	flag.IntVar(&options.stepTime, "steptime", 250, "time per step in milliseconds")
	flag.IntVar(&options.seed, "seed", -1, "seed")
	flag.BoolVar(&options.parallel, "parallel", false, "run in parallel")
	flag.Parse()

	if options.seed < 0 {
		options.seed = time.Now().Nanosecond()
	}

	w := NewWorld(options.width, options.height)
	wScratch := NewWorld(options.width, options.height)

	// stepFn := StepSerial
	// if options.parallel {
	// 	stepFn = StepParallel
	// }

	rand.Seed(int64(options.seed))
	for y := 0; y < w.Height(); y++ {
		for x := 0; x < w.Width(); x++ {
			w.Set(x, y, Cell(rand.Intn(options.target+options.reward)))
		}
	}

	sdl.Main(func() {
		os.Exit(run(w, wScratch))
	})

	// stepFn(w, wScratch, Cell(options.target), Cell(options.reward), Cell(options.penalty), Cell(options.max))
	// w, wScratch = wScratch, w
}

const (
	WindowTitle  = "Go-SDL2 Render"
	WindowWidth  = 800
	WindowHeight = 600
	FrameRate    = 60

	RectWidth  = 20
	RectHeight = 20
	NumRects   = WindowHeight / RectHeight
)

var rects [NumRects]sdl.Rect
var runningMutex sync.Mutex

func run(w, wScratch *World) int {

	var window *sdl.Window
	var renderer *sdl.Renderer
	var err error

	sdl.Do(func() {
		window, err = sdl.CreateWindow(WindowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_OPENGL)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer func() {
		sdl.Do(func() {
			window.Destroy()
		})
	}()

	sdl.Do(func() {
		renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	})
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer func() {
		sdl.Do(func() {
			renderer.Destroy()
		})
	}()

	sdl.Do(func() {
		renderer.Clear()
	})

	for i := range rects {
		rects[i] = sdl.Rect{
			X: int32(rand.Int() % WindowWidth),
			Y: int32(i * WindowHeight / len(rects)),
			W: RectWidth,
			H: RectHeight,
		}
	}

	running := true
	for running {
		sdl.Do(func() {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					runningMutex.Lock()
					running = false
					runningMutex.Unlock()
				}
			}

			renderer.Clear()
			renderer.SetDrawColor(0, 0, 0, 0x20)
			renderer.FillRect(&sdl.Rect{0, 0, WindowWidth, WindowHeight})
		})

		// Do expensive stuff using goroutines
		wg := sync.WaitGroup{}
		for i := range rects {
			wg.Add(1)
			go func(i int) {
				rects[i].X = (rects[i].X + 10) % WindowWidth
				sdl.Do(func() {
					renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
					renderer.DrawRect(&rects[i])
				})
				wg.Done()
			}(i)
		}
		wg.Wait()

		sdl.Do(func() {
			renderer.Present()
			sdl.Delay(1000 / FrameRate)
		})
	}

	return 0
}
