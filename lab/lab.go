package lab

import (
	"image"
	"image/color"

	"gonum.org/v1/netlib/blas/netlib"
)

func init() {
	_ = netlib.Implementation{}
}

type LABImage struct {
	Pix    []float64
	Stride int
	Rect   image.Rectangle
}

func (p *LABImage) ColorModel() color.Model { return &LABModel{} }

func (p *LABImage) Bounds() image.Rectangle { return p.Rect }

func (p *LABImage) At(x, y int) color.Color {
	return p.LABAt(x, y)
}

func (p *LABImage) LABAt(x, y int) LAB {
	if !(image.Point{x, y}.In(p.Rect)) {
		return LAB{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // we store 3 float64's (each 8 bytes)

//		l := float64(uint64(s[0])<<(8*7)|
//			uint64(s[1])<<(8*6)|
//			uint64(s[2])<<(8*5)|
//			uint64(s[3])<<(8*4)|
//			uint64(s[4])<<(8*3)|
//			uint64(s[5])<<(8*2)|
//			uint64(s[6])<<8|
//			uint64(s[7])) / 65535.0
//
//		a := float64(uint64(s[8])<<(8*7)|
//			uint64(s[9])<<(8*6)|
//			uint64(s[10])<<(8*5)|
//			uint64(s[11])<<(8*4)|
//			uint64(s[12])<<(8*3)|
//			uint64(s[13])<<(8*2)|
//			uint64(s[14])<<8|
//			uint64(s[15])) / 65535.0
//
//		b := float64(uint64(s[16])<<(8*7)|
//			uint64(s[17])<<(8*6)|
//			uint64(s[18])<<(8*5)|
//			uint64(s[19])<<(8*4)|
//			uint64(s[20])<<(8*3)|
//			uint64(s[21])<<(8*2)|
//			uint64(s[22])<<8|
//			uint64(s[23])) / 65535.0

	L := s[0]
	A := s[1]
	B := s[2]

	return LAB{L, A, B}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *LABImage) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3 // we store 3 float64's (each 8 bytes)
}

func (p *LABImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := LABModel{}.Convert(c).(LAB)
	s := p.Pix[i : i+3 : i+3] // we store 3 float64's (each 8 bytes)

//		lu64 := uint64(c1.L * 65535.0)
//		au64 := uint64(c1.A * 65535.0)
//		bu64 := uint64(c1.B * 65535.0)
//
//		s[0] = uint8(lu64 >> (8 * 7))
//		s[1] = uint8(lu64 >> (8 * 6))
//		s[2] = uint8(lu64 >> (8 * 5))
//		s[3] = uint8(lu64 >> (8 * 4))
//		s[4] = uint8(lu64 >> (8 * 3))
//		s[5] = uint8(lu64 >> (8 * 2))
//		s[6] = uint8(lu64 >> 8)
//		s[7] = uint8(lu64)
//
//		s[8] = uint8(au64 >> (8 * 7))
//		s[9] = uint8(au64 >> (8 * 6))
//		s[10] = uint8(au64 >> (8 * 5))
//		s[11] = uint8(au64 >> (8 * 4))
//		s[12] = uint8(au64 >> (8 * 3))
//		s[13] = uint8(au64 >> (8 * 2))
//		s[14] = uint8(au64 >> 8)
//		s[15] = uint8(au64)
//
//		s[16] = uint8(bu64 >> (8 * 7))
//		s[17] = uint8(bu64 >> (8 * 6))
//		s[18] = uint8(bu64 >> (8 * 5))
//		s[19] = uint8(bu64 >> (8 * 4))
//		s[20] = uint8(bu64 >> (8 * 3))
//		s[21] = uint8(bu64 >> (8 * 2))
//		s[22] = uint8(bu64 >> 8)
//		s[23] = uint8(bu64)

	s[0] = c1.L
	s[1] = c1.A
	s[2] = c1.B
}

func (p *LABImage) SetLAB(x, y int, c LAB) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // we store 3 float64's (each 8 bytes)

//		lu64 := uint64(c.L * 65535.0)
//		au64 := uint64(c.A * 65535.0)
//		bu64 := uint64(c.B * 65535.0)
//
//		s[0] = uint8(lu64 >> (8 * 7))
//		s[1] = uint8(lu64 >> (8 * 6))
//		s[2] = uint8(lu64 >> (8 * 5))
//		s[3] = uint8(lu64 >> (8 * 4))
//		s[4] = uint8(lu64 >> (8 * 3))
//		s[5] = uint8(lu64 >> (8 * 2))
//		s[6] = uint8(lu64 >> 8)
//		s[7] = uint8(lu64)
//
//		s[8] = uint8(au64 >> (8 * 7))
//		s[9] = uint8(au64 >> (8 * 6))
//		s[10] = uint8(au64 >> (8 * 5))
//		s[11] = uint8(au64 >> (8 * 4))
//		s[12] = uint8(au64 >> (8 * 3))
//		s[13] = uint8(au64 >> (8 * 2))
//		s[14] = uint8(au64 >> 8)
//		s[15] = uint8(au64)
//
//		s[16] = uint8(bu64 >> (8 * 7))
//		s[17] = uint8(bu64 >> (8 * 6))
//		s[18] = uint8(bu64 >> (8 * 5))
//		s[19] = uint8(bu64 >> (8 * 4))
//		s[20] = uint8(bu64 >> (8 * 3))
//		s[21] = uint8(bu64 >> (8 * 2))
//		s[22] = uint8(bu64 >> 8)
//		s[23] = uint8(bu64)

	s[0] = c.L
	s[1] = c.A
	s[2] = c.B
}

// NewLABImage returns a new LABImage image with the given bounds.
func NewLABImage(r image.Rectangle) *LABImage {
	w, h := r.Dx(), r.Dy()
	pix := make([]float64, 3*w*h)       // we store 3 float64's (each 8 bytes)
	return &LABImage{pix, 3 * w, r} // we store 3 float64's (each 8 bytes)
}

func ConvertToLAB(in image.Image) *LABImage {
	bs := in.Bounds()
	out := NewLABImage(bs)

	for y := bs.Min.Y; y < bs.Max.Y; y++ {
		for x := bs.Min.X; x < bs.Max.X; x++ {
			out.Set(x, y, in.At(x, y))
		}
	}

	return out
}
