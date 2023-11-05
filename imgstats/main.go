package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rennis250/colorCosmos/chromStat"
	"github.com/rennis250/colorCosmos/colorImgUtils"
	"github.com/rennis250/colorCosmos/dkl"
	"github.com/rennis250/colorCosmos/lab"
	"github.com/rennis250/colorCosmos/lms"
	"github.com/rennis250/colorCosmos/xyz"
)

func main() {
	monName := flag.String("monName", "", "preferred name of monitor so that output csv file is labelled with it")
	spectraFile := flag.String("spectraFile", "mon_spectra.csv", "CSV file with spectral curves of each primary on monitor")
	chromaFile := flag.String("chromaFile", "mon_chroma.csv", "CSV file with chromaticity coordinates of each primary on monitor")
	gammaFile := flag.String("gammaFile", "mon_chroma.csv", "CSV file with chromaticity coordinates of each primary on monitor")
	maskDir := flag.String("maskDir", "", "directory containing masks to apply to each image, so that localized statistics can be extracted")
	darkThresh := flag.Float64("darkThresh", -1.0, "any pixels with a luminance (CIE1931 Y) less than this value will be excluded from the analysis")

	flag.Parse()

	lms.RGBToLMSFromSpectra(*spectraFile)
	dkl.DKLToRGBFromChroma(*chromaFile)
	lab.RGBToXYZFromChroma(*chromaFile)
	xyz.RGBToXYZFromChroma(*chromaFile)

	gR, gG, gB := colorImgUtils.LoadGamma(*gammaFile)

	masks, err := filepath.Glob(*maskDir + "/*.png")
	if err != nil {
		log.Fatal(err)
	}

	// everything else should be names
	// of files to process
	files := flag.Args()

	for _, maskname := range masks {
		binary_mask := colorImgUtils.LoadBinaryMask(maskname)

		records := make([][]string, len(files)+1)
		records[0] = []string{"img_name",
			"mask_name",
			"lrc",
			"redness",
			"wpld",
			"wprg",
			"wpyv",
			"gwld",
			"gwrg",
			"gwyv",
			"gwl",
			"gwa",
			"gwb",
			"hv",
			"mean_cone_l",
			"mean_cone_m",
			"mean_cone_s",
			"sd_cone_l",
			"sd_cone_m",
			"sd_cone_s",
			"rmc_l",
			"rmc_m",
			"rmc_s",
			"rsd_l",
			"rsd_m",
			"rsd_s",
			"robust_simpler_l",
			"robust_simpler_m",
			"robust_simpler_s",
			"robust_general_l",
			"robust_general_m",
			"robust_general_s",
			"wpL",
			"wpA",
			"wpB"}

		for fc, file := range files {
			fimg, err := os.Open(file)
			if err != nil {
				panic("could not open image file")
			}
			img, err := png.Decode(fimg)
			if err != nil {
				log.Fatal(err)
			}

			imggc := colorImgUtils.GammaCorr(img, gR, gG, gB)

			imgxyz := xyz.ConvertToXYZ(imggc)
			imglms := lms.ConvertToLMS(imggc)
			imgdkl := dkl.ConvertToDKL(imggc)
			imglab := lab.ConvertToLAB(imggc)

			if *darkThresh >= 0 {
				colorImgUtils.ExcludeDarkPixelsFromMask(binary_mask, imgxyz, *darkThresh)
			}

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

			records[fc+1] = []string{
				file,
				maskname,

				strconv.FormatFloat(lrc, 'f', 7, 64),
				strconv.FormatFloat(redness, 'f', 7, 64),

				strconv.FormatFloat(wpdkl.LD, 'f', 7, 64),
				strconv.FormatFloat(wpdkl.RG, 'f', 7, 64),
				strconv.FormatFloat(wpdkl.YV, 'f', 7, 64),

				strconv.FormatFloat(gwdkl.LD, 'f', 7, 64),
				strconv.FormatFloat(gwdkl.RG, 'f', 7, 64),
				strconv.FormatFloat(gwdkl.YV, 'f', 7, 64),

				strconv.FormatFloat(gwlab.L, 'f', 7, 64),
				strconv.FormatFloat(gwlab.A, 'f', 7, 64),
				strconv.FormatFloat(gwlab.B, 'f', 7, 64),

				strconv.FormatFloat(hv, 'f', 7, 64),

				strconv.FormatFloat(lmsmean.L, 'f', 7, 64),
				strconv.FormatFloat(lmsmean.M, 'f', 7, 64),
				strconv.FormatFloat(lmsmean.S, 'f', 7, 64),

				strconv.FormatFloat(lmssd.L, 'f', 7, 64),
				strconv.FormatFloat(lmssd.M, 'f', 7, 64),
				strconv.FormatFloat(lmssd.S, 'f', 7, 64),

				strconv.FormatFloat(rmc[0], 'f', 7, 64),
				strconv.FormatFloat(rmc[1], 'f', 7, 64),
				strconv.FormatFloat(rmc[2], 'f', 7, 64),

				strconv.FormatFloat(rsd[0], 'f', 7, 64),
				strconv.FormatFloat(rsd[1], 'f', 7, 64),
				strconv.FormatFloat(rsd[2], 'f', 7, 64),

				strconv.FormatFloat(tau_simpler[0], 'f', 7, 64),
				strconv.FormatFloat(tau_simpler[1], 'f', 7, 64),
				strconv.FormatFloat(tau_simpler[2], 'f', 7, 64),

				strconv.FormatFloat(tau_general[0], 'f', 7, 64),
				strconv.FormatFloat(tau_general[1], 'f', 7, 64),
				strconv.FormatFloat(tau_general[2], 'f', 7, 64),

				strconv.FormatFloat(wplab.L, 'f', 7, 64),
				strconv.FormatFloat(wplab.A, 'f', 7, 64),
				strconv.FormatFloat(wplab.B, 'f', 7, 64),
			}

			fmt.Println(fc+1, len(records)-1)

			if err := fimg.Close(); err != nil {
				log.Fatal(err)
			}
		}

		var stats_file *os.File
		var err error
		if *darkThresh >= 0 {
			stats_file, err = os.Create(*monName + "_img_stats_dark_pxs_excl.csv")
		} else {
			stats_file, err = os.Create(*monName + "_img_stats.csv")
		}
		if err != nil {
			log.Fatal("An error encountered ::", err)
		}
		w := csv.NewWriter(stats_file)

		for _, record := range records {
			if err := w.Write(record); err != nil {
				log.Fatalln("error write record to csv:", err)
			}
		}

		w.Flush()

		if err := w.Error(); err != nil {
			log.Fatal(err)
		}

		if err := stats_file.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
