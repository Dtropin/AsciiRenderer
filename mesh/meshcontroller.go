package mesh

import (
	"AsciiRenderer/cameracontroller"
	"AsciiRenderer/viewport"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"sync"
)

type MeshController struct {
	meshes []*Mesh
}

func (m *MeshController) Meshes() []*Mesh {
	return m.meshes
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

func (m *MeshController) ProcessVertices(camera *cameracontroller.CameraController, viewPortController *viewport.ViewPortController) {
	view := camera.CameraState().ViewMatrix()
	w, h := viewPortController.GetWindowSize()

	proj := mgl32.Perspective((90*math.Pi)/180.0, float32(w)*0.6/float32(h), 0.1, 100)

	//todo clipping
	for j := 0; j < len(m.meshes); j++ {
		m.meshes[j].UpdateModelMatrix()
		model := m.meshes[j].ModelMatrix
		modelViewMatrix := view.Mul4(model)
		normalModelView := modelViewMatrix.Inv().Transpose()
		for i := 0; i < len(m.meshes[j].RawVertices); i++ {
			var modelSpace = model.Mul4x1(m.meshes[j].RawVertices[i])
			var cameraSpace = view.Mul4x1(modelSpace)
			var clipSpace = proj.Mul4x1(cameraSpace)

			ndcX := clipSpace.X() / clipSpace.W()
			ndcY := clipSpace.Y() / clipSpace.W()

			x := math.Min(float64(w), (float64(ndcX)+1)*0.5*float64(w))
			y := math.Min(float64(h), (1-(float64(ndcY)+1)*0.5)*float64(h))

			m.meshes[j].ProcessedVertices[i] = ProcessedVertex{xScreen: float32(x), yScreen: float32(y), xSource: m.meshes[j].RawVertices[i].X(), ySource: m.meshes[j].RawVertices[i].Y(),
				zSource: m.meshes[j].RawVertices[i].Z(), xCam: cameraSpace.X(), yCam: cameraSpace.Y(), zCam: cameraSpace.Z(), wClip: clipSpace.W()}

			normalVec4 := mgl32.Vec4{m.meshes[j].RawNormals[i].X(), m.meshes[j].RawNormals[i].Y(), m.meshes[j].RawNormals[i].Z(), 0}
			m.meshes[j].ProcessedNormals[i] = normalModelView.Mul4x1(normalVec4).Vec3().Normalize()
		}
	}
}

func (m *MeshController) AddMesh(mesh *Mesh) {
	mesh.ProcessedVertices = make([]ProcessedVertex, len(mesh.RawVertices))
	mesh.ProcessedNormals = make([]mgl32.Vec3, len(mesh.RawNormals))
	m.meshes = append(m.meshes, mesh)

	if m.meshes == nil {
		m.meshes = make([]*Mesh, 1)
	}
}

func (m *MeshController) AddMeshes(meshes []Mesh) {
	for _, ms := range meshes {
		m.AddMesh(&ms)
	}
}
