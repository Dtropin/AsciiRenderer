package mesh

import "github.com/go-gl/mathgl/mgl32"

type Mesh struct {
	RawVertices       []mgl32.Vec4
	RawNormals        []mgl32.Vec4
	ProcessedNormals  []mgl32.Vec3
	ProcessedVertices []ProcessedVertex
	Polys             [][3]int
	PolysNormals      [][3]int
	ModelMatrix       mgl32.Mat4
	Pos               mgl32.Vec3
	Scale             mgl32.Vec3
	AngleRad          float32
}

func (mesh *Mesh) UpdateModelMatrix() {
	mesh.ModelMatrix = mgl32.Translate3D(mesh.Pos.X(), mesh.Pos.Y(), mesh.Pos.Z()).Mul4(mgl32.HomogRotate3D(mesh.AngleRad, [3]float32{1., 0., 0.})).Mul4(mgl32.HomogRotate3D(mesh.AngleRad, [3]float32{0., 1., 0.})).Mul4(mgl32.Scale3D(mesh.Scale.X(), mesh.Scale.Y(), mesh.Scale.Z()))
}

type ProcessedVertex struct {
	xScreen, yScreen, xSource, ySource, zSource, xCam, yCam, zCam, wClip float32
}

func (p *ProcessedVertex) XCam() float32 {
	return p.xCam
}

func (p *ProcessedVertex) YCam() float32 {
	return p.yCam
}

func (p *ProcessedVertex) ZCam() float32 {
	return p.zCam
}

func (p *ProcessedVertex) XSource() float32 {
	return p.xSource
}

func (p *ProcessedVertex) YSource() float32 {
	return p.ySource
}

func (p *ProcessedVertex) ZSource() float32 { return p.zSource }

func (p *ProcessedVertex) XScreen() float32 {
	return p.xScreen
}

func (p *ProcessedVertex) YScreen() float32 {
	return p.yScreen
}

func (p *ProcessedVertex) WClip() float32 {
	return p.wClip
}
