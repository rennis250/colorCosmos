package xyy

import (
	"image"
	"image/color"

	"gonum.org/v1/netlib/blas/netlib"
)

func init() {
	_ = netlib.Implementation{}
}

type XyYImage struct {
	Pix    []float64
	Stride int
	Rect   image.Rectangle
}

func (p *XyYImage) ColorModel() color.Model { return &XyYModel{} }

func (p *XyYImage) Bounds() image.Rectangle { return p.Rect }

func (p *XyYImage) At(x, y int) color.Color {
	return p.XyYAt(x, y)
}

func (p *XyYImage) XyYAt(x, y int) XyY {
	if !(image.Point{x, y}.In(p.Rect)) {
		return XyY{}
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

	//		return XyY{(ld - 0.5) * 2.0, (rg - 0.5) * 2.0, (yv - 0.5) * 2.0}

	x := s[0]
	y := s[1]
	Y := s[2]

	return XyY{x, y, Y}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *XyYImage) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3 // we store 3 float64's (each 8 bytes)
}

func (p *XyYImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := XyYModel{}.Convert(c).(XyY)
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

	s[0] = c1.x
	s[1] = c1.y
	s[2] = c1.Y
}

func (p *XyYImage) SetXyY(x, y int, c XyY) {
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

	s[0] = c.x
	s[1] = c.y
	s[2] = c.Y
}

// NewXyYImage returns a new XyYImage image with the given bounds.
func NewXyYImage(r image.Rectangle) *XyYImage {
	w, h := r.Dx(), r.Dy()
	pix := make([]float64, 3*w*h)   // we store 3 float64's (each 8 bytes)
	return &XyYImage{pix, 3 * w, r} // we store 3 float64's (each 8 bytes)
}

func ConvertToXyY(in image.Image) *XyYImage {
	bs := in.Bounds()
	out := NewXyYImage(bs)

	for y := bs.Min.Y; y < bs.Max.Y; y++ {
		for x := bs.Min.X; x < bs.Max.X; x++ {
			out.Set(x, y, in.At(x, y))
		}
	}

	return out
}
