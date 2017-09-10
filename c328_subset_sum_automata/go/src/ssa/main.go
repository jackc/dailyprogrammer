package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

var options struct {
	width     int
	height    int
	target    int
	reward    int
	penalty   int
	max       int
	stepTime  int
	seed      int
	benchmark int
	parallel  bool
}

func Print(w *World, wr io.Writer) {
	values := [16]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for y := 0; y < w.Height(); y++ {
		for x := 0; x < w.Width(); x++ {
			cell := w.Get(x, y)
			var ch rune
			if cell < 0 {
				ch = '0'
			} else if cell > 15 {
				ch = 'F'
			} else {
				ch = values[cell]
			}

			termbox.SetCell(x*2, y, ch, termbox.ColorWhite, termbox.ColorBlack)
		}
	}

	termbox.Flush()
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
	flag.IntVar(&options.benchmark, "benchmark", 0, "run benchmark of N steps and quit")
	flag.IntVar(&options.seed, "seed", -1, "seed")
	flag.BoolVar(&options.parallel, "parallel", false, "run in parallel")
	flag.Parse()

	if options.seed < 0 {
		options.seed = time.Now().Nanosecond()
	}

	rand.Seed(int64(options.seed))

	w := NewWorld(options.width, options.height)
	wScratch := NewWorld(options.width, options.height)

	for y := 0; y < w.Height(); y++ {
		for x := 0; x < w.Width(); x++ {
			w.Set(x, y, Cell(rand.Intn(options.target+options.reward)))
		}
	}

	stepFn := StepSerial
	if options.parallel {
		stepFn = StepParallel
	}

	if options.benchmark > 0 {
		for i := 0; i < options.benchmark; i++ {
			stepFn(w, wScratch, Cell(options.target), Cell(options.reward), Cell(options.penalty), Cell(options.max))
			w, wScratch = wScratch, w
		}
		fmt.Println(w.Get(0, 0)) // Access results to ensure entire calculation cannot be removed by the optimizer.
		os.Exit(0)
	}

	err := termbox.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(time.Duration(options.stepTime) * time.Millisecond)

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch {
				case ev.Ch == 'q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC:
					return
				}
			}
		case <-ticker.C:
			Print(w, os.Stdout)
			stepFn(w, wScratch, Cell(options.target), Cell(options.reward), Cell(options.penalty), Cell(options.max))
			w, wScratch = wScratch, w
		}
	}
}
