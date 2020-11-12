# OpPi
Go package with a scene management for SDL2 binding of veandco, and some other stuff to make you 2D game creation easier and faster

## Table of contents
* [Technologies](#technologies)
* [Setup](#setup)
* [NextUpdate](#next-update)

## Technologies
Project is created with:
* Golang version 
* https://github.com/veandco/go-sdl2

## Setup
(In progress)
### Main
This is what you main should be, it is the core you need, but feel free to add other stuff just make sur you put all these element in this order

```golang
func main() {
    game := oppi.OpGame{}
    //to init all the stuff OpGame need, it take your path for config's file in param
	if game.Init("./assets/config/") {
		//here you push you scene before enter into the game.loop (important to precise the IDScene here!)
		game.PushScenes(&sTest{IDScene: "test"}, &sTest2{IDScene: "test2"})
        //the main loop
		game.Loop()
    }
    //to clean some stuff from the SDL2
	game.Clean()
}
```

### Scene exemple
And there are exemple of two scene who communicate each other (they use some config file)

First Scene sTest ==>

```golang
type sTest struct {
	IDScene    string
	testSprite oppi.OpSprite
	testAnim   *oppi.OpAnimator
	testAnim2  *oppi.OpAnimator
}

func (sc *sTest) GetFileConfig() string {
	return sc.IDScene
}

func (sc *sTest) Init(gameInfo oppi.OpGameConfig, renderer *sdl.Renderer) {
	scInfo := oppi.OpInfoParser{}
	scInfo.Init(gameInfo.PathConfig + sc.IDScene + ".json")

	sc.testSprite.InitFromFile(renderer, scInfo.Blocks["sprite"])
	sc.testSprite.SetSize(oppi.Convert2i(gameInfo.SizeWindow))
	sc.testAnim = oppi.NewOpAnimatorFromFile(renderer, scInfo.Blocks["animations"].Blocks["jimmySprite"])
	sc.testAnim2 = oppi.NewOpAnimatorFromFile(renderer, scInfo.Blocks["animations"].Blocks["loicSprite"])
}

func (sc *sTest) Reset(gameInfo oppi.OpGameConfig, inputInfo *oppi.OpInput) {
	sc.testAnim.X = int32(gameInfo.SizeWindow.X / 2)
	sc.testAnim.Y = int32(gameInfo.SizeWindow.Y / 2)
}

func (sc *sTest) Update(gameInfo oppi.OpGameConfig, inputInfo *oppi.OpInput, elapsedTime float64) string {
	var force oppi.OpVector2f
	if len(inputInfo.Gamepads) >= 1 {
		force = oppi.GetLeftStick(inputInfo.Gamepads[0], true)
	}
	force.MulForce(elapsedTime * 1000)
	sc.testAnim.Move(oppi.Convert2i(force))

	sc.testAnim.Update(elapsedTime)
	sc.testAnim2.Update(elapsedTime)
	if inputInfo.KeyState[sdl.K_a] {
		return "test2"
	}
	return sc.IDScene
}

func (sc *sTest) Draw(renderer *sdl.Renderer) {
	sc.testAnim.Draw(renderer)
	sc.testAnim2.Draw(renderer)
}

func (sc *sTest) PassInfoToNextScene(nextScene oppi.IOpScene) {
	aller := nextScene.(*sTest2)

	aller.testFromOther = "attends sérieusement !"
}
```

Second Scene sTest2 ==>

```golang
//simple file test to show a scene's exemple
type sTest2 struct {
	IDScene, testFromOther string
	testSprite             oppi.OpSprite
}

func (sc *sTest2) GetFileConfig() string {
	return sc.IDScene
}

func (sc *sTest2) Init(gameInfo oppi.OpGameConfig, renderer *sdl.Renderer) {
	scInfo := oppi.OpInfoParser{}
	scInfo.Init(gameInfo.PathConfig + sc.IDScene + ".json")

	sc.testSprite.InitFromFile(renderer, scInfo.Blocks["sprite"])
	sc.testSprite.SetSize(oppi.Convert2i(gameInfo.SizeWindow))
	fmt.Println("SC 2 :", sc.testFromOther)
}

func (sc *sTest2) Reset(gameInfo oppi.OpGameConfig, inputInfo *oppi.OpInput) {
	fmt.Println("SC 2 :", sc.testFromOther)
}

func (sc *sTest2) Update(gameInfo oppi.OpGameConfig, inputInfo *oppi.OpInput, elapsedTime float64) string {
	if inputInfo.Gamepads[0].Button[sdl.CONTROLLER_BUTTON_A] {
		return "test"
	}
	return sc.IDScene

}

func (sc *sTest2) Draw(renderer *sdl.Renderer) {
	sc.testSprite.Draw(renderer)
}

func (sc *sTest2) PassInfoToNextScene(nextScene oppi.IOpScene) {
}
```

### Config file
Some exemple for you config's file (game.JSON and managerScene.JSON are need good works of the package, make sure to put it into your project)
```JSON
game.JSON =>
{
    "config":{
      "deadZone":"10000.0",
      "sizeWindow":"1600.0,900.0",
      "posWindow":"10.0,10.0",
      "timer":"0.0167"  
	}
}
managerScene.JSON =>
{
    "start":{
        "startingScene":"test"
    }
}
test.JSON =>
{
    "animations":{
        "jimmySprite":{
            "pathTex":"assets/sprite/sprite_sheet_mind_master.png",
            "rectSprite":"50,350,95,135",
            "sizeCase":"0,0,96,138",
            "nbAnim":"16",
            "framePerline":"4",
            "startingLine":"0",
            "timePerFrame":"0.1"    
        }
    }
}
test2.JSON =>
{
    "sprite":{
        "loadedAtStart":"true",
        "pathSprite":"assets/sprite/CDEGEULASS.png",
        "rectSprite":"0,0,300,300"
    }
}
```

## Next update
* add field mouse to OpInput
* add a audio part
* add a manager for assets
* remake the config of OpGame to be more flexible
* make a documentation
