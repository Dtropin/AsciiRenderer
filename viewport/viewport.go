package viewport

import (
	"github.com/gdamore/tcell/v2"
	"sync"
)

type ViewPortController struct {
	screen tcell.Screen //TODO разбить инпут
}

var (
	instance *ViewPortController
	once     sync.Once
)

func Init() *ViewPortController {
	once.Do(func() {
		instance = &ViewPortController{}
		screen, err := tcell.NewScreen()

		if err != nil {
			panic("Can't init viewport: " + err.Error())
		}

		if err := screen.Init(); err != nil {
			panic(err)
		}
		screen.SetStyle(tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorBlack))
		screen.EnablePaste()
		instance.screen = screen
	})
	return instance
}

func (v *ViewPortController) Clear() {
	w, h := v.screen.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v.screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
		}
	}
}

func (v *ViewPortController) PollEvent() tcell.Event {
	return instance.screen.PollEvent()
}

func (v *ViewPortController) SetChar(x, y int, char rune) {
	v.screen.SetContent(x, y, char, nil, tcell.StyleDefault)
}

func (v *ViewPortController) Flush() {
	v.screen.Show()
}

func (v *ViewPortController) GetWindowSize() (int, int) {
	return v.screen.Size()
}

func (v *ViewPortController) Close() {
	v.screen.Fini()
}
