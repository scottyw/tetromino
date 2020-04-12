package ui

import (
	"context"
	"image"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/scottyw/tetromino/pkg/gb"
)

// GLDisplay implements the LCD display using GL
type GLDisplay struct {
	cancelFunc context.CancelFunc
	window     *glfw.Window
	texture    uint32
	width      float32
	height     float32
}

// NewGLDisplay implements an LCD display in GL
func NewGLDisplay(gameboy *gb.Gameboy, cancelFunc context.CancelFunc) (*GLDisplay, error) {
	// initialize glfw
	if err := glfw.Init(); err != nil {
		return nil, err
	}
	// define window width
	var width float32
	var height float32
	if gameboy.Debug() {
		width = 256
		height = 256
	} else {
		width = 160
		height = 144
	}
	// create window
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, 0)
	window, err := glfw.CreateWindow(int(width*3), int(height*3), "Tetromino", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	// For now let's max out speed and worry about locking the framerate later
	glfw.SwapInterval(0)

	// initialize gl
	if err := gl.Init(); err != nil {
		return nil, err
	}
	gl.Enable(gl.TEXTURE_2D)
	window.SetKeyCallback(onKeyFunc(gameboy))
	display := &GLDisplay{
		cancelFunc: cancelFunc,
		window:     window,
		texture:    createTexture(),
		width:      width,
		height:     height,
	}
	return display, nil
}

// Cleanup returns resources to the OS
func (d *GLDisplay) Cleanup() {
	glfw.Terminate()
}

// DisplayFrame draws a frame to the GL window and returns user input
func (d *GLDisplay) DisplayFrame(image *image.RGBA) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BindTexture(gl.TEXTURE_2D, d.texture)
	setTexture(image)
	drawBuffer(d.window, d.width, d.height)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	d.window.SwapBuffers()
	glfw.PollEvents()
	if d.window.ShouldClose() {
		d.cancelFunc()
	}
}

func onKeyFunc(gameboy *gb.Gameboy) func(*glfw.Window, glfw.Key, int, glfw.Action, glfw.ModifierKey) {
	return func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Press && action != glfw.Release {
			return
		}
		switch key {
		case glfw.KeyA:
			gameboy.ButtonAction(gb.Start, action == glfw.Press)
		case glfw.KeyS:
			gameboy.ButtonAction(gb.Select, action == glfw.Press)
		case glfw.KeyZ:
			gameboy.ButtonAction(gb.B, action == glfw.Press)
		case glfw.KeyX:
			gameboy.ButtonAction(gb.A, action == glfw.Press)
		case glfw.KeyUp:
			gameboy.ButtonAction(gb.Up, action == glfw.Press)
		case glfw.KeyDown:
			gameboy.ButtonAction(gb.Down, action == glfw.Press)
		case glfw.KeyLeft:
			gameboy.ButtonAction(gb.Left, action == glfw.Press)
		case glfw.KeyRight:
			gameboy.ButtonAction(gb.Right, action == glfw.Press)
		case glfw.KeyT:
			if action == glfw.Press {
				gameboy.EmulatorAction(gb.TakeScreenshot)
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

func createTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

func setTexture(im *image.RGBA) {
	size := im.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
}

func drawBuffer(window *glfw.Window, width, height float32) {
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / width
	s2 := float32(h) / height
	f := float32(1 - 0)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, height/256.0)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(width/256.0, height/256.0)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(width/256.0, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}
