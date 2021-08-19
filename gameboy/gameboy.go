package gameboy

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/scottyw/tetromino/gameboy/audio"
	"github.com/scottyw/tetromino/gameboy/controller"
	"github.com/scottyw/tetromino/gameboy/cpu"
	"github.com/scottyw/tetromino/gameboy/display"
	"github.com/scottyw/tetromino/gameboy/lcd"
	"github.com/scottyw/tetromino/gameboy/mem"
	"github.com/scottyw/tetromino/gameboy/speakers"
	"github.com/scottyw/tetromino/gameboy/timer"
)

// Config control emulator behaviour
type Config struct {
	RomFilename        string
	DisableVideoOutput bool
	DisableAudioOutput bool
	DebugCPU           bool
	DebugLCD           bool
	SerialWriter       io.Writer
}

// Gameboy represents the Gameboy itself
type Gameboy struct {
	audio      *audio.Audio
	config     Config
	controller *controller.Controller
	dispatch   *cpu.Dispatch
	display    *display.Display
	lcd        *lcd.LCD
	memory     *mem.Memory
	speakers   *speakers.Speakers
	timer      *timer.Timer
}

// NewGameboy returns a new Gameboy
func New(config Config) *Gameboy {
	var rom []byte
	if config.RomFilename == "" {
		rom = make([]byte, 0x8000)
	} else {
		rom = readRomFile(config.RomFilename)
	}
	c := cpu.New(config.DebugCPU)
	controller := controller.New(c.Restart)

	// Create a display
	var d *display.Display
	if !config.DisableVideoOutput {
		d = display.New(controller, config.DebugLCD)
	}

	// Create speakers
	var a *audio.Audio
	var s *speakers.Speakers
	if !config.DisableAudioOutput {
		s = speakers.New()
		a = audio.New(s.Left(), s.Right())
	} else {
		a = audio.New(nil, nil)
	}

	timer := timer.New()
	memory := mem.New(rom, config.SerialWriter, controller, timer, a)
	dispatch := cpu.NewDispatch(c, memory)
	lcd := lcd.New(memory, config.DebugLCD)
	return &Gameboy{
		audio:      a,
		config:     config,
		controller: controller,
		dispatch:   dispatch,
		display:    d,
		lcd:        lcd,
		memory:     memory,
		speakers:   s,
		timer:      timer,
	}
}

func (gb *Gameboy) Cleanup() {
	if gb.speakers != nil {
		gb.speakers.Cleanup()
	}
	if gb.display != nil {
		gb.display.Cleanup()
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

func (gb *Gameboy) runFrame(ctx context.Context) bool {
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
	frame := gb.lcd.FrameEnd()
	if gb.display != nil {
		return gb.display.RenderFrame(frame)
	}
	return false

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
	defer gb.Cleanup()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if gb.runFrame(ctx) {
				return
			}
		}
	}
}
