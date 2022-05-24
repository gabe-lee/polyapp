package polyapp

import (
	"errors"

	geom "github.com/gabe-lee/gengeom"
	math "github.com/gabe-lee/genmath"
)

type GraphicsInterface interface {
	GetWindowSize() Vec2

	XRightYUpZAway() Vec3

	AddRenderer(vertexFlags VertexFlags, shaders []*Shader) (rendererID uint8, err error)
	AddDrawBatch(vertexFlags VertexFlags, textureID uint8, initialSize uint32) (batchID uint8, err error)
	AddTexture(texture *Texture) (textureID uint8, err error)
	AddDrawSurface(size IVec2) (surfaceID uint8, textureID uint8, err error)

	ClearSurface(surfaceID uint8, baseColor Color32)
	ClearSurfaceArea(surfaceID uint8, baseColor Color32, area IRect2D)

	DrawBatch(batchID uint8, surfaceID uint8, rendererID uint8)

	AddVertexToBatch(batchID uint8, vertex Vertex) (vertIndex uint32)
	UpdateVertexInBatch(batchID uint8, vertIndex uint32, vertex Vertex)
	AddIndexesToBatch(batchID uint8, vertIndexes ...uint32) (indexSlice BatchSlice)
	DeleteIndexesFromBatch(indexSlice BatchSlice)
	DeleteVerticesFromBatch(batchSlice BatchSlice)
	ClearBatch(batchID uint8)
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
	Ex160     VertexFlags = 640  // 5 additional 32bit data blocks
	Ex192     VertexFlags = 768  // 6 additional 32bit data blocks
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
	case vf&ExMask == Ex192:
		return 24
	case vf&ExMask == Ex160:
		return 20
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
	Size    Vec2
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

func (b BatchSlice) IdxLen() uint32 {
	return b.IndexEnd - b.IndexStart
}

func (b BatchSlice) VertLen() uint32 {
	return b.VertexEnd - b.VertexStart
}

func (b BatchSlice) Combine(other BatchSlice) (BatchSlice, error) {
	if b.BatchID != other.BatchID || // Must be on same batch
		!(b.IndexEnd == other.IndexStart || b.IndexStart == other.IndexEnd) || // one of the index starts must match one of the index ends
		!(b.VertexStart == other.VertexEnd || b.VertexEnd == other.IndexStart) { // one of the vertex starts must match one of the vertex ends
		return b, errors.New("[PolyApp] Cannot combine batch slices that aren't adjacent and on the same batch")
	}
	if b.IndexEnd == other.IndexStart {
		b.IndexEnd = other.IndexEnd
	} else {
		b.IndexStart = other.IndexStart
	}
	if b.VertexEnd == other.VertexStart {
		b.VertexEnd = other.VertexEnd
	} else {
		b.VertexStart = other.VertexStart
	}
	return b, nil
}

/**************
	LINES
***************/

func (g GraphicsProvider) AddLine2D(batchID uint8, a Vertex, b Vertex, thickness float32, uvThickness float32) (batchSlice BatchSlice) {
	l := Line2D{a.Pos.AsVec2(), b.Pos.AsVec2()}
	u := Line2D{a.UV, b.UV}
	a.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	b.Norm = a.Norm
	l1, l2 := l.PerpLines(thickness / 2)
	u1, u2 := u.PerpLines(uvThickness / 2)
	a.Pos = l1.A().AsVec3()
	a.UV = u1.A()
	a1 := g.AddVertexToBatch(batchID, a)
	a.Pos = l2.A().AsVec3()
	a.UV = u2.A()
	a2 := g.AddVertexToBatch(batchID, a)
	b.Pos = l1.B().AsVec3()
	b.UV = u1.B()
	b1 := g.AddVertexToBatch(batchID, b)
	b.Pos = l2.B().AsVec3()
	b.UV = u2.B()
	b2 := g.AddVertexToBatch(batchID, b)
	batchSlice = g.AddIndexesToBatch(batchID, a1, a2, b1, a1, b2, b1)
	batchSlice.BatchID = batchID
	batchSlice.VertexStart = a1
	batchSlice.VertexEnd = b2 + 1
	return batchSlice
}

func (g GraphicsProvider) UpdateLine2D(batchSlice BatchSlice, a Vertex, b Vertex, thickness float32, uvThickness float32) error {
	if batchSlice.IdxLen() != 6 || batchSlice.VertLen() != 4 {
		return errors.New("[PolyApp] UpdateLine2D(): batch slice provided does not have required dimensions for a line")
	}
	a.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	b.Norm = a.Norm
	l := Line2D{a.Pos.AsVec2(), b.Pos.AsVec2()}
	u := Line2D{a.UV, b.UV}
	l1, l2 := l.PerpLines(thickness / 2)
	u1, u2 := u.PerpLines(uvThickness / 2)
	a.Pos = l1.A().AsVec3()
	a.UV = u1.A()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart, a)
	a.Pos = l2.A().AsVec3()
	a.UV = u2.A()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+1, a)
	b.Pos = l1.B().AsVec3()
	b.UV = u1.B()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+2, a)
	b.Pos = l2.B().AsVec3()
	b.UV = u2.B()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+3, a)
	return nil
}

/**************
	POLYGONS
***************/

func (g GraphicsProvider) AddRegularPolygon2D(batchID uint8, center Vertex, sides uint32, radius float32, shapeRotation float32, uvRadius float32, uvRotation float32) BatchSlice {
	indexes := make([]uint32, sides)
	batchSlice := BatchSlice{}
	center.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	points := geom.PointsOnCircle(shapeRotation*math.DEG_TO_RAD, radius, center.Pos.AsVec2(), sides)
	uvs := geom.PointsOnCircle(uvRotation*math.DEG_TO_RAD, uvRadius, center.UV, sides)
	cenIdx := g.AddVertexToBatch(batchID, center)
	for i := range points {
		center.Pos = points[i].AsVec3()
		center.UV = uvs[i]
		indexes[i] = g.AddVertexToBatch(batchID, center)
		if i > 0 {
			batchSlice, _ = batchSlice.Combine(g.AddIndexesToBatch(batchID, cenIdx, indexes[i-1], indexes[i]))
		}
	}
	batchSlice, _ = batchSlice.Combine(g.AddIndexesToBatch(batchID, cenIdx, indexes[len(indexes)-1], indexes[0]))
	batchSlice.BatchID = batchID
	batchSlice.VertexStart = cenIdx
	batchSlice.VertexEnd = indexes[len(indexes)-1] + 1
	return batchSlice
}

func (g GraphicsProvider) UpdateRegularPolygon2D(batchSlice BatchSlice, center Vertex, sides uint32, radius float32, shapeRotation float32, uvRadius float32, uvRotation float32) error {
	if batchSlice.VertLen() != sides+1 || batchSlice.IdxLen() != sides*3 {
		return errors.New("[PolyApp] UpdateRegularPolygon2D(): batch slice provided does not have required dimensions for a polygon of specified sides")
	}
	center.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	points := geom.PointsOnCircle(shapeRotation*math.DEG_TO_RAD, radius, center.Pos.AsVec2(), sides)
	uvs := geom.PointsOnCircle(uvRotation*math.DEG_TO_RAD, uvRadius, center.UV, sides)
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart, center)
	for i := uint32(0); i < uint32(len(points)); i += 1 {
		center.Pos = points[i].AsVec3()
		center.UV = uvs[i]
		g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+i+1, center)
	}
	return nil
}

func (g GraphicsProvider) AddRegularPolygonRing2D(batchID uint8, center Vertex, sides uint32, innerRadius float32, outerRadius float32, shapeRotation float32, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32) BatchSlice {
	sides = math.Round(sides)
	indexes := make([]uint32, sides*2)
	batchSlice := BatchSlice{}
	center.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	uvs := geom.PointsOnRing(uvRotation*math.DEG_TO_RAD, uvInnerRadius, uvOuterRadius, center.UV, sides)
	points := geom.PointsOnRing(shapeRotation*math.DEG_TO_RAD, innerRadius, outerRadius, center.Pos.AsVec2(), sides)
	for i := range points {
		center.Pos = points[i].AsVec3()
		center.UV = uvs[i]
		indexes[i] = g.AddVertexToBatch(batchID, center)
	}
	for i := 0; i <= len(indexes)-4; i += 2 {
		batchSlice, _ = batchSlice.Combine(g.AddIndexesToBatch(batchID, indexes[i+0], indexes[i+1], indexes[i+2], indexes[i+1], indexes[i+3], indexes[i+2]))
	}
	batchSlice, _ = batchSlice.Combine(g.AddIndexesToBatch(batchID, indexes[len(indexes)-2], indexes[len(indexes)-1], indexes[0], indexes[len(indexes)-1], indexes[1], indexes[0]))
	batchSlice.BatchID = batchID
	batchSlice.VertexStart = indexes[0]
	batchSlice.VertexEnd = indexes[len(indexes)-1] + 1
	return batchSlice
}

func (g GraphicsProvider) UpdateRegularPolygonRing2D(batchSlice BatchSlice, center Vertex, sides uint32, innerRadius float32, outerRadius float32, shapeRotation float32, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32) error {
	if batchSlice.VertLen() != sides*2 || batchSlice.IdxLen() != sides*6 {
		return errors.New("[PolyApp] UpdateRegularPolygonRing2D(): batch slice provided does not have required dimensions for a polygon of specified sides")
	}
	center.Norm = Vec3{0, 0, -g.XRightYUpZAway()[2]}
	uvs := geom.PointsOnRing(uvRotation*math.DEG_TO_RAD, uvInnerRadius, uvOuterRadius, center.UV, sides)
	points := geom.PointsOnRing(shapeRotation*math.DEG_TO_RAD, innerRadius, outerRadius, center.Pos.AsVec2(), sides)
	for i := uint32(0); i < uint32(len(points)); i += 1 {
		center.Pos = points[i].AsVec3()
		center.UV = uvs[i]
		g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+i, center)
	}
	return nil
}

/**************
	CIRCLES
***************/

func (g GraphicsProvider) AddCircleAutoPoints2D(batchID uint8, center Vertex, resolution float32, radius float32, uvRadius float32, uvRotation float32) BatchSlice {
	sides := uint32(math.Ciel(geom.Circumference(radius) / resolution))
	return g.AddRegularPolygon2D(batchID, center, sides, radius, 0, uvRadius, uvRotation)
}
func (g GraphicsProvider) UpdateCircleAutoPoints2D(batchSlice BatchSlice, center Vertex, resolution float32, radius float32, uvRadius float32, uvRotation float32) error {
	sides := uint32(math.Ciel(geom.Circumference(radius) / resolution))
	return g.UpdateRegularPolygon2D(batchSlice, center, sides, radius, 0, uvRadius, uvRotation)
}
func (g GraphicsProvider) AddCircleRingAutoPoints2D(batchID uint8, center Vertex, resolution float32, innerRadius float32, outerRadius float32, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32) BatchSlice {
	sides := uint32(math.Ciel(geom.Circumference(outerRadius) / resolution))
	return g.AddRegularPolygonRing2D(batchID, center, sides, innerRadius, outerRadius, 0, uvInnerRadius, uvOuterRadius, uvRotation)
}
func (g GraphicsProvider) UpdateCircleRingAutoPoints2D(batchSlice BatchSlice, center Vertex, resolution float32, innerRadius float32, outerRadius float32, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32) error {
	sides := uint32(math.Ciel(geom.Circumference(outerRadius) / resolution))
	return g.UpdateRegularPolygonRing2D(batchSlice, center, sides, innerRadius, outerRadius, 0, uvInnerRadius, uvOuterRadius, uvRotation)
}

/**************
	RECTANGLES
***************/

func (g GraphicsProvider) AddQuad2D(batchID uint8, quad Quad2D, color ColorFA, uvQuad Quad2D, extra VertExtra) BatchSlice {
	v := Vertex{
		Pos:   quad.A().AsVec3(),
		Norm:  Vec3{0, 0, g.XRightYUpZAway()[2]},
		UV:    uvQuad.A(),
		Color: color,
		Extra: extra,
	}
	a := g.AddVertexToBatch(batchID, v)
	v.Pos = quad.B().AsVec3()
	v.UV = uvQuad.B()
	b := g.AddVertexToBatch(batchID, v)
	v.Pos = quad.C().AsVec3()
	v.UV = uvQuad.C()
	c := g.AddVertexToBatch(batchID, v)
	v.Pos = quad.D().AsVec3()
	v.UV = uvQuad.D()
	d := g.AddVertexToBatch(batchID, v)
	batchSlice := g.AddIndexesToBatch(batchID, a, b, c)
	batchSlice, _ = batchSlice.Combine(g.AddIndexesToBatch(batchID, c, d, a))
	batchSlice.BatchID = batchID
	batchSlice.VertexStart = a
	batchSlice.VertexEnd = d + 1
	return batchSlice
}
func (g GraphicsProvider) UpdateQuad2D(batchSlice BatchSlice, quad Quad2D, color ColorFA, uvQuad Quad2D, extra VertExtra) error {
	if batchSlice.VertLen() != 4 || batchSlice.IdxLen() != 6 {
		return errors.New("[PolyApp] UpdateQuad2D(): batch slice provided does not have required dimensions for a quad")
	}
	v := Vertex{
		Pos:   quad.A().AsVec3(),
		Norm:  Vec3{0, 0, -g.XRightYUpZAway()[2]},
		UV:    uvQuad.A(),
		Color: color,
		Extra: extra,
	}
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart, v)
	v.Pos = quad.B().AsVec3()
	v.UV = uvQuad.B()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+1, v)
	v.Pos = quad.C().AsVec3()
	v.UV = uvQuad.C()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+2, v)
	v.Pos = quad.D().AsVec3()
	v.UV = uvQuad.D()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+3, v)
	return nil
}
func (g GraphicsProvider) AddRect2D(batchID uint8, rect Rect2D, color ColorFA, uvRect Rect2D, extra VertExtra) BatchSlice {
	quad, uvQuad := rect.Quad(), uvRect.Quad()
	return g.AddQuad2D(batchID, quad, color, uvQuad, extra)
}
func (g GraphicsProvider) UpdateRect2D(batchSlice BatchSlice, rect Rect2D, color ColorFA, uvRect Rect2D, extra VertExtra) error {
	quad, uvQuad := rect.Quad(), uvRect.Quad()
	return g.UpdateQuad2D(batchSlice, quad, color, uvQuad, extra)
}
func (g GraphicsProvider) AddQuadOutline2D(batchID uint8, quadInner Quad2D, quadOuter Quad2D, color ColorFA, uvQuadInner Quad2D, uvQuadOuter Quad2D, extra VertExtra) BatchSlice {
	v := Vertex{
		Pos:   quadInner.A().AsVec3(),
		Norm:  Vec3{0, 0, -g.XRightYUpZAway()[2]},
		UV:    uvQuadInner.A(),
		Color: color,
		Extra: extra,
	}
	ai := g.AddVertexToBatch(batchID, v)
	v.Pos = quadInner.B().AsVec3()
	v.UV = uvQuadInner.B()
	bi := g.AddVertexToBatch(batchID, v)
	v.Pos = quadInner.B().AsVec3()
	v.UV = uvQuadInner.B()
	ci := g.AddVertexToBatch(batchID, v)
	v.Pos = quadInner.B().AsVec3()
	v.UV = uvQuadInner.B()
	di := g.AddVertexToBatch(batchID, v)
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	ao := g.AddVertexToBatch(batchID, v)
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	bo := g.AddVertexToBatch(batchID, v)
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	co := g.AddVertexToBatch(batchID, v)
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	do := g.AddVertexToBatch(batchID, v)
	bSlice := g.AddIndexesToBatch(batchID, ai, ao, bi, ao, bo, bi, bi, bo, ci, bo, co, ci, ci, co, di, co, do, di, di, do, ai, do, ao, ai)
	bSlice.BatchID = batchID
	bSlice.VertexStart = ai
	bSlice.VertexEnd = do + 1
	return bSlice
}
func (g GraphicsProvider) UpdateQuadOutline2D(batchSlice BatchSlice, quadInner Quad2D, quadOuter Quad2D, color ColorFA, uvQuadInner Quad2D, uvQuadOuter Quad2D, extra VertExtra) error {
	if batchSlice.VertLen() != 8 || batchSlice.IdxLen() != 24 {
		return errors.New("[PolyApp] UpdateQuadOutline2D(): batch slice provided does not have required dimensions for a quad outline")
	}
	v := Vertex{
		Pos:   quadInner.A().AsVec3(),
		Norm:  Vec3{0, 0, -g.XRightYUpZAway()[2]},
		UV:    uvQuadInner.A(),
		Color: color,
		Extra: extra,
	}
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart, v)
	v.Pos = quadInner.B().AsVec3()
	v.UV = uvQuadInner.B()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+1, v)
	v.Pos = quadInner.B().AsVec3()
	v.UV = uvQuadInner.B()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+2, v)
	v.Pos = quadInner.B().AsVec3()
	v.UV = uvQuadInner.B()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+3, v)
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+4, v)
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+5, v)
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+6, v)
	v.Pos = quadOuter.B().AsVec3()
	v.UV = uvQuadOuter.B()
	g.UpdateVertexInBatch(batchSlice.BatchID, batchSlice.VertexStart+7, v)
	return nil
}

func (g GraphicsProvider) AddRectOutline2D(batchID uint8, rect Rect2D, thickness float32, color ColorFA, uvRect Rect2D, uvThickness float32, extra VertExtra) BatchSlice {
	innerQuad, uvInnerQuad := rect.Quad(), uvRect.Quad()
	outerQuad := rect.Translate(Vec2{-thickness, -thickness}).Expand(Vec2{2 * thickness, 2 * thickness}).Quad()
	uvOuterQuad := uvRect.Translate(Vec2{-uvThickness, -uvThickness}).Expand(Vec2{2 * uvThickness, 2 * uvThickness}).Quad()
	return g.AddQuadOutline2D(batchID, innerQuad, outerQuad, color, uvInnerQuad, uvOuterQuad, extra)
}

func (g GraphicsProvider) UpdateRectOutline2D(batchSlice BatchSlice, rect Rect2D, thickness float32, color ColorFA, uvRect Rect2D, uvThickness float32, extra VertExtra) error {
	innerQuad, uvInnerQuad := rect.Quad(), uvRect.Quad()
	outerQuad := rect.Translate(Vec2{-thickness, -thickness}).Expand(Vec2{2 * thickness, 2 * thickness}).Quad()
	uvOuterQuad := uvRect.Translate(Vec2{-uvThickness, -uvThickness}).Expand(Vec2{2 * uvThickness, 2 * uvThickness}).Quad()
	return g.UpdateQuadOutline2D(batchSlice, innerQuad, outerQuad, color, uvInnerQuad, uvOuterQuad, extra)
}
