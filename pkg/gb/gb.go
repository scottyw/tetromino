package gb

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"time"

	"github.com/scottyw/tetromino/pkg/gb/cpu"
	"github.com/scottyw/tetromino/pkg/gb/lcd"
	"github.com/scottyw/tetromino/pkg/gb/mem"
	"github.com/scottyw/tetromino/pkg/gb/timer"
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
	Speedup          float64
	DebugCPU         bool
	DebugLCD         bool
	DebugFlowControl bool
	DebugJumps       bool
}

// Gameboy represents the Gameboy itself
type Gameboy struct {
	cpu    *cpu.CPU
	mem    *mem.Memory
	hwr    *mem.HardwareRegisters
	timer  *timer.Timer
	lcd    *lcd.LCD
	start  time.Time
	opts   Options
	cancel func()
	dur    time.Duration
	frame  int
}

// NewGameboy returns a new Gameboy
func NewGameboy(opts Options, cancel func()) Gameboy {
	timer := timer.NewTimer()
	hwr := mem.NewHardwareRegisters(timer, opts.SBWriter)
	cpu := cpu.NewCPU(hwr, opts.DebugCPU, opts.DebugFlowControl, opts.DebugJumps)
	var memory *mem.Memory
	if opts.RomFilename == "" {
		memory = mem.NewMemory(hwr, make([]byte, 0x8000))
	} else {
		memory = mem.NewMemoryFromFile(hwr, opts.RomFilename)
	}
	lcd := lcd.NewLCD(hwr, memory, opts.DebugLCD)
	start := time.Now()
	return Gameboy{
		cpu:    cpu,
		mem:    memory,
		hwr:    hwr,
		timer:  timer,
		lcd:    lcd,
		start:  start,
		opts:   opts,
		cancel: cancel,
		dur:    time.Duration(int(frameDuration / opts.Speedup)),
	}
}

func (gb *Gameboy) runFrame(gui gui, end time.Time) {
	// The Game Boy clock runs at 4.194304MHz
	// Each loop iteration below represents one machine cycle
	// One machine cycle is 4 clock cycles
	// Each LCD frame is 17556 machine cycles
	for mtick := 0; mtick < 17556; mtick++ {
		gb.cpu.ExecuteMachineCycle(gb.mem)
		gb.lcd.EndMachineCycle()
		timerInterruptRequested := gb.timer.EndMachineCycle()
		if timerInterruptRequested {
			gb.hwr.IF |= 0x04
		}
	}
	gb.lcd.FrameEnd()
	if gui != nil {
		gui.DrawFrame(gb.lcd.DebugOffsetFrame)
	}
	time.Sleep(time.Until(end))
	gb.frame++
}

// Run the Gameboy
func (gb *Gameboy) Run(ctx context.Context, gui gui) {
	end := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			end = end.Add(gb.dur)
			gb.runFrame(gui, end)
		}
	}
}

// Time the Gameboy as it runs
func (gb *Gameboy) Time(ctx context.Context, gui gui) {
	end := time.Now()
	for {
		// There are just under 60 frames per second (59.7275) so let's time in blocks of 60 frames
		// On a real Gameboy this would take 1 second
		t0 := time.Now()
		for i := 0; i < 60; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				end = end.Add(gb.dur)
				gb.runFrame(gui, end)
			}
		}
		t1 := time.Now()
		fmt.Println("=========> ", t1.Sub(t0))
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

// Screenshot writes a screenshot to file
func (gb *Gameboy) Screenshot() {
	t := time.Now()
	realfilename := fmt.Sprintf("tetromino-%d%02d%02d-%02d%02d%02d.%9d.png",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
	debugfilename := fmt.Sprintf("tetromino-debug-%d%02d%02d-%02d%02d%02d.%9d.png",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
	debugoffsetfilename := fmt.Sprintf("tetromino-debugoffset-%d%02d%02d-%02d%02d%02d.%9d.png",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
	fmt.Println("Writing screenshot to", realfilename)

	f3, err := os.Create(realfilename)
	if err != nil {
		fmt.Println(err)
	}
	defer f3.Close()
	err = png.Encode(f3, gb.lcd.RealFrame)
	if err != nil {
		fmt.Println(err)
	}

	f2, err := os.Create(debugoffsetfilename)
	if err != nil {
		fmt.Println(err)
	}
	defer f2.Close()
	err = png.Encode(f2, gb.lcd.DebugOffsetFrame)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(debugfilename)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	err = png.Encode(f, gb.lcd.DebugFrame)
	if err != nil {
		fmt.Println(err)
	}

}

// Faster makes the emulator run faster
func (gb *Gameboy) Faster() {
	gb.dur /= 2
}

// Slower makes the emulator run slower
func (gb *Gameboy) Slower() {
	gb.dur *= 2
}

// Debug enabled for the UI
func (gb *Gameboy) Debug() bool {
	return gb.opts.DebugLCD
}

// Shutdown the emulator
func (gb *Gameboy) Shutdown() {
	gb.cancel()
}
