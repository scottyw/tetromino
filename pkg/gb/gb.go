package gb

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/scottyw/tetromino/pkg/gb/audio"
	"github.com/scottyw/tetromino/pkg/gb/controller"
	"github.com/scottyw/tetromino/pkg/gb/cpu"
	"github.com/scottyw/tetromino/pkg/gb/lcd"
	"github.com/scottyw/tetromino/pkg/gb/mem"
	"github.com/scottyw/tetromino/pkg/gb/timer"
)

// Options control emulator behaviour
type Options struct {
	RomFilename string
	DebugCPU    bool
	DebugLCD    bool
	SBWriter    io.Writer
}

// Gameboy represents the Gameboy itself
type Gameboy struct {
	dispatch   *cpu.Dispatch
	memory     *mem.Memory
	controller *controller.Controller
	timer      *timer.Timer
	lcd        *lcd.LCD
	audio      *audio.Audio
	opts       Options
	frame      int
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
	controller := controller.NewController()
	timer := timer.NewTimer()
	audio := audio.NewAudio()
	memory := mem.NewMemory(rom, opts.SBWriter, controller, timer, audio)
	dispatch := cpu.NewDispatch(c, memory)
	lcd := lcd.NewLCD(memory, opts.DebugLCD)
	return &Gameboy{
		dispatch:   dispatch,
		memory:     memory,
		controller: controller,
		timer:      timer,
		lcd:        lcd,
		audio:      audio,
		opts:       opts,
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

func (gb *Gameboy) runFrame() {
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
	gb.frame++

	// The emulator can run a frame much faster than a real Gameboy when running on a modern computer.
	// There is no need to sleep now between frames however, because the audio subsystem consumes
	// samples at a fixed rate from a blocking channel. The emulator can only push audio samples into
	// the channel at the same rate that the "speakers" consume them. Since the "speakers" are
	// consuming the data at the rate of a real Gameboy (in order to make sound play correctly), the
	// rest of the emulator is slowed to the same correct rate. In "fast" mode, the emulator disables
	// the "speakers" meaning there is no constraint on how fast samples are consumed or on how fast
	// the emulator runs.

}

// Run the Gameboy
func (gb *Gameboy) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			gb.runFrame()
		}
	}
}

// Time the Gameboy as it runs
func (gb *Gameboy) Time(ctx context.Context) {
	for {
		// There are just under 60 frames per second (59.7275) so let's time in blocks of 60 frames
		// On a real Gameboy this would take 1 second
		t0 := time.Now()
		for i := 0; i < 60; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				gb.runFrame()
			}
		}
		t1 := time.Now()
		fmt.Println("=========> ", t1.Sub(t0))
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

// RegisterSpeakers registers a real-world audio implementation with the audio subsystem
func (gb *Gameboy) Controller() *controller.Controller {
	return gb.controller
}
