package input

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type InputManager struct {
	mutex sync.RWMutex
	keys  map[glfw.Key]bool
	mouseButtons map[glfw.MouseButton]bool
	mouseX, mouseY float64
}

var globalInputManager = &InputManager{
	keys:	make(map[glfw.Key]bool),
	mouseButtons: make(map[glfw.MouseButton]bool),
}

func KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	globalInputManager.mutex.Lock()
	defer globalInputManager.mutex.Unlock()

	switch action {
	case glfw.Press, glfw.Repeat:
		globalInputManager.keys[key] = true
	case glfw.Release:
		globalInputManager.keys[key] = false
	}
}

func MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	globalInputManager.mutex.Lock()
	defer globalInputManager.mutex.Unlock()

	switch action {
	case glfw.Press:
		globalInputManager.mouseButtons[button] = true
	case glfw.Release:
		globalInputManager.mouseButtons[button] = false
	}
}

func CursorPosCallback(window *glfw.Window, xpos float64, ypos float64) {
	globalInputManager.mutex.Lock()
	defer globalInputManager.mutex.Unlock()

	globalInputManager.mouseX = xpos
	globalInputManager.mouseY = ypos
}

func IsKeyPressed(key glfw.Key) bool {
	globalInputManager.mutex.RLock()
	defer globalInputManager.mutex.RUnlock()

	return globalInputManager.keys[key]
}

func IsMouseButtonPressed(button glfw.MouseButton) bool {
	globalInputManager.mutex.RLock()
	defer globalInputManager.mutex.RUnlock()

	return globalInputManager.mouseButtons[button]
}

func GetMousePosition() (float64, float64) {
	globalInputManager.mutex.RLock()
	defer globalInputManager.mutex.RUnlock()

	return globalInputManager.mouseX, globalInputManager.mouseY
}