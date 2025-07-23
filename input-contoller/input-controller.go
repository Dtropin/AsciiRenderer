package input_contoller

import (
	camera_controller "AsciiRenderer/camera-controller"
	"github.com/nsf/termbox-go"
)

func HandleInputs(cameraController *camera_controller.CameraController) {
	for {
		ev := termbox.PollEvent()

		//TODO как то синхронизовать нажатия клавиш и рендер, иначе теряешь управления при потере фпс
		if ev.Type == termbox.EventKey {
			/*if ev.Ch == 'q' {
				yaw += 0.01
			}
			if ev.Ch == 'e' {
				yaw -= 0.01
			}*/
			switch ev.Key {
			case termbox.KeyArrowLeft:
				//TODO mutex?
				cameraController.AdjustPos(0.1, 0, 0)
			case termbox.KeyArrowRight:
				cameraController.AdjustPos(-0.1, 0, 0)
			case termbox.KeyArrowUp:
				cameraController.AdjustPos(0, 0, 0.1)
			case termbox.KeyArrowDown:
				cameraController.AdjustPos(0, 0, -0.1)
			case termbox.KeyEsc:
				return
			case termbox.KeyCtrlC:
				return
			}

		}
	}
}
