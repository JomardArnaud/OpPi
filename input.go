package oppi

import "github.com/veandco/go-sdl2/sdl"

/*
enum for Keyboard => https://wiki.libsdl.org/SDL_Keycode
	 for button => https://wiki.libsdl.org/SDL_GameControllerButton
	 for axis => https://wiki.libsdl.org/SDL_GameControllerAxis
	 for mouse => https://wiki.libsdl.org/SDL_MouseButtonEvent#Remarks
*/

//OpGamepad is just a layer from sdl.Gamepad with usefull methode
type OpGamepad struct {
	deadZone float64
	Axis     [sdl.CONTROLLER_AXIS_MAX]int
	Button   [sdl.CONTROLLER_BUTTON_MAX]bool
}

func (pad *OpGamepad) empty() {
	for i := 0; i < len(pad.Axis); i++ {
		pad.Axis[i] = 0
	}
	for i := 0; i < len(pad.Button); i++ {
		pad.Button[i] = false
	}
}

//GetLeftStick return a the position of leftStick of the idPad already deadZone checked and can be normalized
func GetLeftStick(gamepad OpGamepad, normalized bool) OpVector2f {
	tmpStick := OpVector2f{float64(gamepad.Axis[sdl.CONTROLLER_AXIS_LEFTX]), float64(gamepad.Axis[sdl.CONTROLLER_AXIS_LEFTY])}
	if Distance(tmpStick, OpVector2f{0.0, 0.0}) < gamepad.deadZone {
		tmpStick = OpVector2f{}
	}
	if normalized {
		tmpStick.NormalizeVector()
	}
	return tmpStick
}

//GetRightStick return a the position of rightStick already deadZone checked and can be normalized
func GetRightStick(gamepad OpGamepad, normalized bool) OpVector2f {
	tmpStick := OpVector2f{float64(gamepad.Axis[sdl.CONTROLLER_AXIS_RIGHTX]), float64(gamepad.Axis[sdl.CONTROLLER_AXIS_RIGHTY])}
	if Distance(tmpStick, OpVector2f{0.0, 0.0}) < gamepad.deadZone {
		tmpStick = OpVector2f{}
	}
	if normalized {
		tmpStick.NormalizeVector()
	}
	return tmpStick
}

//OpInput take input from sdl int the event method and stock into a different array (it serve to be able to set to 0 key after they have been used)
type OpInput struct {
	deadZone float64 //need to check if a global deadZone for all controller works or need more precision
	Gamepads []OpGamepad
	KeyState map[sdl.Keycode]bool
}

func (input *OpInput) init(nDeadZone float64) {
	input.deadZone = nDeadZone

	input.KeyState = make(map[sdl.Keycode]bool)
	for i := 0; i < sdl.NumJoysticks(); i++ {
		input.Gamepads = append(input.Gamepads, OpGamepad{deadZone: input.deadZone})
	}
}

//Empty set all value in input, expect deadZone, to 0 or false
func (input *OpInput) Empty() {
	for i := 0; i < len(input.Gamepads); i++ {
		input.Gamepads[i].empty()
	}
	input.KeyState = make(map[sdl.Keycode]bool)
}

func (input *OpInput) pushGamepad() {
	input.Gamepads = append(input.Gamepads, OpGamepad{deadZone: input.deadZone})
}

func (input *OpInput) deleteGamepad(idPad sdl.JoystickID) {
	input.Gamepads = append(input.Gamepads[:idPad], input.Gamepads[idPad+1:]...)
}
