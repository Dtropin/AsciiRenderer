package mesh_controller

import (
	"AsciiRenderer/camera-controller"
	"AsciiRenderer/mvp"
	"gonum.org/v1/gonum/mat"
	"math"
	"sync"
)

// TODO переименовать в вертекс шейдер?
type MeshController struct {
	vertices          []*mat.Dense
	projectedVertices []int
}

var (
	instance *MeshController
	once     sync.Once
)

func Init() *MeshController {
	once.Do(func() {
		instance = &MeshController{}
	})
	return instance
}

func (m *MeshController) ProcessVertices(camera *camera_controller.CameraController, windowWidth, windowHeight int) {
	xcamera, ycamera, zcamera := camera.GetPos()
	view := mvp.MakeViewMatrix(xcamera, ycamera, zcamera, 0, 0, 0)
	proj := mvp.MakePerspectiveProjection(90, 0.1, float64(windowWidth)/float64(windowHeight), 100)
	model := mvp.MakeModelMatrix(0., 0., 0., 2., 1., 1.)

	//todo clipping
	for i := 0; i < len(m.vertices); i++ {
		var clip_ mat.Dense
		var clip mat.Dense
		var vertex mat.Dense

		vertex.Mul(m.vertices[i], model)
		clip_.Mul(&vertex, view)
		clip.Mul(&clip_, proj)

		x := int(math.Min(float64(windowWidth), (clip.At(0, 0)/clip.At(0, 3)+1)*0.5*float64(windowWidth)))
		y := int(math.Min(float64(windowHeight), (1-(clip.At(0, 1)/clip.At(0, 3)+1)*0.5)*float64(windowHeight)))

		m.projectedVertices[2*i] = x
		m.projectedVertices[2*i+1] = y
	}
}

func (m *MeshController) AddVerticesToMesh(vertices []*mat.Dense) {
	if m.vertices == nil {
		m.vertices = make([]*mat.Dense, len(vertices))
		copy(m.vertices, vertices)
	} else {
		m.vertices = append(m.vertices, vertices...)
	}

	if m.projectedVertices == nil {
		m.projectedVertices = make([]int, len(vertices)*2)
	} else {
		m.projectedVertices = append(m.projectedVertices, make([]int, len(vertices)*2)...)
	}
}

func (m *MeshController) GetProjectedVertices() []int {
	return m.projectedVertices
}
