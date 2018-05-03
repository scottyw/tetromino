package main

import (
	"github.com/scottyw/tetromino/pkg/gb"
)

func main() {
	gameboy := gb.NewGameboy()
	gameboy.Run()
}
