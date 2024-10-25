package app

import (
	"log"
	"os"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/phdavis1027/goecs/entity"
	"github.com/phdavis1027/goecs/render"
	"github.com/phdavis1027/goecs/window"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
)

type App struct {
	name           string
	logger         *log.Logger
	ecs            *entity.ECS
    windowManager  *window.WindowManager
	renderer       *render.Renderer
}

func NewApp(name string, ecsCap int) *App {
	return &App{
		logger:           log.New(os.Stdout, name, log.LstdFlags),
		ecs:              entity.CreateEcsOfCapacity(ecsCap),
    	windowManager:    new(window.WindowManager),
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

func (app *App) Init() error {
	// initialize GL
	err := glfw.Init()
	if err != nil {
		return err
	}

	err = gl.Init()
	if err != nil {
		return err
	}

	err = app.renderer.Init()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Terminate() {
	// Clean up resources
	glfw.Terminate()
}

func (app *App) Run() error {
	// True => make primary window
	app.Init()
	defer app.Terminate()

	// GLFW things
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, app.name, nil, nil)

	if err != nil {
		return err
	}
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {	
		gl.Viewport(0, 0, int32(width), int32(height))
	})
	gl.Viewport(0, 0, WindowWidth, WindowHeight)

	var renderWorkQueue  = make(chan func())
	var stopButton       = make(chan struct{})

	renderer             := render.NewRenderer()
	renderer.Vertices[0]  = -0.5
	renderer.Vertices[1]  = -0.5
	renderer.Vertices[2]  =  0.0
	renderer.Vertices[3]  =  0.5
	renderer.Vertices[4]  = -0.5
	renderer.Vertices[5]  =  0.0
	renderer.Vertices[6]  =  0.5
	renderer.Vertices[7]  =  0.5
	renderer.Vertices[8]  =  0.0
	renderer.Vertices[9]  = -0.5
	renderer.Vertices[10] =  0.5
	renderer.Vertices[11] =  0.0

	renderer.Indices[0]  = 0
	renderer.Indices[1]  = 1
	renderer.Indices[2]  = 2
	renderer.Indices[3]  = 2
	renderer.Indices[4]  = 3
	renderer.Indices[5]  = 0

	renderSystem := func (ecs *entity.ECS, queries []entity.EntityType, entities []roaring64.Bitmap, queriesMut []entity.EntityType, entitiesMut []roaring64.Bitmap) {
		// Render logic here, it sends work to the queue
		renderer.RenderLogic(renderWorkQueue, ecs, queries, entities, queriesMut, entitiesMut)
	}
	app.ecs.RegisterSystem("render", renderSystem)
	app.ecs.RegisterQueries("render", entity.TILE)


	app.logger.Printf("Running app %s\n", app.name)
	go app.Main(stopButton) 

	// Let the render queue take over the main thread 
	for {
		select {
		case work := <-renderWorkQueue:
			work()
		case <-stopButton:
			break
		}
	}

	return nil
}
