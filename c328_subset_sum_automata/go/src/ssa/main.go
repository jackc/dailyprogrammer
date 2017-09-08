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
	cellCount int
	stepTime  int
	seed      int
}

func Print(w *World, wr io.Writer) {
	deadCell := termbox.Cell{Ch: ' ', Fg: termbox.ColorBlack, Bg: termbox.ColorBlack}
	liveCell := termbox.Cell{Ch: '🔶', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for y := 0; y < w.Height(); y++ {
		for x := 0; x < w.Width(); x++ {
			cell := deadCell
			if w.Get(x, y) == 1 {
				cell = liveCell
			}

			termbox.SetCell(x*2, y, cell.Ch, cell.Fg, cell.Bg)
			termbox.SetCell(x*2+1, y, ' ', termbox.ColorBlack, termbox.ColorBlack)
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
	flag.IntVar(&options.cellCount, "livecount", 100, "live cells to start")
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

	for i := 0; i < options.cellCount; i++ {
		w.Set(rand.Intn(w.Width()), rand.Intn(w.Height()), 1)
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
			Step(w, wScratch)
			w, wScratch = wScratch, w
		}
	}
}
