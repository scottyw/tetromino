package gb

import (
	"fmt"
	"time"

	"github.com/scottyw/tetromino/cpu"
	"github.com/scottyw/tetromino/lcd"
	"github.com/scottyw/tetromino/mem"
	"github.com/scottyw/tetromino/ui"
)

// Gameboy represents the Gameboy itself
type Gameboy struct {
	cpu *cpu.CPU
	mem *mem.Memory
	hwr *mem.HardwareRegisters
	lcd *lcd.LCD
	ui  ui.UI
}

// NewGameboy returns a new Gameboy
func NewGameboy() Gameboy {
	hwr := mem.NewHardwareRegisters()
	cpu := cpu.NewCPU(hwr)
	mem := mem.NewMemory(hwr)
	lcd := lcd.NewLCD(hwr, mem)
	ui := ui.NewGL(hwr, cpu)
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
	gb.ui.DrawFrame(gb.lcd)
}

// Run the Gameboy
func (gb Gameboy) Run() {
	for gb.ui.ShouldRun() {
		gb.runFrame()
	}
	gb.ui.Shutdown()
}

// Time the Gameboy as it runs
func (gb Gameboy) Time() {
	for gb.ui.ShouldRun() {
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
