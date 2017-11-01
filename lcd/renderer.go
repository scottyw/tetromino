package lcd

import (
	"fmt"
)

var winWidth, winHeight int = 1024, 256

func renderPixel(x, y, xOffset int, p byte) {
	switch p {
	case 0:
		// renderer.SetDrawColor(255, 255, 255, 255)
	case 1:
		// renderer.SetDrawColor(168, 168, 168, 255)
	case 2:
		// renderer.SetDrawColor(84, 84, 84, 255)
	case 3:
		// renderer.SetDrawColor(0, 0, 0, 255)
	default:
		panic(fmt.Sprintf("Bad pixel: %v", p))
	}
	// renderer.DrawPoint(x+xOffset, y)
}

func renderLine(x, y, xOffset int, a, b byte) {
	for j := 7; j >= 0; j-- {
		p := (a>>uint(j))&1 | ((b>>uint(j))&1)<<1
		renderPixel(x, y, xOffset, p)
		x++
	}
}

func render(mem []byte) {
	var x, y, xOffset int
	for i := 0x0000; i < 0xffff; i += 16 {
		for j := 0; j < 16; j += 2 {
			renderLine(x, y, xOffset, mem[i+j], mem[i+j+1])
			y++
		}
		if i%0x0200 == 0 {
			xOffset += 8
			y = 0
		}
	}
}
