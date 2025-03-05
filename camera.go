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
	camera := &Camera{
		position:          mgl32.Vec3{0.0, 0.0, 3.0},
		direction:         mgl32.Vec3{0.0, 0.0, -1.0},
		projectionInverse: mgl32.Perspective(math.Pi/4.0, 1.0, 1.0, 100.0).Inv(),
	}

	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	camera.viewInverse = mgl32.LookAtV(
		camera.position,
		camera.position.Add(camera.direction),
		upDir,
	).Inv()

	camera.uniform = CameraUniform{
		position:          camera.position,
		viewInverse:       mat4ToArray(camera.viewInverse),
		projectionInverse: mat4ToArray(camera.projectionInverse),
	}

	return camera
}

func (camera *Camera) OnUpdate(window *glfw.Window) {
	if !input.IsMouseButtonPressed(glfw.MouseButton1) {
		input.SetCursorMode(window, glfw.CursorNormal)
		return
	}
	input.SetCursorMode(window, glfw.CursorDisabled)

	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	rightDir := camera.direction.Cross(upDir)

	var speed float32 = 0.1

	// Movement
	if input.IsKeyPressed(glfw.KeyW) {
		camera.position = camera.position.Add(camera.direction.Mul(speed))
	}
	if input.IsKeyPressed(glfw.KeyS) {
		camera.position = camera.position.Sub(camera.direction.Mul(speed))
	}
	if input.IsKeyPressed(glfw.KeyD) {
		camera.position = camera.position.Add(rightDir.Mul(speed))
	}
	if input.IsKeyPressed(glfw.KeyA) {
		camera.position = camera.position.Sub(rightDir.Mul(speed))
	}
	if input.IsKeyPressed(glfw.KeySpace) {
		camera.position = camera.position.Add(upDir.Mul(speed))
	}
	if input.IsKeyPressed(glfw.KeyLeftShift) {
		camera.position = camera.position.Sub(upDir.Mul(speed))
	}

	// Rotation
	delta := input.GetMouseDelta()

	var rotSpeed float32 = 0.002

	if delta.Len() != 0.0 {
		pitchDelta := delta.Y() * rotSpeed
		yawDelta := delta.X() * rotSpeed

		q := mgl32.QuatRotate(-pitchDelta, rightDir).Mul(mgl32.QuatRotate(-yawDelta, upDir)).Normalize()
		camera.direction = q.Rotate(camera.direction)

	}

	camera.viewInverse = mgl32.LookAtV(camera.position, camera.position.Add(camera.direction), upDir).Inv()
	camera.uniform = CameraUniform{
		position:          camera.position,
		viewInverse:       mat4ToArray(camera.viewInverse),
		projectionInverse: mat4ToArray(camera.projectionInverse),
	}

}

func mat4ToArray(mat mgl32.Mat4) [4][4]float32 {
	var arr [4][4]float32
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			arr[col][row] = mat.At(row, col)
		}
	}
	return arr
}
