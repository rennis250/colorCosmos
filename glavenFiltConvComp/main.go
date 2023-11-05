package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"go/build"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/rennis250/colorCosmos/dkl"
	"github.com/rennis250/colorCosmos/lms"
	"gonum.org/v1/gonum/mat"
)

func floorFract(x float64) (float64, float64) {
	floor, frac := math.Modf(x)
	if frac < 0 {
		return floor, -frac
	} else {
		return floor, frac
	}
}

func hash22(p, pp float64) (ox, oy float64) {
	p1 := [3]float64{p, pp, p}
	p2 := [3]float64{0.1031, 0.1030, 0.0973}

	var p3 [3]float64
	for i := range p1 {
		_, p3[i] = floorFract(p1[i] * p2[i])
	}

	p4 := [3]float64{p3[1] + 19.19, p3[2] + 19.19, p3[0] + 19.19}
	s := 0.0
	for i := range p3 {
		s += p3[i] * p4[i]
	}
	p3[0] += s
	p3[1] += s
	p3[2] += s

	_, ox = floorFract((p3[0] + p3[1]) * p3[2])
	_, oy = floorFract((p3[0] + p3[2]) * p3[1])
	return
}

func voronoi(x, y float64) float64 {
	n1, f1 := floorFract(x)
	n2, f2 := floorFract(y)

	m := [3]float64{8.0, 8.0, 8.0}
	for j := -1.0; j < 1.0; j++ {
		for i := -1.0; i < 1.0; i++ {
			ox, oy := hash22(n1+i, n2+j)
			r1 := i - f1 + (0.5 + 0.5*math.Sin(2*math.Pi*ox))
			r2 := j - f2 + (0.5 + 0.5*math.Sin(2*math.Pi*oy))
			d := r1*r1 + r2*r2
			if d < m[0] {
				m[0] = d
				m[1] = ox
				m[2] = oy
			}
		}
	}

	return m[1] + m[2]
}

func loadVoronoiTexture(filename string) (vortex []float64) {
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
		if v, err := strconv.ParseFloat(rec[0], 64); err == nil {
			vortex = append(vortex, v)
		}
	}

	return
}

func main() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	spectraFile := flag.String("spectraFile", gopath+"/src/github.com/rennis250/voronoi_filters/data/mon_spectra.csv", "CSV file with spectral curves of each primary on monitor")
	chromaFile := flag.String("chromaFile", gopath+"/src/github.com/rennis250/voronoi_filters/data/mon_chroma.csv", "CSV file with chromaticity coordinates of each primary on monitor")

	red_refl_filename := flag.String("red", gopath+"/src/github.com/rennis250/voronoi_filters/data/munsell_red_EXTREEEMMMEE.spd", "SSV file with reflectance at different wavelengths for \"red\" filter distribution")

	green_refl_filename := flag.String("green", gopath+"/src/github.com/rennis250/voronoi_filters/data/munsell_green_EXTREEEMMMEE.spd", "SSV file with reflectance at different wavelengths for \"green\" filter distribution")

	blue_refl_filename := flag.String("blue", gopath+"/src/github.com/rennis250/voronoi_filters/data/munsell_blue_EXTREEEMMMEE.spd", "SSV file with reflectance at different wavelengths for \"blue\" filter distribution")

	yellow_refl_filename := flag.String("yellow", gopath+"/src/github.com/rennis250/voronoi_filters/data/munsell_yellow_EXTREEEMMMEE.spd", "SSV file with reflectance at different wavelengths for \"yellow\" filter distribution")

	flag.Parse()

	lms.RGBToLMSFromSpectra(*spectraFile)
	dkl.DKLToRGBFromChroma(*chromaFile)

	cones, wlns := color_img_utils.LoadWavelengthsAndLMS(gopath + "/src/github.com/rennis250/lms/data/linss2_10e_1.csv")

	red_refl := color_img_utils.LoadSpectrum(*red_refl_filename, wlns)
	green_refl := color_img_utils.LoadSpectrum(*green_refl_filename, wlns)
	blue_refl := color_img_utils.LoadSpectrum(*blue_refl_filename, wlns)
	yellow_refl := color_img_utils.LoadSpectrum(*yellow_refl_filename, wlns)

	ld := 0.3714
	rg_mix := 0.3892
	by_mix := 0.8737

	refl := make([]float64, len(red_refl))
	for i := range refl {
		refl[i] = ld * (rg_mix*red_refl[i] + (1-rg_mix)*green_refl[i] + by_mix*blue_refl[i] + (1-by_mix)*yellow_refl[i])
	}

	width, height := 256, 256
	orig_img := image.NewRGBA64(image.Rect(0, 0, width, height))
	vor_glaven_img := image.NewRGBA64(image.Rect(0, 0, width, height))
	vor_filt_img := image.NewRGBA64(image.Rect(0, 0, width, height))

	binary_mask := color_img_utils.MakePieMask(width, height, 0, 0, 60.0, 0.001, 2*math.Pi)

	var filter, bkgd [3]float64
	for i, r := range refl {
		r2 := math.Pow(r, 2.2)
		filter[0] += r2 * cones[i][0]
		filter[1] += r2 * cones[i][1]
		filter[2] += r2 * cones[i][2]

		bkgd[0] += cones[i][0]
		bkgd[1] += cones[i][1]
		bkgd[2] += cones[i][2]
	}

	vortex := loadVoronoiTexture(gopath + "/src/github.com/rennis250/voronoi_filters/data/vor.csv")

	scale := 35.0
	step := scale / float64(width)
	cont := 0.0102
	j := 0
	vor_lms := mat.NewDense(3, width*height, make([]float64, 3*width*height))
	vor_no_filt_lms := mat.NewDense(3, width*height, make([]float64, 3*width*height))
	for y := 0.0; y < scale; y += step {
		i := 0
		for x := 0.0; x < scale; x += step {
			_ = voronoi(x, y) * cont
			idx := j*width + i
			g := vortex[idx] * 3.0 * cont
			if binary_mask[idx] {
				vor_lms.Set(0, idx, g*filter[0])
				vor_lms.Set(1, idx, g*filter[1])
				vor_lms.Set(2, idx, g*filter[2])
			} else {
				vor_lms.Set(0, idx, g*bkgd[0])
				vor_lms.Set(1, idx, g*bkgd[1])
				vor_lms.Set(2, idx, g*bkgd[2])
			}
			vor_no_filt_lms.Set(0, idx, g*bkgd[0])
			vor_no_filt_lms.Set(1, idx, g*bkgd[1])
			vor_no_filt_lms.Set(2, idx, g*bkgd[2])
			i += 1
		}
		j += 1
	}

	vor_rgb := mat.NewDense(3, width*height, make([]float64, 3*width*height))
	vor_rgb.Mul(lms.LMS2RGB, vor_lms)

	vor_no_filt_rgb := mat.NewDense(3, width*height, make([]float64, 3*width*height))
	vor_no_filt_rgb.Mul(lms.LMS2RGB, vor_no_filt_lms)

	vor_no_filt_dkl := mat.NewDense(3, width*height, make([]float64, 3*width*height))
	vor_no_filt_dkl.Mul(dkl.RGB2DKL, vor_no_filt_rgb)

	fmin_M_filt := mat.NewDense(3, 3, []float64{0.8779, 3.5645, -1.1192,
		0.0064, -1.0742, -0.0681,
		-0.0192, 3.4814, 0.8430})
	fmin_t_filt := mat.NewDense(3, 1, []float64{0.0046, 0.0009, 0.0068})

	fmin_M_glaven := mat.NewDense(3, 3, []float64{0.6325, -0.6668, 0.1296,
		-0.2868, 0.3457, -0.0393,
		0.3818, -0.4217, 0.1448})
	fmin_t_glaven := mat.NewDense(3, 1, []float64{-0.3804, -0.2701, 0.3836})

	oog := false
	var r_oog, g_oog, b_oog float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			r_gc := math.Pow(vor_rgb.At(0, idx), 1.0/2.2)
			g_gc := math.Pow(vor_rgb.At(1, idx), 1.0/2.2)
			b_gc := math.Pow(vor_rgb.At(2, idx), 1.0/2.2)

			if r_gc > 1.0 || r_gc < 0.0 || g_gc > 1.0 || g_gc < 0.0 || b_gc > 1.0 || b_gc < 0.0 {
				oog = true
				r_oog = r_gc
				b_oog = b_gc
				g_oog = g_gc
			}

			orig_img.Set(x, y, color.RGBA64{uint16(r_gc * 65535.0),
				uint16(g_gc * 65535.0),
				uint16(b_gc * 65535.0),
				65535})

			var rgb_glaven, rgb_filt [3]float64
			if binary_mask[idx] {
				ldrgyv_glaven := mat.NewDense(3, 1, []float64{vor_no_filt_dkl.At(0, idx),
					vor_no_filt_dkl.At(1, idx),
					vor_no_filt_dkl.At(2, idx)})
				ldrgyv_glaven.Mul(fmin_M_glaven, ldrgyv_glaven)
				ldrgyv_glaven.Add(ldrgyv_glaven, fmin_t_glaven)
				rgb_glaven_vec := mat.NewDense(3, 1, []float64{0, 0, 0})
				rgb_glaven_vec.Mul(dkl.DKL2RGB, ldrgyv_glaven)
				rgb_glaven = [3]float64{math.Pow(rgb_glaven_vec.At(0, 0), 1.0/2.2),
					math.Pow(rgb_glaven_vec.At(1, 0), 1.0/2.2),
					math.Pow(rgb_glaven_vec.At(2, 0), 1.0/2.2)}

				r_gc = rgb_glaven[0]
				g_gc = rgb_glaven[1]
				b_gc = rgb_glaven[2]

				vor_glaven_img.Set(x, y, color.RGBA64{uint16(r_gc * 65535.0),
					uint16(g_gc * 65535.0),
					uint16(b_gc * 65535.0),
					65535})

				ldrgyv_filt := mat.NewDense(3, 1, []float64{vor_no_filt_dkl.At(0, idx),
					vor_no_filt_dkl.At(1, idx),
					vor_no_filt_dkl.At(2, idx)})
				ldrgyv_filt.Mul(fmin_M_filt, ldrgyv_filt)
				ldrgyv_filt.Add(ldrgyv_filt, fmin_t_filt)
				rgb_filt_vec := mat.NewDense(3, 1, []float64{0, 0, 0})
				rgb_filt_vec.Mul(dkl.DKL2RGB, ldrgyv_filt)
				rgb_filt = [3]float64{math.Pow(rgb_filt_vec.At(0, 0), 1.0/2.2),
					math.Pow(rgb_filt_vec.At(1, 0), 1.0/2.2),
					math.Pow(rgb_filt_vec.At(2, 0), 1.0/2.2)}

				r_gc = rgb_filt[0]
				g_gc = rgb_filt[1]
				b_gc = rgb_filt[2]

				vor_filt_img.Set(x, y, color.RGBA64{uint16(r_gc * 65535.0),
					uint16(g_gc * 65535.0),
					uint16(b_gc * 65535.0),
					65535})
			} else {
				vor_glaven_img.Set(x, y, color.RGBA64{uint16(r_gc * 65535.0),
					uint16(g_gc * 65535.0),
					uint16(b_gc * 65535.0),
					65535})

				vor_filt_img.Set(x, y, color.RGBA64{uint16(r_gc * 65535.0),
					uint16(g_gc * 65535.0),
					uint16(b_gc * 65535.0),
					65535})
			}
		}
	}

	if oog {
		fmt.Println("Voronoi - OOG:", r_oog, g_oog, b_oog)
	}

	file, err := os.Create("voronoi_orig.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	png.Encode(file, orig_img)

	file, err = os.Create("voronoi_glaven.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	png.Encode(file, vor_glaven_img)

	file, err = os.Create("voronoi_filter.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	png.Encode(file, vor_filt_img)
}
