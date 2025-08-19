//go:build windows

package inputcontoller

import (
	"AsciiRenderer/cameracontroller"
	"syscall"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

func isKeyPressed(vKey int) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(vKey))
	return ret&0x8000 != 0
}

func MouseEvent(cameraController *cameracontroller.CameraController, xDir, yDir int) {
	cameraController.Yaw(float32(xDir))
	cameraController.Pitch(float32(yDir))
}

func HandleInputKeys(controller *cameracontroller.CameraController) {
	if isKeyPressed(0x41) {
		//a
		controller.Left(0.1)
	}
	if isKeyPressed(0x57) {
		//s
		controller.Forward(0.1)
	}
	if isKeyPressed(0x53) {
		//w
		controller.Back(0.1)
	}
	if isKeyPressed(0x44) {
		//d
		controller.Right(0.1)
	}
	if isKeyPressed(0x52) {
		//r
		controller.Up(0.1)
	}
	if isKeyPressed(0x46) {
		//f
		controller.Down(0.1)
	}
}
