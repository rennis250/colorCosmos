package chromStat

import (
	"math"

	"github.com/rennis250/colorCosmos/lms"
	"gonum.org/v1/gonum/stat"
)

func MeanSdCones(lmsimg *lms.LMSImage, binary_mask []bool) (lmsmean, lmssd lms.LMS) {
	var ls []float64
	var ms []float64
	var ss []float64

	b := lmsimg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			lms := lmsimg.LMSAt(x, y)

			if !binary_mask[y*b.Max.X+x] {
				continue
			}

			ls = append(ls, lms.L)
			ms = append(ms, lms.M)
			ss = append(ss, lms.S)
		}
	}

	lmsmean = lms.LMS{stat.Mean(ls, nil),
		stat.Mean(ms, nil),
		stat.Mean(ss, nil)}

	lmssd = lms.LMS{stat.StdDev(ls, nil),
		stat.StdDev(ms, nil),
		stat.StdDev(ss, nil)}

	return
}

func RatioMeanSdCones(lmsimg *lms.LMSImage, binary_mask []bool) (rmc, rsd [3]float64) {
	var ls_in_mask, ls_out_mask []float64
	var ms_in_mask, ms_out_mask []float64
	var ss_in_mask, ss_out_mask []float64

	b := lmsimg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			lms := lmsimg.LMSAt(x, y)

			if binary_mask[y*b.Max.X+x] {
				ls_in_mask = append(ls_in_mask, lms.L)
				ms_in_mask = append(ms_in_mask, lms.M)
				ss_in_mask = append(ss_in_mask, lms.S)
			} else {
				ls_out_mask = append(ls_out_mask, lms.L)
				ms_out_mask = append(ms_out_mask, lms.M)
				ss_out_mask = append(ss_out_mask, lms.S)
			}
		}
	}

	rmc[0] = stat.Mean(ls_in_mask, nil) / stat.Mean(ls_out_mask, nil)
	rmc[1] = stat.Mean(ms_in_mask, nil) / stat.Mean(ms_out_mask, nil)
	rmc[2] = stat.Mean(ss_in_mask, nil) / stat.Mean(ss_out_mask, nil)

	rsd[0] = stat.StdDev(ls_in_mask, nil) / stat.StdDev(ls_out_mask, nil)
	rsd[1] = stat.StdDev(ms_in_mask, nil) / stat.StdDev(ms_out_mask, nil)
	rsd[2] = stat.StdDev(ss_in_mask, nil) / stat.StdDev(ss_out_mask, nil)

	return
}

// tau is vector of transmittance factors
func GeneralRobustRatio(lmsimg *lms.LMSImage, binary_mask []bool) (tau [3]float64) {
	var ls_in_mask, ls_out_mask []float64
	var ms_in_mask, ms_out_mask []float64
	var ss_in_mask, ss_out_mask []float64

	b := lmsimg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			lms := lmsimg.LMSAt(x, y)

			if binary_mask[y*b.Max.X+x] {
				ls_in_mask = append(ls_in_mask, lms.L)
				ms_in_mask = append(ms_in_mask, lms.M)
				ss_in_mask = append(ss_in_mask, lms.S)
			} else {
				ls_out_mask = append(ls_out_mask, lms.L)
				ms_out_mask = append(ms_out_mask, lms.M)
				ss_out_mask = append(ss_out_mask, lms.S)
			}
		}
	}

	// p is the color of filtered background -> lms_in
	// a is colors of background -> lms_out
	// I is an estimate of the illumination
	I := [3]float64{stat.Mean(ls_out_mask, nil),
		stat.Mean(ms_out_mask, nil),
		stat.Mean(ss_out_mask, nil)}

	// the std of background is used to determine channel which
	// is best for starting the estimate of tau
	std_a := [3]float64{stat.StdDev(ls_out_mask, nil),
		stat.StdDev(ms_out_mask, nil),
		stat.StdDev(ss_out_mask, nil)}

	std_p := [3]float64{stat.StdDev(ls_in_mask, nil),
		stat.StdDev(ms_in_mask, nil),
		stat.StdDev(ss_in_mask, nil)}

	mean_a := [3]float64{stat.Mean(ls_out_mask, nil),
		stat.Mean(ms_out_mask, nil),
		stat.Mean(ss_out_mask, nil)}

	mean_p := [3]float64{stat.Mean(ls_in_mask, nil),
		stat.Mean(ms_in_mask, nil),
		stat.Mean(ss_in_mask, nil)}

	max_std_a := math.Inf(-1)
	maxsdidx := 0
	for i, v := range std_a {
		if v > max_std_a {
			max_std_a = v
			maxsdidx = i
		}
	}

	tau[maxsdidx] = std_p[maxsdidx] / std_a[maxsdidx]

	mu := mean_p[maxsdidx] - tau[maxsdidx]*mean_a[maxsdidx]
	v := (tau[maxsdidx] + mu) * I[maxsdidx]

	// delta is amount of direct reflection
	delta := mu / v

	// once delta is known, we compute other taus from:
	for i := 0; i < len(tau); i++ {
		tau[i] = (mean_p[i] - mu*delta*I[i]) / (mean_a[i] + delta*I[i])
	}

	return
}

// tau is vector of transmittance factors
func SimplerRobustRatio(lmsimg *lms.LMSImage, binary_mask []bool) (tau [3]float64) {
	var ls_in_mask, ls_out_mask []float64
	var ms_in_mask, ms_out_mask []float64
	var ss_in_mask, ss_out_mask []float64

	b := lmsimg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			lms := lmsimg.LMSAt(x, y)

			if binary_mask[y*b.Max.X+x] {
				ls_in_mask = append(ls_in_mask, lms.L)
				ms_in_mask = append(ms_in_mask, lms.M)
				ss_in_mask = append(ss_in_mask, lms.S)
			} else {
				ls_out_mask = append(ls_out_mask, lms.L)
				ms_out_mask = append(ms_out_mask, lms.M)
				ss_out_mask = append(ss_out_mask, lms.S)
			}
		}
	}

	// p is the color of filtered background -> lms_in
	// a is colors of background -> lms_out
	// I is an estimate of the illumination -> not directly used in this form in simpler version
	// I := [3]float64{stat.Mean(ls_out_mask, nil),
	//	stat.Mean(ms_out_mask, nil),
	//	stat.Mean(ss_out_mask, nil)}

	// the std of background is used to determine channel which
	// is best for starting the estimate of tau
	std_a := [3]float64{stat.StdDev(ls_out_mask, nil),
		stat.StdDev(ms_out_mask, nil),
		stat.StdDev(ss_out_mask, nil)}

	std_p := [3]float64{stat.StdDev(ls_in_mask, nil),
		stat.StdDev(ms_in_mask, nil),
		stat.StdDev(ss_in_mask, nil)}

	mean_a := [3]float64{stat.Mean(ls_out_mask, nil),
		stat.Mean(ms_out_mask, nil),
		stat.Mean(ss_out_mask, nil)}

	mean_p := [3]float64{stat.Mean(ls_in_mask, nil),
		stat.Mean(ms_in_mask, nil),
		stat.Mean(ss_in_mask, nil)}

	max_std_a := math.Inf(-1)
	maxsdidx := 0
	for i, v := range std_a {
		if v > max_std_a {
			max_std_a = v
			maxsdidx = i
		}
	}

	tau[maxsdidx] = std_p[maxsdidx] / std_a[maxsdidx]

	// R is RMC of filter
	R := [3]float64{mean_p[0] / mean_a[0], mean_p[1] / mean_a[1], mean_p[2] / mean_a[2]}

	// gamma is clarity of filter
	gamma := tau[maxsdidx] / R[maxsdidx]

	for i := 0; i < len(tau); i++ {
		tau[i] = gamma * R[i]
	}

	return
}
