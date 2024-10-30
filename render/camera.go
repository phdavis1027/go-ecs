package render

import "github.com/go-gl/mathgl/mgl32"

type OrthographicCamera struct {
	ProjectionMatrix     mgl32.Mat4
	ViewMatrix           mgl32.Mat4
	ViewProjectionMatrix mgl32.Mat4
	Position             mgl32.Vec3
}

func NewOrthographicCamera(left, right, bottom, top float32) *OrthographicCamera {
	projection := mgl32.Ortho(left, right, bottom, top, -1, 1)
	return &OrthographicCamera{
		ProjectionMatrix: projection,
		ViewMatrix:      mgl32.Ident4(),
		Position:        mgl32.Vec3{0, 0, 0},
	}
}

func (c *OrthographicCamera) SetPosition(position mgl32.Vec3) {
	c.Position = position
	c.UpdateViewMatrix()
}

func (c *OrthographicCamera) UpdateViewMatrix() {
	c.ViewMatrix = mgl32.Translate3D(-c.Position[0], -c.Position[1], -c.Position[2])
	c.ViewProjectionMatrix = c.ProjectionMatrix.Mul4(c.ViewMatrix)
}
