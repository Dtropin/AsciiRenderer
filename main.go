package main

import (
	"AsciiRenderer/cameracontroller"
	"AsciiRenderer/inputcontoller"
	"AsciiRenderer/mesh"
	"AsciiRenderer/meshesloader"
	"AsciiRenderer/rasterization"
	"AsciiRenderer/viewport"
	"github.com/gdamore/tcell/v2"
	"math"
	"time"
)

func main() {
	viewPortController := viewport.Init()
	defer viewPortController.Close()

	sceneMeshes := meshesloader.LoadObjMeshes([]string{"teapot.obj"})

	meshController := mesh.Init()
	meshController.AddMeshes(sceneMeshes)

	cameraController := cameracontroller.Init()

	colorMap := []rune{
		'$', '@', 'B', '%', '8', '&', 'W', 'M', '#', '*', 'o', 'a', 'h', 'k', 'b', 'd', 'p', 'q', 'w', 'm', 'Z', 'O', '0', 'Q', 'L', 'C', 'J', 'U', 'Y', 'X', 'z', 'c', 'v', 'u', 'n', 'x', 'r', 'j', 'f', 't', '/', '\\', '|', '(', ')', '1', '{', '}', '[', ']', '?', '-', '_', '+', '~', '<', '>', 'i', '!', 'l', 'I', ';', ':', ',', '.', '"', '^', '`', '\'',
	}

	rasterizer := rasterization.Init(colorMap, viewPortController)

	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	inputEvents := make(chan tcell.Event)
	prevMouseX, prevMouseY := 0, 0

	go func() {
		for {
			inputEvents <- viewPortController.PollEvent()
		}
	}()

	for {
		select {
		case <-ticker.C:
			select {
			case ev := <-inputEvents:
				switch evt := ev.(type) {
				case *tcell.EventMouse:
					mouseX, mouseY := evt.Position()
					var xDir, yDir int

					if mouseX != 0 {
						xDir = mouseX - prevMouseX
						if xDir != 0 {
							xDir /= int(math.Abs(float64(xDir)))
						}
					}

					if mouseY != 0 {
						yDir = mouseY - prevMouseY
						if yDir != 0 {
							yDir /= int(math.Abs(float64(yDir)))
						}
					}

					inputcontoller.MouseEvent(cameraController, xDir, yDir)
					prevMouseX = mouseX
					prevMouseY = mouseY
				}
			default:
			}
		}
		w, h := viewPortController.GetWindowSize()
		if prevMouseX <= 0 {
			inputcontoller.MouseEvent(cameraController, -1, 0)
		}
		if prevMouseX >= w-1 {
			inputcontoller.MouseEvent(cameraController, 1, 0)
		}
		if prevMouseY <= 0 {
			inputcontoller.MouseEvent(cameraController, 0, -1)
		}
		if prevMouseY >= h-1 {
			inputcontoller.MouseEvent(cameraController, 0, 1)
		}

		inputcontoller.HandleInputKeys(cameraController)
		viewPortController.Clear()
		meshController.ProcessVertices(cameraController, viewPortController)
		rasterizer.ScanlineRasterization(meshController.Meshes(), viewPortController, cameraController)
		viewPortController.Flush()
	}
}
