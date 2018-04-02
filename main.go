package main

import (
	"github.com/scottyw/tetromino/cpu"
	"github.com/scottyw/tetromino/gb"
	"github.com/scottyw/tetromino/lcd"
	"github.com/scottyw/tetromino/mem"
	"github.com/scottyw/tetromino/ui"
)

func main() {
	hwr := mem.NewHardwareRegisters()
	cpu := cpu.NewCPU(hwr)
	mem := mem.NewMemory(hwr)
	lcd := lcd.NewLCD(hwr, mem)
	ui := ui.NewGL(hwr, cpu)
	gameboy := gb.NewGameboy(cpu, mem, lcd, ui)
	gameboy.Run()
}
