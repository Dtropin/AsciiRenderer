package mvp

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func MakeModelMatrix(x, y, z, xscale, yscale, zscale, angleRad float32) mgl32.Mat4 {
	return mgl32.Translate3D(x, y, z).Mul4(mgl32.HomogRotate3D(angleRad, [3]float32{1., 0., 0.})).Mul4(mgl32.HomogRotate3D(angleRad, [3]float32{0., 1., 0.})).Mul4(mgl32.Scale3D(xscale, yscale, zscale))
}

func MakePerspectiveProjection(fov, aspect, near, far float32) mgl32.Mat4 {
	return mgl32.Perspective((fov*math.Pi)/180.0, aspect, near, far)
}
