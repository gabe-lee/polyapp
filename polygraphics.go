package polyapp

import (
	geom "github.com/gabe-lee/gengeom"
	math "github.com/gabe-lee/genmath"
)

type GraphicsInterface interface {
	GetWindowSize() Vec2

	AddRenderer(vertexSpace VertexSpace, shaders []*Shader) (rendererID uint8, err error)
	AddDrawBatch(vertexSpace VertexSpace, textureID uint8, initialSize uint32) (batchID uint8, err error)
	AddTexture(texture *Texture) (textureID uint8, err error)
	AddDrawSurface(size IVec2) (surfaceID uint8, textureID uint8, err error)

	ClearSurface(surfaceID uint8, baseColor Color32)
	ClearSurfaceArea(surfaceID uint8, baseColor Color32, area IRect2D)

	DrawBatch(batchID uint8, surfaceID uint8, rendererID uint8)

	AddVertexToBatch2D(batchID uint8, position Vec2, color Color32, textureUV Vec2, extra uint32) (vertIndex uint16)
	AddVertexToBatch3D(batchID uint8, position Vec3, color Color32, textureUV Vec2, extra uint32) (vertIndex uint16)
	UpdateVertexInBatch2D(batchID uint8, vertIndex uint16, position Vec2, color Color32, textureUV Vec2, extra uint32)
	UpdateVertexInBatch3D(batchID uint8, vertIndex uint16, position Vec2, color Color32, textureUV Vec2, extra uint32)
	AddIndexesToBatch(batchID uint8, vertIndexes ...uint16) (indexSlice []uint16)
	DeleteIndexesFromBatch(batchID uint8, indexSlice []uint16)
	ClearBatch(batchID uint8)
}

var _ GraphicsInterface = (*GraphicsProvider)(nil)

type GraphicsProvider struct {
	App *App
	GraphicsInterface
}

var ninf = math.NInf32()

var NoUV = Vec2{ninf, ninf}
var NoDraw2D = Vec2{ninf, ninf}
var NoDraw3D = Vec3{ninf, ninf, ninf}

type VertexSpace uint8

const (
	V2D VertexSpace = 2
	V3D VertexSpace = 3
)

type IndexSize uint8

const (
	Idx16 IndexSize = 16
	Idx32 IndexSize = 32
)

type HasTexture uint8

const (
	NoTex  HasTexture = 0
	HasTex HasTexture = 1
	False             = NoTex
	True              = HasTex
)

type ColorSpace uint8

const (
	Col6  ColorSpace = 6
	Col8  ColorSpace = 8
	Col12 ColorSpace = 12
	Col16 ColorSpace = 16
	Col24 ColorSpace = 24
	Col32 ColorSpace = 32
)

type VertexProperties struct {
	VertexSpace
	IndexSize
	ColorSpace
	HasTexture HasTexture
}

type VertexMode uint8

const (
	Pixels        VertexMode = iota // Each vertex is an independant pixel
	Lines                           // Each pair of vertices forms an independant line
	LineStrip                       // Each vertex forms a continuous line with the vertex following it
	LineLoop                        // Each vertex forms a continuous line with the vertex following it, last vertex connects back to first
	Triangles                       // Every 3 vertices form independant triangles
	TriangleStrip                   // Every vertex forms a triangle using the 2 following it with alternating windings
	TriangleFan                     // Every vertex after the first uses the one following it and the very first to form a triangle
)

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

type ShapeInstance struct {
	BatchID         uint8
	BatchIndexStart uint32
	BatchIndexEnd   uint32
}

func (s ShapeInstance) Len() uint32 {
	return s.BatchIndexEnd - s.BatchIndexStart
}

/**************
	LINES
***************/

func (g GraphicsProvider) AddLine2D(batchID uint8, a Vec2, b Vec2, thickness float32, color Color32, uvA Vec2, uvB Vec2, uvThickness float32, extra uint32) {
	l := Line2D{a, b}
	u := Line2D{uvA, uvB}
	l1, l2 := l.PerpLines(thickness / 2)
	u1, u2 := u.PerpLines(uvThickness / 2)
	a1 := g.AddVertexToBatch2D(batchID, l1.A(), color, u1.A(), extra)
	a2 := g.AddVertexToBatch2D(batchID, l2.A(), color, u2.A(), extra)
	b1 := g.AddVertexToBatch2D(batchID, l1.B(), color, u1.B(), extra)
	b2 := g.AddVertexToBatch2D(batchID, l2.B(), color, u2.B(), extra)
	g.AddIndexesToBatch(batchID, a1, a2, b1, a1, b2, b1)
}

/**************
	POLYGONS
***************/

func (g GraphicsProvider) AddRegularPolygon2D(batchID uint8, center Vec2, sides uint32, radius float32, shapeRotation float32, color Color32, uvCenter Vec2, uvRadius float32, uvRotation float32, extra uint32) {
	sides = math.Round(sides)
	indexes := make([]uint16, sides)
	points := geom.PointsOnCircle(shapeRotation*math.DEG_TO_RAD, radius, center, sides)
	uvs := geom.PointsOnCircle(uvRotation*math.DEG_TO_RAD, uvRadius, uvCenter, sides)
	cenIdx := g.AddVertexToBatch2D(batchID, center, color, uvCenter, extra)
	for i := range points {
		indexes[i] = g.AddVertexToBatch2D(batchID, points[i], color, uvs[i], extra)
		if i > 0 {
			g.AddIndexesToBatch(batchID, cenIdx, indexes[i-1], indexes[i])
		}
	}
	g.AddIndexesToBatch(batchID, cenIdx, indexes[len(indexes)-1], indexes[0])
}

func (g GraphicsProvider) AddRegularPolygonRing2D(batchID uint8, center Vec2, sides uint32, innerRadius float32, outerRadius float32, shapeRotation float32, color Color32, uvCenter Vec2, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32, extra uint32) {
	sides = math.Round(sides)
	indexes := make([]uint16, sides*2)
	uvs := geom.PointsOnRing(uvRotation*math.DEG_TO_RAD, uvInnerRadius, uvOuterRadius, uvCenter, sides)
	points := geom.PointsOnRing(shapeRotation*math.DEG_TO_RAD, innerRadius, outerRadius, center, sides)
	for i := range points {
		indexes[i] = g.AddVertexToBatch2D(batchID, points[i], color, uvs[i], extra)
	}
	for i := 0; i <= len(indexes)-4; i += 2 {
		g.AddIndexesToBatch(batchID, indexes[i+0], indexes[i+1], indexes[i+2], indexes[i+1], indexes[i+3], indexes[i+2])
	}
	g.AddIndexesToBatch(batchID, indexes[len(indexes)-2], indexes[len(indexes)-1], indexes[0], indexes[len(indexes)-1], indexes[1], indexes[0])
}

/**************
	CIRCLES
***************/

func (g GraphicsProvider) AddCircleAutoPoints2D(batchID uint8, center Vec2, resolution float32, radius float32, color Color32, uvCenter Vec2, uvRadius float32, uvRotation float32, extra uint32) {
	sides := uint32(math.Ciel(geom.Circumference(radius) / resolution))
	g.AddRegularPolygon2D(batchID, center, sides, radius, 0, color, uvCenter, uvRadius, uvRotation, extra)
}
func (g GraphicsProvider) AddCircleRingAutoPoints2D(batchID uint8, center Vec2, resolution float32, innerRadius float32, outerRadius float32, color Color32, uvCenter Vec2, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32, extra uint32) {
	sides := uint32(math.Ciel(geom.Circumference(outerRadius) / resolution))
	g.AddRegularPolygonRing2D(batchID, center, sides, innerRadius, outerRadius, 0, color, uvCenter, uvInnerRadius, uvOuterRadius, uvRotation, extra)
}

/**************
	RECTANGLES
***************/

func (g GraphicsProvider) AddQuad2D(batchID uint8, quad Quad2D, color Color32, uvQuad Quad2D, extra uint32) {
	a := g.AddVertexToBatch2D(batchID, quad.A(), color, uvQuad.A(), extra)
	b := g.AddVertexToBatch2D(batchID, quad.B(), color, uvQuad.B(), extra)
	c := g.AddVertexToBatch2D(batchID, quad.C(), color, uvQuad.C(), extra)
	d := g.AddVertexToBatch2D(batchID, quad.D(), color, uvQuad.D(), extra)
	g.AddIndexesToBatch(batchID, a, b, c)
	g.AddIndexesToBatch(batchID, c, d, a)
}
func (g GraphicsProvider) AddRect2D(batchID uint8, rect Rect2D, color Color32, uvRect Rect2D, extra uint32) {
	quad, uvQuad := rect.Quad(), uvRect.Quad()
	g.AddQuad2D(batchID, quad, color, uvQuad, extra)
}
func (g GraphicsProvider) AddQuadOutline2D(batchID uint8, quadInner Quad2D, quadOuter Quad2D, color Color32, uvQuadInner Quad2D, uvQuadOuter Quad2D, extra uint32) {
	ai := g.AddVertexToBatch2D(batchID, quadInner.A(), color, uvQuadInner.A(), extra)
	bi := g.AddVertexToBatch2D(batchID, quadInner.B(), color, uvQuadInner.B(), extra)
	ci := g.AddVertexToBatch2D(batchID, quadInner.C(), color, uvQuadInner.C(), extra)
	di := g.AddVertexToBatch2D(batchID, quadInner.D(), color, uvQuadInner.D(), extra)
	ao := g.AddVertexToBatch2D(batchID, quadOuter.A(), color, uvQuadOuter.A(), extra)
	bo := g.AddVertexToBatch2D(batchID, quadOuter.B(), color, uvQuadOuter.B(), extra)
	co := g.AddVertexToBatch2D(batchID, quadOuter.C(), color, uvQuadOuter.C(), extra)
	do := g.AddVertexToBatch2D(batchID, quadOuter.D(), color, uvQuadOuter.D(), extra)
	g.AddIndexesToBatch(batchID, ai, ao, bi, ao, bo, bi, bi, bo, ci, bo, co, ci, ci, co, di, co, do, di, di, do, ai, do, ao, ai)
}

func (g GraphicsProvider) AddRectOutline2D(batchID uint8, rect Rect2D, thickness float32, color Color32, uvRect Rect2D, uvThickness float32, extra uint32) {
	innerQuad, uvInnerQuad := rect.Quad(), uvRect.Quad()
	outerQuad := rect.Translate(Vec2{-thickness, -thickness}).Expand(Vec2{2 * thickness, 2 * thickness}).Quad()
	uvOuterQuad := uvRect.Translate(Vec2{-uvThickness, -uvThickness}).Expand(Vec2{2 * uvThickness, 2 * uvThickness}).Quad()
	g.AddQuadOutline2D(batchID, innerQuad, outerQuad, color, uvInnerQuad, uvOuterQuad, extra)
}
