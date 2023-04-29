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
			strings.TrimPrefix(filename, "testdata/mts-20221022-1430-8d742b9/"),
			"/", "_", -1),
	)
	gameboy.ppu.Screenshot(screenshotFilename)
	if !reflect.DeepEqual(gameboy.cpu.CheckMooneye(), []uint8{3, 5, 8, 13, 21, 34}) {
		t.Errorf("Test ROM failed: %s", filename)
		fmt.Printf("| :boom: fail | %s | [pic](%s) |\n", filename, screenshotFilename)
	} else {
		fmt.Printf("| :green_heart: pass | %s | [pic](%s) |\n", filename, screenshotFilename)
	}
	t.Fail()
}

func TestMooneyeAcceptance(t *testing.T) {
	for _, filename := range []string{
		"testdata/mts-20221022-1430-8d742b9/acceptance/add_sp_e_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/bits/mem_oam.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/bits/reg_f.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/bits/unused_hwio-GS.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/boot_div-dmgABCmgb.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/boot_regs-dmgABC.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/call_cc_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/call_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/div_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/halt_ime0_ei.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/halt_ime0_nointr_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/halt_ime1_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/if_ie_registers.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/instr/daa.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/intr_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/jp_cc_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/jp_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/oam_dma_restart.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/oam_dma_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/oam_dma/basic.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/oam_dma/reg_read.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/pop_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/ret_cc_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/ret_timing.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/reti_timing.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/boot_hwio-dmgABCmgb.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/call_cc_timing2.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/call_timing2.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/di_timing-GS.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ei_sequence.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ei_timing.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/halt_ime1_timing2-GS.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/interrupts/ie_push.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ld_hl_sp_e_timing.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/oam_dma_start.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/oam_dma/sources-GS.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/push_timing.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/rapid_di_ei.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/reti_intr_timing.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/rst_timing.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/serial/boot_sclk_align-dmgABCmgb.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}
}

func TestMooneyePPU(t *testing.T) {
	for _, filename := range []string{
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/hblank_ly_scx_timing-GS.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/intr_1_2_timing-GS.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/intr_2_0_timing.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/intr_2_mode0_timing_sprites.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/intr_2_mode0_timing.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/intr_2_mode3_timing.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/intr_2_oam_ok_timing.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/lcdon_timing-GS.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/lcdon_write_timing-GS.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/stat_irq_blocking.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/stat_lyc_onoff.gb",
		// "testdata/mts-20221022-1430-8d742b9/acceptance/ppu/vblank_stat_intr-GS.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}
}

func TestMooneyeTimer(t *testing.T) {
	for _, filename := range []string{
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/div_write.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/rapid_toggle.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tim00_div_trigger.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tim00.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tim01_div_trigger.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tim01.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tim10_div_trigger.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tim10.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tim11_div_trigger.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tim11.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tima_reload.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tima_write_reloading.gb",
		"testdata/mts-20221022-1430-8d742b9/acceptance/timer/tma_write_reloading.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}
}

func TestMooneyeMBC(t *testing.T) {
	for _, filename := range []string{
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/bits_bank1.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/bits_bank2.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/bits_mode.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/bits_ramg.gb",
		// "testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/multicart_rom_8Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/ram_256kb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/ram_64kb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/rom_16Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/rom_1Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/rom_2Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/rom_4Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/rom_512kb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc1/rom_8Mb.gb",
		// "testdata/mts-20221022-1430-8d742b9/emulator-only/mbc2/bits_ramg.gb",
		// "testdata/mts-20221022-1430-8d742b9/emulator-only/mbc2/bits_romb.gb",
		// "testdata/mts-20221022-1430-8d742b9/emulator-only/mbc2/bits_unused.gb",
		// "testdata/mts-20221022-1430-8d742b9/emulator-only/mbc2/ram.gb",
		// "testdata/mts-20221022-1430-8d742b9/emulator-only/mbc2/rom_1Mb.gb",
		// "testdata/mts-20221022-1430-8d742b9/emulator-only/mbc2/rom_2Mb.gb",
		// "testdata/mts-20221022-1430-8d742b9/emulator-only/mbc2/rom_512kb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc5/rom_16Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc5/rom_1Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc5/rom_2Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc5/rom_32Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc5/rom_4Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc5/rom_512kb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc5/rom_64Mb.gb",
		"testdata/mts-20221022-1430-8d742b9/emulator-only/mbc5/rom_8Mb.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runMooneyeTest(t, filename)
		})
	}
}
