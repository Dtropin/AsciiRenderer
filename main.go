package main

import (
	camera_controller "AsciiRenderer/camera-controller"
	input_contoller "AsciiRenderer/input-contoller"
	"AsciiRenderer/mesh-controller"
	rasterization_contoller "AsciiRenderer/rasterization-contoller"
	render_context "AsciiRenderer/terminal-context"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"time"
)

func main() {
	viewPortController := render_context.Init()
	defer viewPortController.Close()

	var halfSize float32 = 0.5
	rawVertices := []mgl32.Vec4{
		mgl32.Vec4{-halfSize, -halfSize, -halfSize, 1},
		mgl32.Vec4{halfSize, -halfSize, -halfSize, 1},
		mgl32.Vec4{halfSize, halfSize, -halfSize, 1},
		mgl32.Vec4{-halfSize, halfSize, -halfSize, 1},

		mgl32.Vec4{-halfSize, -halfSize, halfSize, 1},
		mgl32.Vec4{halfSize, -halfSize, halfSize, 1},
		mgl32.Vec4{halfSize, halfSize, halfSize, 1},
		mgl32.Vec4{-halfSize, halfSize, halfSize, 1},
	}

	meshController := mesh_controller.Init()
	meshController.AddVerticesToMesh(rawVertices)

	cameraController := camera_controller.Init()
	cameraController.SetPos(0, 0, 5)

	//TODO общую структуру для вершин
	colors := []rune{
		'░', '▒', '▓', '█', '╳', '│', '┼', 'O',
	}

	polys := [][]int{
		{0, 1, 2}, {0, 2, 3}, {1, 5, 6}, {1, 6, 2},
		{5, 4, 7}, {5, 7, 6}, {4, 0, 3}, {4, 3, 7},
		{3, 2, 6}, {3, 6, 7}, {4, 5, 1}, {4, 1, 0},
	}

	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()
	var tick = 0
	for {
		select {
		case <-ticker.C:
			input_contoller.HandleInputKeys(cameraController)
			viewPortController.Clear()
			windowWidth, windowHeight := viewPortController.GetWindowSize()

			//TODO отслеживать изменения окна и если было - пересоздавать буфер
			w, h := viewPortController.GetWindowSize()
			zbuff := make([][]float32, w+1)

			for i := 0; i < w+1; i++ {
				zbuff[i] = make([]float32, h+1)
				for j := 0; j < h+1; j++ {
					zbuff[i][j] = -math.MaxFloat32
				}
			}

			meshController.ProcessVertices(cameraController, windowWidth, windowHeight, tick%360)
			tick = tick + 1
			rasterization_contoller.ScanlineRasterization(polys, meshController.GetProjectedVertices(), zbuff, colors, viewPortController)
			viewPortController.Flush()
		}
	}
}
