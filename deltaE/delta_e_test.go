package deltaE

import (
	"math"
	"testing"

	"github.com/rennis250/colorCosmos/lab"
)

func round(val float64) float64 {
	rounded := val * 10000.0
	return math.Round(rounded) / 10000.0
}

func assertDeltaE(expected float64, lab1, lab2 [3]float64) {
	color_1 := lab.LAB{lab1[0], lab1[1], lab1[2]}
	color_2 := lab.LAB{lab2[0], lab2[1], lab2[2]}

	if round(DE2000(color_1, color_2)) != expected {
		panic("failed DE2000 test...")
	}
}

// Tests taken from Table 1: "CIEDE2000 total color difference test data" of
// "The CIEDE2000 Color-Difference Formula: Implementation Notes,
// Supplementary Test Data, and Mathematical Observations" by Gaurav Sharma,
// Wencheng Wu and Edul N. Dalal.
//
// http://www.ece.rochester.edu/~gsharma/papers/CIEDE2000CRNAFeb05.pdf

func TestDE2000(t *testing.T) {
	assertDeltaE(0.0, [3]float64{0.0, 0.0, 0.0}, [3]float64{0.0, 0.0, 0.0})
	assertDeltaE(0.0, [3]float64{99.5, 0.005, -0.010}, [3]float64{99.5, 0.005, -0.010})
	assertDeltaE(100.0, [3]float64{100.0, 0.005, -0.010}, [3]float64{0.0, 0.0, 0.0})
	assertDeltaE(2.0425,
		[3]float64{50.0000, 2.6772, -79.7751},
		[3]float64{50.0000, 0.0000, -82.7485})
	assertDeltaE(2.8615,
		[3]float64{50.0000, 3.1571, -77.2803},
		[3]float64{50.0000, 0.0000, -82.7485})
	assertDeltaE(3.4412,
		[3]float64{50.0000, 2.8361, -74.0200},
		[3]float64{50.0000, 0.0000, -82.7485})
	assertDeltaE(1.0000,
		[3]float64{50.0000, -1.3802, -84.2814},
		[3]float64{50.0000, 0.0000, -82.7485})
	assertDeltaE(1.0000,
		[3]float64{50.0000, -1.1848, -84.8006},
		[3]float64{50.0000, 0.0000, -82.7485})
	assertDeltaE(1.0000,
		[3]float64{50.0000, -0.9009, -85.5211},
		[3]float64{50.0000, 0.0000, -82.7485})
	assertDeltaE(2.3669,
		[3]float64{50.0000, 0.0000, 0.0000},
		[3]float64{50.0000, -1.0000, 2.0000})
	assertDeltaE(2.3669,
		[3]float64{50.0000, -1.0000, 2.0000},
		[3]float64{50.0000, 0.0000, 0.0000})
	assertDeltaE(7.1792,
		[3]float64{50.0000, 2.4900, -0.0010},
		[3]float64{50.0000, -2.4900, 0.0009})
	assertDeltaE(7.1792,
		[3]float64{50.0000, 2.4900, -0.0010},
		[3]float64{50.0000, -2.4900, 0.0010})
	assertDeltaE(7.2195,
		[3]float64{50.0000, 2.4900, -0.0010},
		[3]float64{50.0000, -2.4900, 0.0011})
	assertDeltaE(7.2195,
		[3]float64{50.0000, 2.4900, -0.0010},
		[3]float64{50.0000, -2.4900, 0.0012})
	assertDeltaE(4.8045,
		[3]float64{50.0000, -0.0010, 2.4900},
		[3]float64{50.0000, 0.0009, -2.4900})
	assertDeltaE(4.7461,
		[3]float64{50.0000, -0.0010, 2.4900},
		[3]float64{50.0000, 0.0011, -2.4900})
	assertDeltaE(4.3065,
		[3]float64{50.0000, 2.5000, 0.0000},
		[3]float64{50.0000, 0.0000, -2.5000})
	assertDeltaE(27.1492,
		[3]float64{50.0000, 2.5000, 0.0000},
		[3]float64{73.0000, 25.0000, -18.0000})
	assertDeltaE(22.8977,
		[3]float64{50.0000, 2.5000, 0.0000},
		[3]float64{61.0000, -5.0000, 29.0000})
	assertDeltaE(31.9030,
		[3]float64{50.0000, 2.5000, 0.0000},
		[3]float64{56.0000, -27.0000, -3.0000})
	assertDeltaE(19.4535,
		[3]float64{50.0000, 2.5000, 0.0000},
		[3]float64{58.0000, 24.0000, 15.0000})
	assertDeltaE(1.0000,
		[3]float64{50.0000, 2.5000, 0.0000},
		[3]float64{50.0000, 3.1736, 0.5854})
	assertDeltaE(1.0000,
		[3]float64{50.0000, 2.5000, 0.0000},
		[3]float64{50.0000, 3.2972, 0.0000})
	assertDeltaE(1.0000,
		[3]float64{50.0000, 2.5000, 0.0000},
		[3]float64{50.0000, 1.8634, 0.5757})
	assertDeltaE(1.0000,
		[3]float64{50.0000, 2.5000, 0.0000},
		[3]float64{50.0000, 3.2592, 0.3350})
	assertDeltaE(1.2644,
		[3]float64{60.2574, -34.0099, 36.2677},
		[3]float64{60.4626, -34.1751, 39.4387})
	assertDeltaE(1.2630,
		[3]float64{63.0109, -31.0961, -5.8663},
		[3]float64{62.8187, -29.7946, -4.0864})
	assertDeltaE(1.8731,
		[3]float64{61.2901, 3.7196, -5.3901},
		[3]float64{61.4292, 2.2480, -4.9620})
	assertDeltaE(1.8645,
		[3]float64{35.0831, -44.1164, 3.7933},
		[3]float64{35.0232, -40.0716, 1.5901})
	assertDeltaE(2.0373,
		[3]float64{22.7233, 20.0904, -46.6940},
		[3]float64{23.0331, 14.9730, -42.5619})
	assertDeltaE(1.4146,
		[3]float64{36.4612, 47.8580, 18.3852},
		[3]float64{36.2715, 50.5065, 21.2231})
	assertDeltaE(1.4441,
		[3]float64{90.8027, -2.0831, 1.4410},
		[3]float64{91.1528, -1.6435, 0.0447})
	assertDeltaE(1.5381,
		[3]float64{90.9257, -0.5406, -0.9208},
		[3]float64{88.6381, -0.8985, -0.7239})
	assertDeltaE(0.6377,
		[3]float64{6.7747, -0.2908, -2.4247},
		[3]float64{5.8714, -0.0985, -2.2286})
	assertDeltaE(0.9082,
		[3]float64{2.0776, 0.0795, -1.1350},
		[3]float64{0.9033, -0.0636, -0.5514})
}
