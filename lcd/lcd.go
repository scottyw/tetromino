package lcd

import "github.com/scottyw/goomba/mem"

// Tick runs the LCD driver for one machine cycle i.e. 4 clock cycles
func Tick(mem mem.Memory, cycle int) {
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
	if ly == uint8(mem.Read(0xff45)) {
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
	mem.Write(0xff41, stat) // STAT register
	mem.Write(0xff44, ly)   // LY register
}
