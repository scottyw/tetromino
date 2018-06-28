package ui

// Button represents a direction pad or button control
type Button int

const (
	// Up on the control pad
	Up = iota
	//Down on the control pad
	Down = iota
	// Left on the control pad
	Left = iota
	// Right on the control pad
	Right = iota
	// A button
	A = iota
	// B button
	B = iota
	// Start button
	Start = iota
	// Select button
	Select = iota
)

// Emulator allows the UI to interact with the emulator itself
type Emulator interface {
	ButtonAction(Button, bool)
	Screenshot()
	Debug() bool
	Shutdown()
}
