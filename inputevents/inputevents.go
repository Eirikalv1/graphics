package inputevents

import "github.com/go-gl/glfw/v3.3/glfw"

var keyState KeyState

type KeyState struct {
	Key    glfw.Key
	Action glfw.Action
	Mods   glfw.ModifierKey
}

func InitKeyState() {
	keyState = KeyState{}
}

func SetKey(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) {
	keyState.Key = key
	keyState.Action = action
	keyState.Mods = mods
}

func GetKey(key glfw.Key) bool {
	return keyState.Key == key
}

func GetAction(action glfw.Action) bool {
	return keyState.Action == action
}

func GetModifierKey(mods glfw.ModifierKey) bool {
	return keyState.Mods == mods
}
