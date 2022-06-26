package toascii


import (
	"image"
	"image/color"
	"math"

	"golang.org/x/exp/constraints"
)

//Converter is a helper struct for those who like chaining function calls
type Converter struct {
	img image.Image
	pixPerCharX int
	pixPerCharY int
	invert bool
}

//NewConverter returns a Converter for a given image with default values
func NewConverter(img image.Image) *Converter {
	return &Converter{img: img}
}

//PixelsPerCharacterWidth sets the number of pixels covered by each character
//across
func (c *Converter) PixelsPerCharacterWidth(x int) *Converter {
	c.pixPerCharX = x
	return c
}

//PixelsPerCharacterHeight sets the number of pixels covered by each character
//from top to bottom
func (c *Converter) PixelsPerCharacterHeight(y int) *Converter {
	c.pixPerCharY = y
	return c
}

//Invert will invert the brightness characters presetn in the converted text.
//Calling Invert an even number of times does nothing. There's no reason to do
//that as that'd be silly. 
func (c *Converter) Invert() *Converter {
	c.invert = !c.invert
	return c
}

//Convert returns the text for the given Converter struct create to this point.
func (c *Converter) Convert() []rune {
	return ImageToText(c.img, c.pixPerCharX, c.pixPerCharY, c.invert)
}

var density = []rune(" .,-=+:;cba!?0123456789$W#@Ñ")

//ImageToText takes an image.Image and returns that image represented by 
//characters of differing "brightness". pixPerCharX/Y represent how many pixels
//each character covers. Given 0 or less these paramters default to 1, meaning 
//one character to represent each pixel. Invert is a flag that will invert the
//level of "brightness" of each character.
func ImageToText(img image.Image, pixPerCharX, pixPerCharY int, invert bool) []rune {
	if pixPerCharX < 1 {
		pixPerCharX = 1
	}

	if pixPerCharY < 1 {
		pixPerCharY = 1
	}

	numCharX := img.Bounds().Dx()/pixPerCharX + 1 //add 1 to allow for \n between lines
	numCharY := img.Bounds().Dy()/pixPerCharY

	text := make([]rune, numCharX*numCharY)
	
	for y := 0; y < numCharY; y++ {
		for x := 0; x < numCharX - 1; x++ {
			pixelGroup := colorBlock(img, x, y, pixPerCharX, pixPerCharY) 
			avgLuma := averagePixelLuma(pixelGroup)

			densityIndex := scale(avgLuma, 0, math.MaxUint16, 0, len(density)-1) //map average luma value to an index in the density slice
			if invert {
				densityIndex = len(density)-1-densityIndex
			}
			text[y * numCharX + x] = density[densityIndex]
		}
		text[(y+1) * numCharX - 1] = '\n'
	}
	return text
}

func colorBlock(img image.Image, blockX, blockY, width, height int) []color.Color {
	x := blockX * width
	y := blockY * height
	colors := make([]color.Color, width*height)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			colors[i*width+j] = img.At(x+j, y+i)
		}
	}
	return colors
}

func averagePixelLuma(colors []color.Color) uint32 {
	var lum uint32 = 0

	for _, pixel := range colors {
		r, g, b, _ := pixel.RGBA()
		lum += luma(r, g, b)
	}
	return lum / uint32(len(colors))
}

func luma(r, g, b uint32) uint32 {
	return uint32(float64(r)*0.2126 + float64(g)*0.7152 + float64(b)*0.0722)
}

type numeric interface {
	constraints.Integer | constraints.Float
}

//scale takes a val of type K between fromLow and fromHigh and scales it to a 
//value of type T between toLow and toHigh.
func scale[K, T numeric](val, fromLow, fromHigh K, toLow, toHigh T) T {
	//switch tolower and tohigher in case of wrong order
	if fromLow > fromHigh {
		fromLow, fromHigh = fromHigh, fromLow
	}

	//switch toLower and toHigher in case of wrong order
	if toLow > toHigh {
		toLow, toHigh = toHigh, toLow
	}

	//if val is lower than low bounds, set val to fromLow
	if val < fromLow  {
		val = fromLow
	} 

	//if val is greater than high bounds set val to fromHigh
	if val > fromHigh {
		val = fromHigh
	}

	floatVal := float64(val)
	floatFromLow := float64(fromLow)
	floatFromHigh := float64(fromHigh)
	floatToLow := float64(toLow)
	floatToHigh := float64(toHigh)

	percentFrom := floatVal / (floatFromHigh - floatFromLow)
	return T(percentFrom*(floatToHigh-floatToLow) + floatToLow)
}
