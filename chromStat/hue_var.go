package chromStat

import (
	"math"
	"math/cmplx"

	"github.com/rennis250/colorCosmos/dkl"
)

func HueVariance(dklimg *dkl.DKLImage, binary_mask []bool) (sh float64) {
	var hues []float64

	b := dklimg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			dkl := dklimg.DKLAt(x, y)

			if !binary_mask[y*b.Max.X+x] {
				continue
			}

			hues = append(hues, math.Atan2(dkl.YV, dkl.RG))
		}
	}

	r := complex(0.0, 0.0)
	for _, v := range hues {
		r += cmplx.Exp(complex(0.0, v))
	}
	mr := cmplx.Abs(r) / float64(len(hues))
	sh = 1.0 - mr

	return
}
