package oppi

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

//OpGameConfig all info may be usefull for other companent
type OpGameConfig struct {
	PathConfig              string
	DeadZone, Buffer, Timer float64
	PosWindow, SizeWindow   OpVector2f
}

//OpGame is the main struct which contains all info to run the game
type OpGame struct {
	//sdl lib stuff
	MainWindow   *sdl.Window
	MainRenderer *sdl.Renderer
	//stuff for the game loop
	game        bool
	config      OpGameConfig
	elapsedTime float64
	Input       OpInput
	clock       time.Time
	scManager   OpSceneManager
}

//Init OpGameConfig with info from parser
func (info *OpGameConfig) Init(parser OpInfoParser) {
	info.DeadZone = OpSetFloat(parser.Blocks["config"].Info["deadZone"])
	info.Buffer = OpSetFloat(parser.Blocks["config"].Info["buffer"])
	info.Timer = OpSetFloat(parser.Blocks["config"].Info["timer"])
	info.PosWindow = OpSetOpVector2f(parser.Blocks["config"].Info["posWindow"])
	info.SizeWindow = OpSetOpVector2f(parser.Blocks["config"].Info["sizeWindow"])
}

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
	game.Input.init(game.config.DeadZone, game.config.Buffer)
	return true
}

//Loop is the OpGame main loop
func (game *OpGame) Loop() bool {
	game.clock = time.Now()
	for game.game {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			game.fillInputBuffer(event)
			game.scManager.event(event)
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
func (game *OpGame) fillInputBuffer(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		game.game = false
		break
	case *sdl.ControllerDeviceEvent:
		switch t.GetType() {
		case sdl.CONTROLLERDEVICEADDED:
			game.Input.pushGamepad()
			break
		case sdl.CONTROLLERDEVICEREMOVED:
			if len(game.Input.Gamepads) == 1 {
				game.Input.Gamepads = nil
				game.Input.GamepadsBuffer = nil
			} else {
				game.Input.deleteGamepad(t.Which)
			}
			break
		}
		break
	// case *sdl.ControllerAxisEvent:
	// 	//standby i need to take time to think about how handle axis movement
	// 	game.Input.Gamepads[t.Which].Axis[t.Axis] = int(t.Value)
	// 	break
	case *sdl.ControllerButtonEvent:
		if t.State != 0 { //push a new buffer only if the button is push
			game.Input.pushNewButtonBuffer(t.Which, t.Button)
		}
		break
	case *sdl.KeyboardEvent:
		if t.GetType() == sdl.KEYUP {
			if t.Keysym.Sym == sdl.K_ESCAPE { //must delete it
				game.game = false
			}
			break
		}
		if t.GetType() == sdl.KEYDOWN {
			if t.Repeat == 0 {
				game.Input.pushNewKeyBuffer(t.Keysym.Sym)
			}
			break
		}
		break
	}
}

func (game *OpGame) update(elapsedTime float64) {
	game.Input.update(elapsedTime)
	game.scManager.update(elapsedTime, game.config, &game.Input)
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
		defer game.Input.Gamepads[i].Close()
	}
	sdl.Quit()
}
