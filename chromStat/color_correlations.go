package chromStat

import (
	"math"

	"github.com/rennis250/colorCosmos/lms"
	"gonum.org/v1/gonum/stat"
)

func LRC(lmsimg *lms.LMSImage, binary_mask []bool) (cr, y float64) {
	var lum []float64
	var r []float64

	b := lmsimg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			lms := lmsimg.LMSAt(x, y)

			if !binary_mask[y*b.Max.X+x] {
				continue
			}

			if lms.L > 0 && lms.M > 0 && lms.S > 0 {
				lum = append(lum, math.Log10(lms.L+lms.M))
				r = append(r, math.Log10(lms.L/(lms.L+lms.M)))
			}
		}
	}

	cr = stat.Correlation(lum, r, nil)
	y = 1.1305*stat.Mean(r, nil) + 0.0063*cr + 0.0202

	return
}

// TODO: need to get formulas for RBC and LBC out of Golz thesis...
