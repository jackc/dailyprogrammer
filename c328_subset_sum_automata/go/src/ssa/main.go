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
	width    int
	height   int
	target   int
	reward   int
	penalty  int
	stepTime int
	seed     int
}

func Print(w *World, wr io.Writer) {
	values := map[Cell]termbox.Cell{
		-9: termbox.Cell{Ch: '9', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
		-8: termbox.Cell{Ch: '8', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
		-7: termbox.Cell{Ch: '7', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
		-6: termbox.Cell{Ch: '6', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
		-5: termbox.Cell{Ch: '5', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
		-4: termbox.Cell{Ch: '4', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
		-3: termbox.Cell{Ch: '3', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
		-2: termbox.Cell{Ch: '2', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
		-1: termbox.Cell{Ch: '1', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
		0:  termbox.Cell{Ch: '0', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		1:  termbox.Cell{Ch: '1', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		2:  termbox.Cell{Ch: '2', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		3:  termbox.Cell{Ch: '3', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		4:  termbox.Cell{Ch: '4', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		5:  termbox.Cell{Ch: '5', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		6:  termbox.Cell{Ch: '6', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		7:  termbox.Cell{Ch: '7', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		8:  termbox.Cell{Ch: '8', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		9:  termbox.Cell{Ch: '9', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		10: termbox.Cell{Ch: 'A', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		11: termbox.Cell{Ch: 'B', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		12: termbox.Cell{Ch: 'C', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		13: termbox.Cell{Ch: 'D', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		14: termbox.Cell{Ch: 'E', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
		15: termbox.Cell{Ch: 'F', Fg: termbox.ColorBlack, Bg: termbox.ColorWhite},
	}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for y := 0; y < w.Height(); y++ {
		for x := 0; x < w.Width(); x++ {
			cell := w.Get(x, y)
			var termCell termbox.Cell
			if cell < -9 {
				termCell = values[-9]
			} else if cell > 15 {
				termCell = values[15]
			} else {
				termCell = values[cell]
			}

			termbox.SetCell(x*2, y, termCell.Ch, termCell.Fg, termCell.Bg)
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
	flag.IntVar(&options.stepTime, "steptime", 250, "time per step in milliseconds")
	flag.IntVar(&options.seed, "seed", -1, "seed")
	flag.Parse()

	if options.seed < 0 {
		options.seed = time.Now().Nanosecond()
	}

	rand.Seed(int64(options.seed))
	err := termbox.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer termbox.Close()

	w := NewWorld(options.width, options.height)
	wScratch := NewWorld(options.width, options.height)

	for y := 0; y < w.Height(); y++ {
		for x := 0; x < w.Width(); x++ {
			w.Set(x, y, Cell(rand.Intn(options.target+options.reward)))
		}
	}

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
			Step(w, wScratch, Cell(options.target), Cell(options.reward), Cell(options.penalty))
			w, wScratch = wScratch, w
		}
	}
}
