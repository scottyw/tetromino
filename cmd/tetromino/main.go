package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"runtime/pprof"

	"github.com/scottyw/tetromino/pkg/gb"
	"github.com/scottyw/tetromino/pkg/ui"
)

func main() {

	// Command line flags
	romFilename := flag.String("f", "", "ROM filename")
	outputSerial := flag.Bool("output-serial-data", false, "When true, data sent to the serial port will be written to console")
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
		SBWriter:         sbWriter,
		DebugCPU:         *debugCPU,
		DebugFlowControl: *debugFlowControl,
		DebugJumps:       *debugJumps,
	}

	// Start running the Gameboy with a GL UI
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	ui := ui.NewGL(cancelFunc, *debugLCD)
	gameboy := gb.NewGameboy(ui, opts)
	if *enableTiming {
		gameboy.Time(ctx)
	} else {
		gameboy.Run(ctx)
	}
}
