package chromStat

import (
	"math"
	"math/rand"
	"time"

	"github.com/rennis250/colorCosmos/lms"
	"gonum.org/v1/gonum/stat"
)

func pxdist(x1, x2 [2]int) int {
	return int(math.Abs(float64(x2[0])-float64(x1[0])) + math.Abs(float64(x2[1])-float64(x1[1])))
}

func vecNorm(a [3]float64) float64 {
	sum := 0.0
	for _, v := range a {
		sum += math.Pow(v, 2)
	}
	return math.Sqrt(sum)
}

func sliceContainsTwoInts(a [][2]int, b [2]int) bool {
	for i := range a {
		if a[i][0] == b[0] && a[i][1] == b[1] {
			return true
		}
	}
	return false
}

func sliceContainsInt(a []int, b int) bool {
	for i := range a {
		if a[i] == b {
			return true
		}
	}
	return false
}

func ConeExcitationRatios(lmsimg1, lmsimg2 lms.LMSImage, binary_mask []bool) (mcer float64) {
	var gi1x []int
	var gi1y []int
	var gi2x []int
	var gi2y []int

	b := lmsimg1.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			lms1 := lmsimg1.LMSAt(x, y)
			lms2 := lmsimg2.LMSAt(x, y)

			if !binary_mask[y*b.Max.X+x] {
				continue
			}

			if lms1.L > 0 && lms1.M > 0 && lms1.S > 0 {
				gi1x = append(gi1x, x)
				gi1y = append(gi1y, y)
			}

			if lms2.L > 0 && lms2.M > 0 && lms2.S > 0 {
				gi2x = append(gi2x, x)
				gi2y = append(gi2y, y)
			}
		}
	}

	var gix []int
	var giy []int
	if len(gi1x) < len(gi2x) {
		gix = gi1x
		giy = gi1y
	} else {
		gix = gi2x
		giy = gi2y
	}

	var r int
	// 10000 comes from one of sergio's papers
	lgix := len(gix)
	if lgix/2 < 10000 {
		if lgix%2 == 0 {
			r = lgix / 2
		} else {
			r = (lgix - 1) / 2
		}
	} else {
		r = 10000
	}

	the_size := len(gix)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	allowedDists := []int{1, 2, 4, 8, 16, 32, 64, 128, 256}
	var grc [][2]int
	var alreadyTested [][2]int
Outer:
	for {
		var rci1, rci2 []int
	Inner:
		for {
			rci1 := rng.Perm(the_size)
			rci2 := rng.Perm(the_size)

			for i := range rci1 {
				if rci1[i] == rci2[i] {
					continue Inner
				}
			}
		}

		rcis := make([][2]int, len(rci1))
		for i := range rci1 {
			rcis[i][0] = rci1[i]
			rcis[i][1] = rci2[i]
		}

		for _, pts := range rcis {
			if sliceContainsTwoInts(alreadyTested, pts) {
				continue
			}
			x := pts[0]
			y := pts[1]
			pd := pxdist([2]int{gix[x], giy[x]}, [2]int{gix[y], giy[y]})
			if sliceContainsInt(allowedDists, pd) {
				grc = append(grc, pts)
				if len(grc) >= r {
					break Outer
				}
			}
		}

		for _, v := range rcis {
			if sliceContainsTwoInts(grc, v) {
				alreadyTested = append(alreadyTested, v)
			}
		}
	}

	var cer []float64
	for _, pts := range grc[:r] {
		lms1 := lmsimg1.LMSAt(gix[pts[0]], giy[pts[0]])
		lms2 := lmsimg1.LMSAt(gix[pts[1]], giy[pts[1]])

		r1 := [3]float64{float64(lms1.L) / float64(lms2.L),
			float64(lms1.M) / float64(lms2.M),
			float64(lms1.S) / float64(lms2.S)}

		lms1 = lmsimg2.LMSAt(gix[pts[0]], giy[pts[0]])
		lms2 = lmsimg2.LMSAt(gix[pts[1]], giy[pts[1]])

		r2 := [3]float64{float64(lms1.L) / float64(lms2.L),
			float64(lms1.M) / float64(lms2.M),
			float64(lms1.S) / float64(lms2.S)}

		r1norm := vecNorm(r1)
		r2norm := vecNorm(r2)
		mr := (r1norm + r2norm) / 2.0

		var r1mr2 [3]float64
		for i := range r1 {
			r1mr2[i] = r1[i] - r2[i]
		}
		tmp := vecNorm(r1mr2) / mr
		if !math.IsInf(tmp, 0) && !math.IsNaN(tmp) {
			cer = append(cer, tmp)
		}
	}

	mcer = stat.Mean(cer, nil)

	return mcer
}
