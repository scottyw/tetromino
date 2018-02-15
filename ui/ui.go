package ui

import (
	"github.com/scottyw/goomba/lcd"
	"github.com/scottyw/goomba/mem"
)

// UI abstracts over the user interface
type UI interface {
	ShouldRun() bool
	DrawFrame(*lcd.LCD, mem.Memory)
	Shutdown()
}
