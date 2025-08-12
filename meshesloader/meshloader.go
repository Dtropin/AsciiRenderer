package meshesloader

import (
	"AsciiRenderer/mesh"
	"bufio"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"os"
	"strconv"
	"strings"
)

func LoadObjMeshes(filenames []string) []mesh.Mesh {
	var result = make([]mesh.Mesh, len(filenames))

	for _, filename := range filenames {
		file, ferr := os.Open(filename)
		if ferr != nil {
			log.Println("Error opening file: %v", ferr)
			continue
		}

		scanner := bufio.NewScanner(file)

		rawVertices := make([]mgl32.Vec4, 0)
		normals := make([]mgl32.Vec4, 0)
		polys := make([][3]int, 0)
		polysNormals := make([][3]int, 0)

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

		loadedMesh := mesh.Mesh{RawVertices: rawVertices, RawNormals: normals, Polys: polys, PolysNormals: polysNormals,
			Pos: mgl32.Vec3{0., -2., 0}, Scale: mgl32.Vec3{1., 1., 1.}, AngleRad: 0.0, Id: filename}

		result = append(result, loadedMesh)

		ferr = file.Close()

		if ferr != nil {
			log.Println("Error closening file: %v", ferr)
			continue
		}

	}

	return result
}
