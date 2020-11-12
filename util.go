package oppi

import (
	"log"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

//OpRgba is a struct represent a color
type OpRgba struct {
	R, G, B, A byte
}

//OpCooldown is a little tool to control time
type OpCooldown struct {
	Tick, Duration float64
}

//Update the tick using the time elapsed
func (cd *OpCooldown) Update(elapsedTime float64) {
	cd.Tick = Clamp(cd.Tick-elapsedTime, 0.0, cd.Tick)
}

//Reset tick to the duration value
func (cd *OpCooldown) Reset() {
	cd.Tick = cd.Duration
}

//OpRect4i blabla
type OpRect4i struct {
	X, Y, W, H int32
}

//OpSetInt convert a string into a int
func OpSetInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return value
}

//OpSetFloat convert a string into a float64
func OpSetFloat(str string) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal(err)
	}
	return value
}

//OpSetByte convert a string into a byte
func OpSetByte(str string) byte {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return byte(value)
}

//OpSetBool convert a string into a bool
func OpSetBool(str string) bool {
	value, err := strconv.ParseBool(str)
	if err != nil {
		log.Fatal(err)
	}
	return value
}

//OpSetOpVector2f convert a string into a OpVector2f
func OpSetOpVector2f(str string) OpVector2f {
	if str != "" {
		s := strings.Split(str, ",")
		return OpVector2f{OpSetFloat(s[0]), OpSetFloat(s[1])}
	}
	return OpVector2f{}
}

//OpSetRgba convert a string into a OpRgba(a color struct basicaly)
func OpSetRgba(str string) OpRgba {
	if str != "" {
		s := strings.Split(str, ",")
		return OpRgba{OpSetByte(s[0]), OpSetByte(s[1]), OpSetByte(s[2]), OpSetByte(s[3])}
	}
	return OpRgba{}
}

//OpSetRect4f convert a string into a OpRect4f
func OpSetRect4f(str string) OpRect4f {
	if str != "" {
		s := strings.Split(str, ",")
		return OpRect4f{OpSetFloat(s[0]), OpSetFloat(s[1]), OpSetFloat(s[2]), OpSetFloat(s[3])}
	}
	return OpRect4f{}
}

//OpSetRect4i convert a string into a OpRect4i
func OpSetRect4i(str string) OpRect4i {
	if str != "" {
		s := strings.Split(str, ",")
		return OpRect4i{int32(OpSetInt(s[0])), int32(OpSetInt(s[1])), int32(OpSetInt(s[2])), int32(OpSetInt(s[3]))}
	}
	return OpRect4i{}
}

//OpSetSdlRect convert a string into a sdl.Rect
func OpSetSdlRect(str string) sdl.Rect {
	if str != "" {
		s := strings.Split(str, ",")
		return sdl.Rect{X: int32(OpSetInt(s[0])), Y: int32(OpSetInt(s[1])), W: int32(OpSetInt(s[2])), H: int32(OpSetInt(s[3]))}
	}
	return sdl.Rect{}
}
