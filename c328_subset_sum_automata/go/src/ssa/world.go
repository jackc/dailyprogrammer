package main

type Cell int16

type World struct {
	cells  []Cell
	width  int
	height int
}

func NewWorld(width, height int) *World {
	cells := make([]Cell, width*height)
	return &World{cells: cells, width: width, height: height}
}

func (w *World) Width() int {
	return w.width
}

func (w *World) Height() int {
	return w.height
}

func (w *World) Set(x, y int, val Cell) {
	idx := w.idxFromCoord(x, y)
	w.cells[idx] = val
}

func (w *World) Get(x, y int) Cell {
	idx := w.idxFromCoord(x, y)
	return w.cells[idx]
}

// idxFromCoord takes x and y coordinates and returns the index in w.cells.
// Coordinates wrap the boundaries of the world. e.g. Given World with a
// width of 10, then an x coordinate of -1 should be equal to 9.
func (w *World) idxFromCoord(x, y int) int {
	x = x % w.width
	if x < 0 {
		x += w.width
	}
	y = y % w.height
	if y < 0 {
		y += w.height
	}

	return y*w.width + x
}

func Step(curr, next *World) {
	for y := 0; y < curr.Height(); y++ {
		for x := 0; x < curr.Width(); x++ {
			var newValue Cell
			neighborCount := countNeighbors(curr, x, y)
			if curr.Get(x, y) == 1 {
				if neighborCount == 2 || neighborCount == 3 {
					newValue = 1
				}
			} else {
				if neighborCount == 3 {
					newValue = 1
				}
			}
			next.Set(x, y, newValue)
		}
	}
}

func countNeighbors(w *World, x, y int) int {
	return int(w.Get(x-1, y-1) +
		w.Get(x, y-1) +
		w.Get(x+1, y-1) +

		w.Get(x-1, y) +
		w.Get(x+1, y) +

		w.Get(x-1, y+1) +
		w.Get(x, y+1) +
		w.Get(x+1, y+1))
}
