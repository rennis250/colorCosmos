package lms

import (
	"encoding/csv"
	"fmt"
	"go/build"
	"image/color"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

var LMS2RGB *mat.Dense
var RGB2LMS *mat.Dense

type LMS struct {
	L, M, S float64
}

func RGBToLMSFromSpectra(calib_file string) {
	var rs []float64
	var gs []float64
	var bs []float64

	recordFile, err := os.Open(calib_file)
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	defer recordFile.Close()
	rdr := csv.NewReader(recordFile)
	records, err := rdr.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	for _, rec := range records {
		if r, err := strconv.ParseFloat(rec[0], 64); err == nil {
			rs = append(rs, r) // TODO: why multiply by 10? rounding issues?
		}
		if g, err := strconv.ParseFloat(rec[1], 64); err == nil {
			gs = append(gs, g)
		}
		if b, err := strconv.ParseFloat(rec[2], 64); err == nil {
			bs = append(bs, b)
		}
	}

	var wlns []float64
	var ls []float64
	var ms []float64
	var ss []float64

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	recordFile, err = os.Open(gopath + "/src/github.com/rennis250/lms/data/linss2_10e_1.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	defer recordFile.Close()
	rdr = csv.NewReader(recordFile)
	records, err = rdr.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	for _, rec := range records {
		if wln, err := strconv.ParseFloat(rec[0], 64); err == nil {
			wlns = append(wlns, wln)
		}
		if l, err := strconv.ParseFloat(rec[1], 64); err == nil {
			ls = append(ls, l)
		}
		if m, err := strconv.ParseFloat(rec[2], 64); err == nil {
			ms = append(ms, m)
		}

		// s is an annoying special case, because the CSV file just
		// has empty values instead of zeros where the curves could
		// not be measured
		s, err := strconv.ParseFloat(rec[3], 64)
		if err == nil {
			ss = append(ss, s)
		} else {
			ss = append(ss, 0.0)
		}
	}

	fmt.Println(rs)
	fmt.Println(gs)
	fmt.Println(bs)

	var r_l, r_m, r_s float64
	var g_l, g_m, g_s float64
	var b_l, b_m, b_s float64

	for x := 10; x < len(rs); x++ {
		r_l += rs[x] * ls[x-10]
		r_m += rs[x] * ms[x-10]
		r_s += rs[x] * ss[x-10]

		g_l += gs[x] * ls[x-10]
		g_m += gs[x] * ms[x-10]
		g_s += gs[x] * ss[x-10]

		b_l += bs[x] * ls[x-10]
		b_m += bs[x] * ms[x-10]
		b_s += bs[x] * ss[x-10]
	}

	RGB2LMS = mat.NewDense(3, 3, []float64{r_l, g_l, b_l,
		r_m, g_m, b_m,
		r_s, g_s, b_s})

	var lms2rgb_temp mat.Dense
	lms2rgb_temp.Inverse(RGB2LMS)

	LMS2RGB = &lms2rgb_temp
}

type LMSModel struct{}

func (m LMSModel) Convert(c color.Color) color.Color {
	if _, ok := c.(LMS); ok {
		return c
	}
	R, G, B, _ := c.RGBA()
	r := float64(R) / 65535.0
	g := float64(G) / 65535.0
	b := float64(B) / 65535.0

	if r > 1.0 || r < 0.0 || g > 1.0 || g < 0.0 || b > 1.0 || b < 0.0 {
		log.Fatalln("LMS - OOG: ", r, g, b)
	}

	var lms mat.VecDense
	lms.MulVec(RGB2LMS, mat.NewVecDense(3, []float64{r, g, b}))
	L := lms.AtVec(0)
	M := lms.AtVec(1)
	S := lms.AtVec(2)

	return LMS{L, M, S}
}

func (l LMS) RGBA() (r, g, b, a uint32) {
	var rgb mat.VecDense
	rgb.MulVec(LMS2RGB, mat.NewVecDense(3, []float64{l.L, l.M, l.S}))

	r = uint32(65535.0 * rgb.AtVec(0))
	g = uint32(65535.0 * rgb.AtVec(1))
	b = uint32(65535.0 * rgb.AtVec(2))
	a = uint32(65535)

	return
}
