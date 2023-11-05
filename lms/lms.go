package lms

import (
	"image"
	"image/color"

	"gonum.org/v1/netlib/blas/netlib"
)

func init() {
	_ = netlib.Implementation{}
}

type LMSImage struct {
	Pix    []float64
	Stride int
	Rect   image.Rectangle
}

func (p *LMSImage) ColorModel() color.Model { return &LMSModel{} }

func (p *LMSImage) Bounds() image.Rectangle { return p.Rect }

func (p *LMSImage) At(x, y int) color.Color {
	return p.LMSAt(x, y)
}

func (p *LMSImage) LMSAt(x, y int) LMS {
	if !(image.Point{x, y}.In(p.Rect)) {
		return LMS{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // we store 3 float64's (each 8 bytes)

//		L := float64(uint64(s[0])<<(8*7)|
//			uint64(s[1])<<(8*6)|
//			uint64(s[2])<<(8*5)|
//			uint64(s[3])<<(8*4)|
//			uint64(s[4])<<(8*3)|
//			uint64(s[5])<<(8*2)|
//			uint64(s[6])<<8|
//			uint64(s[7])) / 65535.0
//
//		M := float64(uint64(s[8])<<(8*7)|
//			uint64(s[9])<<(8*6)|
//			uint64(s[10])<<(8*5)|
//			uint64(s[11])<<(8*4)|
//			uint64(s[12])<<(8*3)|
//			uint64(s[13])<<(8*2)|
//			uint64(s[14])<<8|
//			uint64(s[15])) / 65535.0
//
//		S := float64(uint64(s[16])<<(8*7)|
//			uint64(s[17])<<(8*6)|
//			uint64(s[18])<<(8*5)|
//			uint64(s[19])<<(8*4)|
//			uint64(s[20])<<(8*3)|
//			uint64(s[21])<<(8*2)|
//			uint64(s[22])<<8|
//			uint64(s[23])) / 65535.0

	L := s[0]
	M := s[1]
	S := s[2]

	return LMS{L, M, S}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *LMSImage) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3 // we store 3 float64's (each 8 bytes)
}

func (p *LMSImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := LMSModel{}.Convert(c).(LMS)
	s := p.Pix[i : i+3 : i+3] // we store 3 float64's (each 8 bytes)

//		lu64 := uint64(c1.L * 65535.0)
//		mu64 := uint64(c1.M * 65535.0)
//		su64 := uint64(c1.S * 65535.0)
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
//		s[8] = uint8(mu64 >> (8 * 7))
//		s[9] = uint8(mu64 >> (8 * 6))
//		s[10] = uint8(mu64 >> (8 * 5))
//		s[11] = uint8(mu64 >> (8 * 4))
//		s[12] = uint8(mu64 >> (8 * 3))
//		s[13] = uint8(mu64 >> (8 * 2))
//		s[14] = uint8(mu64 >> 8)
//		s[15] = uint8(mu64)
//
//		s[16] = uint8(su64 >> (8 * 7))
//		s[17] = uint8(su64 >> (8 * 6))
//		s[18] = uint8(su64 >> (8 * 5))
//		s[19] = uint8(su64 >> (8 * 4))
//		s[20] = uint8(su64 >> (8 * 3))
//		s[21] = uint8(su64 >> (8 * 2))
//		s[22] = uint8(su64 >> 8)
//		s[23] = uint8(su64)

	s[0] = c1.L
	s[1] = c1.M
	s[2] = c1.S
}

func (p *LMSImage) SetLMS(x, y int, c LMS) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // we store 3 float64's (each 8 bytes)

//		lu64 := uint64(c.L * 65535.0)
//		mu64 := uint64(c.M * 65535.0)
//		su64 := uint64(c.S * 65535.0)
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
//		s[8] = uint8(mu64 >> (8 * 7))
//		s[9] = uint8(mu64 >> (8 * 6))
//		s[10] = uint8(mu64 >> (8 * 5))
//		s[11] = uint8(mu64 >> (8 * 4))
//		s[12] = uint8(mu64 >> (8 * 3))
//		s[13] = uint8(mu64 >> (8 * 2))
//		s[14] = uint8(mu64 >> 8)
//		s[15] = uint8(mu64)
//
//		s[16] = uint8(su64 >> (8 * 7))
//		s[17] = uint8(su64 >> (8 * 6))
//		s[18] = uint8(su64 >> (8 * 5))
//		s[19] = uint8(su64 >> (8 * 4))
//		s[20] = uint8(su64 >> (8 * 3))
//		s[21] = uint8(su64 >> (8 * 2))
//		s[22] = uint8(su64 >> 8)
//		s[23] = uint8(su64)

	s[0] = c.L
	s[1] = c.M
	s[2] = c.S
}

// NewLMSImage returns a new LMSImage image with the given bounds.
func NewLMSImage(r image.Rectangle) *LMSImage {
	w, h := r.Dx(), r.Dy()
	pix := make([]float64, 3*w*h)       // we store 3 float64's (each 8 bytes)
	return &LMSImage{pix, 3 * w, r} // we store 3 float64's (each 8 bytes)
}

func ConvertToLMS(in image.Image) *LMSImage {
	bs := in.Bounds()
	out := NewLMSImage(bs)

	for y := bs.Min.Y; y < bs.Max.Y; y++ {
		for x := bs.Min.X; x < bs.Max.X; x++ {
			out.Set(x, y, in.At(x, y))
		}
	}

	return out
}
