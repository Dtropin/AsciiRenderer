package mesh

import "github.com/go-gl/mathgl/mgl32"

type Mesh struct {
	RawVertices       []mgl32.Vec4
	RawNormals        []mgl32.Vec4
	ProcessedNormals  []mgl32.Vec3
	ProcessedVertices []ProcessedVertex
	Polys             [][3]int
	PolysNormals      [][3]int
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

func (p *ProcessedVertex) ZSource() float32 {
	return p.zSource
}

func (p *ProcessedVertex) XScreen() float32 {
	return p.xScreen
}

func (p *ProcessedVertex) YScreen() float32 {
	return p.yScreen
}

func (p *ProcessedVertex) WClip() float32 {
	return p.wClip
}
