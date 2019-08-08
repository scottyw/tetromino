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

Tetromino has accurate CPU, timer and MBC1 implementations but there is no sound, no support for other MBCs and sprite support is minimal (no large sprites, palettes or priority).

| Result             | Blargg test                  | Screenshot                                                 |
| ------------------ | ---------------------------- | ---------------------------------------------------------- |
| :green_heart: pass | cpu_instrs/cpu_instrs.gb     | [pic](pkg/gb/testresults/cpu_instrs_cpu_instrs.gb.png)     |
| :boom: fail        | dmg_sound/dmg_sound.gb       | [pic](pkg/gb/testresults/dmg_sound_dmg_sound.gb.png)       |
| :green_heart: pass | halt_bug.gb                  | [pic](pkg/gb/testresults/halt_bug.gb.png)                  |
| :green_heart: pass | instr_timing/instr_timing.gb | [pic](pkg/gb/testresults/instr_timing_instr_timing.gb.png) |
| :green_heart: pass | mem_timing/mem_timing.gb     | [pic](pkg/gb/testresults/mem_timing_mem_timing.gb.png)     |
| :green_heart: pass | mem_timing-2/mem_timing.gb   | [pic](pkg/gb/testresults/mem_timing-2_mem_timing.gb.png)   |
| :boom: fail        | oam_bug/oam_bug.gb           | [pic](pkg/gb/testresults/oam_bug_oam_bug.gb.png)           |

| Result             | Mooneye test                                   | Screenshot                                                                   |
| ------------------ | ---------------------------------------------- | ---------------------------------------------------------------------------- |
| :green_heart: pass | acceptance/add_sp_e_timing.gb                  | [pic](pkg/gb/testresults/acceptance_add_sp_e_timing.gb.png)                  |
| :green_heart: pass | acceptance/bits/mem_oam.gb                     | [pic](pkg/gb/testresults/acceptance_bits_mem_oam.gb.png)                     |
| :green_heart: pass | acceptance/bits/reg_f.gb                       | [pic](pkg/gb/testresults/acceptance_bits_reg_f.gb.png)                       |
| :green_heart: pass | acceptance/bits/unused_hwio-GS.gb              | [pic](pkg/gb/testresults/acceptance_bits_unused_hwio-GS.gb.png)              |
| :green_heart: pass | acceptance/boot_div-dmgABCmgb.gb               | [pic](pkg/gb/testresults/acceptance_boot_div-dmgABCmgb.gb.png)               |
| :green_heart: pass | acceptance/boot_hwio-dmgABCmgb.gb              | [pic](pkg/gb/testresults/acceptance_boot_hwio-dmgABCmgb.gb.png)              |
| :green_heart: pass | acceptance/boot_regs-dmgABC.gb                 | [pic](pkg/gb/testresults/acceptance_boot_regs-dmgABC.gb.png)                 |
| :green_heart: pass | acceptance/call_cc_timing.gb                   | [pic](pkg/gb/testresults/acceptance_call_cc_timing.gb.png)                   |
| :boom: fail        | acceptance/call_cc_timing2.gb                  | [pic](pkg/gb/testresults/acceptance_call_cc_timing2.gb.png)                  |
| :green_heart: pass | acceptance/call_timing.gb                      | [pic](pkg/gb/testresults/acceptance_call_timing.gb.png)                      |
| :boom: fail        | acceptance/call_timing2.gb                     | [pic](pkg/gb/testresults/acceptance_call_timing2.gb.png)                     |
| :boom: fail        | acceptance/di_timing-GS.gb                     | [pic](pkg/gb/testresults/acceptance_di_timing-GS.gb.png)                     |
| :green_heart: pass | acceptance/div_timing.gb                       | [pic](pkg/gb/testresults/acceptance_div_timing.gb.png)                       |
| :boom: fail        | acceptance/ei_sequence.gb                      | [pic](pkg/gb/testresults/acceptance_ei_sequence.gb.png)                      |
| :boom: fail        | acceptance/ei_timing.gb                        | [pic](pkg/gb/testresults/acceptance_ei_timing.gb.png)                        |
| :green_heart: pass | acceptance/halt_ime0_ei.gb                     | [pic](pkg/gb/testresults/acceptance_halt_ime0_ei.gb.png)                     |
| :green_heart: pass | acceptance/halt_ime0_nointr_timing.gb          | [pic](pkg/gb/testresults/acceptance_halt_ime0_nointr_timing.gb.png)          |
| :green_heart: pass | acceptance/halt_ime1_timing.gb                 | [pic](pkg/gb/testresults/acceptance_halt_ime1_timing.gb.png)                 |
| :boom: fail        | acceptance/halt_ime1_timing2-GS.gb             | [pic](pkg/gb/testresults/acceptance_halt_ime1_timing2-GS.gb.png)             |
| :green_heart: pass | acceptance/if_ie_registers.gb                  | [pic](pkg/gb/testresults/acceptance_if_ie_registers.gb.png)                  |
| :green_heart: pass | acceptance/instr/daa.gb                        | [pic](pkg/gb/testresults/acceptance_instr_daa.gb.png)                        |
| :boom: fail        | acceptance/interrupts/ie_push.gb               | [pic](pkg/gb/testresults/acceptance_interrupts_ie_push.gb.png)               |
| :green_heart: pass | acceptance/intr_timing.gb                      | [pic](pkg/gb/testresults/acceptance_intr_timing.gb.png)                      |
| :green_heart: pass | acceptance/jp_cc_timing.gb                     | [pic](pkg/gb/testresults/acceptance_jp_cc_timing.gb.png)                     |
| :green_heart: pass | acceptance/jp_timing.gb                        | [pic](pkg/gb/testresults/acceptance_jp_timing.gb.png)                        |
| :boom: fail        | acceptance/ld_hl_sp_e_timing.gb                | [pic](pkg/gb/testresults/acceptance_ld_hl_sp_e_timing.gb.png)                |
| :green_heart: pass | acceptance/oam_dma/basic.gb                    | [pic](pkg/gb/testresults/acceptance_oam_dma_basic.gb.png)                    |
| :green_heart: pass | acceptance/oam_dma/reg_read.gb                 | [pic](pkg/gb/testresults/acceptance_oam_dma_reg_read.gb.png)                 |
| :green_heart: pass | acceptance/oam_dma/sources-dmgABCmgbS.gb       | [pic](pkg/gb/testresults/acceptance_oam_dma_sources-dmgABCmgbS.gb.png)       |
| :green_heart: pass | acceptance/oam_dma_restart.gb                  | [pic](pkg/gb/testresults/acceptance_oam_dma_restart.gb.png)                  |
| :boom: fail        | acceptance/oam_dma_start.gb                    | [pic](pkg/gb/testresults/acceptance_oam_dma_start.gb.png)                    |
| :green_heart: pass | acceptance/oam_dma_timing.gb                   | [pic](pkg/gb/testresults/acceptance_oam_dma_timing.gb.png)                   |
| :green_heart: pass | acceptance/pop_timing.gb                       | [pic](pkg/gb/testresults/acceptance_pop_timing.gb.png)                       |
| :boom: fail        | acceptance/ppu/hblank_ly_scx_timing-GS.gb      | [pic](pkg/gb/testresults/acceptance_ppu_hblank_ly_scx_timing-GS.gb.png)      |
| :boom: fail        | acceptance/ppu/intr_1_2_timing-GS.gb           | [pic](pkg/gb/testresults/acceptance_ppu_intr_1_2_timing-GS.gb.png)           |
| :boom: fail        | acceptance/ppu/intr_2_0_timing.gb              | [pic](pkg/gb/testresults/acceptance_ppu_intr_2_0_timing.gb.png)              |
| :boom: fail        | acceptance/ppu/intr_2_mode0_timing.gb          | [pic](pkg/gb/testresults/acceptance_ppu_intr_2_mode0_timing.gb.png)          |
| :boom: fail        | acceptance/ppu/intr_2_mode0_timing_sprites.gb  | [pic](pkg/gb/testresults/acceptance_ppu_intr_2_mode0_timing_sprites.gb.png)  |
| :boom: fail        | acceptance/ppu/intr_2_mode3_timing.gb          | [pic](pkg/gb/testresults/acceptance_ppu_intr_2_mode3_timing.gb.png)          |
| :boom: fail        | acceptance/ppu/intr_2_oam_ok_timing.gb         | [pic](pkg/gb/testresults/acceptance_ppu_intr_2_oam_ok_timing.gb.png)         |
| :boom: fail        | acceptance/ppu/lcdon_timing-dmgABCmgbS.gb      | [pic](pkg/gb/testresults/acceptance_ppu_lcdon_timing-dmgABCmgbS.gb.png)      |
| :boom: fail        | acceptance/ppu/lcdon_write_timing-GS.gb        | [pic](pkg/gb/testresults/acceptance_ppu_lcdon_write_timing-GS.gb.png)        |
| :boom: fail        | acceptance/ppu/stat_irq_blocking.gb            | [pic](pkg/gb/testresults/acceptance_ppu_stat_irq_blocking.gb.png)            |
| :boom: fail        | acceptance/ppu/stat_lyc_onoff.gb               | [pic](pkg/gb/testresults/acceptance_ppu_stat_lyc_onoff.gb.png)               |
| :boom: fail        | acceptance/ppu/vblank_stat_intr-GS.gb          | [pic](pkg/gb/testresults/acceptance_ppu_vblank_stat_intr-GS.gb.png)          |
| :boom: fail        | acceptance/push_timing.gb                      | [pic](pkg/gb/testresults/acceptance_push_timing.gb.png)                      |
| :boom: fail        | acceptance/rapid_di_ei.gb                      | [pic](pkg/gb/testresults/acceptance_rapid_di_ei.gb.png)                      |
| :green_heart: pass | acceptance/ret_cc_timing.gb                    | [pic](pkg/gb/testresults/acceptance_ret_cc_timing.gb.png)                    |
| :green_heart: pass | acceptance/ret_timing.gb                       | [pic](pkg/gb/testresults/acceptance_ret_timing.gb.png)                       |
| :boom: fail        | acceptance/reti_intr_timing.gb                 | [pic](pkg/gb/testresults/acceptance_reti_intr_timing.gb.png)                 |
| :green_heart: pass | acceptance/reti_timing.gb                      | [pic](pkg/gb/testresults/acceptance_reti_timing.gb.png)                      |
| :boom: fail        | acceptance/rst_timing.gb                       | [pic](pkg/gb/testresults/acceptance_rst_timing.gb.png)                       |
| :boom: fail        | acceptance/serial/boot_sclk_align-dmgABCmgb.gb | [pic](pkg/gb/testresults/acceptance_serial_boot_sclk_align-dmgABCmgb.gb.png) |
| :green_heart: pass | acceptance/timer/div_write.gb                  | [pic](pkg/gb/testresults/acceptance_timer_div_write.gb.png)                  |
| :green_heart: pass | acceptance/timer/rapid_toggle.gb               | [pic](pkg/gb/testresults/acceptance_timer_rapid_toggle.gb.png)               |
| :green_heart: pass | acceptance/timer/tim00.gb                      | [pic](pkg/gb/testresults/acceptance_timer_tim00.gb.png)                      |
| :green_heart: pass | acceptance/timer/tim00_div_trigger.gb          | [pic](pkg/gb/testresults/acceptance_timer_tim00_div_trigger.gb.png)          |
| :green_heart: pass | acceptance/timer/tim01.gb                      | [pic](pkg/gb/testresults/acceptance_timer_tim01.gb.png)                      |
| :green_heart: pass | acceptance/timer/tim01_div_trigger.gb          | [pic](pkg/gb/testresults/acceptance_timer_tim01_div_trigger.gb.png)          |
| :green_heart: pass | acceptance/timer/tim10.gb                      | [pic](pkg/gb/testresults/acceptance_timer_tim10.gb.png)                      |
| :green_heart: pass | acceptance/timer/tim10_div_trigger.gb          | [pic](pkg/gb/testresults/acceptance_timer_tim10_div_trigger.gb.png)          |
| :green_heart: pass | acceptance/timer/tim11.gb                      | [pic](pkg/gb/testresults/acceptance_timer_tim11.gb.png)                      |
| :green_heart: pass | acceptance/timer/tim11_div_trigger.gb          | [pic](pkg/gb/testresults/acceptance_timer_tim11_div_trigger.gb.png)          |
| :green_heart: pass | acceptance/timer/tima_reload.gb                | [pic](pkg/gb/testresults/acceptance_timer_tima_reload.gb.png)                |
| :green_heart: pass | acceptance/timer/tima_write_reloading.gb       | [pic](pkg/gb/testresults/acceptance_timer_tima_write_reloading.gb.png)       |
| :green_heart: pass | acceptance/timer/tma_write_reloading.gb        | [pic](pkg/gb/testresults/acceptance_timer_tma_write_reloading.gb.png)        |
| :green_heart: pass | emulator-only/mbc1/bits_ram_en.gb              | [pic](pkg/gb/testresults/emulator-only_mbc1_bits_ram_en.gb.png)              |
| :boom: fail        | emulator-only/mbc1/multicart_rom_8Mb.gb        | [pic](pkg/gb/testresults/emulator-only_mbc1_multicart_rom_8Mb.gb.png)        |
| :green_heart: pass | emulator-only/mbc1/ram_256Kb.gb                | [pic](pkg/gb/testresults/emulator-only_mbc1_ram_256Kb.gb.png)                |
| :green_heart: pass | emulator-only/mbc1/ram_64Kb.gb                 | [pic](pkg/gb/testresults/emulator-only_mbc1_ram_64Kb.gb.png)                 |
| :green_heart: pass | emulator-only/mbc1/rom_16Mb.gb                 | [pic](pkg/gb/testresults/emulator-only_mbc1_rom_16Mb.gb.png)                 |
| :green_heart: pass | emulator-only/mbc1/rom_1Mb.gb                  | [pic](pkg/gb/testresults/emulator-only_mbc1_rom_1Mb.gb.png)                  |
| :green_heart: pass | emulator-only/mbc1/rom_2Mb.gb                  | [pic](pkg/gb/testresults/emulator-only_mbc1_rom_2Mb.gb.png)                  |
| :green_heart: pass | emulator-only/mbc1/rom_4Mb.gb                  | [pic](pkg/gb/testresults/emulator-only_mbc1_rom_4Mb.gb.png)                  |
| :green_heart: pass | emulator-only/mbc1/rom_512Kb.gb                | [pic](pkg/gb/testresults/emulator-only_mbc1_rom_512Kb.gb.png)                |
| :green_heart: pass | emulator-only/mbc1/rom_8Mb.gb                  | [pic](pkg/gb/testresults/emulator-only_mbc1_rom_8Mb.gb.png)                  |

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
