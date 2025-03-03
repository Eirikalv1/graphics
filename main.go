package main

import (
	"fmt"
	"graphics/input"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	runtime.LockOSThread()

	window, err := InitWindow()
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	ctx, err := InitGpuContext(window.handle)
	if err != nil {
		panic(err)
	}
	defer ctx.Destroy()

	window.handle.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		ctx.Resize(width, height)
	})

	camera := InitCamera()

	renderer, err := InitRenderer(ctx, camera)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	for !window.handle.ShouldClose() {
		glfw.PollEvents()

		input.UpdateMousePosition(window.handle)
		camera.OnUpdate(window.handle)

		if err := renderer.Render(ctx, camera); err != nil {
			fmt.Println("error occured while rendering: ", err)
			panic(err)
		}
	}
}
