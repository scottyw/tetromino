package ppu

// FF40 - LCDC - LCD Control (R/W)
// Bit 7 - LCD Display Enable             (0=Off, 1=On)
// Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
// Bit 5 - Window Display Enable          (0=Off, 1=On)
// Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
// Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
// Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
// Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
// Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)

// WriteLCDC handles writes to register LCDC
func (ppu *PPU) WriteLCDC(value uint8) {
	// fmt.Printf("> LCDC - 0x%02x\n", value)
	ppu.enabled = value&0x80 > 0
	ppu.highWindowTileMap = value&0x40 > 0
	ppu.windowEnabled = value&0x20 > 0
	ppu.lowTileData = value&0x10 > 0
	ppu.highBgTileMap = value&0x08 > 0
	ppu.spritesLarge = value&0x04 > 0
	ppu.spritesEnabled = value&0x02 > 0
	ppu.bgEnabled = value&0x01 > 0
}

// ReadLCDC handles reads from register LCDC
func (ppu *PPU) ReadLCDC() uint8 {
	var lcdc uint8
	if ppu.enabled {
		lcdc += 0x80
	}
	if ppu.highWindowTileMap {
		lcdc += 0x40
	}
	if ppu.windowEnabled {
		lcdc += 0x20
	}
	if ppu.lowTileData {
		lcdc += 0x10
	}
	if ppu.highBgTileMap {
		lcdc += 0x08
	}
	if ppu.spritesLarge {
		lcdc += 0x04
	}
	if ppu.spritesEnabled {
		lcdc += 0x02
	}
	if ppu.bgEnabled {
		lcdc += 0x01
	}
	// fmt.Printf("< LCDC - 0x%02x\n", lcdc )
	return lcdc
}

// LCD Status Register
// -------------------

// FF41 - STAT - LCDC Status   (R/W)
// Bit 6 - LYC=LY Coincidence Interrupt (1=Enable) (Read/Write)
// Bit 5 - Mode 2 OAM Interrupt         (1=Enable) (Read/Write)
// Bit 4 - Mode 1 V-Blank Interrupt     (1=Enable) (Read/Write)
// Bit 3 - Mode 0 H-Blank Interrupt     (1=Enable) (Read/Write)
// Bit 2 - Coincidence Flag  (0:LYC<>LY, 1:LYC=LY) (Read Only)
// Bit 1-0 - Mode Flag       (Mode 0-3, see below) (Read Only)
// 		  0: During H-Blank
// 		  1: During V-Blank
// 		  2: During Searching OAM-RAM
// 		  3: During Transfering Data to LCD Driver

// WriteSTAT handles writes to register STAT
func (ppu *PPU) WriteSTAT(value uint8) {
	// fmt.Printf("> STAT - 0x%02x\n", value)
	ppu.stat = value
}

// ReadSTAT handles reads from register STAT
func (ppu *PPU) ReadSTAT() uint8 {
	// First bit is always high
	stat := ppu.stat | 0x80
	// fmt.Printf("< STAT - 0x%02x\n", stat )
	return stat
}

// FF42 - SCY - Scroll Y   (R/W)
// FF43 - SCX - Scroll X   (R/W)
// Specifies the position in the 256x256 pixels BG map (32x32 tiles) which is to
// be displayed at the upper/left LCD display position.
// Values in range from 0-255 may be used for X/Y each, the video controller
// automatically wraps back to the upper (left) position in BG map when drawing
// exceeds the lower (right) border of the BG map area.

// WriteSCY handles writes to register SCY
func (ppu *PPU) WriteSCY(value uint8) {
	// fmt.Printf("> SCY - 0x%02x\n", value)
	ppu.scy = value
}

// ReadSCY handles reads from register SCY
func (ppu *PPU) ReadSCY() uint8 {
	scy := ppu.scy
	// fmt.Printf("< SCY - 0x%02x\n", scy )
	return scy
}

// WriteSCX handles writes to register SCX
func (ppu *PPU) WriteSCX(value uint8) {
	// fmt.Printf("> SCX - 0x%02x\n", value)
	ppu.scx = value
}

// ReadSCX handles reads from register SCX
func (ppu *PPU) ReadSCX() uint8 {
	scx := ppu.scx
	// fmt.Printf("< SCX - 0x%02x\n", scx )
	return scx
}

// FF44 - LY - LCDC Y-Coordinate (R)
// The LY indicates the vertical line to which the present data is transferred to
// the LCD Driver. The LY can take on any value between 0 through 153. The values
// between 144 and 153 indicate the V-Blank period. Writing will reset the
// counter.

// WriteLY handles writes to register LY
func (ppu *PPU) WriteLY(value uint8) {
	// fmt.Printf("> LY - 0x%02x\n", value)
	ppu.ly = value
}

// ReadLY handles reads from register LY
func (ppu *PPU) ReadLY() uint8 {
	ly := ppu.ly
	// fmt.Printf("< LY - 0x%02x\n", ly )
	return ly
}

// FF45 - LYC - LY Compare  (R/W)
// The gameboy permanently compares the value of the LYC and LY registers. When
// both values are identical, the coincident bit in the STAT register becomes
// set, and (if enabled) a STAT interrupt is requested.

// WriteLYC handles writes to register LYC
func (ppu *PPU) WriteLYC(value uint8) {
	// fmt.Printf("> LYC - 0x%02x\n", value)
	ppu.lyc = value
}

// ReadLYC handles reads from register LYC
func (ppu *PPU) ReadLYC() uint8 {
	lyc := ppu.lyc
	// fmt.Printf("< LYC - 0x%02x\n", lyc )
	return lyc
}

// FF4A - WY - Window Y Position (R/W)
// FF4B - WX - Window X Position minus 7 (R/W)
// Specifies the upper/left positions of the Window area. (The window is an
// alternate background area which can be displayed above of the normal
// background. OBJs (sprites) may be still displayed above or behinf the window,
// just as for normal BG.)
// The window becomes visible (if enabled) when positions are set in range
// WX=0..166, WY=0..143. A postion of WX=7, WY=0 locates the window at upper
// left, it is then completly covering normal background.

// WriteWY handles writes to register WY
func (ppu *PPU) WriteWY(value uint8) {
	// fmt.Printf("> WY - 0x%02x\n", value)
	ppu.wy = value
}

// ReadWY handles reads from register WY
func (ppu *PPU) ReadWY() uint8 {
	wy := ppu.wy
	// fmt.Printf("< WY - 0x%02x\n", wy )
	return wy
}

// WriteWX handles writes to register WX
func (ppu *PPU) WriteWX(value uint8) {
	// fmt.Printf("> WX - 0x%02x\n", value)
	ppu.wx = value
}

// ReadWX handles reads from register WX
func (ppu *PPU) ReadWX() uint8 {
	wx := ppu.wx
	// fmt.Printf("< WX - 0x%02x\n", wx )
	return wx
}

// FF47 - BGP - BG Palette Data  (R/W) - Non CGB Mode Only
// This register assigns gray shades to the color numbers of the BG and Window
// tiles.
// Bit 7-6 - Shade for Color Number 3
// Bit 5-4 - Shade for Color Number 2
// Bit 3-2 - Shade for Color Number 1
// Bit 1-0 - Shade for Color Number 0
// The four possible gray shades are:
// 0  White
// 1  Light gray
// 2  Dark gray
// 3  Black
// In CGB Mode the Color Palettes are taken from CGB Palette Memory instead.

// WriteBGP handles writes to register BGP
func (ppu *PPU) WriteBGP(value uint8) {
	// fmt.Printf("> BGP - 0x%02x\n", value)
	ppu.bgp = value
}

// ReadBGP handles reads from register BGP
func (ppu *PPU) ReadBGP() uint8 {
	bgp := ppu.bgp
	// fmt.Printf("< BGP - 0x%02x\n", bgp )
	return bgp
}

// FF48 - OBP0 - Object Palette 0 Data (R/W) - Non CGB Mode Only
// This register assigns gray shades for sprite palette 0. It works exactly as
// BGP (FF47), except that the lower two bits aren't used because sprite data 00
// is transparent.

// WriteOBP0 handles writes to register OBP0
func (ppu *PPU) WriteOBP0(value uint8) {
	// fmt.Printf("> OBP0 - 0x%02x\n", value)
	ppu.obp0 = value
}

// ReadOBP0 handles reads from register OBP0
func (ppu *PPU) ReadOBP0() uint8 {
	obp0 := ppu.obp0
	// fmt.Printf("< OBP0 - 0x%02x\n", obp0 )
	return obp0
}

// FF49 - OBP1 - Object Palette 1 Data (R/W) - Non CGB Mode Only
// This register assigns gray shades for sprite palette 1. It works exactly as
// BGP (FF47), except that the lower two bits aren't used because sprite data 00
// is transparent.

// WriteOBP1 handles writes to register OBP1
func (ppu *PPU) WriteOBP1(value uint8) {
	// fmt.Printf("> OBP1 - 0x%02x\n", value)
	ppu.obp1 = value
}

// ReadOBP1 handles reads from register OBP1
func (ppu *PPU) ReadOBP1() uint8 {
	obp1 := ppu.obp1
	// fmt.Printf("< OBP1 - 0x%02x\n", obp1 )
	return obp1
}
