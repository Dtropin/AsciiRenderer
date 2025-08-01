//go:build windows

package input_contoller

import (
	camera_controller "AsciiRenderer/camera-controller"
	"syscall"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

/*type KeyStateHolder struct {
	mu        sync.RWMutex
	keyStates *list.List
	//TODO
	PressedKeys map[rune]bool
}

type KeyState struct {
	char            rune
	lastTimePressed int64
}

func NewKeyStateHolder() *KeyStateHolder {
	return &KeyStateHolder{
		keyStates:   list.New(),
		PressedKeys: make(map[rune]bool),
	}
}*/

// todo интерфейс ивент сурс
/*
func StartConsumingInputEvents(keyStateHolder *KeyStateHolder, viewPortController *render_context.ViewPortController) {

	for {
		switch ev := viewPortController.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyRune {
				fmt.Println(ev.Rune())
				keyStateHolder.mu.Lock()
				if !keyStateHolder.PressedKeys[ev.Rune()] {
					keyStateHolder.PressedKeys[ev.Rune()] = true
					keyStateHolder.keyStates.PushBack(KeyState{ev.Rune(), time.Now().UnixMilli()})
				}
				keyStateHolder.mu.Unlock()
			}
		}
	}
}
*/

func isKeyPressed(vKey int) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(vKey))
	return ret&0x8000 != 0
}

func HandleInputKeys(controller *camera_controller.CameraController) {
	if isKeyPressed(0x41) {
		controller.AdjustPos(-0.1, 0, 0)
	}
	if isKeyPressed(0x57) {
		controller.AdjustPos(0, 0, -0.1)
	}
	if isKeyPressed(0x53) {
		controller.AdjustPos(0, 0, 0.1)
	}
	if isKeyPressed(0x44) {
		controller.AdjustPos(0.1, 0, 0)
	}
}
