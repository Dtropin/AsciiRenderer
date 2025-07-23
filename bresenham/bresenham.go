package bresenham

import (
	"github.com/nsf/termbox-go"
	"math"
)

func DrawLine(sidechar, upchar, diagupchar, diagdownchar byte, x0, y0, x1, y1 int) {
	if x0 == x1 && y0 == y1 {
		return
	}

	dx, dy := math.Abs(float64(x1-x0)), -math.Abs(float64(y1-y0))

	sx := 1

	if x0 >= x1 {
		sx = -1
	}

	sy := 1

	if y0 >= y1 {
		sy = -1
	}

	err := dx + dy

	var dirmap = make(map[[2]int]byte)
	dirmap[[2]int{1, 1}] = diagdownchar
	dirmap[[2]int{1, -1}] = diagupchar
	dirmap[[2]int{-1, -1}] = diagdownchar
	dirmap[[2]int{-1, 1}] = diagupchar
	dirmap[[2]int{0, 1}] = upchar
	dirmap[[2]int{0, -1}] = upchar
	dirmap[[2]int{1, 0}] = sidechar
	dirmap[[2]int{-1, 0}] = sidechar

	var dir = [2]int{0, 0}

	for {
		tmperr := 2 * err

		dir[0] = 0
		dir[1] = 0

		if tmperr >= dy {
			err += dy
			x0 += sx
			dir[0] = sx
		}
		if tmperr <= dx {
			err += dx
			y0 += sy
			dir[1] = sy
		}

		if (x0 == x1 && y0 == y1) || (tmperr < dy && tmperr > dx) {
			return
		}

		//TODO сделать абстракцию через интерфейс
		termbox.SetCell(x0, y0, rune(dirmap[dir]), termbox.ColorWhite, termbox.ColorDefault)
	}
}
