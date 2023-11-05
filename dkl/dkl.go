package dkl

import (
	"image"
	"image/color"

	"gonum.org/v1/netlib/blas/netlib"
)

func init() {
	_ = netlib.Implementation{}
}

type DKLImage struct {
	Pix    []float64
	Stride int
	Rect   image.Rectangle
}

func (p *DKLImage) ColorModel() color.Model { return &DKLModel{} }

func (p *DKLImage) Bounds() image.Rectangle { return p.Rect }

func (p *DKLImage) At(x, y int) color.Color {
	return p.DKLAt(x, y)
}

func (p *DKLImage) DKLAt(x, y int) DKL {
	if !(image.Point{x, y}.In(p.Rect)) {
		return DKL{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // we store 3 float64's (each 8 bytes)

//		ld := float64(uint64(s[0])<<(8*7)|
//			uint64(s[1])<<(8*6)|
//			uint64(s[2])<<(8*5)|
//			uint64(s[3])<<(8*4)|
//			uint64(s[4])<<(8*3)|
//			uint64(s[5])<<(8*2)|
//			uint64(s[6])<<8|
//			uint64(s[7])) / 65535.0
//
//		rg := float64(uint64(s[8])<<(8*7)|
//			uint64(s[9])<<(8*6)|
//			uint64(s[10])<<(8*5)|
//			uint64(s[11])<<(8*4)|
//			uint64(s[12])<<(8*3)|
//			uint64(s[13])<<(8*2)|
//			uint64(s[14])<<8|
//			uint64(s[15])) / 65535.0
//
//		yv := float64(uint64(s[16])<<(8*7)|
//			uint64(s[17])<<(8*6)|
//			uint64(s[18])<<(8*5)|
//			uint64(s[19])<<(8*4)|
//			uint64(s[20])<<(8*3)|
//			uint64(s[21])<<(8*2)|
//			uint64(s[22])<<8|
//			uint64(s[23])) / 65535.0

//		return DKL{(ld - 0.5) * 2.0, (rg - 0.5) * 2.0, (yv - 0.5) * 2.0}

	ld := s[0]
	rg := s[1]
	yv := s[2]

	return DKL{ld, rg, yv}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *DKLImage) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3 // we store 3 float64's (each 8 bytes)
}

func (p *DKLImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := DKLModel{}.Convert(c).(DKL)
	s := p.Pix[i : i+3 : i+3] // we store 3 float64's (each 8 bytes)

//		ldu64 := uint64((c1.LD/2.0 + 0.5) * 65535.0)
//		rgu64 := uint64((c1.RG/2.0 + 0.5) * 65535.0)
//		yvu64 := uint64((c1.YV/2.0 + 0.5) * 65535.0)

//		s[0] = uint8(ldu64 >> (8 * 7))
//		s[1] = uint8(ldu64 >> (8 * 6))
//		s[2] = uint8(ldu64 >> (8 * 5))
//		s[3] = uint8(ldu64 >> (8 * 4))
//		s[4] = uint8(ldu64 >> (8 * 3))
//		s[5] = uint8(ldu64 >> (8 * 2))
//		s[6] = uint8(ldu64 >> 8)
//		s[7] = uint8(ldu64)
//
//		s[8] = uint8(rgu64 >> (8 * 7))
//		s[9] = uint8(rgu64 >> (8 * 6))
//		s[10] = uint8(rgu64 >> (8 * 5))
//		s[11] = uint8(rgu64 >> (8 * 4))
//		s[12] = uint8(rgu64 >> (8 * 3))
//		s[13] = uint8(rgu64 >> (8 * 2))
//		s[14] = uint8(rgu64 >> 8)
//		s[15] = uint8(rgu64)
//
//		s[16] = uint8(yvu64 >> (8 * 7))
//		s[17] = uint8(yvu64 >> (8 * 6))
//		s[18] = uint8(yvu64 >> (8 * 5))
//		s[19] = uint8(yvu64 >> (8 * 4))
//		s[20] = uint8(yvu64 >> (8 * 3))
//		s[21] = uint8(yvu64 >> (8 * 2))
//		s[22] = uint8(yvu64 >> 8)
//		s[23] = uint8(yvu64)

	s[0] = c1.LD
	s[1] = c1.RG
	s[2] = c1.YV
}

func (p *DKLImage) SetDKL(x, y int, c DKL) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // we store 3 float64's (each 8 bytes)

//		ldu64 := uint64((c.LD/2.0 + 0.5) * 65535.0)
//		rgu64 := uint64((c.RG/2.0 + 0.5) * 65535.0)
//		yvu64 := uint64((c.YV/2.0 + 0.5) * 65535.0)
//
//		s[0] = uint8(ldu64 >> (8 * 7))
//		s[1] = uint8(ldu64 >> (8 * 6))
//		s[2] = uint8(ldu64 >> (8 * 5))
//		s[3] = uint8(ldu64 >> (8 * 4))
//		s[4] = uint8(ldu64 >> (8 * 3))
//		s[5] = uint8(ldu64 >> (8 * 2))
//		s[6] = uint8(ldu64 >> 8)
//		s[7] = uint8(ldu64)
//
//		s[8] = uint8(rgu64 >> (8 * 7))
//		s[9] = uint8(rgu64 >> (8 * 6))
//		s[10] = uint8(rgu64 >> (8 * 5))
//		s[11] = uint8(rgu64 >> (8 * 4))
//		s[12] = uint8(rgu64 >> (8 * 3))
//		s[13] = uint8(rgu64 >> (8 * 2))
//		s[14] = uint8(rgu64 >> 8)
//		s[15] = uint8(rgu64)
//
//		s[16] = uint8(yvu64 >> (8 * 7))
//		s[17] = uint8(yvu64 >> (8 * 6))
//		s[18] = uint8(yvu64 >> (8 * 5))
//		s[19] = uint8(yvu64 >> (8 * 4))
//		s[20] = uint8(yvu64 >> (8 * 3))
//		s[21] = uint8(yvu64 >> (8 * 2))
//		s[22] = uint8(yvu64 >> 8)
//		s[23] = uint8(yvu64)

	s[0] = c.LD
	s[1] = c.RG
	s[2] = c.YV
}

// NewDKLImage returns a new DKLImage image with the given bounds.
func NewDKLImage(r image.Rectangle) *DKLImage {
	w, h := r.Dx(), r.Dy()
	pix := make([]float64, 3*w*h)       // we store 3 float64's (each 8 bytes)
	return &DKLImage{pix, 3 * w, r} // we store 3 float64's (each 8 bytes)
}

func ConvertToDKL(in image.Image) *DKLImage {
	bs := in.Bounds()
	out := NewDKLImage(bs)

	for y := bs.Min.Y; y < bs.Max.Y; y++ {
		for x := bs.Min.X; x < bs.Max.X; x++ {
			out.Set(x, y, in.At(x, y))
		}
	}

	return out
}