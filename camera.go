package main

import (
	"graphics/input"
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	position          mgl32.Vec3
	direction         mgl32.Vec3
	viewInverse       mgl32.Mat4
	projectionInverse mgl32.Mat4
	uniform           CameraUniform
}

type CameraUniform struct {
	position          [3]float32
	_                 float32
	projectionInverse [4][4]float32
	viewInverse       [4][4]float32
}

func InitCamera() *Camera {
	return &Camera{
		position:          mgl32.Vec3{0.0, 0.0, 3.0},
		direction:         mgl32.Vec3{0.0, 0.0, -1.0},
		projectionInverse: mgl32.Perspective(math.Pi/4.0, 1.0, 1.0, 100.0).Inv(),
	}
}

func (camera *Camera) OnUpdate(window *glfw.Window) {
	if !input.IsMouseButtonPressed(glfw.MouseButton1) {
		input.SetCursorMode(window, glfw.CursorNormal)
		return
	}
	input.SetCursorMode(window, glfw.CursorNormal)

	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	rightDir := upDir.Cross(camera.direction)

	moved := false
	var speed float32 = 0.2

	// Movement
	if input.IsKeyPressed(glfw.KeyW) {
		camera.position = camera.position.Add(camera.direction.Mul(speed))
		moved = true
	}
	if input.IsKeyPressed(glfw.KeyS) {
		camera.position = camera.position.Sub(camera.direction.Mul(speed))
		moved = true
	}
	if input.IsKeyPressed(glfw.KeyA) {
		camera.position = camera.position.Add(rightDir.Mul(speed))
		moved = true
	}
	if input.IsKeyPressed(glfw.KeyD) {
		camera.position = camera.position.Sub(rightDir.Mul(speed))
		moved = true
	}
	if input.IsKeyPressed(glfw.KeySpace) {
		camera.position = camera.position.Add(upDir.Mul(speed))
		moved = true
	}
	if input.IsKeyPressed(glfw.KeyLeftShift) {
		camera.position = camera.position.Sub(upDir.Mul(speed))
		moved = true
	}

	// Rotation
	delta := input.GetMouseDelta()

	var rotSpeed float32 = 0.001

	if delta.Len() != 0.0 {
		pitchDelta := delta.Y() * rotSpeed
		yawDelta := delta.X() * rotSpeed

		q := mgl32.QuatRotate(-pitchDelta, rightDir).Mul(mgl32.QuatRotate(yawDelta, upDir)).Normalize()
		camera.direction = q.Rotate(camera.direction)

		moved = true
	}

	if moved {
		camera.viewInverse = mgl32.LookAtV(camera.position, camera.position.Add(camera.direction), upDir).Inv()
		camera.uniform = CameraUniform{
			position:          camera.position,
			viewInverse:       mat4ToArray(camera.viewInverse),
			projectionInverse: mat4ToArray(camera.projectionInverse),
		}
	}
}

func mat4ToArray(mat mgl32.Mat4) [4][4]float32 {
	var arr [4][4]float32
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			arr[row][col] = mat.At(row, col)
		}
	}
	return arr
}
