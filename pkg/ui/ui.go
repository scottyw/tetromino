package ui

import (
	"github.com/scottyw/tetromino/pkg/gb/lcd"
)

// UI abstracts over the user interface
type UI interface {
	ShouldRun() bool
	DrawFrame(*lcd.LCD)
	Shutdown()
}
