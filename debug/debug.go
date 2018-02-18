package debug

// CPU returns true if CPU debugging is enabled
func CPU(pc uint16) bool {
	return false
}

// FlowControl returns true if debugging of call/ret/reti/rst is enabled
func FlowControl() bool {
	return false
}

// Jumps returns true if debugging of jp/jr is enabled
func Jumps() bool {
	return false
}
