package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"sync"
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
	holder := image.NewRGBA(b)

	resHolder := make([]*image.Paletted, len(im.Image))

	var wg sync.WaitGroup
	// transforming each frame.
	for index, frame := range im.Image {
		bounds := frame.Bounds()
		prev := holder
		draw.Draw(holder, bounds, frame, bounds.Min, draw.Over)

		wg.Add(1)
		imageToPaletted(&wg, resHolder, index, transform(holder), frame.Palette)

		switch im.Disposal[index] {
		case gif.DisposalBackground:
			holder = image.NewRGBA(b)
		case gif.DisposalPrevious:
			holder = prev
		}
		// on gif.DisposalNone keep use written holder
	}

	wg.Wait()

	im.Image = resHolder
	// Set new height and width into config
	im.Config.Width = im.Image[0].Bounds().Max.X
	im.Config.Height = im.Image[0].Bounds().Max.Y

	return gif.EncodeAll(w, im)
}

func imageToPaletted(wg *sync.WaitGroup, resHolder []*image.Paletted, index int, img image.Image, p color.Palette) {
	defer wg.Done()

	b := img.Bounds()
	pm := image.NewPaletted(b, p)
	// use zero image point
	draw.FloydSteinberg.Draw(pm, b, img, image.Point{})

	resHolder[index] = pm
}

// ImgResTransformer is a function that transforms an image.
type ImgResTransformer func(image.Image) image.Image
