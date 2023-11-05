package lms

import (
	"image"
	"image/color"
	"testing"
)

func TestLMSToRGBAndBack(t *testing.T) {
	RGBToLMSFromSpectra("/home/me/my_docs/science_projects/color_filter_matching/new_structure/color_filter_matching/calibration/eizo_mon_spectra.csv")

	l_exp := 0.07152745477596076
	m_exp := 0.048504029435861865
	s_exp := 0.06381828809453752

	lms_act := LMSModel{}.Convert(color.RGBA64{uint16(32767), // 0.5 * 65535
		uint16(0.2 * 65535),
		uint16(45874), // 0.7 * 65535
		65535}).(LMS)

	if lms_act.L != l_exp || lms_act.M != m_exp || lms_act.S != s_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{l_exp, m_exp, s_exp}, []float64{lms_act.L, lms_act.M, lms_act.S})
	}

	r_exp := 0.49999237048905165
	g_exp := 0.2
	b_exp := 0.6999923704890516

	r_act, g_act, b_act, _ := lms_act.RGBA()
	r_act_f := float64(r_act) / 65535.0
	g_act_f := float64(g_act) / 65535.0
	b_act_f := float64(b_act) / 65535.0

	if r_act_f != r_exp || g_act_f != g_exp || b_act_f != b_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{r_exp, g_exp, b_exp}, []float64{r_act_f, g_act_f, b_act_f})
	}
}

func TestLMSImage(t *testing.T) {
	RGBToLMSFromSpectra("/home/me/my_docs/science_projects/color_filter_matching/new_structure/color_filter_matching/calibration/eizo_mon_spectra.csv")

	lms_img := NewLMSImage(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			lms_img.Set(x, y, color.RGBA64{uint16(32767), // 0.5 * 65535
				uint16(0.2 * 65535),
				uint16(45874), // 0.7 * 65535
				65535})
		}
	}

	l_exp := 0.07152745477596076
	m_exp := 0.048504029435861865
	s_exp := 0.06381828809453752

	lms_act := lms_img.At(5, 5).(LMS)

	if lms_act.L != l_exp || lms_act.M != m_exp || lms_act.S != s_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{l_exp, m_exp, s_exp}, []float64{lms_act.L, lms_act.M, lms_act.S})
	}

	lms_img2 := NewLMSImage(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			lms_img2.SetLMS(x, y, LMS{0.07152745477596076,
				0.048504029435861865,
				0.06381828809453752})
		}
	}

	lms_act2 := lms_img2.LMSAt(5, 5)

	if lms_act2.L != l_exp || lms_act2.M != m_exp || lms_act2.S != s_exp {
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{l_exp, m_exp, s_exp}, []float64{lms_act2.L, lms_act2.M, lms_act2.S})
	}

	rgb_img := image.NewRGBA64(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			rgb_img.Set(x, y, LMS{0.071528,
				0.048504,
				0.063819})
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
		t.Errorf("Test failed, expected: '%v', got: '%v'", []float64{r_exp, g_exp, b_exp}, []float64{r_act_f, g_act_f, b_act_f})
	}
}
