package xyy

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

var XyY2RGB *mat.Dense
var RGB2XyY *mat.Dense

type XyY struct {
	x, y, Y float64
}

func RGBToXyYFromChroma(calib_file string) {
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
	rY := Ys[0]

	gx := xs[1]
	gy := ys[1]
	gY := Ys[1]

	bx := xs[2]
	by := ys[2]
	bY := Ys[2]

	L := []float64{R[0], G[0], B[0]}
	M := []float64{R[1], G[1], B[1]}
	S := []float64{R[2], G[2], B[2]}

	XyY2RGB = mat.NewDense(3, 3, []float64{1.0, 1.0, d_ryv, 1.0, d_grg, d_gyv, 1.0, d_brg, 1.0})

	var RGB2XyY_temp mat.Dense
	RGB2XyY_temp.Inverse(XyY2YB)

	RGB2XyY = &RGB2XyY_temp
}

type XyYModel struct{}

func (m XyYModel) Convert(c color.Color) color.Color {
	if _, ok := c.(XyY); ok {
		return c
	}
	R, G, B, _ := c.RGBA()
	r := float64(R) / 65535.0
	g := float64(G) / 65535.0
	b := float64(B) / 65535.0

	if r > 1.0 || r < 0.0 || g > 1.0 || g < 0.0 || b > 1.0 || b < 0.0 {
		log.Fatalln("XyY - OOG: ", r, g, b)
	}

	var xyy mat.VecDense
	xyy.MulVec(RGB2XyY, mat.NewVecDense(3, []float64{r, g, b}))
	x := xyy.AtVec(0)
	y := xyy.AtVec(1)
	Y := xyy.AtVec(2)

	return XyY{x, y, Y}
}

func (x XyY) RGBA() (r, g, b, a uint32) {
	var rgb mat.VecDense
	rgb.MulVec(XyY2RGB, mat.NewVecDense(3, []float64{x.x, x.y, x.Y}))

	r = uint32(65535.0 * rgb.AtVec(0))
	g = uint32(65535.0 * rgb.AtVec(1))
	b = uint32(65535.0 * rgb.AtVec(2))
	a = uint32(65535)

	return
}
