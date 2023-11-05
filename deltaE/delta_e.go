package deltaE

import (
	"math"

	"github.com/rennis250/colorCosmos/lab"
)

func DE2000(color_1, color_2 lab.LAB) float64 {
	ksub_l := 1.0
	ksub_c := 1.0
	ksub_h := 1.0

	delta_l_prime := color_2.L - color_1.L

	l_bar := (color_1.L + color_2.L) / 2.0

	c1 := math.Sqrt(color_1.A*color_1.A + color_1.B*color_1.B)
	c2 := math.Sqrt(color_2.A*color_2.A + color_2.B*color_2.B)

	c_bar := (c1 + c2) / 2.0

	a_prime_1 := color_1.A +
		(color_1.A/2.0)*
			(1.0-math.Sqrt(math.Pow(c_bar, 7)/(math.Pow(c_bar, 7)+math.Pow(25.0, 7))))
	a_prime_2 := color_2.A +
		(color_2.A/2.0)*
			(1.0-math.Sqrt(math.Pow(c_bar, 7)/(math.Pow(c_bar, 7)+math.Pow(25.0, 7))))

	c_prime_1 := math.Sqrt(a_prime_1*a_prime_1 + color_1.B*color_1.B)
	c_prime_2 := math.Sqrt(a_prime_2*a_prime_2 + color_2.B*color_2.B)

	c_bar_prime := (c_prime_1 + c_prime_2) / 2.0

	delta_c_prime := c_prime_2 - c_prime_1

	l_bar_50 := l_bar - 50.0
	l_bar_50_2 := l_bar_50 * l_bar_50
	s_sub_l := 1.0 + ((0.015 * l_bar_50_2) / math.Sqrt(20.0+l_bar_50_2))

	s_sub_c := 1.0 + 0.045*c_bar_prime

	h_prime_1 := getHPrimeFn(color_1.B, a_prime_1)
	h_prime_2 := getHPrimeFn(color_2.B, a_prime_2)

	delta_h_prime := getDeltaHPrime(c1, c2, h_prime_1, h_prime_2)

	delta_upcase_h_prime := 2.0 * math.Sqrt(c_prime_1*c_prime_2) *
		math.Sin(degreesToRadians(delta_h_prime)/2.0)

	upcase_h_bar_prime := getUpcaseHBarPrime(h_prime_1, h_prime_2, c_prime_1, c_prime_2)

	upcase_t := getUpcaseT(upcase_h_bar_prime)

	s_sub_upcase_h := 1.0 + 0.015*c_bar_prime*upcase_t

	r_sub_t := getRSubT(c_bar_prime, upcase_h_bar_prime)

	lightness := delta_l_prime / (ksub_l * s_sub_l)

	chroma := delta_c_prime / (ksub_c * s_sub_c)

	hue := delta_upcase_h_prime / (ksub_h * s_sub_upcase_h)

	return math.Sqrt(lightness*lightness + chroma*chroma + hue*hue + r_sub_t*chroma*hue)
}

func getHPrimeFn(x, y float64) float64 {
	if x == 0.0 && y == 0.0 {
		return 0.0
	}

	hue_angle := radiansToDegrees(math.Atan2(x, y))

	if hue_angle < 0.0 {
		hue_angle += 360.0
	}

	return hue_angle
}

func getDeltaHPrime(c1, c2, h_prime_1, h_prime_2 float64) float64 {
	if 0.0 == c1 || 0.0 == c2 {
		return 0.0
	}

	if math.Abs(h_prime_1-h_prime_2) <= 180.0 {
		return h_prime_2 - h_prime_1
	}

	if h_prime_2 <= h_prime_1 {
		return h_prime_2 - h_prime_1 + 360.0
	} else {
		return h_prime_2 - h_prime_1 - 360.0
	}
}

func getUpcaseHBarPrime(h_prime_1, h_prime_2, c_prime_1, c_prime_2 float64) float64 {
	if c_prime_1 == 0.0 || c_prime_2 == 0.0 {
		return h_prime_1 + h_prime_2
	}

	if math.Abs(h_prime_1-h_prime_2) > 180.0 {
		if (h_prime_1 + h_prime_2) < 360.0 {
			return (h_prime_1 + h_prime_2 + 360.0) / 2.0
		} else if (h_prime_1 + h_prime_2) >= 360.0 {
			return (h_prime_1 + h_prime_2 - 360.0) / 2.0
		}
	}

	return (h_prime_1 + h_prime_2) / 2.0
}

func getUpcaseT(upcase_h_bar_prime float64) float64 {
	return 1.0 - 0.17*math.Cos(degreesToRadians(upcase_h_bar_prime-30.0)) +
		0.24*math.Cos(degreesToRadians(2.0*upcase_h_bar_prime)) +
		0.32*math.Cos(degreesToRadians(3.0*upcase_h_bar_prime+6.0)) -
		0.20*math.Cos(degreesToRadians(4.0*upcase_h_bar_prime-63.0))
}

func getRSubT(c_bar_prime, upcase_h_bar_prime float64) float64 {
	uhbp25 := (upcase_h_bar_prime - 275.0) / 25.0
	return -2.0 *
		math.Sqrt(math.Pow(c_bar_prime, 7)/(math.Pow(c_bar_prime, 7)+math.Pow(25.0, 7))) *
		math.Sin(degreesToRadians(60.0*math.Exp(-uhbp25*uhbp25)))
}

func radiansToDegrees(radians float64) float64 {
	return radians * (180.0 / math.Pi)
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}
