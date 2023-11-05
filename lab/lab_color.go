package lab

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

var RGB2XYZ *mat.Dense
var XYZ2RGB *mat.Dense
var illumXYZ [3]float64

type LAB struct {
	L, A, B float64
}

func rgbToXYZ(rgb [3]float64) [3]float64 {
	var xyz mat.VecDense
	xyz.MulVec(RGB2XYZ, mat.NewVecDense(3, rgb[:]))
	return [3]float64{xyz.AtVec(0), xyz.AtVec(1), xyz.AtVec(2)}
}

func xyYToXYZ(x, y, Y float64) [3]float64 {
	X := (Y / y) * x
	Z := (Y / y) * (1.0 - y - x)
	return [3]float64{X, Y, Z}
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

	rXYZ := xyYToXYZ(xs[0], ys[0], Ys[0])
	gXYZ := xyYToXYZ(xs[1], ys[1], Ys[1])
	bXYZ := xyYToXYZ(xs[2], ys[2], Ys[2])

	RGB2XYZ = mat.NewDense(3, 3, []float64{rXYZ[0], gXYZ[0], bXYZ[0],
		rXYZ[1], gXYZ[1], bXYZ[1],
		rXYZ[2], gXYZ[2], bXYZ[2]})

	var XYZ2RGB_temp mat.Dense
	XYZ2RGB_temp.Inverse(RGB2XYZ)
	XYZ2RGB = &XYZ2RGB_temp

	illumXYZ = [3]float64{rXYZ[0] + gXYZ[0] + bXYZ[0],
		rXYZ[1] + gXYZ[1] + bXYZ[1],
		rXYZ[2] + gXYZ[2] + bXYZ[2]}

	return
}

type LABModel struct{}

func (m LABModel) Convert(c color.Color) color.Color {
	if _, ok := c.(LAB); ok {
		return c
	}
	R, G, B, _ := c.RGBA()
	r := float64(R) / 65535.0
	g := float64(G) / 65535.0
	b := float64(B) / 65535.0

	if r > 1.0 || r < 0.0 || g > 1.0 || g < 0.0 || b > 1.0 || b < 0.0 {
		log.Fatalln("LAB - OOG: ", r, g, b)
	}

	xyz := rgbToXYZ([3]float64{r, g, b})

	labnx := xyz[0] / illumXYZ[0]
	labny := xyz[1] / illumXYZ[1]
	labnz := xyz[2] / illumXYZ[2]

	C := math.Pow(6.0/29.0, 3.0)
	linmt := 3.0 * math.Pow(6.0/29.0, 2.0)
	linat := 4.0 / 29.0

	if labnx <= C {
		labnx = labnx/linmt + linat
	} else {
		labnx = math.Pow(labnx, 1.0/3.0)
	}

	if labny <= C {
		labny = labny/linmt + linat
	} else {
		labny = math.Pow(labny, 1.0/3.0)
	}

	if labnz <= C {
		labnz = labnz/linmt + linat
	} else {
		labnz = math.Pow(labnz, 1.0/3.0)
	}

	labl := 116.0*labny - 16.0
	laba := 500.0 * (labnx - labny)
	labb := 200.0 * (labny - labnz)

	return LAB{labl, laba, labb}
}

func xyzToRGB(xyz [3]float64) [3]float64 {
	var rgb mat.VecDense
	rgb.MulVec(XYZ2RGB, mat.NewVecDense(3, xyz[:]))
	return [3]float64{rgb.AtVec(0), rgb.AtVec(1), rgb.AtVec(2)}
}

func (l LAB) RGBA() (r, g, b, a uint32) {
	C := 6.0 / 29.0
	linmt := 3.0 * math.Pow(6.0/29.0, 2.0)
	linat := 4.0 / 29.0

	fy := (l.L + 16.0) / 116.0
	fx := (l.A / 500.0) + fy
	fz := fy - (l.B / 200.0)

	var xn float64
	if fx > C {
		xn = math.Pow(fx, 3)
	} else {
		xn = (fx - linat) * linmt
	}

	var yn float64
	if fy > C {
		yn = math.Pow(fy, 3)
	} else {
		yn = (fy - linat) * linmt
	}

	var zn float64
	if fz > C {
		zn = math.Pow(fz, 3)
	} else {
		zn = (fz - linat) * linmt
	}

	X := xn * illumXYZ[0]
	Y := yn * illumXYZ[1]
	Z := zn * illumXYZ[2]

	RGB := xyzToRGB([3]float64{X, Y, Z})
	r = uint32(RGB[0] * 65535.0)
	g = uint32(RGB[1] * 65535.0)
	b = uint32(RGB[2] * 65535.0)
	a = uint32(65535)

	return
}
