package chromStat

import (
	"math"

	"github.com/rennis250/colorCosmos/dkl"
	"github.com/rennis250/colorCosmos/lab"
)

func DKLWhitePoint(dklimg *dkl.DKLImage, binary_mask []bool) (wp dkl.DKL) {
	wpld := -(math.MaxFloat64 - 1) // just to be safe and avoid overflow

	var lds []float64
	var rgs []float64
	var yvs []float64

	b := dklimg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			dkl := dklimg.DKLAt(x, y)

			if !binary_mask[y*b.Max.X+x] {
				continue
			}

			wpld = math.Max(wpld, dkl.LD)

			lds = append(lds, dkl.LD)
			rgs = append(rgs, dkl.RG)
			yvs = append(yvs, dkl.YV)
		}
	}

	wpidx := 0
	for x, v := range lds {
		if v == wpld {
			wpidx = x
		}
	}

	wp = dkl.DKL{wpld, rgs[wpidx], yvs[wpidx]}

	return
}

func LABWhitePoint(labimg *lab.LABImage, binary_mask []bool) (wp lab.LAB) {
	wpl := -(math.MaxFloat64 - 1) // just to be safe and avoid overflow

	var ls []float64
	var as []float64
	var bs []float64

	b := labimg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			lab := labimg.LABAt(x, y)

			if !binary_mask[y*b.Max.X+x] {
				continue
			}

			wpl = math.Max(wpl, lab.L)

			ls = append(ls, lab.L)
			as = append(as, lab.A)
			bs = append(bs, lab.B)
		}
	}

	wpidx := 0
	for x, v := range ls {
		if v == wpl {
			wpidx = x
		}
	}

	wp = lab.LAB{wpl, as[wpidx], bs[wpidx]}

	return
}
