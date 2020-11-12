package oppi

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

//OpGameConfig all info may be usefull for other companent
type OpGameConfig struct {
	PathConfig            string
	DeadZone, Timer       float64
	PosWindow, SizeWindow OpVector2f
}

//Init OpGameConfig with info from parser
func (info *OpGameConfig) Init(parser OpInfoParser) {
	info.DeadZone = OpSetFloat(parser.Blocks["config"].Info["deadZone"])
	info.Timer = OpSetFloat(parser.Blocks["config"].Info["timer"])
	info.PosWindow = OpSetOpVector2f(parser.Blocks["config"].Info["posWindow"])
	info.SizeWindow = OpSetOpVector2f(parser.Blocks["config"].Info["sizeWindow"])
}

//OpGame is the main struct which contains all info to run the game
type OpGame struct {
	game   bool
	config OpGameConfig
	//sdl lib stuff
	MainWindow   *sdl.Window
	MainRenderer *sdl.Renderer
	Gamepads     []*sdl.GameController
	KeyState     []uint8
	//stuff for the game loop
	elapsedTime float64
	infoInput   OpInput
	clock       time.Time
	scManager   OpSceneManager
}

//func Greeting(prefix string, who ...string)

//PushScenes in the OpSceneManager
func (game *OpGame) PushScenes(scenes ...IOpScene) {
	for _, scene := range scenes {
		game.scManager.pushScene(game.config, game.MainRenderer, scene)
	}
}

//Init all setup for running the game
func (game *OpGame) Init(nPathConfig string) bool {
	//main info game init
	game.game = true
	gameParser := OpInfoParser{}
	game.config.PathConfig = nPathConfig
	gameParser.Init(game.config.PathConfig + "game.json")

	game.config.Init(gameParser)
	if game.initSdlContext() == false {
		return false
	}
	game.scManager.init(game.config, game.MainRenderer)
	return true
}

func (game *OpGame) initSdlContext() bool {
	//sdl things init
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//to put setup on window i can do a enum + func getStringFrom enum witcht return a type sdl from id of enum on a array
	game.MainWindow, err = sdl.CreateWindow("Pacific Punch", int32(game.config.PosWindow.X), int32(game.config.PosWindow.Y),
		int32(game.config.SizeWindow.X), int32(game.config.SizeWindow.Y), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return false
	}
	game.MainRenderer, err = sdl.CreateRenderer(game.MainWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return false
	}
	for i := 0; i < sdl.NumJoysticks(); i++ {
		game.Gamepads = append(game.Gamepads, sdl.GameControllerOpen(i))
	}
	game.KeyState = sdl.GetKeyboardState()
	game.infoInput.init(game.config.DeadZone)
	return true
}

//Loop is the OpGame main loop
func (game *OpGame) Loop() bool {
	game.clock = time.Now()
	for game.game {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			game.event(event)
		}
		elapsedTime := time.Since(game.clock).Seconds()
		if elapsedTime > game.config.Timer {
			game.update(elapsedTime)
			game.draw()
			game.clock = time.Now()
		} else {
			sdl.Delay(uint32(game.config.Timer - elapsedTime))
		}
	}
	return game.game
}

//this func need to fill two object one to manage keyboard and one to manage controller (if it not the case i can't let OpGame unsettable)
func (game *OpGame) event(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		game.game = false
		break
	case *sdl.ControllerDeviceEvent:
		switch t.GetType() {
		case sdl.CONTROLLERDEVICEADDED:
			fmt.Print("pad connected")
			game.Gamepads = append(game.Gamepads, sdl.GameControllerOpen(sdl.NumJoysticks()-1))
			game.infoInput.pushGamepad()
			break
		case sdl.CONTROLLERDEVICEREMOVED:
			if len(game.Gamepads) == 1 {
				game.Gamepads = nil
				game.infoInput.Gamepads = nil
			} else {
				game.Gamepads = append(game.Gamepads[:t.Which], game.Gamepads[t.Which+1:]...)
				game.infoInput.deleteGamepad(t.Which)
			}
			break
		}
		break
	case *sdl.ControllerAxisEvent:
		game.infoInput.Gamepads[t.Which].Axis[t.Axis] = int(t.Value)
		break
	case *sdl.ControllerButtonEvent:
		game.infoInput.Gamepads[t.Which].Button[int(t.Button)] = t.State != 0
		break
	case *sdl.KeyboardEvent:
		if t.GetType() == sdl.KEYUP {
			if t.Keysym.Sym == sdl.K_ESCAPE { //must delete it
				game.game = false
			}
			game.infoInput.KeyState[t.Keysym.Sym] = false
			break
		}
		if t.GetType() == sdl.KEYDOWN {
			if t.Repeat == 0 {
				game.infoInput.KeyState[t.Keysym.Sym] = true
			}
			break
		}
		break
	}
}

func (game *OpGame) update(elapsedTime float64) {
	game.scManager.update(elapsedTime, game.config, &game.infoInput)
}

func (game *OpGame) draw() {
	game.MainRenderer.Clear()
	game.scManager.draw(game.MainRenderer)
	game.MainRenderer.Present()
}

//Clean up some stuff from sdl lib
func (game *OpGame) Clean() {
	game.MainWindow.Destroy()
	game.MainRenderer.Destroy()
	for i := 0; i < sdl.NumJoysticks(); i++ {
		defer game.Gamepads[i].Close()
	}
	sdl.Quit()
}
