package dkl

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"github.com/rennis250/colorCosmos/lms"
	"gonum.org/v1/gonum/mat"
)

var DKL2RGB *mat.Dense
var RGB2DKL *mat.Dense

type DKL struct {
	LD, RG, YV float64
}

func solvex(a, b, c, d, e, f float64) float64 {
	return (a*f/d - b) / (c*f/d - e)
}

func solvey(a, b, c, d, e, f float64) float64 {
	return (a*e/c - b) / (d*e/c - f)
}

func cie2lms(x, y float64) [3]float64 {
	cie := mat.NewVecDense(3, []float64{x, y, 1.0 - x - y})
	matrix := mat.NewDense(3, 3,
		[]float64{0.15514, 0.54316, -0.03286, -0.15514, 0.45684, 0.03286, 0.0, 0.0, 0.01608})

	var lms mat.VecDense
	lms.MulVec(matrix, cie)
	l := lms.AtVec(0) / (lms.AtVec(0) + lms.AtVec(1)) // L/(L+M)
	m := lms.AtVec(1) / (lms.AtVec(0) + lms.AtVec(1)) // M/(L+M)
	s := lms.AtVec(2) / (lms.AtVec(0) + lms.AtVec(1)) // S/(L+M)

	return [3]float64{l, m, s}
}

func lumchrm(lumr, lumg, lumb float64, r, g, b [3]float64) {
	lum := []float64{lumr, lumg, lumb}
	bigl := 0.0
	bigm := 0.0
	bigs := 0.0

	for i := 0; i < 3; i++ {
		fmt.Println("Luminances[", i, "] = ", lum[i])
		bigl += lum[i] * r[i]
		bigm += lum[i] * g[i]
		bigs += lum[i] * b[i]
	}

	denom := bigl + bigm
	fmt.Println("bigl = ", bigl, " bigm = ", bigm, " bigs = ", bigs)
	lout := bigl / denom
	sout := bigs / denom
	fmt.Println("L/L+M = ", lout, " S/L+M = ", sout)
}

func DKLToRGBFromChroma(calib_file string) {
	var xs []float64
	var ys []float64
	var Ys []float64

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
		if x, err := strconv.ParseFloat(rec[0], 64); err == nil {
			xs = append(xs, x)
		}
		if y, err := strconv.ParseFloat(rec[1], 64); err == nil {
			ys = append(ys, y)
		}
		if Y, err := strconv.ParseFloat(rec[2], 64); err == nil {
			Ys = append(Ys, Y)
		}
	}

	rx := xs[0]
	ry := ys[0]
	rY := Ys[0] / 2.0

	gx := xs[1]
	gy := ys[1]
	gY := Ys[1] / 2.0

	bx := xs[2]
	by := ys[2]
	bY := Ys[2] / 2.0

	white := []float64{rY, gY, bY}

	R := cie2lms(rx, ry)
	G := cie2lms(gx, gy)
	B := cie2lms(bx, by)

	L := []float64{R[0], G[0], B[0]}
	M := []float64{R[1], G[1], B[1]}
	S := []float64{R[2], G[2], B[2]}

	// lumchrm(white[0], white[1], white[2], L, M, S)

	// fmt.Println("Red Green Axis");
	delta_grg := solvex(white[0]*S[0],
		white[0]*(L[0]+M[0]),
		S[1],
		S[2],
		L[1]+M[1],
		L[2]+M[2])
	delta_brg := solvey(white[0]*S[0],
		white[0]*(L[0]+M[0]),
		S[1],
		S[2],
		L[1]+M[1],
		L[2]+M[2])
	d_grg := -1.0 * delta_grg / white[1]
	d_brg := -1.0 * delta_brg / white[2]
	// lumchrm(0.0, white[1]+delta_grg, white[2]+delta_brg, L, M, S)
	// lumchrm(white[0]*2.0, white[1]-delta_grg, white[2]-delta_brg, L, M, S)
	// fmt.Println("Blue Yellow Axis")
	delta_ryv := solvex(white[2]*L[2], white[2]*M[2], L[0], L[1], M[0], M[1])
	delta_gyv := solvey(white[2]*L[2], white[2]*M[2], L[0], L[1], M[0], M[1])
	d_ryv := -1.0 * delta_ryv / white[0]
	d_gyv := -1.0 * delta_gyv / white[1]
	// lumchrm(white[0]+delta_ryv, white[1]+delta_gyv, 0.0, L, M, S)
	// lumchrm(white[0]-delta_ryv, white[1]-delta_gyv, white[2]*2.0, L, M, S)

	DKL2RGB = mat.NewDense(3, 3, []float64{1.0, 1.0, d_ryv, 1.0, d_grg, d_gyv, 1.0, d_brg, 1.0})

	var RGB2DKL_temp mat.Dense
	RGB2DKL_temp.Inverse(DKL2RGB)

	RGB2DKL = &RGB2DKL_temp
}

type DKLModel struct{}

func (m DKLModel) Convert(c color.Color) color.Color {
	if _, ok := c.(DKL); ok {
		return c
	}
	R, G, B, _ := c.RGBA()
	r := float64(R) / 65535.0
	g := float64(G) / 65535.0
	b := float64(B) / 65535.0

	if r > 1.0 || r < 0.0 || g > 1.0 || g < 0.0 || b > 1.0 || b < 0.0 {
		log.Fatalln("DKL - OOG: ", r, g, b)
	}

	r_s := 2.0 * (r - 0.5)
	g_s := 2.0 * (g - 0.5)
	b_s := 2.0 * (b - 0.5)

	var dkl mat.VecDense
	dkl.MulVec(RGB2DKL, mat.NewVecDense(3, []float64{r_s, g_s, b_s}))
	LD := dkl.AtVec(0)
	RG := dkl.AtVec(1)
	YV := dkl.AtVec(2)

	return DKL{LD, RG, YV}
}

func (d DKL) RGBA() (r, g, b, a uint32) {
	var rgb mat.VecDense
	rgb.MulVec(DKL2RGB, mat.NewVecDense(3, []float64{d.LD, d.RG, d.YV}))

	r = uint32(65535.0 * (rgb.AtVec(0)/2.0 + 0.5))
	g = uint32(65535.0 * (rgb.AtVec(1)/2.0 + 0.5))
	b = uint32(65535.0 * (rgb.AtVec(2)/2.0 + 0.5))
	a = uint32(65535)

	return
}

func FromLMS(lms lms.LMS) DKL {
	l, m, s := lms.L, lms.M, lms.S
	ld, rg, yv := l+m, l-m, 2.0*s-l
	return DKL{ld, rg, yv}
}
