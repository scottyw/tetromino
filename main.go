package main

import (
	"github.com/scottyw/goomba/cpu"
	"github.com/scottyw/goomba/lcd"
	"github.com/scottyw/goomba/mem"
)

func main() {
	cpu := cpu.NewCPU()
	mem := mem.NewMemory()
	// The Game Boy clock runs at 4.194304MHz
	// Each loop iteration represents one machine cycle (i.e. 4 clock cycles)
	for {
		// t0 := time.Now()
		// for i := 0; i < 60; i++ {
		// Each LCD frame is 17556 machine cycles
		for cycle := 0; cycle < 17556; cycle++ {
			lcd.Tick(mem, cycle)
			cpu.Tick(mem)
		}
		// }
		// t1 := time.Now()
		// fmt.Println("=========>", (t1.Sub(t0)))
	}
}
