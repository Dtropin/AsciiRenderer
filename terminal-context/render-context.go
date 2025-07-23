package render_context

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

func Init() {
	err := termbox.Init()

	if err != nil {
		panic("Can't init termbox: " + err.Error())
	}

	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func SetChar(x, y int, char rune) {
	termbox.SetCell(x, y, char, termbox.ColorWhite, termbox.ColorDefault)
}

func Flush() {
	if flusherr := termbox.Flush(); flusherr != nil {
		fmt.Println("Error:", flusherr)
	}
}

func GetWindowSize() (int, int) {
	return termbox.Size()
}

func Close() {
	termbox.Close()
}
