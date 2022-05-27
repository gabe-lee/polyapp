package polyapp

import (
	color "github.com/gabe-lee/color"
	math "github.com/gabe-lee/genmath"
	utils "github.com/gabe-lee/genutils"
	vecs "github.com/gabe-lee/genvecs"
)

var ZeroVec3 = Vec3{0, 0, 0}
var ZeroVec2 = Vec2{0, 0}

type DeepError = utils.DeepError

type Vec2 = vecs.F32Vec2
type IVec2 = vecs.I32Vec2
type Vec3 = vecs.F32Vec3
type IVec3 = vecs.I32Vec3
type Rect2D = vecs.F32AABB2
type IRect2D = vecs.I32AABB2
type Rect3D = vecs.F32AABB3
type IRect3D = vecs.I32AABB3
type Quad2D = vecs.F32Quad2
type IQuad2D = vecs.I32Quad2
type Quad3D = vecs.F32Quad3
type IQuad3D = vecs.I32Quad3
type Line2D = vecs.F32Line2
type Line3D = vecs.F32Line3

type VertExtra = [8]uint32

type ColorFA = color.ColorFA
type ColorF = color.ColorF
type Color64 = color.Color64
type Color48 = color.Color48
type Color32 = color.Color32
type Color24 = color.Color24
type Color16 = color.Color16
type Color8 = color.Color8

type InputState uint8

const (
	UpPosition InputState = iota
	DownPosition
)

type InputAction uint8

const (
	InputUntouched InputAction = iota
	InputPressed
	InputHeld
	InputReleased
	InputHeldRepeat
)

type ImageType uint8

const (
	ImgUnknown ImageType = iota
	ImgPNG
	ImgBMP
	ImgWEBP
)

type BufferZone struct {
	Start uint32
	End   uint32
}

func (b BufferZone) Len() uint32 {
	return b.End - b.Start
}

type BufferZoneLL struct {
	BufferZone
	Next *BufferZoneLL
}

func (b *BufferZoneLL) Insert(zone BufferZone) {
	overlap, start, end := math.CombineRangesIfOverlap(b.Start, b.End, zone.Start, zone.End)
	if overlap {
		b.Start = start
		b.End = end
	} else if b.Next != nil {
		b.Next.Insert(zone)
	} else {
		b.Next = &BufferZoneLL{BufferZone: zone}
	}
	if b.Next != nil {
		overlap, start, end = math.CombineRangesIfOverlap(b.Start, b.End, b.Next.Start, b.Next.End)
		if overlap {
			b.Start = start
			b.End = end
			b.Next = b.Next.Next
		}
	}
}

func (b *BufferZoneLL) Aquire(zoneSize uint32, last *BufferZoneLL) BufferZone {
	if zoneSize <= b.Len() {
		zone := BufferZone{Start: b.Start, End: b.Start + zoneSize}
		b.Start = zone.End
		if b.Len() == 0 && last != nil {
			last.Next = b.Next
		}
		return zone
	}
	if b.Next != nil {
		return b.Next.Aquire(zoneSize, b)
	}
	return BufferZone{}
}
