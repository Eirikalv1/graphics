package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"

	input "graphics/inputevents"
)

func main() {
	runtime.LockOSThread()

	window, err := InitWindow()
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	input.InitKeyState()

	ctx, err := InitGpuContext(window.handle)
	if err != nil {
		panic(err)
	}
	defer ctx.Destroy()

	window.handle.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		ctx.Resize(width, height)
	})

	renderer, err := InitRenderer(ctx)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	for !window.handle.ShouldClose() {
		glfw.PollEvents()

		if err := renderer.Render(ctx); err != nil {
			fmt.Println("error occured while rendering: ", err)
			panic(err)
		}
	}
}
