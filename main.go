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
	normals := make([]mgl32.Vec4, 0)
	colors := make([]rune, 0)
	colors = append(colors, '#') // cuz of numeration of vertices in teapot.obj
	polys := make([][3]int, 0)
	polysNormals := make([][3]int, 0)
	colorMap := []rune{
		'$', '@', 'B', '%', '8', '&', 'W', 'M', '#', '*', 'o', 'a', 'h', 'k', 'b', 'd', 'p', 'q', 'w', 'm', 'Z', 'O', '0', 'Q', 'L', 'C', 'J', 'U', 'Y', 'X', 'z', 'c', 'v', 'u', 'n', 'x', 'r', 'j', 'f', 't', '/', '\\', '|', '(', ')', '1', '{', '}', '[', ']', '?', '-', '_', '+', '~', '<', '>', 'i', '!', 'l', 'I', ';', ':', ',', '.', '"', '^', '`', '\'',
	}
	i := 0
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")

		if strings.HasPrefix(line[0], "vn") {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			z, _ := strconv.ParseFloat(line[3], 32)
			normal := mgl32.Vec4{float32(x), float32(y), float32(z), 0}
			normals = append(normals, normal.Mul(-1))
			continue
		}

		if line[0][0] == 'v' {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			z, _ := strconv.ParseFloat(line[3], 32)
			rawVertices = append(rawVertices, mgl32.Vec4{float32(x), float32(y), float32(z), 1})
			i++
			continue
		}

		if line[0][0] == 'f' {
			v1, _ := strconv.ParseInt(strings.Split(line[1], "//")[0], 10, 32)
			v2, _ := strconv.ParseInt(strings.Split(line[2], "//")[0], 10, 32)
			v3, _ := strconv.ParseInt(strings.Split(line[3], "//")[0], 10, 32)
			vn1, _ := strconv.ParseInt(strings.Split(line[1], "//")[1], 10, 32)
			vn2, _ := strconv.ParseInt(strings.Split(line[2], "//")[1], 10, 32)
			vn3, _ := strconv.ParseInt(strings.Split(line[3], "//")[1], 10, 32)
			polys = append(polys, [3]int{int(v1 - 1), int(v2 - 1), int(v3 - 1)})
			polysNormals = append(polysNormals, [3]int{int(vn1 - 1), int(vn2 - 1), int(vn3 - 1)})
			continue
		}
	}
	ferr = file.Close()
	if ferr != nil {
		log.Fatalf("Error closening file: %v", ferr)
		return
	}

	teapotMesh := mesh.Mesh{RawVertices: rawVertices, RawNormals: normals, Polys: polys, PolysNormals: polysNormals}

	meshController := mesh.Init()
	meshController.AddMesh(&teapotMesh)
	cameraController := cameracontroller.Init()

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

			light := meshController.ProcessVertices(cameraController, windowWidth, windowHeight, tick%360)
			rasterization.ScanlineRasterization(meshController.Meshes()[0], zbuff, viewPortController, light, colorMap)
			viewPortController.Flush()
		}
	}
}
