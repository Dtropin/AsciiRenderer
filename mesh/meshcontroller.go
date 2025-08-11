package mesh

import (
	"AsciiRenderer/cameracontroller"
	"AsciiRenderer/mvp"
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

func (m *MeshController) ProcessVertices(camera *cameracontroller.CameraController, windowWidth, windowHeight int, tick int) {
	xCamera, yCamera, zCamera := camera.GetPos()
	model := mvp.MakeModelMatrix(0., 0., 0., 1., 1., 1., float32(tick)*(math.Pi/180.0))
	view := mvp.MakeViewMatrix(xCamera, yCamera, zCamera)
	proj := mvp.MakePerspectiveProjection(60, float32(windowWidth)*0.6/float32(windowHeight), 0.1, 100)

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

			m.meshes[j].ProjectedVertices[i] = ProjectedVertex{xScreen: float32(x), yScreen: float32(y), xSource: m.meshes[j].RawVertices[i].X(), ySource: m.meshes[j].RawVertices[i].Y(),
				zSource: m.meshes[j].RawVertices[i].Z(), xCam: cameraSpace.X(), yCam: cameraSpace.Y(), zCam: cameraSpace.Z(), wClip: clipSpace.W()}
		}
	}

}

func (m *MeshController) AddMesh(mesh *Mesh) {
	mesh.ProjectedVertices = make([]ProjectedVertex, len(mesh.RawVertices)*3)
	mesh.ProjectedVertices = append(mesh.ProjectedVertices, make([]ProjectedVertex, len(mesh.RawVertices)*3)...)

	mesh.OutlineEdges = make(map[[2]int]int)

	for i := 0; i < len(mesh.Polygons); i++ {
		key0 := [2]int{int(math.Min(float64(mesh.Polygons[i][0]), float64(mesh.Polygons[i][1]))),
			int(math.Max(float64(mesh.Polygons[i][0]), float64(mesh.Polygons[i][1])))}
		key1 := [2]int{int(math.Min(float64(mesh.Polygons[i][1]), float64(mesh.Polygons[i][2]))),
			int(math.Max(float64(mesh.Polygons[i][1]), float64(mesh.Polygons[i][2])))}
		key2 := [2]int{int(math.Min(float64(mesh.Polygons[i][0]), float64(mesh.Polygons[i][2]))),
			int(math.Max(float64(mesh.Polygons[i][0]), float64(mesh.Polygons[i][2])))}
		mesh.OutlineEdges[key0]++
		mesh.OutlineEdges[key1]++
		mesh.OutlineEdges[key2]++
	}

	m.meshes = append(m.meshes, mesh)

	if m.meshes == nil {
		m.meshes = make([]*Mesh, 1)
	}
}
