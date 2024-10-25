package window

import "github.com/veandco/go-sdl2/sdl"

type WindowManager struct {
  primaryWindow *sdl.Window
}

func (wm *WindowManager) OpenWindow(name string, height, width int, setPrimary bool) (*sdl.Window, error) {
	var window *sdl.Window
	var err error

	sdl.Do(func() {
		window, err = sdl.CreateWindow(
			name,
			sdl.WINDOWPOS_UNDEFINED,
			sdl.WINDOWPOS_UNDEFINED,
			int32(width),
			int32(height),
			sdl.WINDOW_SHOWN,
		)
	})

	if err != nil {
		return nil, err
	}

	if setPrimary {
		wm.primaryWindow = window
	}

	return window, nil
}

func (self *WindowManager) DestroyWindows() {
	self.primaryWindow.Destroy()
}
