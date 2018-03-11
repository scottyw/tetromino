package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/scottyw/tetromino/lcd"
	"github.com/scottyw/tetromino/mem"
)

// GL maintains state for the GL UI implementation
type GL struct {
	window  *glfw.Window
	texture uint32
}

// NewGL implements a user interface in GL
func NewGL(memory mem.Memory) UI {
	// initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln(err)
	}

	// create window
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, 0)
	window, err := glfw.CreateWindow(640, 576, "Tetromino", nil, nil)
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

	window.SetKeyCallback(onKeyFunc(memory.JOYP))
	return &GL{
		window:  window,
		texture: createTexture(),
	}
}

// ShouldRun indicates whether the emulator should be running e.g. stop when the GL window is closed
func (glx *GL) ShouldRun() bool {
	return !glx.window.ShouldClose()
}

// Shutdown the GL framework
func (glx *GL) Shutdown() {
	glfw.Terminate()
}

// DrawFrame draws a frame to the GL window
func (glx *GL) DrawFrame(lcd *lcd.LCD) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BindTexture(gl.TEXTURE_2D, glx.texture)
	image := renderFrame(lcd.FrameData())
	setTexture(image)
	drawBuffer(glx.window)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	glx.window.SwapBuffers()
	glfw.PollEvents()
}

func onKeyFunc(joyp *uint8) func(*glfw.Window, glfw.Key, int, glfw.Action, glfw.ModifierKey) {
	return func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// Bit 5 - P15 Select Button Keys      (0=Select)
		// Bit 4 - P14 Select Direction Keys   (0=Select)
		// Bit 3 - P13 Input Down  or Start    (0=Pressed) (Read Only)
		// Bit 2 - P12 Input Up    or Select   (0=Pressed) (Read Only)
		// Bit 1 - P11 Input Left  or Button B (0=Pressed) (Read Only)
		// Bit 0 - P10 Input Right or Button A (0=Pressed) (Read Only)
		if (^*joyp)&0x20 > 0 {
			buttonInput(key, action, joyp)
		}
		if (^*joyp)&0x10 > 0 {
			directionInput(key, action, joyp)
		}
	}
}

func buttonInput(key glfw.Key, action glfw.Action, joyp *uint8) {
	switch key {
	case glfw.KeyA:
		if action == glfw.Press {
			*joyp &^= 0x8
		} else if action == glfw.Release {
			*joyp |= 0x8
		}
	case glfw.KeyS:
		if action == glfw.Press {
			*joyp &^= 0x4
		} else if action == glfw.Release {
			*joyp |= 0x4
		}
	case glfw.KeyZ:
		if action == glfw.Press {
			*joyp &^= 0x2
		} else if action == glfw.Release {
			*joyp |= 0x2
		}
	case glfw.KeyX:
		if action == glfw.Press {
			*joyp &^= 0x1
		} else if action == glfw.Release {
			*joyp |= 0x1
		}
	}
}

func directionInput(key glfw.Key, action glfw.Action, joyp *uint8) {
	switch key {
	case glfw.KeyDown:
		if action == glfw.Press {
			*joyp &^= 0x8
		} else if action == glfw.Release {
			*joyp |= 0x8
		}
	case glfw.KeyUp:
		if action == glfw.Press {
			*joyp &^= 0x4
		} else if action == glfw.Release {
			*joyp |= 0x4
		}
	case glfw.KeyLeft:
		if action == glfw.Press {
			*joyp &^= 0x2
		} else if action == glfw.Release {
			*joyp |= 0x2
		}
	case glfw.KeyRight:
		if action == glfw.Press {
			*joyp &^= 0x1
		} else if action == glfw.Release {
			*joyp |= 0x1
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

func drawBuffer(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / 160
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
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}

func renderPixel(im *image.RGBA, x, y int, pixel uint8) {
	// pixel = pixel % 0x10 // Remove colour offset
	switch pixel {
	case 0x00:
		im.SetRGBA(x, y, color.RGBA{0xff, 0xff, 0xff, 0xff})
	case 0x01:
		im.SetRGBA(x, y, color.RGBA{0xaa, 0xaa, 0xaa, 0xff})
	case 0x02:
		im.SetRGBA(x, y, color.RGBA{0x77, 0x77, 0x77, 0xff})
	case 0x03:
		im.SetRGBA(x, y, color.RGBA{0x33, 0x33, 0x33, 0xff})
	case 0x10:
		im.SetRGBA(x, y, color.RGBA{0xff, 0xaa, 0xaa, 0xff})
	case 0x11:
		im.SetRGBA(x, y, color.RGBA{0xaa, 0x77, 0x77, 0xff})
	case 0x12:
		im.SetRGBA(x, y, color.RGBA{0x77, 0x33, 0x33, 0xff})
	case 0x13:
		im.SetRGBA(x, y, color.RGBA{0x33, 0x00, 0x00, 0xff})
	case 0x20:
		im.SetRGBA(x, y, color.RGBA{0xaa, 0xff, 0xaa, 0xff})
	case 0x21:
		im.SetRGBA(x, y, color.RGBA{0x77, 0xaa, 0x77, 0xff})
	case 0x22:
		im.SetRGBA(x, y, color.RGBA{0x33, 0x77, 0x33, 0xff})
	case 0x23:
		im.SetRGBA(x, y, color.RGBA{0x00, 0x33, 0x00, 0xff})
	case 0x30:
		im.SetRGBA(x, y, color.RGBA{0xaa, 0xaa, 0xff, 0xff})
	case 0x31:
		im.SetRGBA(x, y, color.RGBA{0x77, 0x77, 0xaa, 0xff})
	case 0x32:
		im.SetRGBA(x, y, color.RGBA{0x33, 0x33, 0x77, 0xff})
	case 0x33:
		im.SetRGBA(x, y, color.RGBA{0x00, 0x00, 0x33, 0xff})
	default:
		panic(fmt.Sprintf("Bad pixel: %v", pixel))
	}
}

func renderFrame(data [23040]uint8) *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, 160, 144))
	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			pixel := data[y*160+x]
			renderPixel(im, x, y, pixel)
		}
	}
	return im
}
