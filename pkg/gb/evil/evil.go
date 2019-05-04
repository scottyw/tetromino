package evil

import (
	"fmt"
	"math/rand"
)

const (
	//
	// Tetris shapes
	//
	o  = 0
	l  = 1
	rl = 2
	i  = 3
	z  = 5
	s  = 6
	t  = 7
)

// TotalCycles executed by the emulator
var TotalCycles int64
var previousCycles int64
var previousWells int
var shape byte

// NextShape to give the player
func NextShape(videoRAM [0x2000]byte) byte {
	currentCycles := TotalCycles
	if currentCycles-previousCycles > 100 {
		// Enough cycles have passed that we need a different shape
		wells := evaluateGame(videoRAM)
		shape = pickShape(wells)
		// fmt.Println()
	}
	previousCycles = currentCycles
	return shape
}

func evaluateGame(videoRAM [0x2000]byte) int {

	// Check heights for each column
	cols := make([]int, 12)
	cols[0] = 18
	cols[11] = 18
	for tileX := 2; tileX < 12; tileX++ {
		col := tileX - 1
		for tileY := 0; tileY < 18; tileY++ {
			tileAddr := 32*uint16(tileY) + uint16(tileX)
			tileByte := videoRAM[(0x9800+tileAddr)&0x7fff]
			if tileByte != 0x2f {
				cols[col] = 18 - tileY
				break
			}
		}
	}
	// fmt.Println("HEIGHTS", cols)

	// Now looks for wells i.e. columns withs walls of 3 or more
	var wells int
	for i := 1; i < 11; i++ {
		if cols[i-1]-cols[i] >= 3 && cols[i+1]-cols[i] >= 3 {
			wells++
		}
	}
	// fmt.Println("WELLS", wells)

	return wells
}

func pickShape(wells int) byte {
	x := byte(rand.Intn(8))
	for x == 4 {
		// 4 is not a valid shape so pick again
		x = byte(rand.Intn(8))
	}
	if wells == 1 || wells == 2 {
		// Looks like the player needs a straight
		// Let's make sure they don't get one ...
		fmt.Println("LOL, no straight for you")
		for x == i {
			x = pickShape(wells)
		}
	} else if previousWells > wells {
		// Looks like the player just blocked off a well
		// Let's give them a shot at the straight they needed all along ...
		fmt.Println("Here's a straight, HAHAHAHA")
		x = i
	}
	previousWells = wells
	return x
}
