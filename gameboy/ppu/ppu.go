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
	bgpColour [4]uint8

	// OBP0
	obp0Colour [4]uint8

	// OBP1
	obp1Colour [4]uint8

	// Single-value registers
	ly  uint8
	lyc uint8
	scx uint8
	scy uint8
	wx  uint8
	wy  uint8

	// Internal state
	interrupts     *interrupts.Interrupts
	oam            *oam.OAM
	videoRAM       [0x2000]byte
	frame          *image.RGBA
	spriteOverlaps [40]bool
	ticks          int
	debug          bool
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

	// Find x and y co-ords
	ppu.ly = uint8(ppu.ticks / 114)
	lx := uint8(ppu.ticks % 114)

	// Should we switch to a different mode?
	switch ppu.mode {
	case 2:
		if lx == 20 {
			ppu.mode = 3
		}
	case 3:
		// In reality, this time is dependent on sprite reading delays
		// but we are hardcoding to the minimum for now i.e. 41 mticks
		if lx == 61 {
			ppu.mode = 0
			// If the h-blank interrupt is enabled in stat
			// then the stat interrupt occurs
			if ppu.hlankInterrupt {
				ppu.interrupts.RequestStat()
			}
		}
	case 0:
		// In reality, this time is dependent on the time spent in mode 3
		// but we are hardcoding to the maximum for now i.e. 52
		if lx == 0 {
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
		if ppu.ticks == 0 {
			ppu.mode = 2
		}
	default:
		panic(fmt.Sprintf("unexpected mode during check: %d", ppu.mode))
	}

	// Execute a single tick
	switch ppu.mode {
	case 2:
		ppu.checkOverlappingSprites(lx)
	case 3:
		for offset := 0; offset < 4; offset++ {
			ppu.drawPixel(((lx-20)*4)+uint8(offset), ppu.ly)
		}
		if lx == 0 {
			// Check coincidence flag
			ppu.coincidence = ppu.ly == ppu.lyc
			// If the coincidence interrupt is enabled in stat
			// then the stat interrupt occurs
			if ppu.coincidence && ppu.coincidenceInterrupt {
				ppu.interrupts.RequestStat()
			}
		}
	case 0:
		// Nothing to do
	case 1:
		// Nothing to do
	default:
		panic(fmt.Sprintf("unexpected mode during tick: %d", ppu.mode))
	}

	// Tick the PPU
	ppu.ticks++
	if ppu.ticks == 17556 {
		ppu.ticks = 0
	}

}

func (ppu *PPU) checkOverlappingSprites(lx uint8) {
	ppu.checkOverlappingSprite(lx * 2)
	ppu.checkOverlappingSprite((lx * 2) + 1)
}

func (ppu *PPU) checkOverlappingSprite(sprite uint8) {
	spriteAddr := 0xfe00 + uint16(sprite*4)
	startY := ppu.oam.ReadOAM(spriteAddr)
	ppu.spriteOverlaps[sprite] = startY != 0 && ppu.ly >= startY-16 && ppu.ly < startY-8
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
	ppu.videoRAM[addr-0x8000] = value
}
