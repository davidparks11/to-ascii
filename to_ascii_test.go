package toascii_test

import (
	"image/png"
	"os"
	"testing"

	toascii "github.com/davidparks11/to-ascii"
)

func TestImageToAscii(t *testing.T) {
	inputImgContent, err := os.Open("./testFiles/in.png")
	if err != nil {
		t.Fatal("could not read test image")
	}

	expected, err := os.ReadFile("./testFiles/expected.text")
	if err != nil {
		t.Fatal("could not read expected file")
	}
	
	img, err := png.Decode(inputImgContent)
	if err != nil {
		t.Fatal("could not decode png")
	}

	actual := string(toascii.ImageToText(img, 1, 1, false))
 
	if actual != string(expected) {
		t.Errorf("actual did not match expected\n%s!=\n%s", actual, string(expected))
	}
}

func TestConverter(t *testing.T) {
	inputImgContent, err := os.Open("./testFiles/in.png")
	if err != nil {
		t.Fatal("could not read test image")
	}

	img, err := png.Decode(inputImgContent)
	if err != nil {
		t.Fatal("could not decode png")
	}

	c := toascii.NewConverter(img)
	if string(toascii.ImageToText(img, 1, 1, true)) != string(c.PixelsPerCharacterHeight(1).PixelsPerCharacterWidth(1).Invert().Convert()) {
		t.Error("Convertor did not yield equal result given equivelent options")
	}
}
