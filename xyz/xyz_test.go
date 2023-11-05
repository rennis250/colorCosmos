package xyz

import (
	"image"
	"image/color"
	"testing"
)

func TestXYZToRGBAndBack(t *testing.T) {
	RGBToXYZFromChroma("/home/me/my_docs/science_projects/transparency/color_filter_matching/calibration/eizo_chroma.csv")

	x_exp := 61.53802066179045
	y_exp := 41.487375576409555
	z_exp := 82.5546897289983

	xyz_act := XYZModel{}.Convert(color.RGBA64{uint16(32767), // 0.5 * 65535
		uint16(0.2 * 65535),
		uint16(45874), // 0.7 * 65535
		65535}).(XYZ)

	if xyz_act.X != x_exp || xyz_act.Y != y_exp || xyz_act.Z != z_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{x_exp, y_exp, z_exp}, []float64{xyz_act.X, xyz_act.Y, xyz_act.Z})
	}

	r_exp := 0.49999237048905165
	g_exp := 0.2
	b_exp := 0.6999923704890516

	r_act, g_act, b_act, _ := xyz_act.RGBA()
	r_act_f := float64(r_act) / 65535.0
	g_act_f := float64(g_act) / 65535.0
	b_act_f := float64(b_act) / 65535.0

	if r_act_f != r_exp || g_act_f != g_exp || b_act_f != b_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{r_exp, g_exp, b_exp}, []float64{r_act_f, g_act_f, b_act_f})
	}
}

func TestXYZImage(t *testing.T) {
	RGBToXYZFromChroma("/home/me/my_docs/science_projects/transparency/color_filter_matching/calibration/eizo_chroma.csv")

	xyz_img := NewXYZImage(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			xyz_img.Set(x, y, color.RGBA64{uint16(32767), // 0.5 * 65535
				uint16(0.2 * 65535),
				uint16(45874), // 0.7 * 65535
				65535})
		}
	}

	x_exp := 61.53802066179045
	y_exp := 41.487375576409555
	z_exp := 82.5546897289983

	xyz_act := xyz_img.At(5, 5).(XYZ)

	if xyz_act.X != x_exp || xyz_act.Y != y_exp || xyz_act.Z != z_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{x_exp, y_exp, z_exp}, []float64{xyz_act.X, xyz_act.Y, xyz_act.Z})
	}

	xyz_img2 := NewXYZImage(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			xyz_img2.SetXYZ(x, y, XYZ{61.53802066179045,
				41.487375576409555,
				82.5546897289983})
		}
	}

	xyz_act2 := xyz_img2.XYZAt(5, 5)

	if xyz_act2.X != x_exp || xyz_act2.Y != y_exp || xyz_act2.Z != z_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{x_exp, y_exp, z_exp}, []float64{xyz_act2.X, xyz_act2.Y, xyz_act2.Z})
	}

	rgb_img := image.NewRGBA64(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			rgb_img.Set(x, y, XYZ{61.53802066179045,
				41.487375576409555,
				82.5546897289983})
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
