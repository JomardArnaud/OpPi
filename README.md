# OpPi
Go package with a scene management for SDL2 binding of veandco

this is what you main should be

func main() {
	game := OpGame{}
	if game.Init("./assets/config/") {
		//here you push you scene before enter into the game.loop (important to precise the IDScene here!)
		game.PushScenes(&sTest{IDScene: "test"}, &sTest2{IDScene: "test2"})

		game.Loop()
	}
	game.Clean()
}

and there are exemple of two scene who communicate each other

type sTest2 struct {
	IDScene, testFromOther string
	testSprite             OpSprite
}

func (sc *sTest2) GetFileConfig() string {
	return sc.IDScene
}

func (sc *sTest2) Init(gameInfo OpGameConfig, renderer *sdl.Renderer) {
	scInfo := OpInfoParser{}
	scInfo.Init(gameInfo.pathConfig + sc.IDScene + ".json")

	sc.testSprite.InitFromFile(renderer, scInfo.Blocks["sprite"])
	sc.testSprite.SetSize(Convert2i(gameInfo.sizeWindow))
	fmt.Println("SC 2 :", sc.testFromOther)
}

func (sc *sTest2) Reset(gameInfo OpGameConfig, inputInfo *OpInput) {
	fmt.Println("SC 2 :", sc.testFromOther)
}

func (sc *sTest2) Update(gameInfo OpGameConfig, inputInfo *OpInput, elapsedTime float64) string {
	if inputInfo.Gamepads[0].Button[sdl.CONTROLLER_BUTTON_A] {
		return "test"
	}
	return sc.IDScene

}

func (sc *sTest2) Draw(renderer *sdl.Renderer) {
	sc.testSprite.Draw(renderer)
}

func (sc *sTest2) PassInfoToNextScene(nextScene IOpScene) {
}

type sTest struct {
	IDScene    string
	testSprite OpSprite
	testAnim   *OpAnimator
	testAnim2  *OpAnimator
}

func (sc *sTest) GetFileConfig() string {
	return sc.IDScene
}

func (sc *sTest) Init(gameInfo OpGameConfig, renderer *sdl.Renderer) {
	scInfo := OpInfoParser{}
	scInfo.Init(gameInfo.pathConfig + sc.IDScene + ".json")

	sc.testSprite.InitFromFile(renderer, scInfo.Blocks["sprite"])
	sc.testSprite.SetSize(Convert2i(gameInfo.sizeWindow))
	sc.testAnim = NewOpAnimatorFromFile(renderer, scInfo.Blocks["animations"].Blocks["jimmySprite"])
	sc.testAnim2 = NewOpAnimatorFromFile(renderer, scInfo.Blocks["animations"].Blocks["loicSprite"])
}

func (sc *sTest) Reset(gameInfo OpGameConfig, inputInfo *OpInput) {
	sc.testAnim.X = int32(gameInfo.sizeWindow.X / 2)
	sc.testAnim.Y = int32(gameInfo.sizeWindow.Y / 2)
}

func (sc *sTest) Update(gameInfo OpGameConfig, inputInfo *OpInput, elapsedTime float64) string {
	var force OpVector2f
	if len(inputInfo.Gamepads) >= 1 {
		force = GetLeftStick(inputInfo.Gamepads[0], true)
	}
	force.MulForce(elapsedTime * 1000)
	sc.testAnim.Move(Convert2i(force))

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

func (sc *sTest) PassInfoToNextScene(nextScene IOpScene) {
	aller := nextScene.(*sTest2)

	aller.testFromOther = "attends sérieusement !"
}

next update =>
    add field mouse to OpInput
    add a audio part
    add a manager for assets
