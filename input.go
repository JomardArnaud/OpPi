package oppi

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

/*
enum for Keyboard => https://wiki.libsdl.org/SDL_Keycode
	 for button => https://wiki.libsdl.org/SDL_GameControllerButton
	 for axis => https://wiki.libsdl.org/SDL_GameControllerAxis
	 for mouse => https://wiki.libsdl.org/SDL_MouseButtonEvent#Remarks
*/

//OpGamepad is just a layer from sdl.Gamepad with usefull methode
type OpGamepad struct {
	//Axis   [sdl.CONTROLLER_AXIS_MAX]time.Timer
	Button map[sdl.GameControllerButton][]OpCooldown
}

//OpInput take input from sdl int the event method and stock into a different array (it serve to be able to set to 0 key after they have been used)
type OpInput struct {
	DeadZone, Buffer float64
	Gamepads         []*sdl.GameController
	KeyState         []uint8
	GamepadsBuffer   []OpGamepad
	KeyStateBuffer   map[sdl.Keycode][]OpCooldown
}

func newOpGamepad() OpGamepad {
	dest := OpGamepad{}
	dest.Button = make(map[sdl.GameControllerButton][]OpCooldown)
	return (dest)
}

func (pad *OpGamepad) init() {
	pad.Button = make(map[sdl.GameControllerButton][]OpCooldown)
}

func (pad *OpGamepad) update(buffer, elapsedTime float64) {
	for buttonID, buttonBuffer := range pad.Button {
		for i := sdl.GameControllerButton(0); int(i) < len(buttonBuffer); i++ {
			buttonBuffer[i].Update(elapsedTime)
			if buttonBuffer[i].Tick == 0 {
				pad.Button[buttonID] = buttonBuffer[:int(i)+copy(buttonBuffer[i:], buttonBuffer[i+1:])]
				if len(pad.Button[buttonID]) == 0 {
					pad.Button[buttonID] = nil
				}
			}
		}
	}
}

func (pad *OpGamepad) empty(buffer float64) {
	for i := 0; i < len(pad.Button); i++ {
		for n := 0; n < len(pad.Button[sdl.GameControllerButton(i)]); n++ {
			pad.Button[sdl.GameControllerButton(i)][n].Reset()
		}
	}
}

func (input *OpInput) init(nDeadZone, buffer float64) {
	input.DeadZone = nDeadZone
	input.Buffer = buffer
	fmt.Println("nb gamepad : ", sdl.NumJoysticks())
	input.KeyState = sdl.GetKeyboardState()
	input.KeyStateBuffer = make(map[sdl.Keycode][]OpCooldown)
}

func (input *OpInput) update(elapsedTime float64) {
	for i := 0; i < len(input.GamepadsBuffer); i++ {
		input.GamepadsBuffer[i].update(input.Buffer, elapsedTime)
	}
	for keyID, keyBuffer := range input.KeyStateBuffer {
		for i := sdl.Keycode(0); int(i) < len(keyBuffer); i++ {
			keyBuffer[i].Update(elapsedTime)
			if keyBuffer[i].Tick == 0 {
				input.KeyStateBuffer[keyID] = keyBuffer[:int(i)+copy(keyBuffer[i:], keyBuffer[i+1:])]
				if len(input.KeyStateBuffer[keyID]) == 0 {
					input.KeyStateBuffer[keyID] = nil
				}
			}
		}
	}
	fmt.Println("keyboard buffer = ", input.KeyStateBuffer)
}

//Empty set all value in input, expect deadZone, to 0 or false
func (input *OpInput) Empty() {
	for i := 0; i < len(input.GamepadsBuffer); i++ {
		input.GamepadsBuffer[i].empty(input.Buffer)
	}
	input.KeyStateBuffer = make(map[sdl.Keycode][]OpCooldown)
}

//GetLeftStick return a the position of leftStick of the idPad already deadZone checked and can be normalized
func (input *OpInput) GetLeftStick(ID sdl.GameControllerAxis, normalized bool) OpVector2f {
	tmpStick := OpVector2f{float64(input.Gamepads[ID].Axis(sdl.CONTROLLER_AXIS_LEFTX)), float64(input.Gamepads[ID].Axis(sdl.CONTROLLER_AXIS_LEFTY))}
	if Distance(tmpStick, OpVector2f{0.0, 0.0}) < input.DeadZone {
		tmpStick = OpVector2f{}
	}
	if normalized {
		tmpStick.NormalizeVector()
	}
	return tmpStick
}

//GetRightStick return a the position of rightStick already deadZone checked and can be normalized
func (input *OpInput) GetRightStick(ID sdl.GameControllerAxis, normalized bool) OpVector2f {
	tmpStick := OpVector2f{float64(input.Gamepads[ID].Axis(sdl.CONTROLLER_AXIS_RIGHTX)), float64(input.Gamepads[ID].Axis(sdl.CONTROLLER_AXIS_RIGHTY))}
	if Distance(tmpStick, OpVector2f{0.0, 0.0}) < input.DeadZone {
		tmpStick = OpVector2f{}
	}
	if normalized {
		tmpStick.NormalizeVector()
	}
	return tmpStick
}

func (input *OpInput) pushGamepad() {
	input.Gamepads = append(input.Gamepads, sdl.GameControllerOpen(sdl.NumJoysticks()-1))
	input.GamepadsBuffer = append(input.GamepadsBuffer, newOpGamepad())
}

func (input *OpInput) deleteGamepad(ID sdl.JoystickID) {
	//to avoid memory leaks
	copy(input.Gamepads[ID:], input.Gamepads[ID+1:])
	input.Gamepads[len(input.Gamepads)-1] = nil // or the zero value of T
	input.Gamepads = input.Gamepads[:len(input.Gamepads)-1]

	input.GamepadsBuffer = input.GamepadsBuffer[:int(ID)+copy(input.GamepadsBuffer[ID:], input.GamepadsBuffer[ID+1:])]
}

func (input *OpInput) pushNewButtonBuffer(ID sdl.JoystickID, buttonID uint8) {
	input.GamepadsBuffer[ID].Button[sdl.GameControllerButton(buttonID)] = append(input.GamepadsBuffer[ID].Button[sdl.GameControllerButton(buttonID)],
		OpCooldown{Tick: input.Buffer, Duration: input.Buffer})
}

func (input *OpInput) pushNewAxisBuffer(ID sdl.JoystickID, buttonID uint8) {
	input.GamepadsBuffer[ID].Button[sdl.GameControllerButton(buttonID)] = append(input.GamepadsBuffer[ID].Button[sdl.GameControllerButton(buttonID)],
		OpCooldown{Tick: input.Buffer, Duration: input.Buffer})
}

func (input *OpInput) pushNewKeyBuffer(keyID sdl.Keycode) {
	input.KeyStateBuffer[keyID] = append(input.KeyStateBuffer[keyID],
		OpCooldown{Tick: input.Buffer, Duration: input.Buffer})
}
