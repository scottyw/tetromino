package ppu

import (
	"fmt"
	"image"

	"github.com/scottyw/tetromino/gameboy/interrupts"
	"github.com/scottyw/tetromino/gameboy/oam"
)

type PPU struct {

	//LCDC
	enabled           bool
	highWindowTileMap bool
	windowEnabled     bool
	lowTileData       bool
	highBgTileMap     bool
	spritesLarge      bool
	spritesEnabled    bool
	bgEnabled         bool

	// STAT
	coincidenceInterrupt bool
	oamInterrupt         bool
	vblankInterrupt      bool
	hlankInterrupt       bool
	coincidence          bool
	mode                 uint8

	// BGP
	bgpColour0 uint8
	bgpColour1 uint8
	bgpColour2 uint8
	bgpColour3 uint8

	// OBP0
	obp0Colour1 uint8
	obp0Colour2 uint8
	obp0Colour3 uint8

	// OBP1
	obp1Colour1 uint8
	obp1Colour2 uint8
	obp1Colour3 uint8

	// Single-value registers
	ly  uint8
	lyc uint8
	scx uint8
	scy uint8
	wx  uint8
	wy  uint8

	// Mode 2
	spriteOverlaps [40]bool
	nextSprite     uint16

	// Mode 3
	lx uint8

	// Mode 0
	ticksMode0 int

	// Mode 1
	ticksMode1 int

	// Internal state
	interrupts *interrupts.Interrupts
	oam        *oam.OAM
	videoRAM   [0x2000]byte
	debug      bool

	// Internal LCD state
	tileCache      [384]*[8][8]uint8
	previousBg     [32][32]uint16
	previousWindow [32][32]uint16
	bg             [256][256]uint8
	window         [256][256]uint8
	sprites        [144][160]uint8
	frame          *image.RGBA
}

func New(interrupts *interrupts.Interrupts, oam *oam.OAM, debug bool) *PPU {
	ppu := &PPU{
		interrupts: interrupts,
		oam:        oam,
		frame:      image.NewRGBA(image.Rect(0, 0, 256, 256)),
		debug:      debug,
	}
	ppu.WriteLCDC(0x91)
	ppu.WriteLY(0x00)
	ppu.WriteLYC(0x00)
	ppu.WriteSCX(0x00)
	ppu.WriteSCY(0x00)
	ppu.WriteSTAT(0x00)
	ppu.WriteWX(0x00)
	ppu.WriteWY(0x00)
	ppu.WriteBGP(0xFC)
	ppu.WriteOBP0(0xFF)
	ppu.WriteOBP1(0xFF)
	return ppu
}

// EndMachineCycle updates the LCD driver after each machine cycle i.e. 4 clock cycles
func (ppu *PPU) EndMachineCycle() {

	// Is the lcd enabled?
	if !ppu.enabled {
		return
	}

	// Should we switch to a different mode?
	switch ppu.mode {
	case 2:
		if ppu.nextSprite >= 40 {
			ppu.nextSprite = 0
			ppu.mode = 3
		}
	case 3:
		// In reality, this time is dependent on sprite reading delays
		// but we are hardcoding to the minimum for now
		if ppu.lx >= 168 {
			ppu.lx = 0
			ppu.mode = 0
			// If the h-blank interrupt is enabled in stat
			// then the stat interrupt occurs
			if ppu.hlankInterrupt {
				ppu.interrupts.RequestStat()
			}
		}
	case 0:
		// In reality, this time is dependent on the time spent in mode 3
		// but we are hardcoding to themaximum for now
		if ppu.ticksMode0 >= 208 {
			ppu.ticksMode0 = 0
			ppu.ly++
			if ppu.ly == 144 {
				ppu.mode = 1
				// V-blank interrupt always occurs
				ppu.interrupts.RequestVblank()
				// If the v-blank interrupt is also enabled in stat
				// then the stat interrupt occurs too
				if ppu.vblankInterrupt {
					ppu.interrupts.RequestStat()
				}
			} else {
				ppu.mode = 2
				// If the oam interrupt is enabled in stat
				// then the stat interrupt occurs
				if ppu.oamInterrupt {
					ppu.interrupts.RequestStat()
				}
			}
		}
	case 1:
		if ppu.ticksMode1 >= 456 {
			ppu.ticksMode1 = 0
			ppu.ly++
		}
		if ppu.ly == 154 {
			ppu.ly = 0
			ppu.mode = 2
		}
	default:
		panic(fmt.Sprintf("unexpected mode during check: %d", ppu.mode))
	}

	// Tick the PPU
	switch ppu.mode {
	case 2:
		ppu.checkSpriteOverlap()
	case 3:
		ppu.drawPixel(ppu.lx, ppu.ly)
		if ppu.lx == 0 {
			// Check coincidence flag
			ppu.coincidence = ppu.ly == ppu.lyc
			// If the coincidence interrupt is enabled in stat
			// then the stat interrupt occurs
			if ppu.coincidence && ppu.coincidenceInterrupt {
				ppu.interrupts.RequestStat()
			}
		}
		ppu.lx++
	case 0:
		ppu.ticksMode0++
	case 1:
		ppu.ticksMode1++
	default:
		panic(fmt.Sprintf("unexpected mode during tick: %d", ppu.mode))
	}

}

func (ppu *PPU) checkSpriteOverlap() {
	spriteAddr := 0xfe00 + uint16(ppu.nextSprite*4)
	startY := ppu.oam.ReadOAM(spriteAddr)
	ppu.spriteOverlaps[ppu.nextSprite] = startY != 0 && ppu.ly >= startY-16 && ppu.ly < startY-8
	ppu.nextSprite++
	spriteAddr = 0xfe00 + uint16(ppu.nextSprite*4)
	startY = ppu.oam.ReadOAM(spriteAddr)
	ppu.spriteOverlaps[ppu.nextSprite] = startY != 0 && ppu.ly >= startY-16 && ppu.ly < startY-8
	ppu.nextSprite++
}

func (ppu *PPU) enable() {
	ppu.enabled = true
}

func (ppu *PPU) disable() {
	ppu.enabled = false
	ppu.ly = 0
}

func (ppu *PPU) ReadVideoRAM(addr uint16) uint8 {
	return ppu.videoRAM[addr-0x8000]
}

func (ppu *PPU) WriteVideoRAM(addr uint16, value uint8) {
	if addr < 0x9800 {
		tileNumber := (addr - 0x8000) / 16
		ppu.tileCache[tileNumber] = nil
	}
	ppu.videoRAM[addr-0x8000] = value
}
