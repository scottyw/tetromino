package lcd

import "github.com/scottyw/goomba/mem"

const (
	lcdcReg = 0xff40
	statReg = 0xff41
	lyReg   = 0xff44
)

// FF40 - LCDC - LCD Control (R/W)
// Bit 7 - LCD Display Enable             (0=Off, 1=On)
func lcdDisplayEnable(mem mem.Memory) bool {
	return *mem.Read(lcdcReg)&0x80 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
func highWindowTileMapDisplaySelect(mem mem.Memory) bool {
	return *mem.Read(lcdcReg)&0x40 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 5 - Window Display Enable          (0=Off, 1=On)
func windowDisplayEnable(mem mem.Memory) bool {
	return *mem.Read(lcdcReg)&0x20 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
func lowTileDataSelect(mem mem.Memory) bool {
	return *mem.Read(lcdcReg)&0x10 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
func highBgTileMapDisplaySelect(mem mem.Memory) bool {
	return *mem.Read(lcdcReg)&0x08 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
func largeSpriteSize(mem mem.Memory) bool {
	return *mem.Read(lcdcReg)&0x04 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
func spriteDisplayEnable(mem mem.Memory) bool {
	return *mem.Read(lcdcReg)&0x02 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)
func bgDisplay(mem mem.Memory) bool {
	return *mem.Read(lcdcReg)&0x01 > 0
}
