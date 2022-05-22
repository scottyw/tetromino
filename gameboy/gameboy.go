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
	"github.com/scottyw/tetromino/gameboy/interrupts"
	"github.com/scottyw/tetromino/gameboy/memory"
	"github.com/scottyw/tetromino/gameboy/oam"
	"github.com/scottyw/tetromino/gameboy/ppu"
	"github.com/scottyw/tetromino/gameboy/serial"
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
	cpu        *cpu.CPU
	display    *display.Display
	interrupts *interrupts.Interrupts
	ppu        *ppu.PPU
	mapper     *memory.Mapper
	speakers   *speakers.Speakers
	timer      *timer.Timer
}

// NewGameboy returns a new Gameboy
func New(config Config) *Gameboy {

	// Create interrrupts subsystem
	i := interrupts.New()

	// Create OAM memory
	oam := oam.New()

	// Create speakers
	var a *audio.Audio
	var s *speakers.Speakers
	if !config.DisableAudioOutput {
		s = speakers.New()
		a = audio.New(s.Left(), s.Right())
	} else {
		a = audio.New(nil, nil)
	}

	// Create the PPU
	ppu := ppu.New(i, oam, config.DebugLCD)

	// Create the serial bus subsystem
	serial := serial.New(config.SerialWriter)

	// Create the timer subsystem
	timer := timer.New()

	// Load the ROM file
	rom := readRomFile(config.RomFilename)

	// Create controller

	// FIXME
	controller := controller.New(nil)

	mapper := memory.New(rom, i, oam, ppu, controller, serial, timer, a)

	// Create CPU
	c := cpu.New(i, oam, config.DebugCPU, mapper)

	c.Initialize()

	// Create a display
	var d *display.Display
	if !config.DisableVideoOutput {
		d = display.New(controller, config.DebugLCD)
	}

	return &Gameboy{
		audio:      a,
		config:     config,
		controller: controller,
		cpu:        c,
		display:    d,
		interrupts: i,
		ppu:        ppu,
		mapper:     mapper,
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
		panic("No ROM file specified")
	}
	rom, err := ioutil.ReadFile(romFilename)
	if err != nil {
		panic(fmt.Sprintf("Failed to read the ROM file at \"%s\" (%v)", romFilename, err))
	}
	return rom
}

func (gb *Gameboy) runFrame(ctx context.Context) bool {
	// The Game Boy clock runs at 4.194304MHz
	// One machine cycle is 4 clock cycles
	// Each loop iteration below represents one machine cycle
	// Each LCD frame is 17556 machine cycles
	for mtick := 0; mtick < 17556; mtick++ {
		gb.cpu.ExecuteMachineCycle()
		gb.ppu.EndMachineCycle()
		gb.mapper.EndMachineCycle()
		gb.audio.EndMachineCycle()
		timerInterruptRequested := gb.timer.EndMachineCycle()
		if timerInterruptRequested {
			gb.interrupts.RequestTimer()
		}
	}
	frame := gb.ppu.Frame()
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
	var frames uint8
	// t0 := time.Now().UnixMicro()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if gb.runFrame(ctx) {
				return
			}
		}

		// Show FPS
		if frames == 0 {
			// t1 := time.Now().UnixMicro()
			// microseconds := (t1 - t0) / 256
			// fmt.Printf(">>>>  %0d fps ( %0d μs/frame )\n", 1000000/microseconds, microseconds)
			// t0 = time.Now().UnixMicro()
		}
		frames++

	}
}
