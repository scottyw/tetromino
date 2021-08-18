package lcd

// FF40 - LCDC - LCD Control (R/W)
// Bit 7 - LCD Display Enable             (0=Off, 1=On)
func (lcd *LCD) lcdDisplayEnable() bool {
	return lcd.memory.LCDC&0x80 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
func (lcd *LCD) highWindowTileMapDisplaySelect() bool {
	return lcd.memory.LCDC&0x40 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 5 - Window Display Enable          (0=Off, 1=On)
func (lcd *LCD) windowDisplayEnable() bool {
	return lcd.memory.LCDC&0x20 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
func (lcd *LCD) lowTileDataSelect() bool {
	return lcd.memory.LCDC&0x10 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
func (lcd *LCD) highBgTileMapDisplaySelect() bool {
	return lcd.memory.LCDC&0x08 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
func (lcd *LCD) largeSprites() bool {
	return lcd.memory.LCDC&0x04 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
func (lcd *LCD) spriteDisplayEnable() bool {
	return lcd.memory.LCDC&0x02 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)
func (lcd *LCD) bgDisplayEnable() bool {
	return lcd.memory.LCDC&0x01 > 0
}

// Byte3 - Attributes/Flags:
//   Bit7   OBJ-to-BG Priority (0=OBJ Above BG, 1=OBJ Behind BG color 1-3)
//          (Used for both BG and Window. BG color 0 is always behind OBJ)
//   Bit6   Y flip          (0=Normal, 1=Vertically mirrored)
//   Bit5   X flip          (0=Normal, 1=Horizontally mirrored)
//   Bit4   Palette number  **Non CGB Mode Only** (0=OBP0, 1=OBP1)
//   Bit3   Tile VRAM-Bank  **CGB Mode Only**     (0=Bank 0, 1=Bank 1)
//   Bit2-0 Palette number  **CGB Mode Only**     (OBP0-7)

func spriteBehindBg(att uint8) bool {
	return att&0x80 > 0
}

func spriteYFlip(att uint8) bool {
	return att&0x40 > 0
}

func spriteXFlip(att uint8) bool {
	return att&0x20 > 0
}

func spritePalatteSet(att uint8) bool {
	return att&0x08 > 0
}
