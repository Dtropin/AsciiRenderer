package mesh_controller

import "github.com/go-gl/mathgl/mgl32"

type Mesh struct {
	RawVertices       []mgl32.Vec4
	ProjectedVertices []ProjectedVertex
	Colors            []rune
	Polygons          []Polygon //todo Будто бы лишняя структура
	OutlineEdges      map[[2]int]int
}

type ProjectedVertex struct {
	xScreen, yScreen, xSource, ySource, zSource, xCam, yCam, zCam, wClip float32
}

type Polygon struct {
	VerticesIndices [3]int
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

func (p *ProjectedVertex) XSource() float32 {
	return p.xSource
}

func (p *ProjectedVertex) YSource() float32 {
	return p.ySource
}

func (p *ProjectedVertex) ZSource() float32 {
	return p.zSource
}

func (p *ProjectedVertex) XScreen() float32 {
	return p.xScreen
}

func (p *ProjectedVertex) YScreen() float32 {
	return p.yScreen
}

func (p *ProjectedVertex) WClip() float32 {
	return p.wClip
}
