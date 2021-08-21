package memory

import (
	"fmt"

	"github.com/scottyw/tetromino/gameboy/audio"
	"github.com/scottyw/tetromino/gameboy/controller"
	"github.com/scottyw/tetromino/gameboy/interrupts"
	"github.com/scottyw/tetromino/gameboy/ppu"
	"github.com/scottyw/tetromino/gameboy/serial"
	"github.com/scottyw/tetromino/gameboy/timer"
)

const (

	//
	// Register constants
	//

	JOYP = 0xFF00
	SB   = 0xFF01
	SC   = 0xFF02
	DIV  = 0xFF04
	TIMA = 0xFF05
	TMA  = 0xFF06
	TAC  = 0xFF07
	IF   = 0xFF0F
	NR10 = 0xFF10
	NR11 = 0xFF11
	NR12 = 0xFF12
	NR13 = 0xFF13
	NR14 = 0xFF14
	NR21 = 0xFF16
	NR22 = 0xFF17
	NR23 = 0xFF18
	NR24 = 0xFF19
	NR30 = 0xFF1A
	NR31 = 0xFF1B
	NR32 = 0xFF1C
	NR33 = 0xFF1D
	NR34 = 0xFF1E
	NR41 = 0xFF20
	NR42 = 0xFF21
	NR43 = 0xFF22
	NR44 = 0xFF23
	NR50 = 0xFF24
	NR51 = 0xFF25
	NR52 = 0xFF26
	LCDC = 0xFF40
	STAT = 0xFF41
	SCY  = 0xFF42
	SCX  = 0xFF43
	LY   = 0xFF44
	LYC  = 0xFF45
	DMA  = 0xFF46
	BGP  = 0xFF47
	OBP0 = 0xFF48
	OBP1 = 0xFF49
	WY   = 0xFF4A
	WX   = 0xFF4B
	IE   = 0xFFFF
)

// Memory allows read and write access to memory
type Mapper struct {
	internalRAM [0x2000]byte
	zeroPage    [0x8f]byte
	oam         *[0xa0]byte
	audio       *audio.Audio
	controller  *controller.Controller
	interrupts  *interrupts.Interrupts
	mbc         *mbc
	ppu         *ppu.PPU
	serial      *serial.Serial
	timer       *timer.Timer
	dmaRunning  bool
	dmaCycle    uint16
	dmaBaseAddr uint16
	dmaRead     uint8
}

// NewMemory creates the memory struct and initializes it with ROM contents and default values
func New(rom []byte, oam *[0xa0]byte, interrupts *interrupts.Interrupts, ppu *ppu.PPU, controller *controller.Controller, serial *serial.Serial, timer *timer.Timer, audio *audio.Audio) *Mapper {
	return &Mapper{
		mbc:        newMBC(rom),
		oam:        oam,
		interrupts: interrupts,
		ppu:        ppu,
		controller: controller,
		serial:     serial,
		timer:      timer,
		audio:      audio,
	}
}

func (m *Mapper) EndMachineCycle() {
	if m.dmaRunning {
		m.runDMA()
	}
}

// Read a byte from the chosen memory location
func (m *Mapper) Read(addr uint16) byte {
	switch {
	case addr < 0x8000:
		return m.mbc.read(addr)
	case addr < 0xa000:
		return m.ppu.ReadVideoRAM(addr)
	case addr < 0xc000:
		return m.mbc.read(addr)
	case addr < 0xe000:
		return m.internalRAM[addr-0xc000]
	case addr < 0xfe00:
		return m.internalRAM[addr-0xe000]
	case addr < 0xfea0:
		return m.readOAM(addr)
	case addr < 0xff00:
		// Unusable region
		return 0
	case addr == JOYP:
		return m.controller.ReadJOYP()
	case addr == SB:
		return m.serial.ReadSB()
	case addr == SC:
		return m.serial.ReadSC()
	case addr == DIV:
		return m.timer.ReadDIV()
	case addr == TIMA:
		return m.timer.ReadTIMA()
	case addr == TMA:
		return m.timer.ReadTMA()
	case addr == TAC:
		return m.timer.ReadTAC()
	case addr == IF:
		return m.interrupts.ReadIF()
	case addr == NR10:
		return m.audio.ReadNR10()
	case addr == NR11:
		return m.audio.ReadNR11()
	case addr == NR12:
		return m.audio.ReadNR12()
	case addr == NR13:
		return m.audio.ReadNR13()
	case addr == NR14:
		return m.audio.ReadNR14()
	case addr == NR21:
		return m.audio.ReadNR21()
	case addr == NR22:
		return m.audio.ReadNR22()
	case addr == NR23:
		return m.audio.ReadNR23()
	case addr == NR24:
		return m.audio.ReadNR24()
	case addr == NR30:
		return m.audio.ReadNR30()
	case addr == NR31:
		return m.audio.ReadNR31()
	case addr == NR32:
		return m.audio.ReadNR32()
	case addr == NR33:
		return m.audio.ReadNR33()
	case addr == NR34:
		return m.audio.ReadNR34()
	case addr == NR41:
		return m.audio.ReadNR41()
	case addr == NR42:
		return m.audio.ReadNR42()
	case addr == NR43:
		return m.audio.ReadNR43()
	case addr == NR44:
		return m.audio.ReadNR44()
	case addr == NR50:
		return m.audio.ReadNR50()
	case addr == NR51:
		return m.audio.ReadNR51()
	case addr == NR52:
		return m.audio.ReadNR52()
	case addr < 0xff30:
		// Unused section between audio registers and wave RAM
		return 0xff
	case addr < 0xff40:
		return m.audio.ReadWaveRAM(addr)
	case addr == LCDC:
		return m.ppu.ReadLCDC()
	case addr == STAT:
		return m.ppu.ReadSTAT()
	case addr == SCY:
		return m.ppu.ReadSCY()
	case addr == SCX:
		return m.ppu.ReadSCX()
	case addr == LY:
		return m.ppu.ReadLY()
	case addr == LYC:
		return m.ppu.ReadLYC()
	case addr == DMA:
		return m.readDMA()
	case addr == BGP:
		return m.ppu.ReadBGP()
	case addr == OBP0:
		return m.ppu.ReadOBP0()
	case addr == OBP1:
		return m.ppu.ReadOBP1()
	case addr == WY:
		return m.ppu.ReadWY()
	case addr == WX:
		return m.ppu.ReadWX()
	case addr < 0xff80:
		// Default if a non-hardware register is read
		return 0xff
	case addr < 0xffff:
		return m.zeroPage[addr-0xff80]
	case addr == IE:
		return m.interrupts.ReadIE()
	default:
		panic(fmt.Sprintf("Read failed: 0x%04x", addr))
	}
}

// Write a byte to the chosen memory location
func (m *Mapper) Write(addr uint16, value byte) {
	switch {
	case addr < 0x8000:
		m.mbc.write(addr, value)
	case addr < 0xa000:
		m.ppu.WriteVideoRAM(addr, value)
	case addr < 0xc000:
		m.mbc.write(addr, value)
	case addr < 0xe000:
		m.internalRAM[addr-0xc000] = value
	case addr < 0xfe00:
		m.internalRAM[addr-0xe000] = value
	case addr < 0xfea0:
		m.writeOAM(addr, value)
	case addr < 0xff00:
		// Unusable region
	case addr == JOYP:
		m.controller.WriteJOYP(value)
	case addr == SB:
		m.serial.WriteSB(value)
	case addr == SC:
		m.serial.WriteSC(value)
	case addr == DIV:
		m.timer.WriteDIV(value)
	case addr == TIMA:
		m.timer.WriteTIMA(value)
	case addr == TMA:
		m.timer.WriteTMA(value)
	case addr == TAC:
		m.timer.WriteTAC(value)
	case addr == IF:
		m.interrupts.WriteIF(value)
	case addr == NR10:
		m.audio.WriteNR10(value)
	case addr == NR11:
		m.audio.WriteNR11(value)
	case addr == NR12:
		m.audio.WriteNR12(value)
	case addr == NR13:
		m.audio.WriteNR13(value)
	case addr == NR14:
		m.audio.WriteNR14(value)
	case addr == NR21:
		m.audio.WriteNR21(value)
	case addr == NR22:
		m.audio.WriteNR22(value)
	case addr == NR23:
		m.audio.WriteNR23(value)
	case addr == NR24:
		m.audio.WriteNR24(value)
	case addr == NR30:
		m.audio.WriteNR30(value)
	case addr == NR31:
		m.audio.WriteNR31(value)
	case addr == NR32:
		m.audio.WriteNR32(value)
	case addr == NR33:
		m.audio.WriteNR33(value)
	case addr == NR34:
		m.audio.WriteNR34(value)
	case addr == NR41:
		m.audio.WriteNR41(value)
	case addr == NR42:
		m.audio.WriteNR42(value)
	case addr == NR43:
		m.audio.WriteNR43(value)
	case addr == NR44:
		m.audio.WriteNR44(value)
	case addr == NR50:
		m.audio.WriteNR50(value)
	case addr == NR51:
		m.audio.WriteNR51(value)
	case addr == NR52:
		m.audio.WriteNR52(value)
	case addr < 0xff30:
		// Unused section between audio registers and wave RAM
	case addr < 0xff40:
		m.audio.WriteWaveRAM(addr, value)
	case addr == LCDC:
		m.ppu.WriteLCDC(value)
	case addr == STAT:
		m.ppu.WriteSTAT(value)
	case addr == SCY:
		m.ppu.WriteSCY(value)
	case addr == SCX:
		m.ppu.WriteSCX(value)
	case addr == LY:
		m.ppu.WriteLY(value)
	case addr == LYC:
		m.ppu.WriteLYC(value)
	case addr == DMA:
		m.writeDMA(value)
	case addr == BGP:
		m.ppu.WriteBGP(value)
	case addr == OBP0:
		m.ppu.WriteOBP0(value)
	case addr == OBP1:
		m.ppu.WriteOBP1(value)
	case addr == WY:
		m.ppu.WriteWY(value)
	case addr == WX:
		m.ppu.WriteWX(value)
	case addr < 0xff80:
		// Do nothing if a non-hardware register is written
	case addr < 0xffff:
		m.zeroPage[addr-0xff80] = value
	case addr == IE:
		m.interrupts.WriteIE(value)
	default:
		panic(fmt.Sprintf("Write failed: 0x%04x", addr))
	}
}

// CartRAM returns the contents of cartridge RAM, which is useful for verifing test results
func (m *Mapper) CartRAM() [][0x2000]byte {
	return m.mbc.ram
}

func (m *Mapper) readOAM(addr uint16) uint8 {
	if m.dmaRunning {
		return 0xff
	}
	return m.oam[addr-0xfe00]
}

func (m *Mapper) writeOAM(addr uint16, value uint8) {
	m.oam[addr-0xfe00] = value
}

// ExecuteMachineCycle updates the OAM after a machine cycle
func (m *Mapper) runDMA() {
	if m.dmaCycle == 0 {
		// Setup
	} else if m.dmaCycle == 1 {
		m.dmaRead = m.Read(m.dmaBaseAddr)
	} else if m.dmaCycle == 161 {
		m.oam[159] = m.dmaRead
		m.dmaRunning = false
	} else {
		m.oam[m.dmaCycle-2] = m.dmaRead
		m.dmaRead = m.Read(m.dmaBaseAddr + m.dmaCycle - 1)
	}
	m.dmaCycle++
}

func (m *Mapper) startDMA(value uint8) {
	m.dmaRunning = true
	m.dmaCycle = 0
	m.dmaBaseAddr = uint16(value) << 8
	if m.dmaBaseAddr >= 0xe000 {
		m.dmaBaseAddr -= 0x2000
	}
}

// FF46 - DMA - DMA Transfer and Start Address (W)
// Writing to this register launches a DMA transfer from ROM or RAM to OAM memory
// (sprite attribute table). The written value specifies the transfer source
// address divided by 100h, ie. source & destination are:
// Source:      XX00-XX9F   ;XX in range from 00-F1h
// Destination: FE00-FE9F
// It takes 160 microseconds until the transfer has completed (80 microseconds in
// CGB Double Speed Mode), during this time the CPU can access only HRAM (memory
// at FF80-FFFE). For this reason, the programmer must copy a short procedure
// into HRAM, and use this procedure to start the transfer from inside HRAM, and
// wait until the transfer has finished:
//  ld  (0FF46h),a ;start DMA transfer, a=start address/100h
//  ld  a,28h      ;delay...
// wait:           ;total 5x40 cycles, approx 200ms
//  dec a          ;1 cycle
//  jr  nz,wait    ;4 cycles
// Most programs are executing this procedure from inside of their VBlank
// procedure, but it is possible to execute it during display redraw also,
// allowing to display more than 40 sprites on the screen (ie. for example 40
// sprites in upper half, and other 40 sprites in lower half of the screen).

func (m *Mapper) writeDMA(value uint8) {
	// fmt.Printf("> DMA - 0x%02x\n", value)
	m.startDMA(value)
}

func (m *Mapper) readDMA() uint8 {
	dma := uint8(m.dmaBaseAddr >> 8)
	// fmt.Printf("< DMA - 0x%02x\n", dma)
	return dma
}
