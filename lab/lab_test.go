package lab

import (
	"image"
	"image/color"
	"testing"
)

func TestLABToRGBAndBack(t *testing.T) {
	RGBToXYZFromChroma("/home/me/my_docs/science_projects/color_filter_matching/new_structure/color_filter_matching/calibration/eizo_chroma.csv")

	L_exp := 62.90817233318154
	A_exp := 49.50222733240855
	B_exp := -36.83272867266696

	lab_act := LABModel{}.Convert(color.RGBA64{uint16(32767), // 0.5 * 65535
		uint16(0.2 * 65535),
		uint16(45874), // 0.7 * 65535
		65535}).(LAB)

	if lab_act.L != L_exp || lab_act.A != A_exp || lab_act.B != B_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{L_exp, A_exp, B_exp}, []float64{lab_act.L, lab_act.A, lab_act.B})
	}

	r_exp := 0.49999237048905165
	g_exp := 0.19998474097810331
	b_exp := 0.6999923704890516

	r_act, g_act, b_act, _ := lab_act.RGBA()
	r_act_f := float64(r_act) / 65535.0
	g_act_f := float64(g_act) / 65535.0
	b_act_f := float64(b_act) / 65535.0

	if r_act_f != r_exp || g_act_f != g_exp || b_act_f != b_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{r_exp, g_exp, B_exp}, []float64{r_act_f, g_act_f, b_act_f})
	}
}

func TestLABImage(t *testing.T) {
	RGBToXYZFromChroma("/home/me/my_docs/science_projects/color_filter_matching/new_structure/color_filter_matching/calibration/eizo_chroma.csv")

	lab_img := NewLABImage(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			lab_img.Set(x, y, color.RGBA64{uint16(32767), // 0.5 * 65535
				uint16(0.2 * 65535),
				uint16(45874), // 0.7 * 65535
				65535})
		}
	}

	L_exp := 62.90817233318154
	A_exp := 49.50222733240855
	B_exp := -36.83272867266696

	lab_act := lab_img.At(5, 5).(LAB)

	if lab_act.L != L_exp || lab_act.A != A_exp || lab_act.B != B_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{L_exp, A_exp, B_exp}, []float64{lab_act.L, lab_act.A, lab_act.B})
	}

	lab_img2 := NewLABImage(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			lab_img2.SetLAB(x, y, LAB{62.90817233318154,
				49.50222733240855,
				-36.83272867266696})
		}
	}

	lab_act2 := lab_img2.LABAt(5, 5)

	if lab_act2.L != L_exp || lab_act2.A != A_exp || lab_act2.B != B_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{L_exp, A_exp, B_exp}, []float64{lab_act2.L, lab_act2.A, lab_act2.B})
	}

	rgb_img := image.NewRGBA64(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			rgb_img.Set(x, y, LAB{62.90817233318154,
				49.50222733240855,
				-36.83272867266696})
		}
	}

	r_exp := 0.49999237048905165
	g_exp := 0.19998474097810331
	b_exp := 0.6999923704890516

	r_act, g_act, b_act, _ := rgb_img.At(5, 5).(color.RGBA64).RGBA()
	r_act_f := float64(r_act) / 65535.0
	g_act_f := float64(g_act) / 65535.0
	b_act_f := float64(b_act) / 65535.0

	if r_act_f != r_exp || g_act_f != g_exp || b_act_f != b_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{r_exp, g_exp, B_exp}, []float64{r_act_f, g_act_f, b_act_f})
	}
}
