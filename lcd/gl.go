package lcd

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"runtime"

	"github.com/scottyw/goomba/mem"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(2)

	// we need to keep OpenGL calls on a single thread
	runtime.LockOSThread()
}

func onKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		if key == glfw.KeySpace {
			fmt.Println("Pressed space")
		}
	}
	fmt.Printf("Pressed %v\n", key)
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
	s1 := float32(w) / 256
	s2 := float32(h) / 256
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
	switch pixel {
	case 0:
		im.SetRGBA(x, y, color.RGBA{0xff, 0xff, 0xff, 0xff})
	case 1:
		im.SetRGBA(x, y, color.RGBA{0xa8, 0xa8, 0xa8, 0xff})
	case 2:
		im.SetRGBA(x, y, color.RGBA{0x54, 0x54, 0x54, 0xff})
	case 3:
		im.SetRGBA(x, y, color.RGBA{0x00, 0x00, 0x00, 0xff})
	default:
		panic(fmt.Sprintf("Bad pixel: %v", pixel))
	}
}

func makeImage(mem mem.Memory) *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, 256, 256))
	var i uint16
	for i = 0x0000; i < 0x4000; i += 16 {
		tile := (int(i) - 0x0000) / 16
		xOffset := (tile % 32) * 8
		yOffset := (tile / 32) * 8
		for y := 0; y < 8; y++ {
			a := *mem.Read(i + uint16((y * 2)))
			b := *mem.Read(i + uint16((y*2)+1))
			for x := 0; x < 8; x++ {
				pixel := (a>>uint(7-x))&1 | ((b>>uint(7-x))&1)<<1
				renderPixel(im, xOffset+x, yOffset+y, pixel)
			}
		}
	}
	return im
}

// Run the render frantically, not merely once per frame ...
func Run(mem mem.Memory) {

	// initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	// create window
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, 0)
	window, err := glfw.CreateWindow(512, 512, "goomba", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	window.MakeContextCurrent()

	// initialize gl
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	gl.Enable(gl.TEXTURE_2D)

	window.SetKeyCallback(onKey)
	texture := createTexture()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		image := makeImage(mem)
		setTexture(image)
		drawBuffer(window)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		window.SwapBuffers()
		glfw.PollEvents()
	}

}
