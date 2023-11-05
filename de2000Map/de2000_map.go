package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rennis250/colorCosmos/colorImgUtils"
	"github.com/rennis250/colorCosmos/deltaE"
	"github.com/rennis250/colorCosmos/lab"
)

func main() {
	monName := flag.String("monName", "", "preferred name of monitor so that output csv file is labelled with it")
	chromaFile := flag.String("chromaFile", "mon_chroma.csv", "CSV file with chromaticity coordinates of each primary on monitor")
	gammaFile := flag.String("gammaFile", "mon_chroma.csv", "CSV file with chromaticity coordinates of each primary on monitor")
	maskDir := flag.String("maskDir", "", "directory containing masks to apply to each image, so that localized statistics can be extracted")

	anchorColorL := flag.Float64("anchorColorL", 50.0, "CIELAB L* coordinate of anchor color")
	anchorColorA := flag.Float64("anchorColorA", 0.0, "CIELAB a* coordinate of anchor color")
	anchorColorB := flag.Float64("anchorColorB", 0.0, "CIELAB b* coordinate of anchor color")

	DEDist := flag.Float64("DEDist", 15.0, "threshold distance (in DE2000 units) from anchor color for color in image to be part of colored DE2000 map")

	flag.Parse()

	anchorColorLAB := lab.LAB{*anchorColorL, *anchorColorA, *anchorColorB}

	lab.RGBToXYZFromChroma(*chromaFile)

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

		for _, filename := range files {
			fimg, err := os.Open(filename)
			if err != nil {
				panic("could not open image file")
			}
			img, err := png.Decode(fimg)
			if err != nil {
				log.Fatal(err)
			}

			imggc := colorImgUtils.GammaCorr(img, gR, gG, gB)

			imglab := lab.ConvertToLAB(imggc)

			// prepare a new image to hold the results of the DE2000 mapping
			b := img.Bounds()
			width := b.Max.Y - b.Min.Y
			de2000_img := image.NewRGBA64(img.Bounds())
			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					idx := y*width + x
					if binary_mask[idx] {
						if deltaE.DE2000(imglab.LABAt(x, y), anchorColorLAB) < *DEDist {
							de2000_img.Set(x, y, img.At(x, y))
						} else {
							grayval := color.Gray16Model.Convert(img.At(x, y)).(color.Gray16)

							// reduce the max intensity of the gray values, so that highlights
							// are not accidently mistake for being colored in.
							// also helps to bring out the colored regions better.
							grayval.Y = uint16(float64(grayval.Y) / 1.5)

							de2000_img.Set(x, y, grayval)
						}
					} else {
						de2000_img.Set(x, y, color.RGBA64{0, 0, 0, 65535})
					}
				}
			}

			filenameParts := strings.Split(filename, "/")
			masknameParts := strings.Split(maskname, "/")

			filname_wout_dir := filenameParts[len(filenameParts)-1]
			maskname_wout_dir := masknameParts[len(masknameParts)-1]

			file, err := os.Create("de2000_map_Img_" +
				filname_wout_dir[:len(filname_wout_dir)-4] +
				"_Mask_" +
				maskname_wout_dir[:len(maskname_wout_dir)-4] +
				"_L_" +
				strconv.FormatFloat(*anchorColorL, 'f', 7, 64) +
				"_A_" +
				strconv.FormatFloat(*anchorColorA, 'f', 7, 64) +
				"_B_" +
				strconv.FormatFloat(*anchorColorB, 'f', 7, 64) +
				"_Monitor_" +
				*monName +
				".png")
			if err != nil {
				log.Fatalln(err)
			}
			png.Encode(file, de2000_img)
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}

			if err := fimg.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
