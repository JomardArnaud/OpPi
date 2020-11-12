package oppi

import (
	"github.com/veandco/go-sdl2/sdl"
)

//IOpScene is interface for scene
type IOpScene interface {
	GetFileConfig() string
	Init(OpGameConfig, *sdl.Renderer)
	Reset(OpGameConfig, *OpInput)
	Update(OpGameConfig, *OpInput, float64) string
	Draw(*sdl.Renderer)
	PassInfoToNextScene(IOpScene)
}

//OpSceneManager contains and manage scene's game
type OpSceneManager struct {
	idScene, idPrevScene string
	allScene             map[string]IOpScene //be carefull to not use twice the same name for a scene
}

func (manager *OpSceneManager) pushScene(gameInfo OpGameConfig, renderer *sdl.Renderer, nSc IOpScene) {
	manager.allScene[nSc.GetFileConfig()] = nSc
	manager.allScene[nSc.GetFileConfig()].Init(gameInfo, renderer)
}

func (manager *OpSceneManager) init(gameInfo OpGameConfig, renderer *sdl.Renderer) {
	infoManager := OpInfoParser{}
	infoManager.Init(gameInfo.pathConfig + "managerScene.json")

	manager.idScene = infoManager.Blocks["start"].Info["startingScene"] //gameInfo.startingScene
	manager.allScene = make(map[string]IOpScene)
}

func (manager *OpSceneManager) update(elapsedTime float64, gameInfo OpGameConfig, infoInput *OpInput) {
	manager.idPrevScene = manager.idScene
	manager.idScene = manager.allScene[manager.idScene].Update(gameInfo, infoInput, elapsedTime)
	if manager.idPrevScene != manager.idScene {
		manager.allScene[manager.idPrevScene].PassInfoToNextScene(manager.allScene[manager.idScene])
		manager.allScene[manager.idScene].Reset(gameInfo, infoInput)
		infoInput.Empty()
	}
}

func (manager *OpSceneManager) draw(renderer *sdl.Renderer) {
	manager.allScene[manager.idScene].Draw(renderer)
}
