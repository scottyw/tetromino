package main

import (
	"github.com/scottyw/tetromino/cpu"
	"github.com/scottyw/tetromino/gb"
	"github.com/scottyw/tetromino/lcd"
	"github.com/scottyw/tetromino/mem"
	"github.com/scottyw/tetromino/ui"
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
