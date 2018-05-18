package lcd

import "github.com/scottyw/tetromino/pkg/gb/mem"

// FF40 - LCDC - LCD Control (R/W)
// Bit 7 - LCD Display Enable             (0=Off, 1=On)
func lcdDisplayEnable(hwr *mem.HardwareRegisters) bool {
	return hwr.LCDC&0x80 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
func highWindowTileMapDisplaySelect(hwr *mem.HardwareRegisters) bool {
	return hwr.LCDC&0x40 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 5 - Window Display Enable          (0=Off, 1=On)
func windowDisplayEnable(hwr *mem.HardwareRegisters) bool {
	return hwr.LCDC&0x20 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
func lowTileDataSelect(hwr *mem.HardwareRegisters) bool {
	return hwr.LCDC&0x10 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
func highBgTileMapDisplaySelect(hwr *mem.HardwareRegisters) bool {
	return hwr.LCDC&0x08 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
func largeSpriteSize(hwr *mem.HardwareRegisters) bool {
	return hwr.LCDC&0x04 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
func spriteDisplayEnable(hwr *mem.HardwareRegisters) bool {
	return hwr.LCDC&0x02 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)
func bgDisplay(hwr *mem.HardwareRegisters) bool {
	return hwr.LCDC&0x01 > 0
}
