package oppi

import (
	"image/png"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

//OpSprite is a struct here to help the management of *sdl.texture
type OpSprite struct {
	tex                     *sdl.Texture
	rectTexture, rectSprite sdl.Rect
}

//InitFromFile from a JSON file
func (sprite *OpSprite) InitFromFile(renderer *sdl.Renderer, spriteBlock block) {
	var sizeTex OpVector2i
	sprite.tex, sizeTex = FillTexture(renderer, spriteBlock.Info["pathSprite"])
	sprite.rectSprite = OpSetSdlRect(spriteBlock.Info["rectSprite"])
	sprite.rectTexture = sdl.Rect{X: 0, Y: 0, W: sizeTex.X, H: sizeTex.Y}
}

//Init the values sprite
func (sprite *OpSprite) Init(nTex *sdl.Texture, nRectTexture sdl.Rect, posSprite, sizeSprite OpVector2i) {
	sprite.tex = nTex
	sprite.rectTexture = nRectTexture
	sprite.rectSprite = sdl.Rect{X: posSprite.X, Y: posSprite.Y, W: sizeSprite.X, H: sizeSprite.Y}
}

//Move the position's sprite
func (sprite *OpSprite) Move(force OpVector2i) {
	sprite.rectSprite.X += force.X
	sprite.rectSprite.Y += force.Y
}

//SetPosition of the sprite
func (sprite *OpSprite) SetPosition(nPos OpVector2i) {
	sprite.rectSprite.X = nPos.X
	sprite.rectSprite.Y = nPos.Y
}

//SetSize of the sprite
func (sprite *OpSprite) SetSize(nSize OpVector2i) {
	sprite.rectSprite.W = nSize.X
	sprite.rectSprite.H = nSize.Y
}

//Draw the sprite of the renderer
func (sprite *OpSprite) Draw(renderer *sdl.Renderer) {
	renderer.Copy(sprite.tex, &sprite.rectTexture, &sprite.rectSprite)
}

//function related to Sprite and Texture

//FillTexture is vast
func FillTexture(renderer *sdl.Renderer, filename string) (*sdl.Texture, OpVector2i) {

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
	return tex, OpVector2i{int32(w), int32(h)}
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)
	return tex
}
