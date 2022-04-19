package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/scottyw/tetromino/gameboy"
)

func main() {

	// Command line flags
	fast := flag.Bool("fast", false, "When true, Tetromino runs the emulator as fast as possible (audio support is disabled)")
	debugCPU := flag.Bool("debugcpu", false, "When true, CPU debugging is enabled")
	debugLCD := flag.Bool("debuglcd", false, "When true, colour-based LCD debugging is enabled")
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
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			pprof.StopCPUProfile()
			os.Exit(0)
		}()
	}

	rom := flag.Arg(0)
	if rom == "" {
		fmt.Println("No ROM filename was specified")
		os.Exit(1)
	}

	// Fast mode requires audio to be disabled
	config := gameboy.Config{
		RomFilename:        rom,
		DisableVideoOutput: false,
		DisableAudioOutput: *fast,
		DebugCPU:           *debugCPU,
		DebugLCD:           *debugLCD,
	}

	// Create the Gameboy emulator
	gameboy := gameboy.New(config)

	// Start running the emulator
	gameboy.Run(context.Background())

}
