package gb

import (
	"context"
	"fmt"
	"image"
	"io"
	"time"

	"github.com/scottyw/tetromino/pkg/gb/cpu"
	"github.com/scottyw/tetromino/pkg/gb/lcd"
	"github.com/scottyw/tetromino/pkg/gb/mem"
	"github.com/scottyw/tetromino/pkg/ui"
)

const frameDuration = float64(16742706)

type gui interface {
	DrawFrame(image *image.RGBA)
}

// Options control emulator behaviour
type Options struct {
	RomFilename      string
	SBWriter         io.Writer
	Speed            float64
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
	start time.Time
	opts  Options
	frame int
}

// NewGameboy returns a new Gameboy
func NewGameboy(opts Options) Gameboy {
	hwr := mem.NewHardwareRegisters(opts.SBWriter)
	cpu := cpu.NewCPU(hwr, opts.DebugCPU, opts.DebugFlowControl, opts.DebugJumps)
	mem := mem.NewMemory(hwr, opts.RomFilename)
	lcd := lcd.NewLCD(hwr, mem, opts.DebugLCD)
	start := time.Now()
	return Gameboy{cpu: cpu,
		mem:   mem,
		hwr:   hwr,
		lcd:   lcd,
		start: start,
		opts:  opts,
	}
}

func (gb *Gameboy) runFrame(gui gui) {
	// The Game Boy clock runs at 4.194304MHz
	// There are 4 clock cycles to a "machine cycle" giving 1048576 machine cycles per second
	// Each loop iteration below represents one machine cycle (i.e. 4 clock cycles)
	// Each LCD frame is 17556 machine cycles
	for tick := 0; tick < 17556; tick++ {
		gb.lcd.Tick(gb.frame > 500)
		gb.cpu.Tick(gb.mem)
		gb.hwr.Tick()
	}
	if gui != nil {
		gui.DrawFrame(gb.lcd.Frame)
	}
	gb.frame++
	expectedFrameEndTime := gb.start.Add(time.Duration(gb.frame * int(frameDuration*gb.opts.Speed)))
	sleepDuration := time.Until(expectedFrameEndTime)
	time.Sleep(sleepDuration)

	// FIXME

	// if ui.UserInput().InputRecv {
	// 	gb.cpu.Start()
	// }
	// FIXME maybe make this a flaggable feature?
	// if gb.frame > 500 {
	// 	time.Sleep(1000 * time.Millisecond)
	// }
}

// Run the Gameboy
func (gb *Gameboy) Run(ctx context.Context, gui gui) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			gb.runFrame(gui)
		}
	}
}

// Time the Gameboy as it runs
func (gb *Gameboy) Time(ctx context.Context, gui gui) {
	for {
		// There are just under 60 frames per second (59.7275) so let's time in blocks of 60 frames
		// On a real Gameboy this would take 1 second
		t0 := time.Now()
		for i := 0; i < 60; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				gb.runFrame(gui)
			}
		}
		t1 := time.Now()
		fmt.Println("=========>", (t1.Sub(t0)))
	}
}

// ButtonAction turns UI key presses into emulator button presses
func (gb *Gameboy) ButtonAction(b ui.Button, pressed bool) {
	// Start the CPU in case it was stopped waiting for input
	gb.cpu.Start()
	// Bit 3 - P13 Input Down  or Start    (0=Pressed) (Read Only)
	// Bit 2 - P12 Input Up    or Select   (0=Pressed) (Read Only)
	// Bit 1 - P11 Input Left  or Button B (0=Pressed) (Read Only)
	// Bit 0 - P10 Input Right or Button A (0=Pressed) (Read Only)
	if b == ui.Start {
		if pressed {
			gb.hwr.ButtonInput &^= 0x8
		} else {
			gb.hwr.ButtonInput |= 0x8
		}
	} else if b == ui.Select {
		if pressed {
			gb.hwr.ButtonInput &^= 0x4
		} else {
			gb.hwr.ButtonInput |= 0x4
		}
	}
	if b == ui.B {
		if pressed {
			gb.hwr.ButtonInput &^= 0x2
		} else {
			gb.hwr.ButtonInput |= 0x2
		}
	}
	if b == ui.A {
		if pressed {
			gb.hwr.ButtonInput &^= 0x1
		} else {
			gb.hwr.ButtonInput |= 0x1
		}
	}
	if b == ui.Down {
		if pressed {
			gb.hwr.DirectionInput &^= 0x8
		} else {
			gb.hwr.DirectionInput |= 0x8
		}
	} else if b == ui.Up {
		if pressed {
			gb.hwr.DirectionInput &^= 0x4
		} else {
			gb.hwr.DirectionInput |= 0x4
		}
	}
	if b == ui.Left {
		if pressed {
			gb.hwr.DirectionInput &^= 0x2
		} else {
			gb.hwr.DirectionInput |= 0x2
		}
	} else if b == ui.Right {
		if pressed {
			gb.hwr.DirectionInput &^= 0x1
		} else {
			gb.hwr.DirectionInput |= 0x1
		}
	}
}

// Debug enabled for the UI
func (gb *Gameboy) Debug() bool {
	return gb.opts.DebugLCD
}
