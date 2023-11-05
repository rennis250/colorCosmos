package xyz

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

var XYZ2RGB *mat.Dense
var RGB2XYZ *mat.Dense

type XYZ struct {
	X, Y, Z float64
}

func xyYtoXYZ(x, y, Yin float64) (X, Yout, Z float64) {
	X = (Yin / y) * x
	Yout = Yin
	Z = (Yin / y) * (1 - y - x)
	return
}

func RGBToXYZFromChroma(calib_file string) {
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

	rX, rY, rZ := xyYtoXYZ(rx, ry, rY)
	gX, gY, gZ := xyYtoXYZ(gx, gy, gY)
	bX, bY, bZ := xyYtoXYZ(bx, by, bY)

	RGB2XYZ = mat.NewDense(3, 3, []float64{rX, gX, bX, rY, gY, bY, rZ, gZ, bZ})

	var XYZ2RGB_temp mat.Dense
	XYZ2RGB_temp.Inverse(RGB2XYZ)

	XYZ2RGB = &XYZ2RGB_temp
}

type XYZModel struct{}

func (m XYZModel) Convert(c color.Color) color.Color {
	if _, ok := c.(XYZ); ok {
		return c
	}
	R, G, B, _ := c.RGBA()
	r := float64(R) / 65535.0
	g := float64(G) / 65535.0
	b := float64(B) / 65535.0

	if r > 1.0 || r < 0.0 || g > 1.0 || g < 0.0 || b > 1.0 || b < 0.0 {
		log.Fatalln("XYZ - OOG: ", r, g, b)
	}

	var xyz mat.VecDense
	xyz.MulVec(RGB2XYZ, mat.NewVecDense(3, []float64{r, g, b}))
	X := xyz.AtVec(0)
	Y := xyz.AtVec(1)
	Z := xyz.AtVec(2)

	return XYZ{X, Y, Z}
}

func (x XYZ) RGBA() (r, g, b, a uint32) {
	var rgb mat.VecDense
	rgb.MulVec(XYZ2RGB, mat.NewVecDense(3, []float64{x.X, x.Y, x.Z}))

	r = uint32(65535.0 * rgb.AtVec(0))
	g = uint32(65535.0 * rgb.AtVec(1))
	b = uint32(65535.0 * rgb.AtVec(2))
	a = uint32(65535)

	return
}
