package main

import (
	"github.com/scottyw/goomba/cpu"
	"github.com/scottyw/goomba/gb"
	"github.com/scottyw/goomba/lcd"
	"github.com/scottyw/goomba/mem"
	"github.com/scottyw/goomba/ui"
)

func main() {
	gameboy := gb.NewGameboy(
		cpu.NewCPU(),
		mem.NewMemory(),
		lcd.NewLCD(),
		ui.NewGL(),
	)
	gameboy.Run()
}
