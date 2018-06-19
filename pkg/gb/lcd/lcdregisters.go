package lcd

// FF40 - LCDC - LCD Control (R/W)
// Bit 7 - LCD Display Enable             (0=Off, 1=On)
func (lcd *LCD) lcdDisplayEnable() bool {
	return lcd.hwr.LCDC&0x80 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
func (lcd *LCD) highWindowTileMapDisplaySelect() bool {
	return lcd.hwr.LCDC&0x40 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 5 - Window Display Enable          (0=Off, 1=On)
func (lcd *LCD) windowDisplayEnable() bool {
	return lcd.hwr.LCDC&0x20 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
func (lcd *LCD) lowTileDataSelect() bool {
	return lcd.hwr.LCDC&0x10 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
func (lcd *LCD) highBgTileMapDisplaySelect() bool {
	return lcd.hwr.LCDC&0x08 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
func (lcd *LCD) largeSprites() bool {
	return lcd.hwr.LCDC&0x04 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
func (lcd *LCD) spriteDisplayEnable() bool {
	return lcd.hwr.LCDC&0x02 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)
func (lcd *LCD) bgDisplayEnable() bool {
	return lcd.hwr.LCDC&0x01 > 0
}
