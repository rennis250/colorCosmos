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

	"github.com/rennis250/colorCosmos/chromStat"
	"github.com/rennis250/colorCosmos/colorImgUtils"
	"github.com/rennis250/colorCosmos/dkl"
	"github.com/rennis250/colorCosmos/lab"
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
	gammaFile := flag.String("gammaFile", "mon_chroma.csv", "CSV file with chromaticity coordinates of each primary on monitor")

	printStatsFlag := flag.Bool("printStats", true, "print out some image statistics for the filter?")
	saveImgFlag := flag.Bool("saveImg", false, "save an image of the flat filter overlaying an achromatic voronoi background?")

	red_refl_filename := flag.String("red", gopath+"/src/github.com/rennis250/voronoi_filters/data/munsell_red_EXTREEEMMMEE.spd", "SSV file with reflectance at different wavelengths for \"red\" filter distribution")

	green_refl_filename := flag.String("green", gopath+"/src/github.com/rennis250/voronoi_filters/data/munsell_green_EXTREEEMMMEE.spd", "SSV file with reflectance at different wavelengths for \"green\" filter distribution")

	blue_refl_filename := flag.String("blue", gopath+"/src/github.com/rennis250/voronoi_filters/data/munsell_blue_EXTREEEMMMEE.spd", "SSV file with reflectance at different wavelengths for \"blue\" filter distribution")

	yellow_refl_filename := flag.String("yellow", gopath+"/src/github.com/rennis250/voronoi_filters/data/munsell_yellow_EXTREEEMMMEE.spd", "SSV file with reflectance at different wavelengths for \"yellow\" filter distribution")

	ld := flag.Float64("ld", 1.0, "How thick/light should the filter be?")
	rg_mix := flag.Float64("rg", 0.5, "More green (1.0) or red (0.0)?")
	by_mix := flag.Float64("by", 0.5, "More yellow (1.0) or blue (0.0)?")

	flag.Parse()

	lms.RGBToLMSFromSpectra(*spectraFile)

	gR, gG, gB := colorImgUtils.LoadGamma(*gammaFile)

	cones, wlns := colorImgUtils.LoadWavelengthsAndLMS(gopath + "/src/github.com/rennis250/lms/data/linss2_10e_1.csv")

	red_refl := colorImgUtils.LoadSpectrum(*red_refl_filename, wlns)
	green_refl := colorImgUtils.LoadSpectrum(*green_refl_filename, wlns)
	blue_refl := colorImgUtils.LoadSpectrum(*blue_refl_filename, wlns)
	yellow_refl := colorImgUtils.LoadSpectrum(*yellow_refl_filename, wlns)

	refl := make([]float64, len(red_refl))
	for i := range refl {
		refl[i] = *ld * (*rg_mix*red_refl[i] + (1-*rg_mix)*green_refl[i] + *by_mix*blue_refl[i] + (1-*by_mix)*yellow_refl[i])
	}

	width, height := 256, 256
	img := image.NewRGBA64(image.Rect(0, 0, width, height))

	binary_mask := colorImgUtils.MakePieMask(width, height, 0, 0, 60.0, 0.001, 2*math.Pi)

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
			i += 1
		}
		j += 1
	}

	vor_rgb := mat.NewDense(3, width*height, make([]float64, 3*width*height))
	vor_rgb.Mul(lms.LMS2RGB, vor_lms)

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

			img.Set(x, y, color.RGBA64{uint16(r_gc * 65535.0),
				uint16(g_gc * 65535.0),
				uint16(b_gc * 65535.0),
				65535})
		}
	}

	if oog {
		fmt.Println("Voronoi - OOG:", r_oog, g_oog, b_oog)
	}

	if *saveImgFlag {
		file, err := os.Create("voronoi_filter_ld_" +
			strconv.FormatFloat(*ld, 'f', 7, 64) +
			"_rg_" +
			strconv.FormatFloat(*rg_mix, 'f', 7, 64) +
			"_by_" +
			strconv.FormatFloat(*by_mix, 'f', 7, 64) +
			".png")
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		png.Encode(file, img)
	}

	if *printStatsFlag {
		lms.RGBToLMSFromSpectra(*spectraFile)
		dkl.DKLToRGBFromChroma(*chromaFile)
		lab.RGBToXYZFromChroma(*chromaFile)

		imggc := colorImgUtils.GammaCorr(img, gR, gG, gB)

		imglms := lms.ConvertToLMS(imggc)
		imgdkl := dkl.ConvertToDKL(imggc)
		imglab := lab.ConvertToLAB(imggc)

		lrc, redness := chromStat.LRC(imglms, binary_mask)
		wpdkl := chromStat.DKLWhitePoint(imgdkl, binary_mask)
		wplab := chromStat.LABWhitePoint(imglab, binary_mask)
		gwdkl := chromStat.DKLGrayWorld(imgdkl, binary_mask)
		gwlab := chromStat.LABGrayWorld(imglab, binary_mask)
		hv := chromStat.HueVariance(imgdkl, binary_mask)
		lmsmean, lmssd := chromStat.MeanSdCones(imglms, binary_mask)
		rmc, rsd := chromStat.RatioMeanSdCones(imglms, binary_mask)
		tau_simpler := chromStat.GeneralRobustRatio(imglms, binary_mask)
		tau_general := chromStat.SimplerRobustRatio(imglms, binary_mask)
		// _, _ = chromStat.FFTColorConstancy(imgdkl, 5, binary_mask)

		fmt.Print(strconv.FormatFloat(lrc, 'f', 7, 64) + "," +
			strconv.FormatFloat(redness, 'f', 7, 64) + "," +

			strconv.FormatFloat(wpdkl.LD, 'f', 7, 64) + "," +
			strconv.FormatFloat(wpdkl.RG, 'f', 7, 64) + "," +
			strconv.FormatFloat(wpdkl.YV, 'f', 7, 64) + "," +

			strconv.FormatFloat(gwdkl.LD, 'f', 7, 64) + "," +
			strconv.FormatFloat(gwdkl.RG, 'f', 7, 64) + "," +
			strconv.FormatFloat(gwdkl.YV, 'f', 7, 64) + "," +

			strconv.FormatFloat(gwlab.L, 'f', 7, 64) + "," +
			strconv.FormatFloat(gwlab.A, 'f', 7, 64) + "," +
			strconv.FormatFloat(gwlab.B, 'f', 7, 64) + "," +

			strconv.FormatFloat(hv, 'f', 7, 64) + "," +

			strconv.FormatFloat(lmsmean.L, 'f', 7, 64) + "," +
			strconv.FormatFloat(lmsmean.M, 'f', 7, 64) + "," +
			strconv.FormatFloat(lmsmean.S, 'f', 7, 64) + "," +

			strconv.FormatFloat(lmssd.L, 'f', 7, 64) + "," +
			strconv.FormatFloat(lmssd.M, 'f', 7, 64) + "," +
			strconv.FormatFloat(lmssd.S, 'f', 7, 64) + "," +

			strconv.FormatFloat(rmc[0], 'f', 7, 64) + "," +
			strconv.FormatFloat(rmc[1], 'f', 7, 64) + "," +
			strconv.FormatFloat(rmc[2], 'f', 7, 64) + "," +

			strconv.FormatFloat(rsd[0], 'f', 7, 64) + "," +
			strconv.FormatFloat(rsd[1], 'f', 7, 64) + "," +
			strconv.FormatFloat(rsd[2], 'f', 7, 64) + "," +

			strconv.FormatFloat(tau_simpler[0], 'f', 7, 64) + "," +
			strconv.FormatFloat(tau_simpler[1], 'f', 7, 64) + "," +
			strconv.FormatFloat(tau_simpler[2], 'f', 7, 64) + "," +

			strconv.FormatFloat(tau_general[0], 'f', 7, 64) + "," +
			strconv.FormatFloat(tau_general[1], 'f', 7, 64) + "," +
			strconv.FormatFloat(tau_general[2], 'f', 7, 64) + "," +

			strconv.FormatFloat(wplab.L, 'f', 7, 64) + "," +
			strconv.FormatFloat(wplab.A, 'f', 7, 64) + "," +
			strconv.FormatFloat(wplab.B, 'f', 7, 64))
	}
}
