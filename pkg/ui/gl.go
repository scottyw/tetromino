package ui

import (
	"image"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

// GL maintains state for the GL UI implementation
type GL struct {
	emu     Emulator
	window  *glfw.Window
	texture uint32
	width   float32
}

// NewGL implements a user interface in GL
func NewGL(emu Emulator) *GL {
	// initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln(err)
	}
	// define window width
	var width float32
	if emu.Debug() {
		width = 256
	} else {
		width = 160
	}
	// create window
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, 0)
	window, err := glfw.CreateWindow(int(width*3), 144*3, "Tetromino", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	window.MakeContextCurrent()

	// For now let's max out speed and worry about locking the framerate later
	glfw.SwapInterval(0)

	// initialize gl
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	gl.Enable(gl.TEXTURE_2D)
	window.SetKeyCallback(onKeyFunc(emu))
	return &GL{
		emu:     emu,
		window:  window,
		texture: createTexture(),
		width:   width,
	}
}

// DrawFrame draws a frame to the GL window and returns user input
func (glx *GL) DrawFrame(image *image.RGBA) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BindTexture(gl.TEXTURE_2D, glx.texture)
	setTexture(image)
	drawBuffer(glx.window, glx.width)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	glx.window.SwapBuffers()
	glfw.PollEvents()
	if glx.window.ShouldClose() {
		glfw.Terminate()
		glx.emu.Shutdown()
	}
}

func onKeyFunc(emu Emulator) func(*glfw.Window, glfw.Key, int, glfw.Action, glfw.ModifierKey) {
	return func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Press && action != glfw.Release {
			return
		}

		// Bit 3 - P13 Input Down  or Start    (0=Pressed) (Read Only)
		// Bit 2 - P12 Input Up    or Select   (0=Pressed) (Read Only)
		// Bit 1 - P11 Input Left  or Button B (0=Pressed) (Read Only)
		// Bit 0 - P10 Input Right or Button A (0=Pressed) (Read Only)
		if key == glfw.KeyA {
			emu.ButtonAction(Start, action == glfw.Press)
		}
		if key == glfw.KeyS {
			emu.ButtonAction(Select, action == glfw.Press)
		}
		if key == glfw.KeyZ {
			emu.ButtonAction(B, action == glfw.Press)
		}
		if key == glfw.KeyX {
			emu.ButtonAction(A, action == glfw.Press)
		}
		if key == glfw.KeyUp {
			emu.ButtonAction(Up, action == glfw.Press)
		} else if key == glfw.KeyDown {
			emu.ButtonAction(Down, action == glfw.Press)
		}
		if key == glfw.KeyLeft {
			emu.ButtonAction(Left, action == glfw.Press)
		} else if key == glfw.KeyRight {
			emu.ButtonAction(Right, action == glfw.Press)
		}

		// Emulator controls
		if key == glfw.KeyT && action == glfw.Press {
			emu.Screenshot()
		}
		if key == glfw.KeyP && action == glfw.Press {
			emu.Faster()
		}
		if key == glfw.KeyO && action == glfw.Press {
			emu.Slower()
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

func drawBuffer(window *glfw.Window, width float32) {
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / width
	s2 := float32(h) / 144
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
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(width/256.0, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(width/256.0, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}
