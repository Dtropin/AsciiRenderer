package main

import (
	"AsciiRenderer/cameracontroller"
	"AsciiRenderer/inputcontoller"
	"AsciiRenderer/mesh"
	"AsciiRenderer/meshesloader"
	"AsciiRenderer/rasterization"
	"AsciiRenderer/viewport"
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

	var tick = 0

	for {
		select {
		case <-ticker.C:
			tick = inputcontoller.HandleInputKeys(tick, cameraController)
			viewPortController.Clear()
			meshController.ProcessVertices(cameraController, viewPortController)
			rasterizer.ScanlineRasterization(meshController.Meshes(), viewPortController, cameraController)
			viewPortController.Flush()
		}
	}
}
