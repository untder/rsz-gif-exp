package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
)

// processImage :nodoc:
func processImage(w io.Writer, r io.Reader, transform ImgResTransformer) error {
	if transform == nil {
		_, err := io.Copy(w, r)
		return err
	}

	// decode using "image/gif" package
	im, err := gif.DecodeAll(r)
	if err != nil {
		return err
	}

	// get first frame
	firstFrame := im.Image[0].Bounds()
	b := image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy())
	// store first frame in holder
	img := image.NewRGBA(b)

	resImgPal := make([]*image.Paletted, len(im.Image))
	// transforming each frame.
	for index, frame := range im.Image {
		bounds := frame.Bounds()
		prev := img
		draw.Draw(img, bounds, frame, bounds.Min, draw.Over)
		resImgPal[index] = imageToPaletted(transform(img), frame.Palette)

		switch im.Disposal[index] {
		case gif.DisposalBackground:
			img = image.NewRGBA(b)
		case gif.DisposalPrevious:
			img = prev
		}
	}

	im.Image = resImgPal

	// Set new height and width into config
	im.Config.Width = im.Image[0].Bounds().Max.X
	im.Config.Height = im.Image[0].Bounds().Max.Y

	return gif.EncodeAll(w, im)
}

func imageToPaletted(img image.Image, p color.Palette) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, p)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
}

// ImgResTransformer is a function that transforms an image.
type ImgResTransformer func(image.Image) image.Image
