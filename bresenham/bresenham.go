package bresenham

import (
	render_context "AsciiRenderer/terminal-context"
	"math"
)

func DrawLine(viewPortController *render_context.ViewPortController, downsidechar, upsidechar, upchar, diagupchar, diagdownchar rune, x0, y0, x1, y1 int) {
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

	var dirmap = make(map[[2]int]rune)
	dirmap[[2]int{1, 1}] = diagdownchar
	dirmap[[2]int{1, -1}] = diagupchar
	dirmap[[2]int{-1, -1}] = diagdownchar
	dirmap[[2]int{-1, 1}] = diagupchar
	dirmap[[2]int{0, 1}] = upchar
	dirmap[[2]int{0, -1}] = upchar
	dirmap[[2]int{1, 0}] = downsidechar
	dirmap[[2]int{-1, 0}] = downsidechar

	var dir = [2]int{0, 0}
	var prevchar *rune

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

		var char = dirmap[dir]

		if char == downsidechar {
			if prevchar != nil {
				if *prevchar != downsidechar {
					//TODO map
					if *prevchar == upsidechar {
						char = upsidechar
					}
					if *prevchar == upchar {
						char = downsidechar
					}
					if *prevchar == diagupchar {
						char = upsidechar
					}
					if *prevchar == diagdownchar {
						char = downsidechar
					}
				}
			}
		}

		//TODO сделать абстракцию через интерфейс
		viewPortController.SetChar(x0, y0, char)
		prevchar = &char
	}
}
