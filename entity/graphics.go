package entity

import "github.com/go-gl/mathgl/mgl32"

type RenderableQuadComponent struct {
	Height, Width float32
	Color         mgl32.Vec4
}

func NewRenderableQuadComponent(height, width float32, color mgl32.Vec4) *RenderableQuadComponent {
	return &RenderableQuadComponent{
		height: height,
		width:  width,
		color:  color,
	}
}


