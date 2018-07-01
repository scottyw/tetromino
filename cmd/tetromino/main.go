package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"

	"github.com/scottyw/tetromino/pkg/gb"
	"github.com/scottyw/tetromino/pkg/ui"
)

func main() {

	fmt.Println()
	fmt.Println("Welcome to Tetromino!")
	fmt.Println()
	fmt.Println("Arrows keys : Up/Down/Left/Right")
	fmt.Println("A : Start")
	fmt.Println("S : Select")
	fmt.Println("Z : B button")
	fmt.Println("X : A button")
	fmt.Println()
	fmt.Println("T : Take screenshot")
	fmt.Println()

	// Command line flags
	romFilename := flag.String("f", "", "ROM filename")
	outputSerial := flag.Bool("output-serial-data", false, "When true, data sent to the serial port will be written to console")
	speed := flag.Float64("speed", 1000, "The speed at which to run as a percentage e.g. 100 for normal speed, 200 for double speed")
	debugCPU := flag.Bool("debug-cpu", false, "When true, CPU debugging is enabled")
	debugFlowControl := flag.Bool("debug-flow", false, "When true, flow control debugging is enabled")
	debugJumps := flag.Bool("debug-jumps", false, "When true, jump debugging is enabled")
	debugLCD := flag.Bool("debug-lcd", false, "When true, LCD colour-based debugging is enabled")
	enableTiming := flag.Bool("enable-timing", false, "When true, timing is output every 60 frames")
	enableProfiling := flag.Bool("enable-profiling", false, "When true, CPU profiling data is written to 'cpuprofile.pprof'")
	flag.Parse()

	// CPU profiling
	if *enableProfiling {
		f, err := os.Create("cpuprofile.pprof")
		if err != nil {
			log.Fatalf("Failed to write cpuprofile.pprof: %v", err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Set up options
	var sbWriter io.Writer
	if *outputSerial {
		sbWriter = os.Stdout
	}
	opts := gb.Options{
		RomFilename:      *romFilename,
		Speed:            100 / *speed,
		SBWriter:         sbWriter,
		DebugCPU:         *debugCPU,
		DebugLCD:         *debugLCD,
		DebugFlowControl: *debugFlowControl,
		DebugJumps:       *debugJumps,
	}

	// Start running the Gameboy with a GL UI
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gameboy := gb.NewGameboy(opts, cancel)
	gui := ui.NewGL(&gameboy)
	if *enableTiming {
		gameboy.Time(ctx, gui)
	} else {
		gameboy.Run(ctx, gui)
	}
}
