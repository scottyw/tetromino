package options

// DebugCPU returns true if CPU debugging is enabled
func DebugCPU(pc uint16) bool {
	return false
}

// DebugFlowControl returns true if debugging of call/ret/reti/rst is enabled
func DebugFlowControl() bool {
	return false
}

// DebugJumps returns true if debugging of jp/jr is enabled
func DebugJumps() bool {
	return false
}

// DebugLCD returns true if LCD debug colouring is enabled
func DebugLCD() bool {
	return false
}
