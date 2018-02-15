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
	*mem.Read(statReg) = stat
	*mem.Read(lyReg) = ly
}

// FrameData returns the frame data as a 256x256 array of bytes where each element is a colour value between 0 and 3
func (lcd *LCD) FrameData(mem mem.Memory) [65536]uint8 {
	lcd.drawTiles(mem, highBgTileMapDisplaySelect)
	if windowDisplayEnable(mem) {
		lcd.drawTiles(mem, highWindowTileMapDisplaySelect)
	}
	return lcd.data
}

// Returns 16 bytes representing one 8x8 tile
func tileData(mem mem.Memory, tile uint16, displaySelect func(mem.Memory) bool) []byte {
	var tileAddr uint16
	if displaySelect(mem) {
		tileAddr = 0x9c00 + tile
	} else {
		tileAddr = 0x9800 + tile
	}
	tileIndex := mem.Read(tileAddr)
	if lowTileDataSelect(mem) {
		return mem.ReadRegion(uint16(0x8000+uint(*tileIndex)), 16)
	}
	return mem.ReadRegion(uint16(0x9000+int(*tileIndex)), 16)
}

func (lcd *LCD) drawTiles(mem mem.Memory, displaySelect func(mem.Memory) bool) {
	var x, y, row, col uint16
	var pixel uint8
	for y = 0; y < 32; y++ {
		for x = 0; x < 32; x++ {
			tile := tileData(mem, y*32+x, displaySelect)
			for row = 0; row < 8; row++ {
				a := tile[row*2]
				b := tile[row*2+1]
				for col = 0; col < 8; col++ {
					pixel = (a>>uint(7-col))&1 | ((b>>uint(7-col))&1)<<1
					index := (((y * 8) + row) * 256) + ((x * 8) + col)
					lcd.data[index] = pixel
				}
			}
		}
	}
}
