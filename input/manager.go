package input

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type InputManager struct {
	mutex             sync.RWMutex
	keys              map[glfw.Key]bool
	mouseButtons      map[glfw.MouseButton]bool
	mousePosition     mgl32.Vec2
	prevMousePosition mgl32.Vec2
}

var globalInputManager = &InputManager{
	keys:         make(map[glfw.Key]bool),
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

	currentPos := mgl32.Vec2{float32(xpos), float32(ypos)}
	globalInputManager.prevMousePosition = globalInputManager.mousePosition
	globalInputManager.mousePosition = currentPos
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

func GetMouseDelta() mgl32.Vec2 {
	globalInputManager.mutex.RLock()
	defer globalInputManager.mutex.RUnlock()

	return globalInputManager.mousePosition.Sub(globalInputManager.prevMousePosition)
}

func GetMousePosition() mgl32.Vec2 {
	globalInputManager.mutex.RLock()
	defer globalInputManager.mutex.RUnlock()

	return globalInputManager.mousePosition
}
