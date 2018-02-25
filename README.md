## Welcome to Tetromino

Tetromino is a Game Boy emulator written in Go.

### Screenshots

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
![Tetris](screenshots/tetris-title.png)

### Progress so far ...

Tetromino isn't functional yet but can run Tetris as far as the title screen. This probably seems deeply unimpressive but I can tell you that the moment I first saw that screen appear successfully was bliss ...

Things that work:

* Maybe 80% of the instructions are implemented, including flags
* The fetch-execute architecture is pretty sound, though I would like to tweak instruction dispatch
* V-Blank interrupt works
* LCD is wired up enought to display backgrounds (with some caveats below)

Here is a sample of the todo list:

* Implement the rest of the instructions
    * About three quarters are implemented now
    * Testing of what's there is not terrible but not great, and really just what fell out of TDD
* Implement controls
    * That's mostly a wiring job. A fiddly wiring job but nonetheless hooking GFLW keyboard input into the register which looks after that
* Implement more LCD stuff
    * Sprite support
    * Window support
    * I don't even think the background upper tile map is working ...
    * Implement the proper 160x144 display window - right now I display the full 256x256 video ram
* More on interrupts
    * Implement the rest of the interrupts
    * Multiple interrupts isn't right at the moment either
* Full on validation with test ROMs
    * Use blargg tests or similar to have more confidence that it's all solid
* Sound
    * Not even close to starting this one ...

### Running

You'll need a ROM. I highly recommend Tetris. For now, create an environment variable called ```ROM_FILENAME``` that contains the ROM filename.

    export ROM_FILENAME=/roms/tetris.gb

Run like this:

    go run main.go

### Common Errors

Tetris will corrupt the display and crash as soon as it hits demo mode. Every other ROM will probably immediately crash when it hits an unimplemented instruction. Errors will be very common in other words. Still, it's a bit of fun, isn't it?

### Dependencies

#### Go dependencies

Install these Go dependencies:

    github.com/go-gl/gl/v2.1/gl
    github.com/go-gl/glfw/v3.1/glfw

#### GLFW dependencies

Tetromino uses [GLFW](http://www.glfw.org) for graphics so you might need to install some OS-specific packages.

> * GLFW C library source is included and built automatically as part of the Go package. But you need to make sure you have dependencies of GLFW:
> 	* On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
> 	* On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev` and `xorg-dev` packages.
> 	* On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel` packages.

See [this page](https://github.com/go-gl/glfw) if those instructions don't work for you.

### References and Thanks

You can find a huge amount of great information about the Game Boy out there and many people have shared their work for others to build on. Thanks to everyone who has shared their experiences, code and documentation.

These ones have been invaluable for me:
* http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf
* http://bgb.bircd.org/pandocs.htm
* http://cturt.github.io/cinoop.html
* https://github.com/fogleman/nes
* http://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html
* https://github.com/lmmendes/game-boy-opcodes
