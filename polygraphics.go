package polyapp

import (
	geom "github.com/gabe-lee/gengeom"
	math "github.com/gabe-lee/genmath"
)

type GraphicsInterface interface {
	GetWindowSize() Vec2

	AddRenderer(shaders []*Shader) (rendererID uint8)
	AddDrawBatch(initialSize uint32) (batchID uint8)
	AddTexture(texture *Texture) (textureID uint8)
	AddDrawSurface(size Vec2) (surfaceID uint8, textureID uint8)

	ClearDrawSurface(surfaceID uint8, baseColor Color)
	ClearDrawSurfaceArea(surfaceID uint8, baseColor Color, area Rect2D)

	DrawBatchIndexedTriangles2D(batchID uint8, surfaceID uint8, rendererID uint8)
	DrawBatchIndexedTriangles3D(batchID uint8, surfaceID uint8, rendererID uint8)

	AddVertexToBatch2D(batchID uint8, position Vec2, color Color, textureID uint8, textureUV Vec2, flags uint16) (index uint16)
	AddVertexToBatch3D(batchID uint8, position Vec3, color Color, textureID uint8, textureUV Vec2, flags uint16) (index uint16)
	AddIndexesToBatch(batchID uint8, indexes ...uint16)
	ClearBatch(batchID uint8)
}

var _ GraphicsInterface = (*GraphicsProvider)(nil)

type GraphicsProvider struct {
	GraphicsInterface
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

/**************
	LINES
***************/

func (g GraphicsProvider) AddLine2D(batchID uint8, a Vec2, b Vec2, thickness float32, color Color, textureID uint8, uvA Vec2, uvB Vec2, uvThickness float32, flags uint16) {
	l := Line2D{a, b}
	u := Line2D{uvA, uvB}
	l1, l2 := l.PerpLines(thickness / 2)
	u1, u2 := u.PerpLines(uvThickness / 2)
	a1 := g.AddVertexToBatch2D(batchID, l1.A(), color, textureID, u1.A(), flags)
	a2 := g.AddVertexToBatch2D(batchID, l2.A(), color, textureID, u2.A(), flags)
	b1 := g.AddVertexToBatch2D(batchID, l1.B(), color, textureID, u1.B(), flags)
	b2 := g.AddVertexToBatch2D(batchID, l2.B(), color, textureID, u2.B(), flags)
	g.AddIndexesToBatch(batchID, a1, a2, b1, a1, b2, b1)
}

/**************
	POLYGONS
***************/

func (g GraphicsProvider) AddRegularPolygon2D(batchID uint8, center Vec2, sides uint32, radius float32, shapeRotation float32, color Color, textureID uint8, uvCenter Vec2, uvRadius float32, uvRotation float32, flags uint16) {
	sides = math.Round(sides)
	indexes := make([]uint16, sides)
	points := geom.PointsOnCircle(shapeRotation*math.DEG_TO_RAD, radius, center, sides)
	uvs := geom.PointsOnCircle(uvRotation*math.DEG_TO_RAD, uvRadius, uvCenter, sides)
	cenIdx := g.AddVertexToBatch2D(batchID, center, color, textureID, uvCenter, flags)
	for i := range points {
		indexes[i] = g.AddVertexToBatch2D(batchID, points[i], color, textureID, uvs[i], flags)
		if i > 0 {
			g.AddIndexesToBatch(batchID, cenIdx, indexes[i-1], indexes[i])
		}
	}
	g.AddIndexesToBatch(batchID, cenIdx, indexes[len(indexes)-1], indexes[0])
}

func (g GraphicsProvider) AddRegularPolygonRing2D(batchID uint8, center Vec2, sides uint32, innerRadius float32, outerRadius float32, shapeRotation float32, color Color, textureID uint8, uvCenter Vec2, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32, flags uint16) {
	sides = math.Round(sides)
	indexes := make([]uint16, sides*2)
	uvs := geom.PointsOnRing(uvRotation*math.DEG_TO_RAD, uvInnerRadius, uvOuterRadius, uvCenter, sides)
	points := geom.PointsOnRing(shapeRotation*math.DEG_TO_RAD, innerRadius, outerRadius, center, sides)
	for i := range points {
		indexes[i] = g.AddVertexToBatch2D(batchID, points[i], color, textureID, uvs[i], flags)
	}
	for i := 0; i <= len(indexes)-4; i += 2 {
		g.AddIndexesToBatch(batchID, indexes[i+0], indexes[i+1], indexes[i+2], indexes[i+1], indexes[i+3], indexes[i+2])
	}
	g.AddIndexesToBatch(batchID, indexes[len(indexes)-2], indexes[len(indexes)-1], indexes[0], indexes[len(indexes)-1], indexes[1], indexes[0])
}

/**************
	CIRCLES
***************/

func (g GraphicsProvider) AddCircleAutoPoints2D(batchID uint8, center Vec2, resolution float32, radius float32, color Color, textureID uint8, uvCenter Vec2, uvRadius float32, uvRotation float32, flags uint16) {
	sides := uint32(math.Ciel(geom.Circumference(radius) / resolution))
	g.AddRegularPolygon2D(batchID, center, sides, radius, 0, color, textureID, uvCenter, uvRadius, uvRotation, flags)
}
func (g GraphicsProvider) AddCircleRingAutoPoints2D(batchID uint8, center Vec2, resolution float32, innerRadius float32, outerRadius float32, color Color, textureID uint8, uvCenter Vec2, uvInnerRadius float32, uvOuterRadius float32, uvRotation float32, flags uint16) {
	sides := uint32(math.Ciel(geom.Circumference(outerRadius) / resolution))
	g.AddRegularPolygonRing2D(batchID, center, sides, innerRadius, outerRadius, 0, color, textureID, uvCenter, uvInnerRadius, uvOuterRadius, uvRotation, flags)
}

/**************
	RECTANGLES
***************/

func (g GraphicsProvider) AddQuad2D(batchID uint8, quad Quad2D, color Color, textureID uint8, uvQuad Quad2D, flags uint16) {
	a := g.AddVertexToBatch2D(batchID, quad.A(), color, textureID, uvQuad.A(), flags)
	b := g.AddVertexToBatch2D(batchID, quad.B(), color, textureID, uvQuad.B(), flags)
	c := g.AddVertexToBatch2D(batchID, quad.C(), color, textureID, uvQuad.C(), flags)
	d := g.AddVertexToBatch2D(batchID, quad.D(), color, textureID, uvQuad.D(), flags)
	g.AddIndexesToBatch(batchID, a, b, c)
	g.AddIndexesToBatch(batchID, c, d, a)
}
func (g GraphicsProvider) AddRect2D(batchID uint8, rect Rect2D, color Color, textureID uint8, uvRect Rect2D, flags uint16) {
	quad, uvQuad := rect.Quad(), uvRect.Quad()
	g.AddQuad2D(batchID, quad, color, textureID, uvQuad, flags)
}
func (g GraphicsProvider) AddQuadOutline2D(batchID uint8, quadInner Quad2D, quadOuter Quad2D, color Color, textureID uint8, uvQuadInner Quad2D, uvQuadOuter Quad2D, flags uint16) {
	ai := g.AddVertexToBatch2D(batchID, quadInner.A(), color, textureID, uvQuadInner.A(), flags)
	bi := g.AddVertexToBatch2D(batchID, quadInner.B(), color, textureID, uvQuadInner.B(), flags)
	ci := g.AddVertexToBatch2D(batchID, quadInner.C(), color, textureID, uvQuadInner.C(), flags)
	di := g.AddVertexToBatch2D(batchID, quadInner.D(), color, textureID, uvQuadInner.D(), flags)
	ao := g.AddVertexToBatch2D(batchID, quadOuter.A(), color, textureID, uvQuadOuter.A(), flags)
	bo := g.AddVertexToBatch2D(batchID, quadOuter.B(), color, textureID, uvQuadOuter.B(), flags)
	co := g.AddVertexToBatch2D(batchID, quadOuter.C(), color, textureID, uvQuadOuter.C(), flags)
	do := g.AddVertexToBatch2D(batchID, quadOuter.D(), color, textureID, uvQuadOuter.D(), flags)
	g.AddIndexesToBatch(batchID, ai, ao, bi, ao, bo, bi, bi, bo, ci, bo, co, ci, ci, co, di, co, do, di, di, do, ai, do, ao, ai)
}
func (g GraphicsProvider) AddRectOutline2D(batchID uint8, rect Rect2D, thickness float32, color Color, textureID uint8, uvRect Rect2D, uvThickness float32, flags uint16) {
	innerQuad, uvInnerQuad := rect.Quad(), uvRect.Quad()
	outerQuad := rect.Translate(Vec2{-thickness, -thickness}).Expand(Vec2{2 * thickness, 2 * thickness}).Quad()
	uvOuterQuad := uvRect.Translate(Vec2{-uvThickness, -uvThickness}).Expand(Vec2{2 * uvThickness, 2 * uvThickness}).Quad()
	g.AddQuadOutline2D(batchID, innerQuad, outerQuad, color, textureID, uvInnerQuad, uvOuterQuad, flags)
}

// // Sprite Instance
// func (s *SystemSolution) DrawSpriteInstanceTinted(sInst *SpriteInstance, pos Vec2, color *Color) {
// 	frame := sInst.GetFrame()
// 	source := frame.texRect
// 	destPos := frame.drawOffset.Add(pos)
// 	dest := NewRect2D(destPos, source.Size())
// 	s.DrawFromTexComplete(frame.texIndex, source, dest, color, 0, Vec2{}, true)
// }
// func (s *SystemSolution) DrawSpriteInstanceDestRectTinted(sInst *SpriteInstance, dest Rect2D, color *Color) {
// 	frame := sInst.GetFrame()
// 	source := frame.texRect
// 	scale := dest.Size().Div(source.Size())
// 	destFinal := dest.TranslateCopy(frame.drawOffset.Mult(scale))
// 	s.DrawFromTexComplete(frame.texIndex, source, destFinal, color, 0, Vec2{}, true)
// }
