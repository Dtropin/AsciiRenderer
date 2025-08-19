package mesh

import (
	"AsciiRenderer/cameracontroller"
	"AsciiRenderer/viewport"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"runtime"
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
	view := camera.CameraState().ViewMatrix
	w, h := viewPortController.GetWindowSize()

	proj := mgl32.Perspective((60*math.Pi)/180.0, float32(w)*0.6/float32(h), 0.1, 100)

	//todo clipping
	numCPU := runtime.NumCPU()
	waitGroup := sync.WaitGroup{}
	channel := make(chan chunkProcessingResult, numCPU*len(m.meshes))
	for j := 0; j < len(m.meshes); j++ {
		m.meshes[j].UpdateModelMatrix()
		model := m.meshes[j].ModelMatrix
		modelViewMatrix := view.Mul4(model)
		normalModelView := modelViewMatrix.Inv().Transpose()
		chunkSize := len(m.meshes[j].RawVertices) / numCPU
		for i := 0; i < numCPU; i++ {
			waitGroup.Add(1)
			start := i * chunkSize
			end := start + chunkSize
			if i == numCPU-1 {
				end = len(m.meshes[j].RawVertices)
			}
			go processVerticesChunk(model, view, proj, normalModelView, w, h, j, m.meshes[j].RawVertices, m.meshes[j].RawNormals, start, end, &waitGroup, channel)
		}
	}
	waitGroup.Wait()
	close(channel)
	for result := range channel {
		resultMesh := m.meshes[result.meshIdx]
		copy(resultMesh.ProcessedVertices[result.start:], result.localProcessedVertices)
		copy(resultMesh.ProcessedNormals[result.start:], result.localProcessedNormals)
	}
}

type chunkProcessingResult struct {
	start                  int
	meshIdx                int
	localProcessedVertices []ProcessedVertex
	localProcessedNormals  []mgl32.Vec3
}

func processVerticesChunk(model,
	view,
	proj,
	normalModelView mgl32.Mat4,
	w, h int,
	meshIdx int,
	rawVertices []mgl32.Vec4,
	rawNormals []mgl32.Vec4,
	start, end int,
	syncGroup *sync.WaitGroup,
	resultChannel chan chunkProcessingResult) {
	defer syncGroup.Done()
	localProcessedVertices := make([]ProcessedVertex, end-start)
	localProcessedNormals := make([]mgl32.Vec3, end-start)
	for i := start; i < end; i++ {
		var modelSpace = model.Mul4x1(rawVertices[i])
		var cameraSpace = view.Mul4x1(modelSpace)
		var clipSpace = proj.Mul4x1(cameraSpace)

		ndcX := clipSpace.X() / clipSpace.W()
		ndcY := clipSpace.Y() / clipSpace.W()

		x := math.Min(float64(w), (float64(ndcX)+1)*0.5*float64(w))
		y := math.Min(float64(h), (1-(float64(ndcY)+1)*0.5)*float64(h))

		localProcessedVertices[i-start] = ProcessedVertex{xScreen: float32(x), yScreen: float32(y), xSource: rawVertices[i].X(), ySource: rawVertices[i].Y(),
			zSource: rawVertices[i].Z(), xCam: cameraSpace.X(), yCam: cameraSpace.Y(), zCam: cameraSpace.Z(), wClip: clipSpace.W()}

		normalVec4 := mgl32.Vec4{rawNormals[i].X(), rawNormals[i].Y(), rawNormals[i].Z(), 0}
		localProcessedNormals[i-start] = normalModelView.Mul4x1(normalVec4).Vec3().Normalize()
	}

	resultChannel <- chunkProcessingResult{start: start, meshIdx: meshIdx, localProcessedNormals: localProcessedNormals, localProcessedVertices: localProcessedVertices}
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
