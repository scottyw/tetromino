## Welcome to Tetromino

[![Build Status](https://travis-ci.org/scottyw/tetromino.svg?branch=master)](https://travis-ci.org/scottyw/tetromino) [![Go Report Card](https://goreportcard.com/badge/github.com/scottyw/tetromino)](https://goreportcard.com/report/github.com/scottyw/tetromino)

Tetromino is a Game Boy emulator written in Go.

### Tetris
![Tetris demo](https://github.com/scottyw/tetromino/blob/master/screenshots/tetris/Large%20GIF%20(480x432).gif)

### Super Mario Land
![Mario demo](https://github.com/scottyw/tetromino/blob/master/screenshots/mario/Large%20GIF%20(480x432).gif)

### Pokemon Red
![Pokemon demo](https://github.com/scottyw/tetromino/blob/master/screenshots/pokemon/Large%20GIF%20(480x432).gif)

### Debugging the LCD

Tetromino has some fun LCD debugging that uses colour to differentiate sprites (in blue) from background (white when on-screen and red when off-screen) to show how scrolling is implemented.

![Debugging Mario](https://github.com/scottyw/tetromino/blob/master/screenshots/mario-debug/Large%20GIF%20(766x434).gif)

### Running

You'll need a ROM file which you can specify like this:

    go run cmd/tetromino/main.go /roms/tetris.gb

Other options exist including enabling debug. List them like this:

    go run cmd/tetromino/main.go -help

Note that flag parsing follows the Go language rules and so flags must be specified before the ROM filename e.g.

    go run cmd/tetromino/main.go --debuglcd /roms/tetris.gb

### Controls

Arrows keys : Up/Down/Left/Right
A : Start
S : Select
Z : B button
X : A button
T : Take screenshot

### Tests

Tetromino is pretty functional and passes the blargg CPU and timing tests, plus some mooneye tests. The CPU, timer and MBC1 implementations are accurate but there is no sound, no support for other MBCs and sprite support is minimal (no large sprites, palettes or priority). Most of the sub-instruction timing isn't accurate either.

![Blargg CPU tests pass](https://github.com/scottyw/tetromino/blob/master/screenshots/blargg-cpu-tests.png)
![Blargg timing tests pass](https://github.com/scottyw/tetromino/blob/master/screenshots/blargg-timing-tests.png)

### Dependencies

Tetromino uses Go modules and requires Go 1.11 or later. You may need to enable module support like this:

    export GO111MODULE=on

When you run Tetromino or the tests, the dependencies will be fetched automatically.

#### GLFW dependencies

Tetromino uses [GLFW](http://www.glfw.org) for graphics so you might need to install some OS-specific packages.

> * GLFW C library source is included and built automatically as part of the Go package. But you need to make sure you have dependencies of GLFW:
> 	* On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
> 	* On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev` and `xorg-dev` packages.
> 	* On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel` packages.

See [this page](https://github.com/go-gl/glfw) if those instructions don't work for you.

### References and Thanks

You can find a huge amount of great information about the Game Boy out there and many people have shared their work for others to build on. Thanks to everyone who has shared their experiences, code and documentation.

Incredible test roms from:
* blargg
* Gekkio (https://github.com/Gekkio/mooneye-gb)

These resources have been invaluable for me:
* http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf
* http://bgb.bircd.org/pandocs.htm
* http://cturt.github.io/cinoop.html
* https://github.com/fogleman/nes
* http://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html
* https://github.com/lmmendes/game-boy-opcodes
