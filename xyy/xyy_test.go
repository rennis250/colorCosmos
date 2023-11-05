package dkl

import (
	"image"
	"image/color"
	"testing"
)

func TestDKLToRGBAndBack(t *testing.T) {
	DKLToRGBFromChroma("/home/me/my_docs/science_projects/color_filter_matching/new_structure/color_filter_matching/calibration/eizo_chroma.csv")

	ld_exp := -0.37046195482623934
	rg_exp := 0.3307311869955214
	yv_exp := 0.7603233593684466

	dkl_act := DKLModel{}.Convert(color.RGBA64{uint16(32767), // 0.5 * 65535
		uint16(0.2 * 65535),
		uint16(45874), // 0.7 * 65535
		65535}).(DKL)

	if dkl_act.LD != ld_exp || dkl_act.RG != rg_exp || dkl_act.YV != yv_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{ld_exp, rg_exp, yv_exp}, []float64{dkl_act.LD, dkl_act.RG, dkl_act.YV})
	}

	r_exp := 0.49999237048905165
	g_exp := 0.2
	b_exp := 0.6999923704890516

	r_act, g_act, b_act, _ := dkl_act.RGBA()
	r_act_f := float64(r_act) / 65535.0
	g_act_f := float64(g_act) / 65535.0
	b_act_f := float64(b_act) / 65535.0

	if r_act_f != r_exp || g_act_f != g_exp || b_act_f != b_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{r_exp, g_exp, b_exp}, []float64{r_act_f, g_act_f, b_act_f})
	}
}

func TestDKLImage(t *testing.T) {
	DKLToRGBFromChroma("/home/me/my_docs/science_projects/color_filter_matching/new_structure/color_filter_matching/calibration/eizo_chroma.csv")

	dkl_img := NewDKLImage(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			dkl_img.Set(x, y, color.RGBA64{uint16(32767), // 0.5 * 65535
				uint16(0.2 * 65535),
				uint16(45874), // 0.7 * 65535
				65535})
		}
	}

	ld_exp := -0.37046195482623934
	rg_exp := 0.3307311869955214
	yv_exp := 0.7603233593684466

	dkl_act := dkl_img.At(5, 5).(DKL)

	if dkl_act.LD != ld_exp || dkl_act.RG != rg_exp || dkl_act.YV != yv_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{ld_exp, rg_exp, yv_exp}, []float64{dkl_act.LD, dkl_act.RG, dkl_act.YV})
	}

	dkl_img2 := NewDKLImage(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			dkl_img2.SetDKL(x, y, DKL{-0.37046195482623934,
				0.3307311869955214,
				0.7603233593684466})
		}
	}

	dkl_act2 := dkl_img2.DKLAt(5, 5)

	if dkl_act2.LD != ld_exp || dkl_act2.RG != rg_exp || dkl_act2.YV != yv_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{ld_exp, rg_exp, yv_exp}, []float64{dkl_act2.LD, dkl_act2.RG, dkl_act2.YV})
	}

	rgb_img := image.NewRGBA64(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			rgb_img.Set(x, y, DKL{-0.37046195482623934,
				0.3307311869955214,
				0.7603233593684466})
		}
	}

	r_exp := 0.49999237048905165
	g_exp := 0.2
	b_exp := 0.6999923704890516

	r_act, g_act, b_act, _ := rgb_img.At(5, 5).(color.RGBA64).RGBA()
	r_act_f := float64(r_act) / 65535.0
	g_act_f := float64(g_act) / 65535.0
	b_act_f := float64(b_act) / 65535.0

	if r_act_f != r_exp || g_act_f != g_exp || b_act_f != b_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{r_exp, g_exp, b_exp}, []float64{r_act_f, g_act_f, b_act_f})
	}
}
