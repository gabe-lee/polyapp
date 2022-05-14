package polyapp

import (
	color "github.com/gabe-lee/gencolor"
	vecs "github.com/gabe-lee/genvecs"
)

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
type Color32 = color.U8x4Color

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

type Shader = poly.Shader
type Texture = poly.Tex