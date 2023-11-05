package chromStat

import (
	"math"

	"github.com/rennis250/colorCosmos/dkl"
	"github.com/runningwild/go-fftw/fftw"
	"gonum.org/v1/gonum/stat"
)

// The Gaussian Function.
// `r` is the standard deviation.
func gaussian(x, y, xc, yc, r float64) float64 {
	rsq := math.Pow(r, 2)
	rsq2 := 2.0 * rsq
	denom := 1.0 / math.Sqrt(math.Pi*rsq2)
	return denom * math.Exp(-(math.Pow(x-xc, 2)+math.Pow(y-yc, 2))/rsq2)
}

func FFTColorConstancy(dklimg *dkl.DKLImage, sigma float64, objmask []bool) (lowpassmean, highpassmean dkl.DKL) {
	bs := dklimg.Bounds()
	w, h := bs.Max.X-bs.Min.X, bs.Max.Y-bs.Min.Y
	xc, yc := float64(w/2), float64(h/2)

	ld_low := fftw.NewArray2(w, h)
	rg_low := fftw.NewArray2(w, h)
	yv_low := fftw.NewArray2(w, h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dkl := dklimg.DKLAt(x+bs.Min.X, y+bs.Min.Y)

			ld_low.Set(x, y, complex(dkl.LD, 0))
			rg_low.Set(x, y, complex(dkl.RG, 0))
			yv_low.Set(x, y, complex(dkl.YV, 0))
		}
	}

	fftw.NewPlan2(ld_low, ld_low, fftw.Forward, fftw.Estimate).Execute().Destroy()
	fftw.NewPlan2(rg_low, rg_low, fftw.Forward, fftw.Estimate).Execute().Destroy()
	fftw.NewPlan2(yv_low, yv_low, fftw.Forward, fftw.Estimate).Execute().Destroy()

	dc := [3]float64{real(ld_low.At(int(xc), int(yc))),
		real(rg_low.At(int(xc), int(yc))),
		real(yv_low.At(int(xc), int(yc)))}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			g := gaussian(float64(x), float64(y), xc, yc, sigma)
			gc := complex(g, 0)

			ldg := ld_low.At(x, y) * gc
			rgg := rg_low.At(x, y) * gc
			yvg := yv_low.At(x, y) * gc

			ld_low.Set(x, y, ldg)
			rg_low.Set(x, y, rgg)
			yv_low.Set(x, y, yvg)
		}
	}

	fftw.NewPlan2(ld_low, ld_low, fftw.Backward, fftw.Estimate).Execute().Destroy()
	fftw.NewPlan2(rg_low, rg_low, fftw.Backward, fftw.Estimate).Execute().Destroy()
	fftw.NewPlan2(yv_low, yv_low, fftw.Backward, fftw.Estimate).Execute().Destroy()

	ld_temp, rg_temp, yv_temp := make([]float64, w*h), make([]float64, w*h), make([]float64, w*h)
	for i := range ld_low.Elems {
		ld_temp[i] = real(ld_low.Elems[i])
		rg_temp[i] = real(rg_low.Elems[i])
		yv_temp[i] = real(yv_low.Elems[i])
	}

	lowpassmean = dkl.DKL{stat.Mean(ld_temp, nil),
		stat.Mean(rg_temp, nil),
		stat.Mean(yv_temp, nil)}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dkl := dklimg.DKLAt(x+bs.Min.X, y+bs.Min.Y)

			ld_temp[y*w+x] = dkl.LD - real(ld_low.At(x, y)) + dc[0]
			rg_temp[y*w+x] = dkl.RG - real(rg_low.At(x, y)) + dc[1]
			yv_temp[y*w+x] = dkl.YV - real(yv_low.At(x, y)) + dc[2]
		}
	}

	highpassmean = dkl.DKL{stat.Mean(ld_temp, nil),
		stat.Mean(rg_temp, nil),
		stat.Mean(yv_temp, nil)}

	return
}
