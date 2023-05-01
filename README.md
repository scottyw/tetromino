## Welcome to Tetromino

[![Build Status](https://travis-ci.org/scottyw/tetromino.svg?branch=master)](https://travis-ci.org/scottyw/tetromino) [![Go Report Card](https://goreportcard.com/badge/github.com/scottyw/tetromino)](https://goreportcard.com/report/github.com/scottyw/tetromino)

Tetromino is a Game Boy emulator written in Go.

### Tetris
![Tetris demo](https://github.com/scottyw/tetromino/blob/main/screenshots/tetris/Large%20GIF%20(480x432).gif)
### Super Mario Land
![Mario demo](https://github.com/scottyw/tetromino/blob/main/screenshots/mario/Large%20GIF%20(480x432).gif)

### Pokemon Red
![Pokemon demo](https://github.com/scottyw/tetromino/blob/main/screenshots/pokemon/Large%20GIF%20(480x432).gif)

### Debugging the LCD

Tetromino has some fun LCD debugging that uses colour to differentiate sprites (in blue) from background (white when on-screen and red when off-screen) to show how scrolling is implemented.

![Mario running with LCD debugging enabled](https://github.com/scottyw/tetromino/blob/main/screenshots/mario-debug/Large%20GIF%20(766x434).gif)

### Building and Running

Build Tetromino like this. You may need to install some OS-specific packages to support video and sound - see below for details.

    go build

To run, you'll need a ROM file which you can specify like this:

    tetromino tetris.gb

Other options exist including enabling debug. List them like this:

    tetromino --help

Flags must be specified before the ROM filename e.g.

    tetromino --debuglcd /roms/tetris.gb

### Controls

Arrows keys : Up/Down/Left/Right
A : Start
S : Select
Z : B button
X : A button
T : Take screenshot

### Tests

Tetromino has accurate CPU, timer, sound and MBC1 implementations (though no support for other MBCs). 

#### Blargg Tests

All Blargg tests pass.

| Result             | Blargg test                  | Screenshot                                                 |
| ------------------ | ---------------------------- | ---------------------------------------------------------- |
| :green_heart: pass | cpu_instrs/cpu_instrs.gb     | [pic](pkg/gb/testresults/cpu_instrs_cpu_instrs.gb.png)     |
| :green_heart: pass | dmg_sound/dmg_sound.gb       | [pic](pkg/gb/testresults/dmg_sound_dmg_sound.gb.png)       |
| :green_heart: pass | halt_bug.gb                  | [pic](pkg/gb/testresults/halt_bug.gb.png)                  |
| :green_heart: pass | instr_timing/instr_timing.gb | [pic](pkg/gb/testresults/instr_timing_instr_timing.gb.png) |
| :green_heart: pass | mem_timing/mem_timing.gb     | [pic](pkg/gb/testresults/mem_timing_mem_timing.gb.png)     |
| :green_heart: pass | mem_timing-2/mem_timing.gb   | [pic](pkg/gb/testresults/mem_timing-2_mem_timing.gb.png)   |
| :green_heart: pass | oam_bug/oam_bug.gb           | [pic](pkg/gb/testresults/oam_bug_oam_bug.gb.png)           |

#### Mooneye Tests

Some Mooneye tests pass (65 of 94).

| Result             | Mooneye test                                                         | Screenshot                                                      |
| ------------------ | -------------------------------------------------------------------- | --------------------------------------------------------------- |
| :green_heart: pass | acceptance/add_sp_e_timing.gb | [pic](testresults/acceptance_add_sp_e_timing.gb.png) |
| :green_heart: pass | acceptance/bits/mem_oam.gb | [pic](testresults/acceptance_bits_mem_oam.gb.png) |
| :green_heart: pass | acceptance/bits/reg_f.gb | [pic](testresults/acceptance_bits_reg_f.gb.png) |
| :green_heart: pass | acceptance/bits/unused_hwio-GS.gb | [pic](testresults/acceptance_bits_unused_hwio-GS.gb.png) |
| :green_heart: pass | acceptance/boot_div-dmgABCmgb.gb | [pic](testresults/acceptance_boot_div-dmgABCmgb.gb.png) |
| :green_heart: pass | acceptance/boot_regs-dmgABC.gb | [pic](testresults/acceptance_boot_regs-dmgABC.gb.png) |
| :green_heart: pass | acceptance/call_cc_timing.gb | [pic](testresults/acceptance_call_cc_timing.gb.png) |
| :green_heart: pass | acceptance/call_timing.gb | [pic](testresults/acceptance_call_timing.gb.png) |
| :green_heart: pass | acceptance/div_timing.gb | [pic](testresults/acceptance_div_timing.gb.png) |
| :green_heart: pass | acceptance/halt_ime0_ei.gb | [pic](testresults/acceptance_halt_ime0_ei.gb.png) |
| :green_heart: pass | acceptance/halt_ime0_nointr_timing.gb | [pic](testresults/acceptance_halt_ime0_nointr_timing.gb.png) |
| :green_heart: pass | acceptance/halt_ime1_timing.gb | [pic](testresults/acceptance_halt_ime1_timing.gb.png) |
| :green_heart: pass | acceptance/if_ie_registers.gb | [pic](testresults/acceptance_if_ie_registers.gb.png) |
| :green_heart: pass | acceptance/instr/daa.gb | [pic](testresults/acceptance_instr_daa.gb.png) |
| :green_heart: pass | acceptance/intr_timing.gb | [pic](testresults/acceptance_intr_timing.gb.png) |
| :green_heart: pass | acceptance/jp_cc_timing.gb | [pic](testresults/acceptance_jp_cc_timing.gb.png) |
| :green_heart: pass | acceptance/jp_timing.gb | [pic](testresults/acceptance_jp_timing.gb.png) |
| :green_heart: pass | acceptance/oam_dma_restart.gb | [pic](testresults/acceptance_oam_dma_restart.gb.png) |
| :green_heart: pass | acceptance/oam_dma_timing.gb | [pic](testresults/acceptance_oam_dma_timing.gb.png) |
| :green_heart: pass | acceptance/oam_dma/basic.gb | [pic](testresults/acceptance_oam_dma_basic.gb.png) |
| :green_heart: pass | acceptance/oam_dma/reg_read.gb | [pic](testresults/acceptance_oam_dma_reg_read.gb.png) |
| :green_heart: pass | acceptance/pop_timing.gb | [pic](testresults/acceptance_pop_timing.gb.png) |
| :green_heart: pass | acceptance/ret_cc_timing.gb | [pic](testresults/acceptance_ret_cc_timing.gb.png) |
| :green_heart: pass | acceptance/ret_timing.gb | [pic](testresults/acceptance_ret_timing.gb.png) |
| :green_heart: pass | acceptance/reti_timing.gb | [pic](testresults/acceptance_reti_timing.gb.png) |
| :green_heart: pass | acceptance/timer/div_write.gb | [pic](testresults/acceptance_timer_div_write.gb.png) |
| :green_heart: pass | acceptance/timer/rapid_toggle.gb | [pic](testresults/acceptance_timer_rapid_toggle.gb.png) |
| :green_heart: pass | acceptance/timer/tim00_div_trigger.gb | [pic](testresults/acceptance_timer_tim00_div_trigger.gb.png) |
| :green_heart: pass | acceptance/timer/tim00.gb | [pic](testresults/acceptance_timer_tim00.gb.png) |
| :green_heart: pass | acceptance/timer/tim01_div_trigger.gb | [pic](testresults/acceptance_timer_tim01_div_trigger.gb.png) |
| :green_heart: pass | acceptance/timer/tim01.gb | [pic](testresults/acceptance_timer_tim01.gb.png) |
| :green_heart: pass | acceptance/timer/tim10_div_trigger.gb | [pic](testresults/acceptance_timer_tim10_div_trigger.gb.png) |
| :green_heart: pass | acceptance/timer/tim10.gb | [pic](testresults/acceptance_timer_tim10.gb.png) |
| :green_heart: pass | acceptance/timer/tim11_div_trigger.gb | [pic](testresults/acceptance_timer_tim11_div_trigger.gb.png) |
| :green_heart: pass | acceptance/timer/tim11.gb | [pic](testresults/acceptance_timer_tim11.gb.png) |
| :green_heart: pass | acceptance/timer/tima_reload.gb | [pic](testresults/acceptance_timer_tima_reload.gb.png) |
| :green_heart: pass | acceptance/timer/tima_write_reloading.gb | [pic](testresults/acceptance_timer_tima_write_reloading.gb.png) |
| :green_heart: pass | acceptance/timer/tma_write_reloading.gb | [pic](testresults/acceptance_timer_tma_write_reloading.gb.png) |
| :green_heart: pass | emulator-only/mbc1/bits_bank1.gb | [pic](testresults/emulator-only_mbc1_bits_bank1.gb.png) |
| :green_heart: pass | emulator-only/mbc1/bits_bank2.gb | [pic](testresults/emulator-only_mbc1_bits_bank2.gb.png) |
| :green_heart: pass | emulator-only/mbc1/bits_mode.gb | [pic](testresults/emulator-only_mbc1_bits_mode.gb.png) |
| :green_heart: pass | emulator-only/mbc1/bits_ramg.gb | [pic](testresults/emulator-only_mbc1_bits_ramg.gb.png) |
| :green_heart: pass | emulator-only/mbc1/ram_256kb.gb | [pic](testresults/emulator-only_mbc1_ram_256kb.gb.png) |
| :green_heart: pass | emulator-only/mbc1/ram_64kb.gb | [pic](testresults/emulator-only_mbc1_ram_64kb.gb.png) |
| :green_heart: pass | emulator-only/mbc1/rom_16Mb.gb | [pic](testresults/emulator-only_mbc1_rom_16Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc1/rom_1Mb.gb | [pic](testresults/emulator-only_mbc1_rom_1Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc1/rom_2Mb.gb | [pic](testresults/emulator-only_mbc1_rom_2Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc1/rom_4Mb.gb | [pic](testresults/emulator-only_mbc1_rom_4Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc1/rom_512kb.gb | [pic](testresults/emulator-only_mbc1_rom_512kb.gb.png) |
| :green_heart: pass | emulator-only/mbc1/rom_8Mb.gb | [pic](testresults/emulator-only_mbc1_rom_8Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc2/bits_ramg.gb | [pic](testresults/emulator-only_mbc2_bits_ramg.gb.png) |
| :green_heart: pass | emulator-only/mbc2/bits_romb.gb | [pic](testresults/emulator-only_mbc2_bits_romb.gb.png) |
| :green_heart: pass | emulator-only/mbc2/bits_unused.gb | [pic](testresults/emulator-only_mbc2_bits_unused.gb.png) |
| :green_heart: pass | emulator-only/mbc2/ram.gb | [pic](testresults/emulator-only_mbc2_ram.gb.png) |
| :green_heart: pass | emulator-only/mbc2/rom_1Mb.gb | [pic](testresults/emulator-only_mbc2_rom_1Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc2/rom_2Mb.gb | [pic](testresults/emulator-only_mbc2_rom_2Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc2/rom_512kb.gb | [pic](testresults/emulator-only_mbc2_rom_512kb.gb.png) |
| :green_heart: pass | emulator-only/mbc5/rom_16Mb.gb | [pic](testresults/emulator-only_mbc5_rom_16Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc5/rom_1Mb.gb | [pic](testresults/emulator-only_mbc5_rom_1Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc5/rom_2Mb.gb | [pic](testresults/emulator-only_mbc5_rom_2Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc5/rom_32Mb.gb | [pic](testresults/emulator-only_mbc5_rom_32Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc5/rom_4Mb.gb | [pic](testresults/emulator-only_mbc5_rom_4Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc5/rom_512kb.gb | [pic](testresults/emulator-only_mbc5_rom_512kb.gb.png) |
| :green_heart: pass | emulator-only/mbc5/rom_64Mb.gb | [pic](testresults/emulator-only_mbc5_rom_64Mb.gb.png) |
| :green_heart: pass | emulator-only/mbc5/rom_8Mb.gb | [pic](testresults/emulator-only_mbc5_rom_8Mb.gb.png) |

### Dependencies

Tetromino uses [GLFW](http://www.glfw.org) for video and [PortAudio](http://www.portaudio.com) for sound so you might need to install some OS-specific packages.

#### GLFW

> * GLFW C library source is included and built automatically as part of the Go package. But you need to make sure you have dependencies of GLFW:
> 	* On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
> 	* On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev` and `xorg-dev` packages.
> 	* On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel` packages.

See [this page](https://github.com/go-gl/glfw) if those instructions don't work for you.

#### PortAudio

On macOS: `brew install portaudio`

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
