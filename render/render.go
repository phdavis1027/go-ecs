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
	Vao      uint32
	Vbo      uint32
}

func NewRenderer() *Renderer {
	return &Renderer{
		Vertices: make([]float32, 4),
		Indices:  make([]uint32, 6),
	}
}

func loadShader(shaderType uint32, path string) (uint32, error) {

}

// WARNING: MUST BE CALLED FROM THE MAIN THREAD
func (self *Renderer) Init() {
	gl.GenVertexArrays(1, &self.Vao)
	gl.BindVertexArray(self.Vao)

	gl.GenBuffers(1, &self.Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, self.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Vertices)*4, gl.Ptr(self.Vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(POSITION_ATTRIB, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(POSITION_ATTRIB)
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
		gl.BindVertexArray(self.Vao)

	})
}
