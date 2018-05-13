package options

import (
	"flag"
)

// RomFilename indicates the file containing the ROM to be run
var RomFilename = flag.String("rom", "", "ROM filename")

// DebugCPU is true if CPU debugging is enabled
var DebugCPU = flag.Bool("debug-cpu", false, "When true, CPU debugging is enabled")

// DebugFlowControl is true if debugging of call/ret/reti/rst is enabled
var DebugFlowControl = flag.Bool("debug-flow", false, "When true, flow control debugging is enabled")

// DebugJumps is true if debugging of jp/jr is enabled
var DebugJumps = flag.Bool("debug-jumps", false, "When true, jump debugging is enabled")

// DebugLCD is true if LCD debug colouring is enabled
var DebugLCD = flag.Bool("debug-lcd", false, "When true, LCD colour-based debugging is enabled")

// ShowSerialData is true if serial data should be written to console
var ShowSerialData = flag.Bool("show-serial-data", false, "When true, data sent to the serial port will be written to console")

func init() {
	flag.Parse()
}
