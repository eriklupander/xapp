package imageprocessor

import (
	"bytes"
	"github.com/disintegration/gift"
	"image"
	"image/jpeg"
)

type ImageProcessor interface {
	Hipsterize(src image.Image, buf *bytes.Buffer) error
}

type GiftImageProcessor struct {
	filterChain *gift.GIFT
	resampler   gift.Resampling
}

func NewGiftImageProcessor(resampler gift.Resampling) *GiftImageProcessor {
	if resampler == nil {
		resampler = gift.LanczosResampling
	}
	filterChain := gift.New(
		gift.Saturation(20),
		gift.Contrast(20),
	)
	return &GiftImageProcessor{filterChain: filterChain, resampler: resampler}
}

// Hipsterize performs a sepia conversion
func (gip *GiftImageProcessor) Hipsterize(src image.Image, buf *bytes.Buffer) error {

	dst := image.NewRGBA(gip.filterChain.Bounds(src.Bounds()))

	// Use Draw func to apply the filters to src and store the result in dst:
	gip.filterChain.Draw(dst, src)

	return jpeg.Encode(buf, dst, nil)
}
