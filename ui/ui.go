package ui

import (
	"github.com/scottyw/tetromino/lcd"
	"github.com/scottyw/tetromino/mem"
)

// UI abstracts over the user interface
type UI interface {
	ShouldRun() bool
	DrawFrame(*lcd.LCD, mem.Memory)
	Shutdown()
}
