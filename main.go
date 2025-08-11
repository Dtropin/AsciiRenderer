package main

import (
	"AsciiRenderer/cameracontroller"
	"AsciiRenderer/inputcontoller"
	"AsciiRenderer/mesh"
	"AsciiRenderer/rasterization"
	"AsciiRenderer/viewport"
	"bufio"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	viewPortController := viewport.Init()
	defer viewPortController.Close()

	file, ferr := os.Open("teapot.obj")
	if ferr != nil {
		log.Fatalf("Error opening file: %v", ferr)
		return
	}
	scanner := bufio.NewScanner(file)
	rawVertices := make([]mgl32.Vec4, 0)
	colors := make([]rune, 0)
	colors = append(colors, '#') // cuz of numeration of vertices in teapot.obj
	polys := make([][3]int, 0)
	colormap := []rune{
		'░', '░', '░', '░',
	}
	i := 0
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if line[0][0] == 'v' {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			z, _ := strconv.ParseFloat(line[3], 32)
			rawVertices = append(rawVertices, mgl32.Vec4{float32(x), float32(y), float32(z), 1})
			colors = append(colors, colormap[i%len(colormap)])
		}
		if line[0][0] == 'f' {
			v1, _ := strconv.ParseInt(line[1], 10, 32)
			v2, _ := strconv.ParseInt(line[2], 10, 32)
			v3, _ := strconv.ParseInt(line[3], 10, 32)
			polys = append(polys, [3]int{int(v1 - 1), int(v2 - 1), int(v3 - 1)})
		}
		i++
	}
	ferr = file.Close()
	if ferr != nil {
		log.Fatalf("Error closening file: %v", ferr)
		return
	}

	teapotMesh := mesh.Mesh{RawVertices: rawVertices, Polygons: polys, Colors: colors}

	meshController := mesh.Init()
	meshController.AddMesh(&teapotMesh)

	cameraController := cameracontroller.Init()
	cameraController.SetPos(0, 2, 10)

	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()
	var tick = 0
	var zbuff [][]float32

	for {
		select {
		case <-ticker.C:
			tick = inputcontoller.HandleInputKeys(tick, cameraController)
			viewPortController.Clear()
			windowWidth, windowHeight := viewPortController.GetWindowSize()

			//todo zbuff obj?
			if zbuff == nil || len(zbuff) != windowWidth || len(zbuff[0]) != windowHeight {
				zbuff = make([][]float32, windowWidth+1)

				for i := 0; i < windowWidth+1; i++ {
					zbuff[i] = make([]float32, windowHeight+1)
					for j := 0; j < windowHeight+1; j++ {
						zbuff[i][j] = -math.MaxFloat32
					}
				}
			}

			meshController.ProcessVertices(cameraController, windowWidth, windowHeight, tick%360)
			rasterization.ScanlineRasterization(meshController.Meshes()[0], zbuff, viewPortController)
			viewPortController.Flush()
		}
	}
}
