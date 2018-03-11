package main

import (
	"github.com/scottyw/tetromino/cpu"
	"github.com/scottyw/tetromino/gb"
	"github.com/scottyw/tetromino/lcd"
	"github.com/scottyw/tetromino/mem"
	"github.com/scottyw/tetromino/ui"
)

func main() {
	cpu := cpu.NewCPU()
	mem := mem.NewMemory()
	lcd := lcd.NewLCD()
	ui := ui.NewGL(mem)
	gameboy := gb.NewGameboy(cpu, mem, lcd, ui)
	gameboy.Run()
}
