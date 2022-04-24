package gameboy

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func runMooneyeTest(t *testing.T, filename string) {
	config := Config{
		RomFilename:        filename,
		DisableVideoOutput: true,
		DisableAudioOutput: true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	gameboy := New(config)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if gameboy.cpu.CheckMooneye() != nil {
					cancel()
					return
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	gameboy.Run(ctx)
	<-ctx.Done()
	screenshotFilename := fmt.Sprintf(
		"testresults/%s.png",
		strings.Replace(
			strings.TrimPrefix(filename, "testdata/mooneye-gb_hwtests/"),
			"/", "_", -1),
	)
	gameboy.ppu.Screenshot(screenshotFilename)
	if !reflect.DeepEqual(gameboy.cpu.CheckMooneye(), []uint8{0, 3, 5, 8, 13, 21, 34}) {
		t.Errorf("Test ROM failed: %s", filename)
		// fmt.Printf("| :boom: fail | %s | [pic](%s) |\n", filename, screenshotFilename)
		// } else {
		// fmt.Printf("| :green_heart: pass | %s | [pic](%s) |\n", filename, screenshotFilename)
	}

}

func TestMooneyeBoot(t *testing.T) {

	for _, filename := range []string{
		"testdata/mooneye-gb_hwtests/acceptance/boot_div-dmgABCmgb.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/boot_hwio-dmgABCmgb.gb",
		"testdata/mooneye-gb_hwtests/acceptance/boot_regs-dmgABC.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyeTiming(t *testing.T) {

	for _, filename := range []string{
		"testdata/mooneye-gb_hwtests/acceptance/add_sp_e_timing.gb",
		"testdata/mooneye-gb_hwtests/acceptance/call_cc_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/call_cc_timing2.gb",
		"testdata/mooneye-gb_hwtests/acceptance/call_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/call_timing2.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/di_timing-GS.gb",
		"testdata/mooneye-gb_hwtests/acceptance/div_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ei_sequence.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ei_timing.gb",
		"testdata/mooneye-gb_hwtests/acceptance/halt_ime0_ei.gb",
		"testdata/mooneye-gb_hwtests/acceptance/halt_ime0_nointr_timing.gb",
		"testdata/mooneye-gb_hwtests/acceptance/halt_ime1_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/halt_ime1_timing2-GS.gb",
		"testdata/mooneye-gb_hwtests/acceptance/if_ie_registers.gb",
		"testdata/mooneye-gb_hwtests/acceptance/intr_timing.gb",
		"testdata/mooneye-gb_hwtests/acceptance/jp_cc_timing.gb",
		"testdata/mooneye-gb_hwtests/acceptance/jp_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ld_hl_sp_e_timing.gb",
		"testdata/mooneye-gb_hwtests/acceptance/pop_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/push_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/rapid_di_ei.gb",
		"testdata/mooneye-gb_hwtests/acceptance/ret_cc_timing.gb",
		"testdata/mooneye-gb_hwtests/acceptance/ret_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/reti_intr_timing.gb",
		"testdata/mooneye-gb_hwtests/acceptance/reti_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/rst_timing.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyeBits(t *testing.T) {

	for _, filename := range []string{
		"testdata/mooneye-gb_hwtests/acceptance/bits/mem_oam.gb",
		"testdata/mooneye-gb_hwtests/acceptance/bits/reg_f.gb",
		"testdata/mooneye-gb_hwtests/acceptance/bits/unused_hwio-GS.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyeDAA(t *testing.T) {

	for _, filename := range []string{
		"testdata/mooneye-gb_hwtests/acceptance/instr/daa.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyeInterrupts(t *testing.T) {

	for _, filename := range []string{
		// "testdata/mooneye-gb_hwtests/acceptance/interrupts/ie_push.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyeOAMDMA(t *testing.T) {

	for _, filename := range []string{
		"testdata/mooneye-gb_hwtests/acceptance/oam_dma_restart.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/oam_dma_start.gb",
		"testdata/mooneye-gb_hwtests/acceptance/oam_dma_timing.gb",
		"testdata/mooneye-gb_hwtests/acceptance/oam_dma/basic.gb",
		"testdata/mooneye-gb_hwtests/acceptance/oam_dma/reg_read.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/oam_dma/sources-GS.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyePPU(t *testing.T) {

	for _, filename := range []string{
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/hblank_ly_scx_timing-GS.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/intr_1_2_timing-GS.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/intr_2_0_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/intr_2_mode0_timing_sprites.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/intr_2_mode0_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/intr_2_mode3_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/intr_2_oam_ok_timing.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/lcdon_timing-GS.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/lcdon_write_timing-GS.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/stat_irq_blocking.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/stat_lyc_onoff.gb",
		// "testdata/mooneye-gb_hwtests/acceptance/ppu/vblank_stat_intr-GS.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyeSerial(t *testing.T) {

	for _, filename := range []string{
		// "testdata/mooneye-gb_hwtests/acceptance/serial/boot_sclk_align-dmgABCmgb.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyeTimer(t *testing.T) {

	for _, filename := range []string{
		"testdata/mooneye-gb_hwtests/acceptance/timer/div_write.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/rapid_toggle.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tim00_div_trigger.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tim00.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tim01_div_trigger.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tim01.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tim10_div_trigger.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tim10.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tim11_div_trigger.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tim11.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tima_reload.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tima_write_reloading.gb",
		"testdata/mooneye-gb_hwtests/acceptance/timer/tma_write_reloading.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyeMBC1(t *testing.T) {

	for _, filename := range []string{
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/bits_bank1.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/bits_bank2.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/bits_mode.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/bits_ramg.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc1/multicart_rom_8Mb.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/ram_256Kb.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/ram_64Kb.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/rom_16Mb.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/rom_1Mb.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/rom_2Mb.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/rom_4Mb.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/rom_512Kb.gb",
		"testdata/mooneye-gb_hwtests/emulator-only/mbc1/rom_8Mb.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyeMBC2(t *testing.T) {

	for _, filename := range []string{
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc2/bits_ramg.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc2/bits_romb.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc2/bits_unused.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc2/ram.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc2/rom_1Mb.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc2/rom_2Mb.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc2/rom_512kb.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}

func TestMooneyeMBC5(t *testing.T) {

	for _, filename := range []string{
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc5/rom_16Mb.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc5/rom_1Mb.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc5/rom_2Mb.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc5/rom_32Mb.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc5/rom_4Mb.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc5/rom_512kb.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc5/rom_64Mb.gb",
		// "testdata/mooneye-gb_hwtests/emulator-only/mbc5/rom_8Mb.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}

}
