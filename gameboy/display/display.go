package display

import (
	"fmt"
	"image"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/scottyw/tetromino/gameboy/controller"
)

// Display implements the LCD display using GL
type Display struct {
	window *glfw.Window
}

// New implements an LCD display in GL
func New(controller *controller.Controller, onInput func(), debug bool) *Display {

	if err := glfw.Init(); err != nil {
		panic(fmt.Sprintf("Failed to create display: %v", err))
	}

	var width float32
	var height float32
	if debug {
		width = 256
		height = 256
	} else {
		width = 160
		height = 144
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, 0)
	window, err := glfw.CreateWindow(int(width*3), int(height*3), "Tetromino", nil, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create display: %v", err))
	}
	window.MakeContextCurrent()

	// Max out speed - framerate is controlled by the audio subsystem
	glfw.SwapInterval(0)

	if err := gl.Init(); err != nil {
		panic(fmt.Sprintf("Failed to create display: %v", err))
	}
	gl.Enable(gl.TEXTURE_2D)
	window.SetKeyCallback(onKeyFunc(controller, onInput))

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	display := &Display{
		window: window,
	}
	return display
}

// Cleanup returns resources to the OS
func (d *Display) Cleanup() {
	glfw.Terminate()
}

// RenderFrame draws a frame to the GL window and returns user input
func (d *Display) RenderFrame(image *image.RGBA) bool {
	size := image.Rect.Size()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(image.Pix))
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(1, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(1, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-1, 1)
	gl.End()
	d.window.SwapBuffers()
	glfw.PollEvents()
	return d.window.ShouldClose()
}

func onKeyFunc(c *controller.Controller, onInput func()) func(*glfw.Window, glfw.Key, int, glfw.Action, glfw.ModifierKey) {
	return func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Press && action != glfw.Release {
			return
		}
		switch key {
		case glfw.KeyA:
			c.ButtonAction(controller.Start, action == glfw.Press)
			onInput()
		case glfw.KeyS:
			c.ButtonAction(controller.Select, action == glfw.Press)
			onInput()
		case glfw.KeyZ:
			c.ButtonAction(controller.B, action == glfw.Press)
			onInput()
		case glfw.KeyX:
			c.ButtonAction(controller.A, action == glfw.Press)
			onInput()
		case glfw.KeyUp:
			c.ButtonAction(controller.Up, action == glfw.Press)
			onInput()
		case glfw.KeyDown:
			c.ButtonAction(controller.Down, action == glfw.Press)
			onInput()
		case glfw.KeyLeft:
			c.ButtonAction(controller.Left, action == glfw.Press)
			onInput()
		case glfw.KeyRight:
			c.ButtonAction(controller.Right, action == glfw.Press)
			onInput()
		case glfw.KeyT:
			if action == glfw.Press {
				c.EmulatorAction(controller.TakeScreenshot)
			}
		}
	}
}

func init() {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(2)

	// we need to keep OpenGL calls on a single thread
	runtime.LockOSThread()
}
