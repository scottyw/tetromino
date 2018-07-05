package lcd

import (
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/scottyw/tetromino/pkg/gb/mem"
)

const (
	bit7 = 1 << iota
	bit6 = 1 << iota
	bit5 = 1 << iota
	bit4 = 1 << iota
	bit3 = 1 << iota
	bit2 = 1 << iota
	bit1 = 1 << iota
	bit0 = 1 << iota
)

var (
	gray = []color.RGBA{
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0xaa, 0xaa, 0xaa, 0xff},
		color.RGBA{0x77, 0x77, 0x77, 0xff},
		color.RGBA{0x33, 0x33, 0x33, 0xff},
	}

	red = []color.RGBA{
		color.RGBA{0xff, 0xaa, 0xaa, 0xff},
		color.RGBA{0xdd, 0x77, 0x77, 0xff},
		color.RGBA{0xaa, 0x33, 0x33, 0xff},
		color.RGBA{0x55, 0x00, 0x00, 0xff},
	}

	green = []color.RGBA{
		color.RGBA{0xaa, 0xff, 0xaa, 0xff},
		color.RGBA{0x77, 0xdd, 0x77, 0xff},
		color.RGBA{0x33, 0xaa, 0x33, 0xff},
		color.RGBA{0x00, 0x55, 0x00, 0xff},
	}

	blue = []color.RGBA{
		color.RGBA{0xaa, 0xaa, 0xff, 0xff},
		color.RGBA{0x77, 0x77, 0xdd, 0xff},
		color.RGBA{0x33, 0x33, 0xaa, 0xff},
		color.RGBA{0x00, 0x00, 0x55, 0xff},
	}
)

// LCD represents the LCD display of the Gameboy
type LCD struct {
	hwr            *mem.HardwareRegisters
	videoRAM       *[0x2000]byte
	oam            *[0xa0]byte
	tileCache      [384]*[8][8]uint8
	previousBg     [32][32]uint16
	previousWindow [32][32]uint16
	bg             [256][256]uint8
	window         [256][256]uint8
	sprites        [144][160]uint8
	Frame          *image.RGBA
	tick           int
	debug          bool
}

// NewLCD returns the configured LCD
func NewLCD(hwr *mem.HardwareRegisters, memory *mem.Memory, debug bool) *LCD {
	lcd := LCD{
		hwr:      hwr,
		videoRAM: &memory.VideoRAM,
		oam:      &memory.OAM,
		Frame:    image.NewRGBA(image.Rect(0, 0, 256, 256)),
		debug:    debug,
	}
	memory.WriteNotification = &lcd
	return &lcd
}

// WriteToVideoRAM implements memory write notification§
func (lcd *LCD) WriteToVideoRAM(addr uint16) {
	if addr < 0x9800 {
		tileNumber := (addr - 0x8000) / 16
		lcd.tileCache[tileNumber] = nil
	}
}

// Tick runs the LCD driver for one machine cycle i.e. 4 clock cycles
func (lcd *LCD) Tick() {

	// Is the LCD enabled?
	if !lcd.lcdDisplayEnable() {
		lcd.hwr.LY = 0
		lcd.tick = 0
		return
	}

	// Where are we on the LCD?
	lcd.hwr.LY = uint8(lcd.tick / 114)
	x := lcd.tick % 114
	lcd.tick++
	if lcd.tick >= 17556 {
		lcd.tick = 0
	}

	// Set mode on STAT register
	switch {
	case x == 0 && lcd.hwr.LY == 144:
		// V-Blank period starts
		lcd.hwr.STAT = (lcd.hwr.STAT & 0xfc) | 0x01
		// V-Blank interrupt
		lcd.hwr.IF |= 0x01
		// Is LCD STAT interrupt enabled?
		if lcd.hwr.STAT&0x10 > 0 {
			lcd.hwr.IF |= 0x02
		}
	case x == 0 && lcd.hwr.LY < 144:
		// OAM period starts
		lcd.hwr.STAT = (lcd.hwr.STAT & 0xfc) | 0x02
		// Is LCD STAT interrupt enabled?
		if lcd.hwr.STAT&0x20 > 0 {
			lcd.hwr.IF |= 0x02
		}
	case x == 20 && lcd.hwr.LY < 144:
		// LCD data transfer period starts
		lcd.hwr.STAT = (lcd.hwr.STAT & 0xfc) | 0x03
	case x == 63 && lcd.hwr.LY < 144:
		// H-Blank period starts
		lcd.hwr.STAT = (lcd.hwr.STAT & 0xfc)
		// Is LCD STAT interrupt enabled?
		if lcd.hwr.STAT&0x08 > 0 {
			lcd.hwr.IF |= 0x02
		}
		// Render LCD line
		lcd.updateLcdLine(lcd.hwr.LY)
	}

	// Check coincidence flag
	if x == 0 {
		if lcd.hwr.LY == lcd.hwr.LYC {
			lcd.hwr.STAT |= 0x04
			// Is LCD STAT interrupt enabled?
			if lcd.hwr.STAT&0x40 > 0 {
				lcd.hwr.IF |= 0x02
			}
		} else {
			lcd.hwr.STAT &^= 0x04
		}
	}
}

// TakeSnapshot writes the current contents of LCD to a file
func (lcd *LCD) TakeSnapshot() {
	file, err := os.Create("snapshot.gob")
	if err != nil {
		fmt.Printf("Failed to save LCD snapshot: %v\n", err)
		return
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	encoder.Encode(lcd)
}

// LoadSnapshot updates the current LCD based on the file contents
func (lcd *LCD) LoadSnapshot(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to load LCD snapshot: %v\n", err)
		return
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(lcd)
}

func (lcd *LCD) readVideoRAM(memoryAddr uint16) byte {
	return lcd.videoRAM[memoryAddr&0x7fff]
}

func (lcd *LCD) readTile(tileNumber uint16) (*[8][8]uint8, bool) {
	tile := lcd.tileCache[tileNumber]
	if tile != nil {
		return tile, true
	}
	tile = &[8][8]uint8{}
	startAddr := uint16(0x8000 + (tileNumber * 16))
	for y := uint16(0); y < 8; y++ {
		a := lcd.readVideoRAM(startAddr + y*2)
		b := lcd.readVideoRAM(startAddr + y*2 + 1)
		for x := uint16(0); x < 8; x++ {
			tile[y][0] = (a&bit0)>>7 | (b&bit0)>>6
			tile[y][1] = (a&bit1)>>6 | (b&bit1)>>5
			tile[y][2] = (a&bit2)>>5 | (b&bit2)>>4
			tile[y][3] = (a&bit3)>>4 | (b&bit3)>>3
			tile[y][4] = (a&bit4)>>3 | (b&bit4)>>2
			tile[y][5] = (a&bit5)>>2 | (b&bit5)>>1
			tile[y][6] = (a&bit6)>>1 | (b & bit6)
			tile[y][7] = (a & bit7) | (b&bit7)<<1
		}
	}
	lcd.tileCache[tileNumber] = tile
	return tile, false
}

func (lcd *LCD) updateTiles(lcdY uint8, offsetAddr uint16, layer *[256][256]uint8, previousTiles *[32][32]uint16) {
	lowTileData := lcd.lowTileDataSelect()
	tileY := lcdY / 8
	for tileX := 0; tileX < 32; tileX++ {
		var tileNumber uint16
		tileAddr := 32*uint16(tileY) + uint16(tileX)
		tileByte := lcd.readVideoRAM(offsetAddr + tileAddr)
		if !lowTileData {
			tileNumber = uint16(256 + int(int8(tileByte)))
		} else {
			tileNumber = uint16(tileByte)
		}
		tile, cacheHit := lcd.readTile(tileNumber)
		if !cacheHit || tileNumber != previousTiles[tileY][tileX] {
			lcdX := uint8(tileX * 8)
			for y := uint8(0); y < 8; y++ {
				for x := uint8(0); x < 8; x++ {
					layer[lcdY+y][lcdX+x] = tile[y][x]
				}
			}
		}
		previousTiles[tileY][tileX] = tileNumber
	}
}

func (lcd *LCD) updateBG(lcdY, scy uint8) {
	if !lcd.bgDisplayEnable() {
		return
	}
	var offsetAddr uint16
	if lcd.highBgTileMapDisplaySelect() {
		offsetAddr = 0x9c00
	} else {
		offsetAddr = 0x9800
	}
	lcd.updateTiles(lcdY+scy, offsetAddr, &lcd.bg, &lcd.previousBg)
}

func (lcd *LCD) updateWindow(lcdY uint8) {
	if !lcd.windowDisplayEnable() {
		return
	}
	var offsetAddr uint16
	if lcd.highWindowTileMapDisplaySelect() {
		offsetAddr = 0x9c00
	} else {
		offsetAddr = 0x9800
	}
	lcd.updateTiles(lcdY, offsetAddr, &lcd.window, &lcd.previousWindow)
}

func (lcd *LCD) readSpriteTile(tileNumber uint16, att uint8) *[8][8]uint8 {
	tile := lcd.tileCache[tileNumber]
	tile, _ = lcd.readTile(uint16(tileNumber))
	if spriteXFlip(att) {
		new := [8][8]uint8{}
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				new[y][x] = tile[y][7-x]
			}
		}
		tile = &new
	}
	if spriteYFlip(att) {
		new := [8][8]uint8{}
		for y := 0; y < 8; y++ {
			new[y] = tile[7-y]
		}
		tile = &new
	}
	return tile
}

func (lcd *LCD) updateSprites(lcdY uint8) {
	if lcdY >= 144 {
		return
	}
	if lcd.largeSprites() {
		panic(fmt.Sprintf("Large sprites are not supported"))
	}
	lcd.sprites[lcdY] = [160]uint8{}
	for sprite := 0; sprite < 40; sprite++ {
		spriteAddr := sprite * 4
		startY := lcd.oam[spriteAddr]
		if startY == 0 || lcdY < startY-16 || lcdY >= startY-8 {
			continue
		}
		startX := lcd.oam[spriteAddr+1]
		if startX == 0 || startX >= 168 {
			continue
		}
		tileNumber := lcd.oam[spriteAddr+2]
		attributes := lcd.oam[spriteAddr+3]
		tile := lcd.readSpriteTile(uint16(tileNumber), attributes)
		spriteX := startX - 8
		for tileX := uint8(0); tileX < 8; tileX++ {
			lcdX := spriteX + tileX
			if lcdX >= 0 && lcdX < 160 {
				lcd.sprites[lcdY][lcdX] = tile[lcdY-startY+16][tileX]
			}
		}
	}
}

func (lcd *LCD) renderPixel(x, y, scx, scy, wx, wy uint8, debug bool) color.RGBA {
	if lcd.spriteDisplayEnable() {
		if x < 160 && y < 144 {
			pixel := lcd.sprites[y][x]
			if pixel > 0 {
				if debug {
					return blue[pixel]
				}
				return gray[pixel]
			}
		}
	}
	if lcd.windowDisplayEnable() {
		// Use WX/WY to shift the visible pixels
		if x >= wx && y >= wy {
			pixel := lcd.window[y-wy][x-wx]
			if debug {
				return green[pixel]
			}
			return gray[pixel]
		}
	}
	if lcd.bgDisplayEnable() {
		// Use SCX/SCY to shift the visible pixels
		pixel := lcd.bg[y+scy][x+scx]
		if debug && (x >= 160 || y >= 144) {
			return red[pixel]
		}
		return gray[pixel]
	}
	return gray[0]
}

func (lcd *LCD) renderLine(y, scy uint8) {
	scx := lcd.hwr.SCX
	wx := lcd.hwr.WX
	wy := lcd.hwr.WY
	if lcd.debug {
		for x := 0; x < 256; x++ {
			pixel := lcd.renderPixel(uint8(x)-scx, y-scy, scx, scy, wx, wy, true)
			lcd.Frame.SetRGBA(x, int(y), pixel)
		}
	} else {
		for x := 0; x < 160; x++ {
			pixel := lcd.renderPixel(uint8(x), y, scx, scy, wx, wy, false)
			lcd.Frame.SetRGBA(x, int(y), pixel)
		}
	}
}

func (lcd *LCD) updateLcdLine(y uint8) {
	scy := lcd.hwr.SCY
	lcd.updateBG(y, scy)
	lcd.updateWindow(y)
	lcd.updateSprites(y)
	lcd.renderLine(y, scy)
}

// FrameEnd writes any remaining VRAM lines to the GUI for debugging
func (lcd *LCD) FrameEnd() {
	if lcd.debug {
		for y := lcd.hwr.SCY + 144; y < lcd.hwr.SCY || y >= lcd.hwr.SCY+144; y++ {
			lcd.updateLcdLine(y)
		}
	}
}
