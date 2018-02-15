package ui

import "github.com/scottyw/goomba/lcd"

// UI abstracts over the user interface
type UI interface {
	ShouldRun() bool
	DrawFrame(*lcd.LCD)
	Shutdown()
}
