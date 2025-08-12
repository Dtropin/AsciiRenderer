package cameracontroller

import (
	"github.com/go-gl/mathgl/mgl32"
	"sync"
)

type CameraController struct {
	cameraState *CameraState
}

func (m *CameraController) CameraState() *CameraState {
	return m.cameraState
}

type CameraState struct {
	pos        mgl32.Vec3
	lookAt     mgl32.Vec3
	up         mgl32.Vec3
	forward    mgl32.Vec3
	right      mgl32.Vec3
	viewMatrix mgl32.Mat4
}

func (m *CameraController) updateViewMatrix() {
	//TODO sync
	m.cameraState.forward = m.cameraState.lookAt.Sub(m.cameraState.pos).Normalize()
	m.cameraState.right = m.cameraState.forward.Cross(m.cameraState.up).Normalize()
	m.cameraState.viewMatrix = mgl32.LookAtV(m.cameraState.pos, m.cameraState.lookAt, m.cameraState.up)
}

func Init() *CameraController {
	once.Do(func() {
		cameraPosition := mgl32.Vec3{0, 0, 10}
		lookAtPoint := mgl32.Vec3{0, 0, 0}
		upVector := mgl32.Vec3{0, 1, 0}

		instance = &CameraController{&CameraState{
			pos:    cameraPosition,
			lookAt: lookAtPoint,
			up:     upVector,
		}}
		instance.updateViewMatrix()
	})
	return instance
}

func (m *CameraController) Forward(speed float32) {
	m.cameraState.pos = m.cameraState.pos.Add(m.cameraState.forward.Mul(speed))
	m.cameraState.lookAt = m.cameraState.lookAt.Add(m.cameraState.forward.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) Up(speed float32) {
	m.cameraState.pos = m.cameraState.pos.Add(m.cameraState.up.Mul(speed))
	m.cameraState.lookAt = m.cameraState.lookAt.Add(m.cameraState.up.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) Down(speed float32) {
	m.cameraState.pos = m.cameraState.pos.Sub(m.cameraState.up.Mul(speed))
	m.cameraState.lookAt = m.cameraState.lookAt.Sub(m.cameraState.up.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) Back(speed float32) {
	m.cameraState.pos = m.cameraState.pos.Sub(m.cameraState.forward.Mul(speed))
	m.cameraState.lookAt = m.cameraState.lookAt.Sub(m.cameraState.forward.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) Right(speed float32) {
	m.cameraState.pos = m.cameraState.pos.Add(m.cameraState.right.Mul(speed))
	m.cameraState.lookAt = m.cameraState.lookAt.Add(m.cameraState.right.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) Left(speed float32) {
	m.cameraState.pos = m.cameraState.pos.Sub(m.cameraState.right.Mul(speed))
	m.cameraState.lookAt = m.cameraState.lookAt.Sub(m.cameraState.right.Mul(speed))
	m.updateViewMatrix()
}

func (c CameraState) ViewMatrix() mgl32.Mat4 {
	return c.viewMatrix
}

func (c CameraState) Up() mgl32.Vec3 {
	return c.up
}

func (c CameraState) Right() mgl32.Vec3 {
	return c.right
}

func (c CameraState) Forward() mgl32.Vec3 {
	return c.forward
}

func (c CameraState) LookAt() mgl32.Vec3 {
	return c.lookAt
}

func (c CameraState) Pos() mgl32.Vec3 {
	return c.pos
}

var (
	instance *CameraController
	once     sync.Once
)
