package oppi

import (
	"image/png"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

//OpAnimator manage a Texture to draw a suite's sprite on renderer
type OpAnimator struct {
	sdl.Rect
	caseTexRect                                       sdl.Rect
	animeTex                                          *sdl.Texture
	nbAnim, idFrame, line, startingLine, framePerLine int32
	timePerFrame, timeFrame                           float64
}

//Move the pos's sprite with a force
func (anime *OpAnimator) Move(force OpVector2i) {
	anime.X += force.X
	anime.Y += force.Y
}

//SetPosition of the sprite
func (anime *OpAnimator) SetPosition(nPos OpVector2i) {
	anime.X = nPos.X
	anime.Y = nPos.Y
}

//SetSize of the sprite
func (anime *OpAnimator) SetSize(nSize OpVector2i) {
	anime.W = nSize.X
	anime.H = nSize.Y
}

//Update the rect's texture to select the corect part's image to display
func (anime *OpAnimator) Update(elapsedTime float64) {
	anime.timeFrame = Clamp(anime.timeFrame-elapsedTime, 0.0, anime.timePerFrame)

	if anime.timeFrame == 0.0 {
		anime.idFrame++
		if anime.idFrame%anime.framePerLine == 0 {
			anime.line++
		}
		if anime.idFrame == anime.nbAnim {
			anime.idFrame = 0
			anime.line = anime.startingLine
		}
		anime.timeFrame = anime.timePerFrame
	}
	anime.caseTexRect = sdl.Rect{X: (anime.idFrame % 4) * anime.caseTexRect.W, Y: anime.line * anime.caseTexRect.H, W: anime.caseTexRect.W, H: anime.caseTexRect.H}

}

//Draw rect's texture on a sdl.Renderer
func (anime *OpAnimator) Draw(renderer *sdl.Renderer) {
	renderer.Copy(anime.animeTex, &anime.caseTexRect, &anime.Rect)
}

//ImgFileToTexture fill a sdl.Texture with each pixel of a img
func ImgFileToTexture(renderer *sdl.Renderer, filename string) *sdl.Texture {

	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	pIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[pIndex] = byte(r / 256)
			pIndex++
			pixels[pIndex] = byte(g / 256)
			pIndex++
			pixels[pIndex] = byte(b / 256)
			pIndex++
			pixels[pIndex] = byte(a / 256)
			pIndex++
		}
	}
	tex := pixelsToTexture(renderer, pixels, w, h)
	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	return tex
}

//NewOpAnimatorFromFile return a ready to use OpAnimator using a block of .JSON to setup
func NewOpAnimatorFromFile(renderer *sdl.Renderer, infoAnim block) *OpAnimator {
	animeTex := ImgFileToTexture(renderer, infoAnim.Info["pathTex"])
	return &OpAnimator{OpSetSdlRect(infoAnim.Info["rectSprite"]), OpSetSdlRect(infoAnim.Info["sizeCase"]), animeTex,
		int32(OpSetInt(infoAnim.Info["nbAnim"])), 0, int32(OpSetInt(infoAnim.Info["startingLine"])), int32(OpSetInt(infoAnim.Info["startingLine"])),
		int32(OpSetInt(infoAnim.Info["framePerline"])), OpSetFloat(infoAnim.Info["timePerFrame"]), OpSetFloat(infoAnim.Info["timePerFrame"])}
}
