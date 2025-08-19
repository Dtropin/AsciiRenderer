package cameracontroller

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"sync"
)

var (
	instance *CameraController
	once     sync.Once
)

type CameraController struct {
	cameraState *CameraState
}

type CameraState struct {
	Pos        mgl32.Vec3
	LookAt     mgl32.Vec3
	Up         mgl32.Vec3
	Forward    mgl32.Vec3
	Right      mgl32.Vec3
	ViewMatrix mgl32.Mat4
	Yaw        float32
	Pitch      float32
	mutex      sync.RWMutex
}

func (m *CameraController) CameraState() *CameraState {
	m.cameraState.mutex.RLock()
	defer m.cameraState.mutex.RUnlock()
	return m.cameraState
}

func Init() *CameraController {
	once.Do(func() {
		cameraPosition := mgl32.Vec3{0, 0, 10}
		lookAtPoint := mgl32.Vec3{0, 0, 0}
		upVector := mgl32.Vec3{0, 1, 0}

		instance = &CameraController{&CameraState{
			Pos:    cameraPosition,
			LookAt: lookAtPoint,
			Up:     upVector,
			Pitch:  0.,
			Yaw:    -90.,
			mutex:  sync.RWMutex{},
		}}
		instance.updateViewMatrix()
	})
	return instance
}

func (m *CameraController) Yaw(delta float32) {
	m.cameraState.mutex.Lock()
	defer m.cameraState.mutex.Unlock()
	m.cameraState.Yaw += delta
	m.updateViewMatrix()
}

func (m *CameraController) Pitch(delta float32) {
	m.cameraState.mutex.Lock()
	defer m.cameraState.mutex.Unlock()
	m.cameraState.Pitch -= delta
	if m.cameraState.Pitch > 89.0 {
		m.cameraState.Pitch = 89.0
	}
	if m.cameraState.Pitch < -89.0 {
		m.cameraState.Pitch = -89.0
	}
	m.updateViewMatrix()
}

func (m *CameraController) Forward(speed float32) {
	m.cameraState.mutex.Lock()
	defer m.cameraState.mutex.Unlock()
	m.cameraState.Pos = m.cameraState.Pos.Add(m.cameraState.Forward.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) Up(speed float32) {
	m.cameraState.mutex.Lock()
	defer m.cameraState.mutex.Unlock()
	m.cameraState.Pos = m.cameraState.Pos.Add(m.cameraState.Up.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) Down(speed float32) {
	m.cameraState.mutex.Lock()
	defer m.cameraState.mutex.Unlock()
	m.cameraState.Pos = m.cameraState.Pos.Sub(m.cameraState.Up.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) Back(speed float32) {
	m.cameraState.mutex.Lock()
	defer m.cameraState.mutex.Unlock()
	m.cameraState.Pos = m.cameraState.Pos.Sub(m.cameraState.Forward.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) Right(speed float32) {
	m.cameraState.mutex.Lock()
	defer m.cameraState.mutex.Unlock()
	m.cameraState.Pos = m.cameraState.Pos.Add(m.cameraState.Right.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) Left(speed float32) {
	m.cameraState.mutex.Lock()
	defer m.cameraState.mutex.Unlock()
	m.cameraState.Pos = m.cameraState.Pos.Sub(m.cameraState.Right.Mul(speed))
	m.updateViewMatrix()
}

func (m *CameraController) updateViewMatrix() {
	forwardVector := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(m.cameraState.Yaw))) * math.Cos(float64(mgl32.DegToRad(m.cameraState.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(m.cameraState.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(m.cameraState.Yaw))) * math.Cos(float64(mgl32.DegToRad(m.cameraState.Pitch)))),
	}.Normalize()

	m.cameraState.LookAt = m.cameraState.Pos.Add(forwardVector)
	m.cameraState.Forward = forwardVector
	m.cameraState.Right = m.cameraState.Forward.Cross(m.cameraState.Up).Normalize()
	m.cameraState.ViewMatrix = mgl32.LookAtV(m.cameraState.Pos, m.cameraState.LookAt, m.cameraState.Up)
}
