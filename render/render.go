package render

import (
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
// WARNING: Assumes a valid OpenGL context has been created
func (self *Renderer) Init() error {
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)

	gl.GenVertexArrays(1, &self.Vao)
	gl.BindVertexArray(self.Vao)

	gl.GenBuffers(1, &self.Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, self.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Vertices)*4, gl.Ptr(self.Vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(POSITION_ATTRIB, 3, gl.FLOAT, false, 12, nil)
	gl.EnableVertexArrayAttrib(self.Vbo, POSITION_ATTRIB)

	gl.CreateBuffers(1, &self.IBuf)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, self.IBuf)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(self.Indices)*4, gl.Ptr(self.Indices), gl.DYNAMIC_DRAW)

	program, err := LoadShaderProgram("/home/phillipdavis/everyday/dev/go/go-ecs/render/shaders/passthrough.vert.glsl", "/home/phillipdavis/everyday/dev/go/go-ecs/render/shaders/solid_color.glsl.frag")
	if err != nil {
		panic(err)	
	}
	self.Program = program

	return nil
}

func (self *Renderer) RenderLogic(
	window *glfw.Window,
	workQueue chan func(),
	ecs *entity.ECS,
	queries []entity.EntityType,
	entities []roaring64.Bitmap,
	queriesMut []entity.EntityType,
	entitiesMut []roaring64.Bitmap) {

	DoOn(workQueue, func() {
		h, w := window.GetSize()

		self.Vertices[0]  = float32(w/2.0 - w/4.0)
		self.Vertices[1]  = float32(h/2.0 + h/4.0)
		self.Vertices[2]  = 0.0

		self.Vertices[3]  = float32(w/2.0 - w/4.0) 
		self.Vertices[4]  = float32(h/2.0 - h/4.0) 
		self.Vertices[5]  = 0.0

		self.Vertices[6]  = float32(w/2.0 + w/4.0) 
		self.Vertices[7]  = float32(h/2.0 - h/4.0) 
		self.Vertices[8]  = 0.0

		self.Vertices[9]  = float32(w/2.0 + w/4.0) 
		self.Vertices[10] = float32(h/2.0 + h/4.0) 
		self.Vertices[11] = 0.0
		// Render the tiles
		window.MakeContextCurrent()	
		glfw.PollEvents()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(self.Program)
		
		/// resLoc := gl.GetUniformLocation(self.Program, gl.Str("res\000"))

		/// gl.Uniform3f( resLoc, float32(w), float32(h), 0.0)

		orth := mgl32.Ortho(0, float32(w), float32(h), 0, -1, 1)

		gl.BindVertexArray(self.Vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, self.Vbo)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(self.Vertices)*4, gl.Ptr(self.Vertices))

		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		window.SwapBuffers()
	})
}
