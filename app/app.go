package app

import (
	"log"
	"os"
	"time"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/phdavis1027/goecs/entity"
	"github.com/phdavis1027/goecs/render"
)

const (
	WindowWidth  = 400 
	WindowHeight = 400
)

type App struct {
	name           string
	logger         *log.Logger
	ecs            *entity.ECS
	renderer       *render.Renderer
}

func NewApp(name string, ecsCap int) *App {
	return &App{
		logger:           log.New(os.Stdout, name, log.LstdFlags),
		ecs:              entity.CreateEcsOfCapacity(ecsCap),
		name:             name,
		renderer:         render.NewRenderer(),
	}
}

func (app *App) Main(stopButton chan(struct {})) error {
	// Main loop logic here
	var err error
	err = app.ecs.CompileSchedule(true)
	if err != nil {
		return err
	}

	err = app.ecs.RunSchedule()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Init() (*glfw.Window, error) {
	// initialize GL
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)


	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, app.name, nil, nil)

	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {	
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	err = gl.Init()

	if err != nil {
		panic(err)
	}


	err = app.renderer.Init()
	if err != nil {
		panic(err)
	}

	return window, nil
}

func (app *App) Terminate(window *glfw.Window)  {
	// Clean up resources
	glfw.Terminate()
	window.Destroy()	
}

// WARNING: Must run on main thread
func (app *App) Run() error {

	// GLFW things
	var renderWorkQueue  = make(chan func())
	var stopButton       = make(chan struct{})

	app.renderer.Vertices[0]  = 100 
	app.renderer.Vertices[1]  = 300 
	app.renderer.Vertices[2]  = 0.0

	app.renderer.Vertices[3]  = 100 
	app.renderer.Vertices[4]  = 100 
	app.renderer.Vertices[5]  = 0.0

	app.renderer.Vertices[6]  = 300 
	app.renderer.Vertices[7]  = 100 
	app.renderer.Vertices[8]  = 0.0

	app.renderer.Vertices[9]  = 300 
	app.renderer.Vertices[10] = 300 
	app.renderer.Vertices[11] = 0.0

	app.renderer.Indices[0]  = 0
	app.renderer.Indices[1]  = 1
	app.renderer.Indices[2]  = 2
	app.renderer.Indices[3]  = 2
	app.renderer.Indices[4]  = 3
	app.renderer.Indices[5]  = 0

	window, err := app.Init()
	gl.Viewport(0, 0, WindowWidth, WindowHeight)
	if err != nil {
		panic(err)
	}	
	defer app.Terminate(window)

	n := 0

	renderSystem := func (ecs *entity.ECS, queries []entity.EntityType, entities []roaring64.Bitmap, queriesMut []entity.EntityType, entitiesMut []roaring64.Bitmap) {
		// Render logic here, it sends work to the queue
		app.renderer.RenderLogic(window, 
							 renderWorkQueue, 
							 ecs, 
							 queries, 
							 entities, 
							 queriesMut, 
							 entitiesMut)

		n++
	}
	app.ecs.RegisterSystem("render", renderSystem)
	app.ecs.RegisterQueries("render", entity.TILE)

	app.logger.Printf("Running app %s\n", app.name)
	go app.Main(stopButton) 


	timeout := make(chan struct{})
	go func() {
		time.Sleep(20 * time.Second)
		timeout <- struct{}{}
	}()

	// Let the render queue take over the main thread 
	// This can also also us to handle other plugins that require thread-local state
	for !window.ShouldClose() {
		select {
		case work := <-renderWorkQueue:
			work()
		case <-stopButton:
			break
		case <-timeout:
			return nil	
		}
	}

	return nil
}
