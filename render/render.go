package render

import (
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/phdavis1027/goecs/entity"
)

// Index -> Entity Map:
// 0 - Tiles

// Vertex Attribute Enum
const (
	POSITION_ATTRIB = iota
)

type Renderer struct {
	Vertices []float32
	Indices  []uint32
	IBuf     uint32
	Vao      uint32
	Vbo      uint32
	Program  uint32
}

func NewRenderer() *Renderer {
	return &Renderer{
		Vertices: make([]float32, 12),
		Indices:  make([]uint32, 6),
	}
}

// WARNING: MUST BE CALLED FROM THE MAIN THREAD
func (self *Renderer) Init() error {
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)

	gl.CreateVertexArrays(1, &self.Vao)
	gl.BindVertexArray(self.Vao)

	gl.CreateBuffers(1, &self.Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, self.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Vertices)*4, gl.Ptr(self.Vertices), gl.STATIC_DRAW)

	gl.EnableVertexArrayAttrib(self.Vbo, POSITION_ATTRIB)
	gl.VertexAttribPointer(POSITION_ATTRIB, 3, gl.FLOAT, false, 0, nil)

	gl.CreateBuffers(1, &self.IBuf)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, self.IBuf)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(self.Indices)*4, gl.Ptr(self.Indices), gl.STATIC_DRAW)

	program, err := LoadShaderProgram("/home/phillipdavis/everyday/dev/go/go-ecs/render/shaders/passthrough.glsl.vert", "/home/phillipdavis/everyday/dev/go/go-ecs/render/shaders/fragment.glsl")
	if err != nil {
		return err
	}
	self.Program = program

	gl.UseProgram(program)

	return nil
}

func (self *Renderer) RenderLogic(
	workQueue chan func(),
	ecs *entity.ECS,
	queries []entity.EntityType,
	entities []roaring64.Bitmap,
	queriesMut []entity.EntityType,
	entitiesMut []roaring64.Bitmap) {

	DoOn(workQueue, func() {
		// Render the tiles
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(self.Program)
		gl.BindVertexArray(self.Vao)

		gl.DrawElements(gl.TRIANGLES, int32(len(self.Indices)), gl.UNSIGNED_INT, nil)
	})
}
