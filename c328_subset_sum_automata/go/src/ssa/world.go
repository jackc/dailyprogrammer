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

func Step(curr, next *World, target, reward, penalty Cell) {
	for y := 0; y < curr.Height(); y++ {
		for x := 0; x < curr.Width(); x++ {
			newValue := curr.Get(x, y)
			if detectSubsetSum(curr, x, y, target) {
				newValue += reward
			} else {
				newValue -= penalty
			}
			next.Set(x, y, newValue)
		}
	}
}

func detectSubsetSum(w *World, x, y int, want Cell) bool {
	cells := [8]Cell{
		w.Get(x-1, y-1),
		w.Get(x, y-1),
		w.Get(x+1, y-1),

		w.Get(x-1, y),
		w.Get(x+1, y),

		w.Get(x-1, y+1),
		w.Get(x, y+1),
		w.Get(x+1, y+1),
	}

	for i := 1; i < 256; i++ {
		sum := Cell(0)
		for j := uint(0); j < uint(len(cells)); j++ {
			if i&(1<<j) != 0 {
				sum += cells[j]
			}
		}

		if sum == want {
			return true
		}
	}

	return false
}
