package ui

// UserInput captures user input for D-pad and buttons
type UserInput struct {
	InputRecv      bool
	DirectionInput uint8
	ButtonInput    uint8
}

// UI  abstracts over the user interface
type UI interface {
	HandleFrame([256 * 144]uint8)
	UserInput() *UserInput
	Shutdown()
}
