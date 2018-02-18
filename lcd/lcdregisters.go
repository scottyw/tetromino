package lcd

import "github.com/scottyw/goomba/mem"

// FF40 - LCDC - LCD Control (R/W)
// Bit 7 - LCD Display Enable             (0=Off, 1=On)
func lcdDisplayEnable(memory mem.Memory) bool {
	return *memory.Read(mem.LCDC)&0x80 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
func highWindowTileMapDisplaySelect(memory mem.Memory) bool {
	return *memory.Read(mem.LCDC)&0x40 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 5 - Window Display Enable          (0=Off, 1=On)
func windowDisplayEnable(memory mem.Memory) bool {
	return *memory.Read(mem.LCDC)&0x20 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
func lowTileDataSelect(memory mem.Memory) bool {
	return *memory.Read(mem.LCDC)&0x10 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
func highBgTileMapDisplaySelect(memory mem.Memory) bool {
	return *memory.Read(mem.LCDC)&0x08 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
func largeSpriteSize(memory mem.Memory) bool {
	return *memory.Read(mem.LCDC)&0x04 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
func spriteDisplayEnable(memory mem.Memory) bool {
	return *memory.Read(mem.LCDC)&0x02 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)
func bgDisplay(memory mem.Memory) bool {
	return *memory.Read(mem.LCDC)&0x01 > 0
}
