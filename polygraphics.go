package polyapp

import (
	geom "github.com/gabe-lee/gengeom"
	math "github.com/gabe-lee/genmath"
	utils "github.com/gabe-lee/genutils"
)

type GraphicsInterface interface {
	XRightYUpZAway() Vec3

	AddRenderer(vertexFlags VertexFlags, shaders []*Shader) (rendererID uint8, err DeepError)
	AddDrawBatch(vertexFlags VertexFlags, textureID uint8, initialSize uint32) (batchID uint8, err DeepError)
	AddTexture(texture *Texture) (textureID uint8, err DeepError)
	AddDrawSurface(size IVec2, mipMaps uint32) (surfaceID uint8, textureID uint8, err DeepError)

	ClearSurface(surfaceID uint8, baseColor ColorFA) DeepError
	ClearSurfaceArea(surfaceID uint8, baseColor ColorFA, area IRect2D) DeepError

	AllocateShapeInBatch(batchID uint8, shape BatchShape) (BatchSlice, DeepError)
	UpdateVertexInBatch(batchID uint8, vertIndex uint32, vertex Vertex) DeepError
	DeleteShapeFromBatch(batchSlice BatchSlice) DeepError

	DrawBatch(batchID uint8, surfaceID uint8, rendererID uint8) DeepError
	ClearBatch(batchID uint8) DeepError
}

var _ GraphicsInterface = (*GraphicsProvider)(nil)

type GraphicsProvider struct {
	GraphicsInterface
}

var ninf = math.NInf32()

var NoUV = Vec2{ninf, ninf}
var NoVert = Vec3{ninf, ninf, ninf}
var NoColor = ColorFA{ninf, ninf, ninf, ninf}
var NoExtra = VertExtra{}
var NoNorm = Vec3{ninf, ninf, ninf}

// Describes the types of vertex attributes present on a draw batch or renderer.
//
// Zero value defaults to: 2D Positions + 16 bit indexes + Traingle draw mode + No texture + No Color + No Extra data blocks + No Camera
//
// Vertex attribute layout should follow this order: Position -> Normals -> UVs -> Color -> Extra
type VertexFlags uint16

const (
	Pos2D     VertexFlags = 0  // 2D Vertex space
	Pos3D     VertexFlags = 1  // 3D Vertex space
	PosMask   VertexFlags = 1  // Mask for checking vertex space
	Idx16     VertexFlags = 0  // Indexes are uint16
	Idx32     VertexFlags = 2  // Indexes ar uint32
	IdxMask   VertexFlags = 2  // Mask for checking index size
	NoTex     VertexFlags = 0  // No texture (no UV coordinates)
	HasTex    VertexFlags = 4  // Uses Texture with UV coordinates
	TexMask   VertexFlags = 4  // Mask for checking texture use
	NoCol     VertexFlags = 0  // No color channel (uniform color)
	Col8      VertexFlags = 8  // 2bit RGBA channels
	Col16     VertexFlags = 16 // 4bit RGBA channels
	Col24     VertexFlags = 24 // 8bit RGB channels
	Col32     VertexFlags = 32 // 8bit RGBA channels
	Col48     VertexFlags = 40 // 16bit RGB channels
	Col64     VertexFlags = 48 // 16bit RGBA channels
	ColF      VertexFlags = 56 // float32 RGB channels
	ColFA     VertexFlags = 64 // float32 RGBA channels
	_col9     VertexFlags = 72
	_col10    VertexFlags = 80
	_col11    VertexFlags = 88
	_col12    VertexFlags = 96
	_col13    VertexFlags = 104
	_col14    VertexFlags = 112
	_col15    VertexFlags = 120
	ColMask   VertexFlags = 120  // Mask for checking color mode
	NoEx      VertexFlags = 0    // No aditional 32bit data blocks
	Ex32      VertexFlags = 128  // 1 additional 32bit data block
	Ex64      VertexFlags = 256  // 2 additional 32bit data blocks
	Ex96      VertexFlags = 384  // 3 additional 32bit data blocks
	Ex128     VertexFlags = 512  // 4 additional 32bit data blocks
	Ex192     VertexFlags = 640  // 6 additional 32bit data blocks
	Ex224     VertexFlags = 768  // 7 additional 32bit data blocks
	Ex256     VertexFlags = 896  // 8 additional 32bit data blocks
	ExMask    VertexFlags = 896  // Mask for checking number of extra data blocks
	Tris      VertexFlags = 0    // Every 3 Vertices are an independant triangle
	Lines     VertexFlags = 1024 // Every 2 vertices are an independant line
	Pixels    VertexFlags = 2048 // Every vertex is an independant point
	_draw4    VertexFlags = 3072
	DrawMask  VertexFlags = 3072 // Mask for checking draw mode
	NoCam     VertexFlags = 0    // No Camera Projection (Draws as if draw surface IS the camera, no transform)
	Cam2D     VertexFlags = 4096 // 2D Camera projection
	Cam3D     VertexFlags = 8192 // 3D Camera projection
	_cam4D    VertexFlags = 12288
	CamMask   VertexFlags = 12288 // Mask for checking camera mode
	NoNorms   VertexFlags = 0     // No vertex Normals
	Norms     VertexFlags = 16384 // Includes Vertex normals
	NormsMask VertexFlags = 16384 // Mask for checking if uses vertex normals
	_un4      VertexFlags = 32768

	VertexAttributeMask  VertexFlags = PosMask | ColMask | IdxMask | TexMask | ExMask | NormsMask // Mask describing layout of vertex attributes and indexes
	UniformAttributeMask VertexFlags = CamMask | DrawMask                                         // Mask decribing rendering uniforms and draw mode
)

func (vf VertexFlags) SameAttributes(other VertexFlags) bool {
	return vf&VertexAttributeMask == other&VertexAttributeMask
}

func (vf VertexFlags) SameUniforms(other VertexFlags) bool {
	return vf&UniformAttributeMask == other&UniformAttributeMask
}

func (vf VertexFlags) PositionOffset() uint32 {
	return 0
}
func (vf VertexFlags) PositionSize() uint32 {
	if vf&PosMask == Pos3D {
		return 12
	}
	return 8
}

func (vf VertexFlags) NormalOffset() uint32 {
	return vf.PositionSize()
}
func (vf VertexFlags) NormalSize() uint32 {
	if vf&NormsMask == NormsMask {
		if vf&PosMask == Pos3D {
			return 12
		}
		return 8
	}
	return 0
}

func (vf VertexFlags) UVOffset() uint32 {
	return vf.NormalOffset() + vf.NormalSize()
}
func (vf VertexFlags) UVSize() uint32 {
	if vf&TexMask == HasTex {
		return 8
	}
	return 0
}

func (vf VertexFlags) ColorOffset() uint32 {
	return vf.UVOffset() + vf.UVSize()
}
func (vf VertexFlags) ColorSize() uint32 {
	switch {
	case vf&ColMask == ColFA:
		return 16
	case vf&ColMask == ColF:
		return 12
	case vf&ColMask == Col64:
		return 8
	case vf&ColMask == Col48:
		return 6
	case vf&ColMask == Col32:
		return 4
	case vf&ColMask == Col24:
		return 3
	case vf&ColMask == Col16:
		return 2
	case vf&ColMask == Col8:
		return 1
	default:
		return 0
	}
}

func (vf VertexFlags) ExOffset() uint32 {
	return vf.ColorOffset() + vf.ColorSize()
}
func (vf VertexFlags) ExSize() uint32 {
	switch {
	case vf&ExMask == Ex256:
		return 32
	case vf&ExMask == Ex224:
		return 28
	case vf&ExMask == Ex192:
		return 24
	case vf&ExMask == Ex128:
		return 16
	case vf&ExMask == Ex96:
		return 12
	case vf&ExMask == Ex64:
		return 8
	case vf&ExMask == Ex32:
		return 4
	default:
		return 0
	}
}

func (vf VertexFlags) Stride() uint32 {
	sum := uint32(0)
	sum += vf.PositionSize()
	sum += vf.NormalSize()
	sum += vf.UVSize()
	sum += vf.ColorSize()
	sum += vf.ExSize()
	return sum
}

type Vertex struct {
	Pos   Vec3
	Norm  Vec3
	UV    Vec2
	Color ColorFA
	Extra VertExtra
}

type ShaderType uint8

const (
	ShaderVertex ShaderType = iota
	ShaderTessControl
	ShaderTessEval
	ShaderGeometry
	ShaderFragment
	ShaderCompute
)

type Texture struct {
	Data    []byte
	File    string
	ImgType ImageType
	Size    IVec2
	MipMaps uint32
	ID      uint32
	TexUnit uint32
}

type Shader struct {
	SType ShaderType
	Code  string
	Data  []byte
	File  string
}

type BatchSlice struct {
	BatchID     uint8
	IndexStart  uint32
	IndexEnd    uint32
	VertexStart uint32
	VertexEnd   uint32
}

type BatchShape struct {
	VertCount uint32
	Indexes   []uint32
}

func (b BatchSlice) IdxLen() uint32 {
	return b.IndexEnd - b.IndexStart
}

func (b BatchSlice) VertLen() uint32 {
	return b.VertexEnd - b.VertexStart
}

/**************
	LINES
***************/

func (g GraphicsProvider) AddLine2D(batchID uint8, a Vertex, b Vertex, thickness float32, uvThickness float32) (batchSlice BatchSlice, err DeepError) {
	dErr := utils.NewDeepError("[PolyApp] AddLine2D():")
	dErr.IsErr = false
	bSlice, err := g.AllocateShapeInBatch(batchID, BatchShape{
		VertCount: 4,
		Indexes:   []uint32{0, 1, 2, 0, 3, 2},
	})
	if err.IsErr {
		dErr.AddChildDeepError(err)
		return bSlice, err
	}
	dErr.AddChildDeepError(g.UpdateLine2D(bSlice, a, b, thickness, uvThickness))
	return batchSlice, dErr
}

func (g GraphicsProvider) UpdateLine2D(batchSlice BatchSlice, a Vertex, b Vertex, thickness float32, uvThickness float32) DeepError {
	if batchSlice.IdxLen() != 6 || batchSlice.VertLen() != 4 {
		return utils.NewDeepError("[PolyApp] UpdateLine2D(): batch slice provided does not have required dimensions for a line")
	}
	dErr := utils.NewDeepError("[PolyApp] UpdateLine2D():")
	dErr.IsErr = false
	a.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	b.Norm = a.Norm
	l := Line2D{a.Pos.AsVec2(), b.Pos.AsVec2()}
	u := Line2D{a.UV, b.UV}
	l1, l2 := l.PerpLines(thickness / 2)
	u1, u2 := u.PerpLines(uvThickness / 2)
	a.Pos = l1.A().AsVec3()
	a.UV = u1.A()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart, a))
	a.Pos = l2.A().AsVec3()
	a.UV = u2.A()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+1, a))
	b.Pos = l1.B().AsVec3()
	b.UV = u1.B()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+2, a))
	b.Pos = l2.B().AsVec3()
	b.UV = u2.B()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+3, a))
	return dErr
}

/**************
	TRIANGLES
***************/

func (g GraphicsProvider) AddTriangle2D(batchID uint8, a Vertex, b Vertex, c Vertex) (BatchSlice, DeepError) {
	dErr := utils.NewDeepError("[PolyApp] AddTriangle2D():")
	dErr.IsErr = false
	bSlice, err := g.AllocateShapeInBatch(batchID, BatchShape{
		VertCount: 3,
		Indexes:   []uint32{0, 1, 2},
	})
	if err.IsErr {
		dErr.AddChildDeepError(err)
		return bSlice, dErr
	}
	a.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	b.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	c.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	err.AddChildDeepError(g.UpdateTriangle2D(bSlice, a, b, c))
	return bSlice, err
}

func (g GraphicsProvider) UpdateTriangle2D(batchSlice BatchSlice, a Vertex, b Vertex, c Vertex) DeepError {
	if batchSlice.VertLen() != 3 || batchSlice.IdxLen() != 3 {
		return utils.NewDeepError("[PolyApp] UpdateTriangle2D(): batch slice provided does not have required dimensions for a triangle")
	}
	dErr := utils.NewDeepError("[PolyApp] UpdateTriangle2D():")
	dErr.IsErr = false
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart, a))
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+1, b))
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+2, c))
	return dErr
}

/**************
	POLYGONS
***************/

func (g GraphicsProvider) AddRegularPolygon2D(batchID uint8, center Vertex, sides uint32, radius float32, shapeRotation float32, uvRadius float32, uvRotation float32) (BatchSlice, DeepError) {
	dErr := utils.NewDeepError("[PolyApp] AddRegularPolygon2D():")
	dErr.IsErr = false
	iCount := 3 * sides
	vCount := sides + 1
	idx := make([]uint32, iCount)
	_ = idx[iCount-1]
	for i, v := uint32(0), uint32(1); i < iCount; i, v = i+3, v+1 {
		idx[i] = 0
		idx[i+1] = v
		idx[i+2] = v + 1
	}
	idx[iCount-1] = 1
	bSlice, err := g.AllocateShapeInBatch(batchID, BatchShape{
		VertCount: vCount,
		Indexes:   idx,
	})
	if err.IsErr {
		dErr.AddChildDeepError(err)
		return bSlice, dErr
	}
	dErr.AddChildDeepError(g.UpdateRegularPolygon2D(bSlice, center, sides, radius, shapeRotation, uvRadius, uvRotation))
	return bSlice, dErr
}

func (g GraphicsProvider) UpdateRegularPolygon2D(batchSlice BatchSlice, center Vertex, sides uint32, radius float32, shapeRotation float32, uvRadius float32, uvRotation float32) DeepError {
	if batchSlice.VertLen() != sides+1 || batchSlice.IdxLen() != sides*3 {
		return utils.NewDeepError("[PolyApp] UpdateRegularPolygon2D(): batch slice provided does not have required dimensions for a polygon of specified sides")
	}
	dErr := utils.NewDeepError("[PolyApp] UpdateRegularPolygon2D():")
	dErr.IsErr = false
	center.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	points := geom.PointsOnCircle(shapeRotation*math.DEG_TO_RAD, radius, center.Pos.AsVec2(), sides)
	uvs := geom.PointsOnCircle(uvRotation*math.DEG_TO_RAD, uvRadius, center.UV, sides)
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart, center))
	for i := uint32(0); i < uint32(len(points)); i += 1 {
		center.Pos = points[i].AsVec3()
		center.UV = uvs[i]
		dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+i+1, center))
	}
	return dErr
}

func (g GraphicsProvider) AddRegularPolygonRing2D(batchID uint8, center Vertex, sides uint32, innerRadius float32, outerRadius float32, shapeRotation float32, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32) (BatchSlice, DeepError) {
	dErr := utils.NewDeepError("[PolyApp] AddRegularPolygonRing2D():")
	dErr.IsErr = false
	iCount := 6 * sides
	vCount := 2 * sides
	idx := make([]uint32, iCount)
	_ = idx[iCount-1]
	for i, v := uint32(0), uint32(0); i < iCount; i, v = i+6, v+2 {
		idx[i] = v
		idx[i+1] = v + 1
		idx[i+2] = v + 2
		idx[i+3] = v
		idx[i+4] = v + 3
		idx[i+5] = v + 2
	}
	idx[iCount-2] = 1
	idx[iCount-1] = 0
	bSlice, err := g.AllocateShapeInBatch(batchID, BatchShape{
		VertCount: vCount,
		Indexes:   idx,
	})
	if err.IsErr {
		dErr.AddChildDeepError(err)
		return bSlice, err
	}
	dErr.AddChildDeepError(g.UpdateRegularPolygonRing2D(bSlice, center, sides, innerRadius, outerRadius, shapeRotation, uvInnerRadius, uvOuterRadius, uvRotation))
	return bSlice, err
}

func (g GraphicsProvider) UpdateRegularPolygonRing2D(batchSlice BatchSlice, center Vertex, sides uint32, innerRadius float32, outerRadius float32, shapeRotation float32, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32) DeepError {
	if batchSlice.VertLen() != sides*2 || batchSlice.IdxLen() != sides*6 {
		return utils.NewDeepError("[PolyApp] UpdateRegularPolygonRing2D(): batch slice provided does not have required dimensions for a polygon of specified sides")
	}
	dErr := utils.NewDeepError("[PolyApp] UpdateRegularPolygonRing2D():")
	dErr.IsErr = false
	center.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	uvs := geom.PointsOnRing(uvRotation*math.DEG_TO_RAD, uvInnerRadius, uvOuterRadius, center.UV, sides)
	points := geom.PointsOnRing(shapeRotation*math.DEG_TO_RAD, innerRadius, outerRadius, center.Pos.AsVec2(), sides)
	for i := uint32(0); i < uint32(len(points)); i += 1 {
		center.Pos = points[i].AsVec3()
		center.UV = uvs[i]
		dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+i, center))
	}
	return dErr
}

/**************
	CIRCLES
***************/

func (g GraphicsProvider) AddCircleAutoPoints2D(batchID uint8, center Vertex, resolution float32, radius float32, uvRadius float32, uvRotation float32) (BatchSlice, DeepError) {
	dErr := utils.NewDeepError("[PolyApp] AddCircleAutoPoints2D():")
	dErr.IsErr = false
	sides := uint32(math.Ciel(geom.Circumference(radius) / resolution))
	bs, err := g.AddRegularPolygon2D(batchID, center, sides, radius, 0, uvRadius, uvRotation)
	dErr.AddChildDeepError(err)
	return bs, dErr
}
func (g GraphicsProvider) UpdateCircleAutoPoints2D(batchSlice BatchSlice, center Vertex, resolution float32, radius float32, uvRadius float32, uvRotation float32) DeepError {
	dErr := utils.NewDeepError("[PolyApp] UpdateCircleAutoPoints2D():")
	dErr.IsErr = false
	sides := uint32(math.Ciel(geom.Circumference(radius) / resolution))
	dErr.AddChildDeepError(g.UpdateRegularPolygon2D(batchSlice, center, sides, radius, 0, uvRadius, uvRotation))
	return dErr
}
func (g GraphicsProvider) AddCircleRingAutoPoints2D(batchID uint8, center Vertex, resolution float32, innerRadius float32, outerRadius float32, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32) (BatchSlice, DeepError) {
	dErr := utils.NewDeepError("[PolyApp] AddCircleRingAutoPoints2D():")
	dErr.IsErr = false
	sides := uint32(math.Ciel(geom.Circumference(outerRadius) / resolution))
	bs, err := g.AddRegularPolygonRing2D(batchID, center, sides, innerRadius, outerRadius, 0, uvInnerRadius, uvOuterRadius, uvRotation)
	dErr.AddChildDeepError(err)
	return bs, dErr
}
func (g GraphicsProvider) UpdateCircleRingAutoPoints2D(batchSlice BatchSlice, center Vertex, resolution float32, innerRadius float32, outerRadius float32, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32) DeepError {
	dErr := utils.NewDeepError("[PolyApp] UpdateCircleRingAutoPoints2D():")
	dErr.IsErr = false
	sides := uint32(math.Ciel(geom.Circumference(outerRadius) / resolution))
	dErr.AddChildDeepError(g.UpdateRegularPolygonRing2D(batchSlice, center, sides, innerRadius, outerRadius, 0, uvInnerRadius, uvOuterRadius, uvRotation))
	return dErr
}

/**************
	RECTANGLES
***************/

func (g GraphicsProvider) AddQuad2D(batchID uint8, quad Quad2D, color ColorFA, uvQuad Quad2D, extra VertExtra) (BatchSlice, DeepError) {
	dErr := utils.NewDeepError("[PolyApp] AddQuad2D():")
	dErr.IsErr = false
	bSlice, err := g.AllocateShapeInBatch(batchID, BatchShape{
		VertCount: 4,
		Indexes:   []uint32{0, 1, 2, 2, 3, 0},
	})
	if err.IsErr {
		dErr.AddChildDeepError(err)
		return bSlice, dErr
	}
	dErr.AddChildDeepError(g.UpdateQuad2D(bSlice, quad, color, uvQuad, extra))
	return bSlice, dErr
}
func (g GraphicsProvider) UpdateQuad2D(batchSlice BatchSlice, quad Quad2D, color ColorFA, uvQuad Quad2D, extra VertExtra) DeepError {
	if batchSlice.VertLen() != 4 || batchSlice.IdxLen() != 6 {
		return utils.NewDeepError("[PolyApp] UpdateQuad2D(): batch slice provided does not have required dimensions for a quad")
	}
	dErr := utils.NewDeepError("[PolyApp] UpdateQuad2D():")
	dErr.IsErr = false
	v := Vertex{
		Pos:   quad.A().AsVec3(),
		Norm:  Vec3{0, 0, -g.XRightYUpZAway()[2]},
		UV:    uvQuad.A(),
		Color: color,
		Extra: extra,
	}
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart, v))
	v.Pos = quad.B().AsVec3()
	v.UV = uvQuad.B()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+1, v))
	v.Pos = quad.C().AsVec3()
	v.UV = uvQuad.C()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+2, v))
	v.Pos = quad.D().AsVec3()
	v.UV = uvQuad.D()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+3, v))
	return dErr
}
func (g GraphicsProvider) AddRect2D(batchID uint8, rect Rect2D, color ColorFA, uvRect Rect2D, extra VertExtra) (BatchSlice, DeepError) {
	dErr := utils.NewDeepError("[PolyApp] AddRect2D():")
	dErr.IsErr = false
	quad, uvQuad := rect.Quad(), uvRect.Quad()
	bs, err := g.AddQuad2D(batchID, quad, color, uvQuad, extra)
	dErr.AddChildDeepError(err)
	return bs, dErr
}
func (g GraphicsProvider) UpdateRect2D(batchSlice BatchSlice, rect Rect2D, color ColorFA, uvRect Rect2D, extra VertExtra) DeepError {
	dErr := utils.NewDeepError("[PolyApp] UpdateRect2D():")
	dErr.IsErr = false
	quad, uvQuad := rect.Quad(), uvRect.Quad()
	dErr.AddChildDeepError(g.UpdateQuad2D(batchSlice, quad, color, uvQuad, extra))
	return dErr
}
func (g GraphicsProvider) AddQuadOutline2D(batchID uint8, quadInner Quad2D, quadOuter Quad2D, color ColorFA, uvQuadInner Quad2D, uvQuadOuter Quad2D, extra VertExtra) (BatchSlice, DeepError) {
	dErr := utils.NewDeepError("[PolyApp] AddQuadOutline2D():")
	dErr.IsErr = false
	bSlice, err := g.AllocateShapeInBatch(batchID, BatchShape{
		VertCount: 8,
		Indexes:   []uint32{0, 1, 2, 1, 3, 2, 2, 3, 4, 3, 5, 4, 4, 5, 6, 5, 7, 6, 6, 7, 0, 7, 1, 0},
	})
	if err.IsErr {
		dErr.AddChildDeepError(err)
		return bSlice, dErr
	}
	dErr.AddChildDeepError(g.UpdateQuadOutline2D(bSlice, quadInner, quadOuter, color, uvQuadInner, uvQuadOuter, extra))
	return bSlice, dErr
}
func (g GraphicsProvider) UpdateQuadOutline2D(batchSlice BatchSlice, quadInner Quad2D, quadOuter Quad2D, color ColorFA, uvQuadInner Quad2D, uvQuadOuter Quad2D, extra VertExtra) DeepError {
	if batchSlice.VertLen() != 8 || batchSlice.IdxLen() != 24 {
		return utils.NewDeepError("[PolyApp] UpdateQuadOutline2D(): batch slice provided does not have required dimensions for a quad outline")
	}
	dErr := utils.NewDeepError("[PolyApp] UpdateQuadOutline2D():")
	dErr.IsErr = false
	v := Vertex{
		Pos:   quadInner.A().AsVec3(),
		Norm:  Vec3{0, 0, -g.XRightYUpZAway()[2]},
		UV:    uvQuadInner.A(),
		Color: color,
		Extra: extra,
	}
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart, v))
	v.Pos = quadInner.B().AsVec3()
	v.UV = uvQuadInner.B()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+1, v))
	v.Pos = quadInner.B().AsVec3()
	v.UV = uvQuadInner.B()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+2, v))
	v.Pos = quadInner.B().AsVec3()
	v.UV = uvQuadInner.B()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+3, v))
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+4, v))
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+5, v))
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+6, v))
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	dErr.AddChildDeepError(g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+7, v))
	return dErr
}

func (g GraphicsProvider) AddRectOutline2D(batchID uint8, rect Rect2D, thickness float32, color ColorFA, uvRect Rect2D, uvThickness float32, extra VertExtra) (BatchSlice, DeepError) {
	dErr := utils.NewDeepError("[PolyApp] AddRectOutline2D():")
	dErr.IsErr = false
	innerQuad, uvInnerQuad := rect.Quad(), uvRect.Quad()
	outerQuad := rect.Translate(Vec2{-thickness, -thickness}).Expand(Vec2{2 * thickness, 2 * thickness}).Quad()
	uvOuterQuad := uvRect.Translate(Vec2{-uvThickness, -uvThickness}).Expand(Vec2{2 * uvThickness, 2 * uvThickness}).Quad()
	bs, err := g.AddQuadOutline2D(batchID, innerQuad, outerQuad, color, uvInnerQuad, uvOuterQuad, extra)
	dErr.AddChildDeepError(err)
	return bs, dErr
}

func (g GraphicsProvider) UpdateRectOutline2D(batchSlice BatchSlice, rect Rect2D, thickness float32, color ColorFA, uvRect Rect2D, uvThickness float32, extra VertExtra) DeepError {
	dErr := utils.NewDeepError("[PolyApp] UpdateRectOutline2D():")
	dErr.IsErr = false
	innerQuad, uvInnerQuad := rect.Quad(), uvRect.Quad()
	outerQuad := rect.Translate(Vec2{-thickness, -thickness}).Expand(Vec2{2 * thickness, 2 * thickness}).Quad()
	uvOuterQuad := uvRect.Translate(Vec2{-uvThickness, -uvThickness}).Expand(Vec2{2 * uvThickness, 2 * uvThickness}).Quad()
	dErr.AddChildDeepError(g.UpdateQuadOutline2D(batchSlice, innerQuad, outerQuad, color, uvInnerQuad, uvOuterQuad, extra))
	return dErr
}
