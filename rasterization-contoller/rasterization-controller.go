package rasterization_contoller

import (
	mesh_controller "AsciiRenderer/mesh-controller"
	render_context "AsciiRenderer/terminal-context"
	"math"
	"sort"
)

// TODO сделать обьектом
func ScanlineRasterization(polys [][]int, projectedVertices []mesh_controller.ProjectedVertex, zbuff [][]float32, colors []rune, viewPortController *render_context.ViewPortController) {
	for i := 0; i < len(polys); i++ {
		p0 := projectedVertices[polys[i][0]]
		p1 := projectedVertices[polys[i][1]]
		p2 := projectedVertices[polys[i][2]]

		minY, maxY := int(math.Min(float64(p0.YScreen()), math.Min(float64(p1.YScreen()), float64(p2.YScreen())))),
			int(math.Max(float64(p0.YScreen()), math.Max(float64(p1.YScreen()), float64(p2.YScreen()))))

		if minY < 0 {
			continue
		}

		for y := minY; y < maxY; y++ {
			var intersections = make([]Intersection, 0)
			intersections = addEdgeIntersection(intersections, &p0, &p1, float32(y))
			intersections = addEdgeIntersection(intersections, &p0, &p2, float32(y))
			intersections = addEdgeIntersection(intersections, &p1, &p2, float32(y))

			sort.Slice(intersections, func(a, b int) bool {
				return intersections[a].interpolatedX < intersections[b].interpolatedX
			})

			if len(intersections) == 2 {
				xStart := int(intersections[0].interpolatedX)
				xEnd := int(intersections[1].interpolatedX)

				for x := xStart; x <= xEnd; x++ {
					u, v, w := barycentric(x, y, p0.XScreen(), p0.YScreen(), p1.XScreen(), p1.YScreen(), p2.XScreen(), p2.YScreen())
					wp := u*(1/p0.WClip()) + v*(1/p1.WClip()) + w*(1/p2.WClip())
					if wp > zbuff[x][y] {
						zbuff[x][y] = wp
						px3d := (u*p0.XCam()/p0.WClip() + v*p1.XCam()/p1.WClip() + w*p2.XCam()/p2.WClip()) / wp
						py3d := (u*p0.YCam()/p0.WClip() + v*p1.YCam()/p1.WClip() + w*p2.YCam()/p2.WClip()) / wp
						pz3d := (u*p0.ZCam()/p0.WClip() + v*p1.ZCam()/p1.WClip() + w*p2.ZCam()/p2.WClip()) / wp
						d0 := (p0.XCam()-px3d)*(p0.XCam()-px3d) + (p0.YCam()-py3d)*(p0.YCam()-py3d) + (p0.ZCam()-pz3d)*(p0.ZCam()-pz3d)
						d1 := (p1.XCam()-px3d)*(p1.XCam()-px3d) + (p1.YCam()-py3d)*(p1.YCam()-py3d) + (p1.ZCam()-pz3d)*(p1.ZCam()-pz3d)
						d2 := (p2.XCam()-px3d)*(p2.XCam()-px3d) + (p2.YCam()-py3d)*(p2.YCam()-py3d) + (p2.ZCam()-pz3d)*(p2.ZCam()-pz3d)
						var char rune
						if d0 < d1 && d0 < d2 {
							char = colors[polys[i][0]]
						} else if d1 < d0 && d1 < d2 {
							char = colors[polys[i][1]]
						} else {
							char = colors[polys[i][2]]
						}
						viewPortController.SetChar(x, y, char)
					}
				}
			}
		}
	}
}

func barycentric(x, y int, p0x, p0y, p1x, p1y, p2x, p2y float32) (float32, float32, float32) {
	v0x := p1x - p0x
	v0y := p1y - p0y
	v1x := p2x - p0x
	v1y := p2y - p0y
	v2x := float32(x) - p0x
	v2y := float32(y) - p0y

	dot00 := v0x*v0x + v0y*v0y
	dot01 := v0x*v1x + v0y*v1y
	dot02 := v0x*v2x + v0y*v2y
	dot11 := v1x*v1x + v1y*v1y
	dot12 := v1x*v2x + v1y*v2y

	denom := 1 / (dot00*dot11 - dot01*dot01)
	v := (dot11*dot02 - dot01*dot12) * denom
	w := (dot00*dot12 - dot01*dot02) * denom
	u := 1 - v - w
	return u, v, w
}

type Intersection struct {
	interpolatedX, interpolatedZClip, interpolatedWClip float32
}

func addEdgeIntersection(intersections []Intersection, p0 *mesh_controller.ProjectedVertex, p1 *mesh_controller.ProjectedVertex, y float32) []Intersection {
	if (p0.YScreen() <= y && p1.YScreen() >= y) || (p1.YScreen() <= y && p0.YScreen() >= y) {
		alpha := (y - p0.YScreen()) / (p1.YScreen() - p0.YScreen())
		interpolatedX := p0.XScreen() + (p1.XScreen()-p0.XScreen())*alpha
		interpolatedZClip := p0.ZClip() + (p1.ZClip()-p0.ZClip())*alpha
		interpolatedWClip := p0.WClip() + (p1.WClip()-p0.WClip())*alpha
		return append(intersections, Intersection{interpolatedX: interpolatedX, interpolatedZClip: interpolatedZClip, interpolatedWClip: interpolatedWClip})
	}

	return intersections
}
