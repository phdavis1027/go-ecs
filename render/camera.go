package render

import "github.com/go-gl/mathgl/mgl32"

type OrthographicCamera struct {
	ProjectionMatrix     mgl32.Mat4
	ViewMatrix           mgl32.Mat4
	ViewProjectionMatrix mgl32.Mat4
	Position             mgl32.Vec3
}

func NewOrthographicCamera(width, height, resolution float32) *OrthographicCamera {
	projection := mgl32.Ortho(0, width, height, 0, 0, resolution)
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

func (c *OrthographicCamera) ApplyMovement(translation mgl32.Vec3) {
	c.SetPosition(c.Position.Add(translation))
}

func (c *OrthographicCamera) UpdateViewMatrix() {
    tranform := mgl32.Translate3D(c.Position[0], c.Position[1], c.Position[2])
	inv := tranform.Inv()
	c.ViewProjectionMatrix = c.ProjectionMatrix.Mul4(inv)
}
