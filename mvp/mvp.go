package mvp

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

func MakeModelMatrix(x, y, z, xscale, yscale, zscale float64) *mat.Dense {
	return mat.NewDense(4, 4, []float64{xscale, 0, 0, x,
		0, yscale, 0, y,
		0, 0, zscale, z,
		0, 0, 0, 1,
	})
}

func MakeViewMatrix(x, y, z, pitch, yaw, roll float64) *mat.Dense {
	Rx := mat.NewDense(4, 4, []float64{
		1, 0, 0, 0,
		0, math.Cos(pitch), -math.Sin(pitch), 0,
		0, math.Sin(pitch), math.Cos(pitch), 0,
		0, 0, 0, 1,
	})

	Ry := mat.NewDense(4, 4, []float64{
		math.Cos(yaw), 0, math.Sin(yaw), 0,
		0, 1, 0, 0,
		-math.Sin(yaw), 0, math.Cos(yaw), 0,
		0, 0, 0, 1,
	})

	Rz := mat.NewDense(4, 4, []float64{
		math.Cos(roll), -math.Sin(roll), 0, 0,
		math.Sin(roll), math.Cos(roll), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})

	var R, temp mat.Dense
	temp.Mul(Rx, Ry)
	R.Mul(&temp, Rz)

	T := mat.NewDense(4, 4, []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		-x, -y, -z, 1,
	})

	var view mat.Dense
	view.Mul(T, &R)

	return &view
}

func MakePerspectiveProjection(fov, near, aspect, far float64) *mat.Dense {
	t := 1 / (aspect * math.Tan(fov*0.5*math.Pi/180))

	return mat.NewDense(4, 4, []float64{
		t, 0, 0, 0,
		0, t * aspect, 0, 0,
		0, 0, -(far + near) / (far - near), -(2 * far * near) / (far - near),
		0, 0, -1, 0,
	})
}
