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

func StepSerial(curr, next *World, target, reward, penalty Cell, max Cell) {
	for y := 0; y < curr.Height(); y++ {
		stepRow(curr, next, target, reward, penalty, max, y)
	}
}

func StepParallel(curr, next *World, target, reward, penalty Cell, max Cell) {
	doneChan := make(chan struct{})

	for y := 0; y < curr.Height(); y++ {
		go func(y int) {
			stepRow(curr, next, target, reward, penalty, max, y)
			doneChan <- struct{}{}
		}(y)
	}

	for y := 0; y < curr.Height(); y++ {
		<-doneChan
	}
}

func stepRow(curr, next *World, target, reward, penalty Cell, max Cell, y int) {
	for x := 0; x < curr.Width(); x++ {
		newValue := curr.Get(x, y)
		if detectSubsetSum(curr, x, y, target) {
			newValue += reward
		} else {
			newValue -= penalty
		}

		if newValue < 0 {
			newValue = 0
		} else if newValue > max {
			newValue = max
		}

		next.Set(x, y, newValue)
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

	i := 0
	j := 1

	for j < 256 {
		iSum := Cell(0)

		if i&(1<<0) != 0 {
			iSum += cells[0]
		}
		if i&(1<<1) != 0 {
			iSum += cells[1]
		}
		if i&(1<<2) != 0 {
			iSum += cells[2]
		}
		if i&(1<<3) != 0 {
			iSum += cells[3]
		}
		if i&(1<<4) != 0 {
			iSum += cells[4]
		}
		if i&(1<<5) != 0 {
			iSum += cells[5]
		}
		if i&(1<<6) != 0 {
			iSum += cells[6]
		}
		if i&(1<<7) != 0 {
			iSum += cells[7]
		}

		if iSum == want {
			return true
		}

		jSum := Cell(0)

		if j&(1<<0) != 0 {
			jSum += cells[0]
		}
		if j&(1<<1) != 0 {
			jSum += cells[1]
		}
		if j&(1<<2) != 0 {
			jSum += cells[2]
		}
		if j&(1<<3) != 0 {
			jSum += cells[3]
		}
		if j&(1<<4) != 0 {
			jSum += cells[4]
		}
		if j&(1<<5) != 0 {
			jSum += cells[5]
		}
		if j&(1<<6) != 0 {
			jSum += cells[6]
		}
		if j&(1<<7) != 0 {
			jSum += cells[7]
		}

		if jSum == want {
			return true
		}

		i += 2
		j += 2
	}

	return false
}
