package chromStat

import (
	"github.com/rennis250/colorCosmos/dkl"
	"github.com/rennis250/colorCosmos/lab"
	"gonum.org/v1/gonum/stat"
)

func DKLGrayWorld(dklimg *dkl.DKLImage, binary_mask []bool) (gw dkl.DKL) {
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

			lds = append(lds, dkl.LD)
			rgs = append(rgs, dkl.RG)
			yvs = append(yvs, dkl.YV)
		}
	}

	gw = dkl.DKL{stat.Mean(lds, nil), stat.Mean(rgs, nil), stat.Mean(yvs, nil)}

	return
}

func LABGrayWorld(labimg *lab.LABImage, binary_mask []bool) (gw lab.LAB) {
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

			ls = append(ls, lab.L)
			as = append(as, lab.A)
			bs = append(bs, lab.B)
		}
	}

	gw = lab.LAB{stat.Mean(ls, nil), stat.Mean(as, nil), stat.Mean(bs, nil)}

	return
}
