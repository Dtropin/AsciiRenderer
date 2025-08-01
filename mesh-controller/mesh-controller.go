package mesh_controller

import (
	"AsciiRenderer/camera-controller"
	"AsciiRenderer/mvp"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"sync"
)

// TODO переименовать в вертекс шейдер?
type MeshController struct {
	vertices          []mgl32.Vec4 //Таким образом будет задействовано кеширование цпу
	projectedVertices []ProjectedVertex
}

type ProjectedVertex struct {
	xScreen, yScreen, xClip, yClip, zClip, xCam, yCam, zCam, wClip float32
}

func (p *ProjectedVertex) XCam() float32 {
	return p.xCam
}

func (p *ProjectedVertex) YCam() float32 {
	return p.yCam
}

func (p *ProjectedVertex) ZCam() float32 {
	return p.zCam
}

func (p *ProjectedVertex) XClip() float32 {
	return p.xClip
}

func (p *ProjectedVertex) YClip() float32 {
	return p.yClip
}

func (p *ProjectedVertex) XScreen() float32 {
	return p.xScreen
}

func (p *ProjectedVertex) YScreen() float32 {
	return p.yScreen
}

func (p *ProjectedVertex) ZClip() float32 {
	return p.zClip
}

func (p *ProjectedVertex) WClip() float32 {
	return p.wClip
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

func (m *MeshController) ProcessVertices(camera *camera_controller.CameraController, windowWidth, windowHeight int, tick int) {
	xcamera, ycamera, zcamera := camera.GetPos()
	model := mvp.MakeModelMatrix(0., 0., 0., 1., 1., 1., float32(tick)*(math.Pi/200.0))
	view := mvp.MakeViewMatrix(xcamera, ycamera, zcamera, 0)
	proj := mvp.MakePerspectiveProjection(60, float32(windowWidth)*0.6/float32(windowHeight), 0.1, 100)

	//todo clipping
	for i := 0; i < len(m.vertices); i++ {
		var modelSpace = model.Mul4x1(m.vertices[i])
		var cameraSpace = view.Mul4x1(modelSpace)
		var clipSpace = proj.Mul4x1(cameraSpace)

		ndcX := clipSpace.X() / clipSpace.W()
		ndcY := clipSpace.Y() / clipSpace.W()

		x := math.Min(float64(windowWidth), (float64(ndcX)+1)*0.5*float64(windowWidth))
		y := math.Min(float64(windowHeight), (1-(float64(ndcY)+1)*0.5)*float64(windowHeight))

		m.projectedVertices[i] = ProjectedVertex{xScreen: float32(x), yScreen: float32(y), xClip: clipSpace.X(), yClip: clipSpace.Y(), zClip: clipSpace.Z(),
			xCam: cameraSpace.X(), yCam: cameraSpace.Y(), zCam: cameraSpace.Z(), wClip: clipSpace.W()}
	}
}

func (m *MeshController) AddVerticesToMesh(vertices []mgl32.Vec4) {
	if m.vertices == nil {
		m.vertices = make([]mgl32.Vec4, len(vertices))
		copy(m.vertices, vertices)
	} else {
		m.vertices = append(m.vertices, vertices...)
	}

	if m.projectedVertices == nil {
		m.projectedVertices = make([]ProjectedVertex, len(vertices)*3)
	} else {
		m.projectedVertices = append(m.projectedVertices, make([]ProjectedVertex, len(vertices)*3)...)
	}
}

func (m *MeshController) GetProjectedVertices() []ProjectedVertex {
	return m.projectedVertices
}
