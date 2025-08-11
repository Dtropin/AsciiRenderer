package main

import (
	"AsciiRenderer/cameracontroller"
	"AsciiRenderer/inputcontoller"
	"AsciiRenderer/mesh"
	rasterization_contoller "AsciiRenderer/rasterization"
	render_context "AsciiRenderer/viewport"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"time"
)

func main() {
	viewPortController := render_context.Init()
	defer viewPortController.Close()

	rawVertices := []mgl32.Vec4{
		// Передняя грань (красный)
		{-0.5, -0.5, 0.5, 1},
		{0.5, -0.5, 0.5, 1},
		{0.5, 0.5, 0.5, 1},
		{-0.5, 0.5, 0.5, 1},

		// Задняя грань (зеленый)
		{-0.5, -0.5, -0.5, 1},
		{0.5, -0.5, -0.5, 1},
		{0.5, 0.5, -0.5, 1},
		{-0.5, 0.5, -0.5, 1},

		// Левая грань (синий)
		{-0.5, -0.5, -0.5, 1},
		{-0.5, -0.5, 0.5, 1},
		{-0.5, 0.5, 0.5, 1},
		{-0.5, 0.5, -0.5, 1},

		// Правая грань (желтый)
		{0.5, -0.5, 0.5, 1},
		{0.5, -0.5, -0.5, 1},
		{0.5, 0.5, -0.5, 1},
		{0.5, 0.5, 0.5, 1},

		// Верхняя грань (пурпурный)
		{-0.5, 0.5, 0.5, 1},
		{0.5, 0.5, 0.5, 1},
		{0.5, 0.5, -0.5, 1},
		{-0.5, 0.5, -0.5, 1},

		// Нижняя грань (голубой)
		{-0.5, -0.5, -0.5, 1},
		{0.5, -0.5, -0.5, 1},
		{0.5, -0.5, 0.5, 1},
		{-0.5, -0.5, 0.5, 1},
	}

	colors := []rune{
		'░', '░', '░', '░', '▒', '▒', '▒', '▒', '▓', '▓', '▓', '▓', '█', '█', '█', '█', '▒', '▒', '▒', '▓', '░', '░', '░', '░',
	}

	var polys = []mesh.Polygon{
		{VerticesIndices: [3]int{0, 1, 2}},
		{VerticesIndices: [3]int{0, 2, 3}},
		{VerticesIndices: [3]int{4, 5, 6}},
		{VerticesIndices: [3]int{4, 6, 7}},
		{VerticesIndices: [3]int{8, 9, 10}},
		{VerticesIndices: [3]int{8, 10, 11}},
		{VerticesIndices: [3]int{12, 13, 14}},
		{VerticesIndices: [3]int{12, 14, 15}},
		{VerticesIndices: [3]int{16, 17, 18}},
		{VerticesIndices: [3]int{16, 18, 19}},
		{VerticesIndices: [3]int{20, 21, 22}},
		{VerticesIndices: [3]int{20, 22, 23}},
	}

	cubeMesh := mesh.Mesh{RawVertices: rawVertices, Polygons: polys, Colors: colors}

	meshController := mesh.Init()
	meshController.AddMesh(&cubeMesh)

	cameraController := cameracontroller.Init()
	cameraController.SetPos(0, 0, 2)

	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()
	var tick = 0
	for {
		select {
		case <-ticker.C:
			tick = inputcontoller.HandleInputKeys(tick, cameraController)
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
			rasterization_contoller.ScanlineRasterization(meshController.Meshes()[0], zbuff, viewPortController)
			viewPortController.Flush()
		}
	}
}
