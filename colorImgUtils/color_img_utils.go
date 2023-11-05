package colorImgUtils

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/rennis250/colorCosmos/xyz"
)

func LoadGamma(filename string) (gR, gG, gB float64) {
	recordFile, err := os.Open(filename)
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
		gR, err = strconv.ParseFloat(rec[0], 64)
		if err != nil {
			fmt.Println("Gamma exponent should be a number ::", err)
		}
		gG, err = strconv.ParseFloat(rec[1], 64)
		if err != nil {
			fmt.Println("Gamma exponent should be a number ::", err)
		}
		gB, err = strconv.ParseFloat(rec[2], 64)
		if err != nil {
			fmt.Println("Gamma exponent should be a number ::", err)
		}
	}

	return gR, gG, gB
}

func linearInterp(y1, y2, mu float64) float64 {
	return y1*(1-mu) + y2*mu
}

// linearly interpolates so that the input spectrum "lines up" with the reference
// wavelength scale
func LoadSpectrum(spectrumName string, ref_wlns []float64) (spect_interp []float64) {
	recordFile, err := os.Open(spectrumName)
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	defer recordFile.Close()
	rdr := csv.NewReader(recordFile)
	rdr.Comma = ' '
	records, err := rdr.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	var input_wlns, input_power []float64
	for _, rec := range records {
		if rw, err := strconv.ParseFloat(rec[0], 64); err == nil {
			input_wlns = append(input_wlns, rw)
		}
		if r, err := strconv.ParseFloat(rec[1], 64); err == nil {
			input_power = append(input_power, r)
		}
	}

	for _, w := range ref_wlns {
		j := sort.SearchFloat64s(input_wlns, w)
		if j == 0 {
			spect_interp = append(spect_interp, input_power[j-1])
		} else if j == len(input_power) {
			spect_interp = append(spect_interp, 0.0)
		} else {
			mid := w - input_wlns[j-1]
			end := input_wlns[j] - input_wlns[j-1]
			mu := mid / end
			spect_interp = append(spect_interp, linearInterp(input_power[j-1], input_power[j], mu))
		}
	}

	return
}

func LoadWavelengthsAndLMS(calib_name string) (cones [][3]float64, wlns []float64) {
	recordFile, err := os.Open(calib_name)
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	defer recordFile.Close()
	rdr := csv.NewReader(recordFile)
	records, err := rdr.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(records); i += 14 {
		if wln, err := strconv.ParseFloat(records[i][0], 64); err == nil {
			wlns = append(wlns, wln)
		}

		l, err1 := strconv.ParseFloat(records[i][1], 64)
		m, err2 := strconv.ParseFloat(records[i][2], 64)
		s, err3 := strconv.ParseFloat(records[i][3], 64)
		if err1 == nil && err2 == nil && err3 == nil {
			cones = append(cones, [3]float64{l, m, s})
		} else {
			cones = append(cones, [3]float64{l, m, 0.0})
		}
	}

	return
}

func MakePieMask(w, h, xc, yc int, r_max, theta_start, theta_max float64) []bool {
	binary_mask := make([]bool, w*h)

	xcf := float64(xc)
	ycf := float64(yc)
	wf := float64(w)

	for r := 1.0; r < r_max; r++ {
		for theta := theta_start; theta < theta_max; theta += 0.001 {
			y := r*math.Sin(theta) - ycf
			x := r*math.Cos(theta) - xcf
			if math.Sqrt(x*x+y*y) < wf {
				i := int(x) + w/2
				j := int(y) + h/2
				binary_mask[j*w+i] = true
			}
		}
	}

	return binary_mask
}

func LoadBinaryMask(maskname string) (binary_mask []bool) {
	fmask, err := os.Open(maskname)
	if err != nil {
		panic("could not open mask file")
	}
	defer fmask.Close()
	mask, err := png.Decode(fmask)
	if err != nil {
		log.Fatal(err)
	}
	b := mask.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			mr, mg, mb, _ := mask.At(x, y).RGBA()
			if mr != 0 && mg != 0 && mb != 0 {
				binary_mask = append(binary_mask, true)
			} else {
				binary_mask = append(binary_mask, false)
			}
		}
	}
	return
}

func ExcludeDarkPixelsFromMask(binary_mask []bool, imgxyz *xyz.XYZImage, darkThresh float64) {
	b := imgxyz.Bounds()
	width := b.Max.X - b.Min.X
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			XYZ := imgxyz.XYZAt(x, y)
			if XYZ.Y < darkThresh {
				idx := (y-b.Min.Y)*width + (x - b.Min.X)
				binary_mask[idx] = false
			}
		}
	}
}

func GammaCorr(in image.Image, gR, gG, gB float64) *image.RGBA64 {
	bs := in.Bounds()
	out := image.NewRGBA64(bs)

	for y := bs.Min.Y; y < bs.Max.Y; y++ {
		for x := bs.Min.X; x < bs.Max.X; x++ {
			r, g, b, _ := in.At(x, y).RGBA()
			r_s := float64(r) / 65535
			g_s := float64(g) / 65535
			b_s := float64(b) / 65535

			r_lin := uint16(math.Pow(r_s, gR) * 65535)
			g_lin := uint16(math.Pow(g_s, gG) * 65535)
			b_lin := uint16(math.Pow(b_s, gB) * 65535)
			out.Set(x, y, color.RGBA64{r_lin, g_lin, b_lin, 65535})
		}
	}

	return out
}
