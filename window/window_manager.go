package window

import "github.com/go-gl/glfw/v3.3/glfw"

type Window struct {
  height int
  width int
}

type WindowManager struct {
  primaryWindow Window
}

func (wm *WindowManager) Init(height, width int) error {
  wm.primaryWindow = Window{height: height, width: width}

  err := glfw.Init()
  if err != nil {
    return err
  }
  
  window, err := glfw.CreateWindow(width, height, "Primary", nil, nil)
  if err != nil {
    return err
  }
  window.MakeContextCurrent()

  for !window.ShouldClose() {
    window.SwapBuffers()
    glfw.PollEvents()
  }

  return nil
}
