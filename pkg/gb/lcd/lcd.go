package lcd

import (
	"encoding/gob"
	"fmt"
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
	fgSprites      [144][160]uint8
	bgSprites      [144][160]uint8
	data           [144][256]uint8
	tick           int
	debug          bool
}

// NewLCD returns the configured LCD
func NewLCD(hwr *mem.HardwareRegisters, memory *mem.Memory, debug bool) *LCD {
	lcd := LCD{
		hwr:      hwr,
		videoRAM: &memory.VideoRAM,
		oam:      &memory.OAM,
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
func (lcd *LCD) Tick(foo bool) {

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
		lcd.updateLcdLine()
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

// FrameData returns the frame data as a 256x144 array of bytes where each element is a colour value between 0 and 3
func (lcd *LCD) FrameData() [144][256]uint8 {
	return lcd.data
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

func (lcd *LCD) updateTiles(offsetAddr uint16, layer *[256][256]uint8, previousTiles *[32][32]uint16) {
	lowTileData := lcd.lowTileDataSelect()
	for j := uint16(0); j < 32; j++ {
		for i := uint16(0); i < 32; i++ {
			var tileNumber uint16
			tileAddr := 32*j + i
			tileByte := lcd.readVideoRAM(offsetAddr + tileAddr)
			if !lowTileData {
				tileNumber = uint16(256 + int(int8(tileByte)))
			} else {
				tileNumber = uint16(tileByte)
			}
			tile, cacheHit := lcd.readTile(tileNumber)
			if !cacheHit || tileNumber != previousTiles[j][i] {
				bgX := i * 8
				bgY := j * 8
				for y := uint16(0); y < 8; y++ {
					for x := uint16(0); x < 8; x++ {
						layer[bgY+y][bgX+x] = tile[y][x]
					}
				}
			}
			previousTiles[j][i] = tileNumber
		}
	}
}

func (lcd *LCD) updateBG() {
	if !lcd.bgDisplayEnable() {
		return
	}
	var offsetAddr uint16
	if lcd.highBgTileMapDisplaySelect() {
		offsetAddr = 0x9c00
	} else {
		offsetAddr = 0x9800
	}
	lcd.updateTiles(offsetAddr, &lcd.bg, &lcd.previousBg)
}

func (lcd *LCD) updateWindow() {
	if !lcd.windowDisplayEnable() {
		return
	}
	var offsetAddr uint16
	if lcd.highWindowTileMapDisplaySelect() {
		offsetAddr = 0x9c00
	} else {
		offsetAddr = 0x9800
	}
	lcd.updateTiles(offsetAddr, &lcd.window, &lcd.previousWindow)
}

func (lcd *LCD) updateSprites() {
	if lcd.largeSprites() {
		panic(fmt.Sprintf("Large sprites are not supported"))
	}
	lcd.fgSprites = [144][160]uint8{}
	lcd.bgSprites = [144][160]uint8{}
	for sprite := 0; sprite < 40; sprite++ {
		spriteAddr := sprite * 4
		startY := uint16(lcd.oam[spriteAddr])
		if startY == 0 || startY > 160 {
			continue
		}
		startX := uint16(lcd.oam[spriteAddr+1])
		if startX == 0 || startX > 168 {
			continue
		}
		tileNumber := lcd.oam[spriteAddr+2]
		tile, _ := lcd.readTile(uint16(tileNumber))
		spriteX := startX - 8
		spriteY := startY - 16
		for y := uint16(0); y < 8; y++ {
			for x := uint16(0); x < 8; x++ {
				if spriteX+x >= 0 &&
					spriteX+x < 160 &&
					spriteY+y >= 0 &&
					spriteY+y < 144 {
					// FIXME background sprites added to bgSprites based on attributes
					lcd.fgSprites[spriteY+y][spriteX+x] = tile[y][x]
				}
			}
		}
	}
}

func (lcd *LCD) renderPixel(x, y, scx, scy, wx, wy uint8) uint8 {
	if x < 160 && y < 144 {
		pixel := lcd.fgSprites[y][x]
		if pixel > 0 {
			return pixel
		}
	}
	if lcd.windowDisplayEnable() {
		// Use WX/WY to shift the visible pixels
		if x >= wx && y >= wy {
			return lcd.window[y-wy][x-wx]
		}
	}
	if lcd.bgDisplayEnable() {
		// Use SCX/SCY to shift the visible pixels
		return lcd.bg[y+scy][x+scx]
	}
	return 0
}

func (lcd *LCD) renderPixelForDebug(x, y, scx, scy, wx, wy uint8) uint8 {
	var pixel uint8
	// Check for a foreground sprite ...
	if x-scx < 160 && y-scy < 144 {
		pixel = lcd.fgSprites[y-scy][x-scx]
		if pixel > 0 {
			return pixel + 0x30
		}
	}
	if lcd.windowDisplayEnable() {
		// Use WX/WY to shift the visible pixels
		if x-scx < 160 && y-scy < 144 {
			if x >= wx && y >= wy {
				return lcd.window[y-wy-scy][x-wx-scx] + 0x20
			}
		}
	}
	if lcd.bgDisplayEnable() {
		// Check for backround using SCX/SCY to colour pixels to illusrate the offset
		if x-scx >= 160 || y-scy >= 144 {
			pixel += 0x10
		}
		pixel += lcd.bg[y][x]
	}
	return pixel
}

func (lcd *LCD) renderLine() {
	y := lcd.hwr.LY
	scx := lcd.hwr.SCX
	scy := lcd.hwr.SCY
	wx := lcd.hwr.WX
	wy := lcd.hwr.WY
	if lcd.debug {
		for x := 0; x < 256; x++ {
			lcd.data[y][x] = lcd.renderPixelForDebug(uint8(x), y, scx, scy, wx, wy)
		}
	} else {
		for x := 0; x < 160; x++ {
			lcd.data[y][x] = lcd.renderPixel(uint8(x), y, scx, scy, wx, wy)
		}
	}
}

func (lcd *LCD) updateLcdLine() {
	lcd.updateBG()
	lcd.updateWindow()
	lcd.updateSprites()
	lcd.renderLine()
}
