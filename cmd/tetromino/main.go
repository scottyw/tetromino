package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/scottyw/tetromino/pkg/gb"
	"github.com/scottyw/tetromino/pkg/ui"
)

func main() {

	// Command line flags
	fast := flag.Bool("fast", true, "When true, Tetromino runs the emulator as fast as possible (true by default)")
	silent := flag.Bool("silent", true, "When true, Tetromino disables all sound output (true by default)")
	debugCPU := flag.Bool("debugcpu", false, "When true, CPU debugging is enabled")
	debugTimer := flag.Bool("debugtimer", false, "When true, timer debugging is enabled")
	debugLCD := flag.Bool("debuglcd", false, "When true, colour-based LCD debugging is enabled")
	enableTiming := flag.Bool("timing", false, "When true, timing is output every 60 frames")
	enableProfiling := flag.Bool("profiling", false, "When true, CPU profiling data is written to 'cpuprofile.pprof'")
	flag.Parse()

	// CPU profiling
	if *enableProfiling {
		f, err := os.Create("cpuprofile.pprof")
		if err != nil {
			log.Fatalf("Failed to write cpuprofile.pprof: %v", err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatalf("Failed to start cpu profile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}

	rom := flag.Arg(0)
	if rom == "" {
		fmt.Println("No ROM filename was specified")
		os.Exit(1)
	}

	opts := gb.Options{
		RomFilename: rom,
		Fast:        *fast,
		Silent:      *silent,
		DebugCPU:    *debugCPU,
		DebugTimer:  *debugTimer,
		DebugLCD:    *debugLCD,
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
