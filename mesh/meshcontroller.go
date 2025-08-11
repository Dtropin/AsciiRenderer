package mesh

import (
	"AsciiRenderer/cameracontroller"
	"AsciiRenderer/mvp"
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

func (m *MeshController) ProcessVertices(camera *cameracontroller.CameraController, windowWidth, windowHeight int, tick int) mgl32.Vec3 { //TODO Вынести матрицы view модели и тд отдельно и в параметры передавать
	xCamera, yCamera, zCamera := camera.GetPos()
	model := mvp.MakeModelMatrix(0., 0., 0., 1., 1., 1., float32(tick)*(math.Pi/180.0))
	view := mvp.MakeViewMatrix(xCamera, yCamera, zCamera)
	proj := mvp.MakePerspectiveProjection(60, float32(windowWidth)*0.6/float32(windowHeight), 0.1, 100)
	normalModelView := model.Mul4(view).Inv().Transpose()
	//todo clipping
	for j := 0; j < len(m.meshes); j++ {
		for i := 0; i < len(m.meshes[j].RawVertices); i++ {
			var modelSpace = model.Mul4x1(m.meshes[j].RawVertices[i])
			var cameraSpace = view.Mul4x1(modelSpace)
			var clipSpace = proj.Mul4x1(cameraSpace)

			ndcX := clipSpace.X() / clipSpace.W()
			ndcY := clipSpace.Y() / clipSpace.W()

			x := math.Min(float64(windowWidth), (float64(ndcX)+1)*0.5*float64(windowWidth))
			y := math.Min(float64(windowHeight), (1-(float64(ndcY)+1)*0.5)*float64(windowHeight))

			m.meshes[j].ProjectedVertices[i] = ProcessedVertex{xScreen: float32(x), yScreen: float32(y), xSource: m.meshes[j].RawVertices[i].X(), ySource: m.meshes[j].RawVertices[i].Y(),
				zSource: m.meshes[j].RawVertices[i].Z(), xCam: cameraSpace.X(), yCam: cameraSpace.Y(), zCam: cameraSpace.Z(), wClip: clipSpace.W()}
			m.meshes[j].ProcessedNormals[i] = normalModelView.Mul4x1(m.meshes[j].RawNormals[i])
		}
	}
	return mgl32.Vec3{view.At(0, 2), view.At(1, 2), view.At(2, 2)}.Normalize()
}

func (m *MeshController) AddMesh(mesh *Mesh) {
	mesh.ProjectedVertices = make([]ProcessedVertex, len(mesh.RawVertices))
	mesh.ProcessedNormals = make([]mgl32.Vec4, len(mesh.RawNormals))
	m.meshes = append(m.meshes, mesh)

	if m.meshes == nil {
		m.meshes = make([]*Mesh, 1)
	}
}
