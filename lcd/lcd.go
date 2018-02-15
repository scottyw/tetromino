package lcd

import (
	"github.com/scottyw/goomba/mem"
)

// LCD represents the LCD display of the Gameboy
type LCD struct {
	data [65536]uint8
}

// NewLCD returns the configured LCD
func NewLCD() *LCD {
	return &LCD{}
}

// Tick runs the LCD driver for one machine cycle i.e. 4 clock cycles
func (lcd *LCD) Tick(mem mem.Memory, cycle int) {
	ly := uint8(cycle / 114)
	lyRemainder := cycle % 114
	var stat uint8
	// Set mode on stat register
	switch {
	case ly >= 144:
		stat = 1
	case lyRemainder < 20:
		stat = 2
	case lyRemainder < 63:
		stat = 3
	case lyRemainder < 114:
		stat = 0
	default:
		panic("LCD driver error setting mode")
	}
	// Set coincidence flag and coincidence interrupt on stat register
	if ly == uint8(*mem.Read(0xff45)) {
		stat |= 0x44
	} else {
		stat &^= 0x44
	}
	// Set interrupts on stat register
	switch {
	case ly == 144:
		stat |= 0x10
	case lyRemainder == 0:
		stat |= 0x20
	case lyRemainder == 63:
		stat |= 0x08
	}
	*mem.Read(0xff41) = stat // STAT register
	*mem.Read(0xff44) = ly   // LY register
}

// FrameData returns the frame data as a 256x256 array of bytes where each element is a colour value between 0 and 3
func (lcd *LCD) FrameData() [65536]uint8 {
	lcd.drawBackground()
	return lcd.data
}

func (lcd *LCD) drawBackground() {

	// a := *mem.Read(i + uint16((y * 2)))
	// 		b := *mem.Read(i + uint16((y*2)+1))
	// 		for x := 0; x < 8; x++ {
	// 			pixel := (a>>uint(7-x))&1 | ((b>>uint(7-x))&1)<<1
	// 		}

	for i := 0; i < 65536; i++ {
		lcd.data[i] = uint8(i % 4)
	}

}
