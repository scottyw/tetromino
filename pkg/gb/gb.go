package gb

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/scottyw/tetromino/pkg/gb/cpu"
	"github.com/scottyw/tetromino/pkg/gb/lcd"
	"github.com/scottyw/tetromino/pkg/gb/mem"
	"github.com/scottyw/tetromino/pkg/ui"
)

const frameDuration = float64(16742706)

// Options control emulator behaviour
type Options struct {
	RomFilename      string
	SBWriter         io.Writer
	DebugCPU         bool
	DebugLCD         bool
	DebugFlowControl bool
	DebugJumps       bool
}

// Gameboy represents the Gameboy itself
type Gameboy struct {
	cpu   *cpu.CPU
	mem   *mem.Memory
	hwr   *mem.HardwareRegisters
	lcd   *lcd.LCD
	ui    ui.UI
	start time.Time
	frame int
}

// NewGameboy returns a new Gameboy
func NewGameboy(ui ui.UI, opts Options) Gameboy {
	hwr := mem.NewHardwareRegisters(ui.UserInput(), opts.SBWriter)
	cpu := cpu.NewCPU(hwr, opts.DebugCPU, opts.DebugFlowControl, opts.DebugJumps)
	mem := mem.NewMemory(hwr, opts.RomFilename)
	lcd := lcd.NewLCD(hwr, mem, opts.DebugLCD)
	start := time.Now()
	return Gameboy{cpu: cpu, mem: mem, hwr: hwr, lcd: lcd, ui: ui, start: start}
}

func (gb *Gameboy) runFrame(speed float64) {
	// The Game Boy clock runs at 4.194304MHz
	// There are 4 clock cycles to a "machine cycle" giving 1048576 machine cycles per second
	// Each loop iteration below represents one machine cycle (i.e. 4 clock cycles)
	// Each LCD frame is 17556 machine cycles
	for cycle := 0; cycle < 17556; cycle++ {
		gb.lcd.Tick(cycle, gb.frame > 500)
		gb.cpu.Tick(gb.mem)
		gb.hwr.Tick()
	}
	gb.ui.HandleFrame(gb.lcd.FrameData())
	if gb.ui.UserInput().InputRecv {
		gb.cpu.Start()
	}
	gb.frame++
	expectedFrameEndTime := gb.start.Add(time.Duration(gb.frame * int(frameDuration*speed)))
	sleepDuration := time.Until(expectedFrameEndTime)
	time.Sleep(sleepDuration)
	// FIXME maybe make this a flaggable feature?
	// if gb.frame > 500 {
	// 	time.Sleep(1000 * time.Millisecond)
	// }
}

// Run the Gameboy
func (gb *Gameboy) Run(ctx context.Context, speed float64) {
	for {
		select {
		case <-ctx.Done():
			gb.ui.Shutdown()
			return
		default:
			gb.runFrame(speed)
		}
	}
}

// Time the Gameboy as it runs
func (gb *Gameboy) Time(ctx context.Context, speed float64) {
	for {
		// There are just under 60 frames per second (59.7275) so let's time in blocks of 60 frames
		// On a real Gameboy this would take 1 second
		t0 := time.Now()
		for i := 0; i < 60; i++ {
			select {
			case <-ctx.Done():
				gb.ui.Shutdown()
				return
			default:
				gb.runFrame(speed)
			}
		}
		t1 := time.Now()
		fmt.Println("=========>", (t1.Sub(t0)))
	}
}
