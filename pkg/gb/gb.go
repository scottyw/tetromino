package gb

import (
	"fmt"
	"io"
	"time"

	"github.com/scottyw/tetromino/pkg/gb/cpu"
	"github.com/scottyw/tetromino/pkg/gb/lcd"
	"github.com/scottyw/tetromino/pkg/gb/mem"
	"github.com/scottyw/tetromino/pkg/ui"
)

// Options control emulator behaviour
type Options struct {
	RomFilename      string
	SBWriter         io.Writer
	DebugCPU         bool
	DebugFlowControl bool
	DebugJumps       bool
}

// Gameboy represents the Gameboy itself
type Gameboy struct {
	cpu *cpu.CPU
	mem *mem.Memory
	hwr *mem.HardwareRegisters
	lcd *lcd.LCD
	ui  ui.UI
}

// NewGameboy returns a new Gameboy
func NewGameboy(ui ui.UI, opts Options) Gameboy {
	hwr := mem.NewHardwareRegisters(ui.UserInput(), opts.SBWriter)
	cpu := cpu.NewCPU(hwr)
	mem := mem.NewMemory(hwr, opts.RomFilename)
	lcd := lcd.NewLCD(hwr, mem)
	return Gameboy{cpu: cpu, mem: mem, hwr: hwr, lcd: lcd, ui: ui}
}

func (gb Gameboy) runFrame() {
	// The Game Boy clock runs at 4.194304MHz
	// There are 4 clock cycles to a "machine cycle" giving 1048576 machine cycles per second
	// Each loop iteration below represents one machine cycle (i.e. 4 clock cycles)
	// Each LCD frame is 17556 machine cycles
	for cycle := 0; cycle < 17556; cycle++ {
		gb.lcd.Tick(cycle)
		gb.cpu.Tick(gb.mem)
		gb.hwr.Tick()
	}
	gb.ui.HandleFrame(gb.lcd.FrameData())
}

// Run the Gameboy
func (gb Gameboy) Run() {
	for gb.ui.KeepRunning() {
		gb.runFrame()
	}
	gb.ui.Shutdown()
}

// Time the Gameboy as it runs
func (gb Gameboy) Time() {
	for gb.ui.KeepRunning() {
		// There are just under 60 frames per second (59.7275) so let's time in blocks of 60 frames
		// On a real Gameboy this would take 1 second
		t0 := time.Now()
		for i := 0; i < 60; i++ {
			gb.runFrame()
		}
		t1 := time.Now()
		fmt.Println("=========>", (t1.Sub(t0)))
	}
	gb.ui.Shutdown()
}
