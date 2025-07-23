package main

import (
	"AsciiRenderer/bresenham"
	camera_controller "AsciiRenderer/camera-controller"
	input_contoller "AsciiRenderer/input-contoller"
	"AsciiRenderer/mesh-controller"
	"AsciiRenderer/terminal-context"
	"gonum.org/v1/gonum/mat"
	"time"
)

type Edge [2]int

func main() {
	render_context.Init()
	defer render_context.Close()

	halfSize := 0.5
	rawVertices := []*mat.Dense{
		mat.NewDense(1, 4, []float64{-halfSize, -halfSize, -halfSize, 1}),
		mat.NewDense(1, 4, []float64{halfSize, -halfSize, -halfSize, 1}), // 1: правый-нижний-задний
		mat.NewDense(1, 4, []float64{halfSize, halfSize, -halfSize, 1}),  // 2: правый-верхний-задний
		mat.NewDense(1, 4, []float64{-halfSize, halfSize, -halfSize, 1}), // 3: левый-верхний-задний
		mat.NewDense(1, 4, []float64{-halfSize, -halfSize, halfSize, 1}), // 4: левый-нижний-передний
		mat.NewDense(1, 4, []float64{halfSize, -halfSize, halfSize, 1}),  // 5: правый-нижний-передний
		mat.NewDense(1, 4, []float64{halfSize, halfSize, halfSize, 1}),   // 6: правый-верхний-передний
		mat.NewDense(1, 4, []float64{-halfSize, halfSize, halfSize, 1}),  // 7: левый-верхний-передний
	}

	meshController := mesh_controller.Init()
	meshController.AddVerticesToMesh(rawVertices)

	cameraController := camera_controller.Init()
	cameraController.SetPos(0, 0, -15)

	go input_contoller.HandleInputs(cameraController)

	edges := []Edge{
		{0, 1}, {1, 2}, {2, 3}, {3, 0}, // Нижняя грань
		{4, 5}, {5, 6}, {6, 7}, {7, 4}, // Верхняя грань
		{0, 4}, {1, 5}, {2, 6}, {3, 7}, // Вертикальные рёбра
	}

	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			render_context.Clear()
			windowWidth, windowHeight := render_context.GetWindowSize()
			meshController.ProcessVertices(cameraController, windowWidth, windowHeight)
			projectedVertices := meshController.GetProjectedVertices()

			for i := 0; i < len(projectedVertices); i = i + 2 {
				render_context.SetChar(projectedVertices[i], projectedVertices[i+1], '*')
			}

			//todo rasterization-contoller
			for _, edge := range edges {
				bresenham.DrawLine('-', '|', '/', '\\', projectedVertices[edge[0]*2], projectedVertices[edge[0]*2+1], projectedVertices[edge[1]*2], projectedVertices[edge[1]*2+1])
			}

			render_context.Flush()
		}
	}
}
