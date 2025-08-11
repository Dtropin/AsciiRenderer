package cameracontroller

import "sync"

type CameraController struct {
	cameraState *CameraState
}

type CameraState struct {
	x, y, z, pitch, yaw, roll float32
}

var (
	instance *CameraController
	once     sync.Once
)

func Init() *CameraController {
	once.Do(func() {
		instance = &CameraController{&CameraState{}}
	})
	return instance
}

func (m *CameraController) GetPos() (x, y, z float32) {
	x, y, z = m.cameraState.x, m.cameraState.y, m.cameraState.z
	return
}

func (m *CameraController) SetPos(x, y, z float32) {
	m.cameraState.x = x
	m.cameraState.y = y
	m.cameraState.z = z
}

func (m *CameraController) AdjustPos(dx, dy, dz float32) {
	m.cameraState.x += dx
	m.cameraState.y += dy
	m.cameraState.z += dz
}
