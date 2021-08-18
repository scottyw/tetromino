package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/scottyw/tetromino/gameboy"
	"github.com/scottyw/tetromino/gameboy/display"
	"github.com/scottyw/tetromino/gameboy/speakers"
)

func main() {

	// Command line flags
	fast := flag.Bool("fast", false, "When true, Tetromino runs the emulator as fast as possible (audio support is disabled)")
	debugCPU := flag.Bool("debugcpu", false, "When true, CPU debugging is enabled")
	debugLCD := flag.Bool("debuglcd", false, "When true, colour-based LCD debugging is enabled")
	enableTiming := flag.Bool("timing", false, "When true, timing is output every 60 frames")
	enableProfiling := flag.Bool("profiling", false, "When true, CPU profiling data is written to 'cpuprofile.pprof'")
	flag.Parse()

	// CPU profiling
	if *enableProfiling {
		f, err := os.Create("cpuprofile.pprof")
		if err != nil {
			log.Printf("Failed to write cpuprofile.pprof: %v", err)
			return
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Printf("Failed to start cpu profile: %v", err)
			return
		}
		defer pprof.StopCPUProfile()
	}

	rom := flag.Arg(0)
	if rom == "" {
		fmt.Println("No ROM filename was specified")
		os.Exit(1)
	}

	opts := gameboy.Options{
		RomFilename: rom,
		DebugCPU:    *debugCPU,
		DebugLCD:    *debugLCD,
	}

	// Run context
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Create the Gameboy emulator
	gameboy := gameboy.NewGameboy(opts)

	// Create a display
	display, err := display.NewGLDisplay(gameboy, cancelFunc)
	if err != nil {
		log.Printf("Failed to create display: %v", err)
		return
	}
	defer display.Cleanup()
	gameboy.RegisterDisplay(display)

	// Create speakers if we are not running in fast mode
	if !*fast {
		speakers, err := speakers.NewPortaudioSpeakers()
		if err != nil {
			log.Printf("Failed to create speakers: %v", err)
			return
		}
		defer speakers.Cleanup()
		gameboy.RegisterSpeakers(speakers)
	}

	// Start running the emulator
	if *enableTiming {
		gameboy.Time(ctx)
	} else {
		gameboy.Run(ctx)
	}

}
