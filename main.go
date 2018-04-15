package main

import (
	"github.com/scottyw/tetromino/gb"
)

func main() {
	gameboy := gb.NewGameboy()
	gameboy.Run()
}
