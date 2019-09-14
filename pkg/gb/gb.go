package gb

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/scottyw/tetromino/pkg/gb/audio"
	"github.com/scottyw/tetromino/pkg/gb/cpu"
	"github.com/scottyw/tetromino/pkg/gb/lcd"
	"github.com/scottyw/tetromino/pkg/gb/mem"
	"github.com/scottyw/tetromino/pkg/gb/timer"
)

const frameDuration = float64(16742706)

// Button represents a direction pad or button control
type Button int

const (
	// Up on the control pad
	Up = iota
	//Down on the control pad
	Down = iota
	// Left on the control pad
	Left = iota
	// Right on the control pad
	Right = iota
	// A button
	A = iota
	// B button
	B = iota
	// Start button
	Start = iota
	// Select button
	Select = iota
)

// Action  represents emulator controls
type Action int

const (
	// TakeScreenshot of the current LCD
	TakeScreenshot = iota
	// RunFaster speeds the emulator up
	RunFaster = iota
	// RunSlower slows the emulator down
	RunSlower = iota
)

// Options control emulator behaviour
type Options struct {
	RomFilename string
	Fast        bool
	DebugCPU    bool
	DebugTimer  bool
	DebugLCD    bool
	SBWriter    io.Writer
}

// Gameboy represents the Gameboy itself
type Gameboy struct {
	dispatch *cpu.Dispatch
	memory   *mem.Memory
	timer    *timer.Timer
	lcd      *lcd.LCD
	audio    *audio.Audio
	start    time.Time
	opts     Options
	dur      time.Duration
	frame    int
}

// NewGameboy returns a new Gameboy
func NewGameboy(opts Options) *Gameboy {
	var rom []byte
	if opts.RomFilename == "" {
		rom = make([]byte, 0x8000)
	} else {
		rom = readRomFile(opts.RomFilename)
	}
	c := cpu.NewCPU(opts.DebugCPU)
	timer := timer.NewTimer(opts.DebugTimer)
	audio := audio.NewAudio()
	memory := mem.NewMemory(rom, opts.SBWriter, timer, audio)
	dispatch := cpu.NewDispatch(c, memory)
	lcd := lcd.NewLCD(memory, opts.DebugLCD)
	start := time.Now()
	duration := frameDuration
	if opts.Fast {
		duration = 0
	}
	return &Gameboy{
		dispatch: dispatch,
		memory:   memory,
		timer:    timer,
		lcd:      lcd,
		audio:    audio,
		start:    start,
		opts:     opts,
		dur:      time.Duration(duration),
	}
}

func readRomFile(romFilename string) []byte {
	var rom []byte
	if romFilename == "" {
		panic(fmt.Sprintf("No ROM file specified"))
	}
	rom, err := ioutil.ReadFile(romFilename)
	if err != nil {
		panic(fmt.Sprintf("Failed to read the ROM file at \"%s\" (%v)", romFilename, err))
	}
	return rom
}

func (gb *Gameboy) runFrame(end time.Time) {
	// The Game Boy clock runs at 4.194304MHz
	// Each loop iteration below represents one machine cycle
	// One machine cycle is 4 clock cycles
	// Each LCD frame is 17556 machine cycles
	for mtick := 0; mtick < 17556; mtick++ {
		gb.dispatch.ExecuteMachineCycle()
		gb.memory.ExecuteMachineCycle()
		gb.lcd.EndMachineCycle()
		gb.audio.EndMachineCycle()
		timerInterruptRequested := gb.timer.EndMachineCycle()
		if timerInterruptRequested {
			gb.memory.IF |= 0x04
		}
	}
	gb.lcd.FrameEnd()
	time.Sleep(time.Until(end))
	gb.frame++
}

// Run the Gameboy
func (gb *Gameboy) Run(ctx context.Context) {
	end := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			end = end.Add(gb.dur)
			gb.runFrame(end)
		}
	}
}

// Time the Gameboy as it runs
func (gb *Gameboy) Time(ctx context.Context) {
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
				gb.runFrame(end)
			}
		}
		t1 := time.Now()
		fmt.Println("=========> ", t1.Sub(t0))
	}
}

// ButtonAction turns UI key presses into emulator button presses corresponding to the Gameboy controls
func (gb *Gameboy) ButtonAction(button Button, pressed bool) {

	// Start the CPU in case it was stopped waiting for input
	gb.dispatch.Start()

	// Bit 3 - P13 Input Down  or Start    (0=Pressed) (Read Only)
	// Bit 2 - P12 Input Up    or Select   (0=Pressed) (Read Only)
	// Bit 1 - P11 Input Left  or Button B (0=Pressed) (Read Only)
	// Bit 0 - P10 Input Right or Button A (0=Pressed) (Read Only)

	switch button {

	// FIXME it shouldn't be possible to press left and right at once or up and down at once

	case Start:
		if pressed {
			gb.memory.ButtonInput &^= 0x8
		} else {
			gb.memory.ButtonInput |= 0x8
		}

	case Select:
		if pressed {
			gb.memory.ButtonInput &^= 0x4
		} else {
			gb.memory.ButtonInput |= 0x4
		}

	case B:
		if pressed {
			gb.memory.ButtonInput &^= 0x2
		} else {
			gb.memory.ButtonInput |= 0x2
		}

	case A:
		if pressed {
			gb.memory.ButtonInput &^= 0x1
		} else {
			gb.memory.ButtonInput |= 0x1
		}

	case Down:
		if pressed {
			gb.memory.DirectionInput &^= 0x8
		} else {
			gb.memory.DirectionInput |= 0x8
		}

	case Up:
		if pressed {
			gb.memory.DirectionInput &^= 0x4
		} else {
			gb.memory.DirectionInput |= 0x4
		}

	case Left:
		if pressed {
			gb.memory.DirectionInput &^= 0x2
		} else {
			gb.memory.DirectionInput |= 0x2
		}

	case Right:
		if pressed {
			gb.memory.DirectionInput &^= 0x1
		} else {
			gb.memory.DirectionInput |= 0x1
		}
	}
}

// EmulatorAction turns UI key presses into actions controlling the emulator itself
func (gb *Gameboy) EmulatorAction(action Action) {
	switch action {
	case TakeScreenshot:
		t := time.Now()
		filename := fmt.Sprintf("tetromino-%d%02d%02d-%02d%02d%02d.png",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
		fmt.Println("Writing screenshot to", filename)
		gb.lcd.Screenshot(filename)
	case RunFaster:
		gb.dur /= 2
	case RunSlower:
		gb.dur *= 2
	}
}

// Debug enabled for the UI
func (gb *Gameboy) Debug() bool {
	return gb.opts.DebugLCD
}

// RegisterDisplay registers a real-world display implementation with the LCD subsystem
func (gb *Gameboy) RegisterDisplay(display lcd.Display) {
	gb.lcd.RegisterDisplay(display)
}

// RegisterSpeakers registers a real-world audio implementation with the audio subsystem
func (gb *Gameboy) RegisterSpeakers(speakers audio.Speakers) {
	gb.audio.RegisterSpeakers(speakers)
}
