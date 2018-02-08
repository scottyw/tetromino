package main

import (
	"fmt"
	"time"

	"github.com/scottyw/goomba/cpu"
	"github.com/scottyw/goomba/lcd"
	"github.com/scottyw/goomba/mem"
)

func runFrame(cpu cpu.CPU, mem mem.Memory) {
	// The Game Boy clock runs at 4.194304MHz
	// There are 4 clock cycles to a "machine cycle" giving 1048576 machine cycles per second
	// Each loop iteration below represents one machine cycle (i.e. 4 clock cycles)
	// Each LCD frame is 17556 machine cycles
	for cycle := 0; cycle < 17556; cycle++ {
		lcd.Tick(mem, cycle)
		cpu.Tick(mem)
	}
}

func runCPU(cpu cpu.CPU, mem mem.Memory) {
	for {
		runFrame(cpu, mem)
	}
}

func timeCPU(cpu cpu.CPU, mem mem.Memory) {
	for {
		// There are just under 60 frames per second (59.7275) so let's time in blocks of 60 frames
		// On a real Gameboy this would take 1 second
		t0 := time.Now()
		for i := 0; i < 60; i++ {
			runFrame(cpu, mem)
		}
		t1 := time.Now()
		fmt.Println("=========>", (t1.Sub(t0)))
	}
}

func main() {
	cpu := cpu.NewCPU()
	mem := mem.NewMemory()
	runCPU(cpu, mem)
	//	lcd.Run(mem)
}
