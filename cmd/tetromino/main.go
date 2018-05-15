package main

import (
	"flag"
	"io"
	"os"

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
	flag.Parse()

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
	ui := ui.NewGL(*debugLCD)
	gameboy := gb.NewGameboy(ui, opts)
	gameboy.Run()
}
