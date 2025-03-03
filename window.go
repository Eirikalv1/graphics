package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"graphics/input"
)

type Window struct {
	handle *glfw.Window
}

func InitWindow() (window *Window, err error) {
	defer func() {
		if err != nil {
			window.Destroy()
			window = nil
		}
	}()

	window = &Window{}

	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	window.handle, err = glfw.CreateWindow(600, 600, "title", nil, nil)
	if err != nil {
		return window, err
	}

	window.handle.SetKeyCallback(input.KeyCallback)
	window.handle.SetMouseButtonCallback(input.MouseButtonCallback)

	return window, err
}

func (window *Window) Destroy() {
	if window.handle != nil {
		window.handle.Destroy()
	}

	glfw.Terminate()
}
