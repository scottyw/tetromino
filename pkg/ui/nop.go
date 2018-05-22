package ui

// Nop implements a UI with no real input/output, usually used for testing
type Nop struct {
	input *UserInput
}

// NewNop implements a no-op user interface
func NewNop() UI {
	input := UserInput{
		DirectionInput: 0x0f,
		ButtonInput:    0x0f,
	}
	return &Nop{
		input: &input,
	}
}

// UserInput returns a data structure containing user input
func (n *Nop) UserInput() *UserInput {
	return n.input
}

// HandleFrame ...
func (n *Nop) HandleFrame(lcd [23040]uint8) {
	// Do nothing
}

// Shutdown ...
func (n *Nop) Shutdown() {
	// Do nothing
}
