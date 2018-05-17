package ui

// Nop implements a UI with no real input/output, usually used for testing
type Nop struct {
	input *UserInput
	debug bool
}

// NewNop implements a no-op user interface
func NewNop(debug bool) UI {
	input := UserInput{
		DirectionInput: 0x0f,
		ButtonInput:    0x0f,
	}
	return &Nop{
		input: &input,
		debug: debug,
	}
}

// UserInput returns a data structure containing user input
func (n *Nop) UserInput() *UserInput {
	return n.input
}

// KeepRunning indicates whether the emulator should be running
func (n *Nop) KeepRunning() bool {
	return true
}

// Shutdown ...
func (n *Nop) Shutdown() {
	// Do nothing
}

// HandleFrame ...
func (n *Nop) HandleFrame(lcd [23040]uint8) {
	// Do nothing
}
