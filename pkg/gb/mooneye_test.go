package gb

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func runMooneyeTest(t *testing.T, filename string) {
	opts := Options{
		RomFilename: "testdata/mooneye-gb_hwtests/" + filename,
		Fast:        true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	gameboy := NewGameboy(opts)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if gameboy.dispatch.Mooneye {
					cancel()
					return
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	gameboy.Run(ctx)
	<-ctx.Done()
	screenshotFilename := fmt.Sprintf("testresults/%s.png", strings.Replace(filename, "/", "_", -1))
	gameboy.lcd.Screenshot(screenshotFilename)
	if gameboy.dispatch.TestA() != 0 || !gameboy.dispatch.Mooneye {
		t.Errorf("Test ROM failed: %s", filename)
		// fmt.Printf("| :boom: fail | %s | [pic](pkg/gb/%s) |\n", filename, screenshotFilename)
	} else {
		// fmt.Printf("| :green_heart: pass | %s | [pic](pkg/gb/%s) |\n", filename, screenshotFilename)
	}

}

func TestMooneye00(t *testing.T) { runMooneyeTest(t, "acceptance/add_sp_e_timing.gb") }

func TestMooneye01(t *testing.T) { runMooneyeTest(t, "acceptance/bits/mem_oam.gb") }

func TestMooneye02(t *testing.T) { runMooneyeTest(t, "acceptance/bits/reg_f.gb") }

func TestMooneye03(t *testing.T) { runMooneyeTest(t, "acceptance/bits/unused_hwio-GS.gb") }

func TestMooneye06(t *testing.T) { runMooneyeTest(t, "acceptance/boot_div-dmgABCmgb.gb") }

func TestMooneye10(t *testing.T) { runMooneyeTest(t, "acceptance/boot_hwio-dmgABCmgb.gb") }

func TestMooneye12(t *testing.T) { runMooneyeTest(t, "acceptance/boot_regs-dmgABC.gb") }

func TestMooneye16(t *testing.T) { runMooneyeTest(t, "acceptance/call_cc_timing.gb") }

// func TestMooneye17(t *testing.T) { runMooneyeTest(t, "acceptance/call_cc_timing2.gb") }

func TestMooneye18(t *testing.T) { runMooneyeTest(t, "acceptance/call_timing.gb") }

// func TestMooneye19(t *testing.T) { runMooneyeTest(t, "acceptance/call_timing2.gb") }

// func TestMooneye20(t *testing.T) { runMooneyeTest(t, "acceptance/di_timing-GS.gb") }

func TestMooneye21(t *testing.T) { runMooneyeTest(t, "acceptance/div_timing.gb") }

// func TestMooneye22(t *testing.T) { runMooneyeTest(t, "acceptance/ei_sequence.gb") }

// func TestMooneye23(t *testing.T) { runMooneyeTest(t, "acceptance/ei_timing.gb") }

func TestMooneye24(t *testing.T) { runMooneyeTest(t, "acceptance/halt_ime0_ei.gb") }

func TestMooneye25(t *testing.T) { runMooneyeTest(t, "acceptance/halt_ime0_nointr_timing.gb") }

func TestMooneye26(t *testing.T) { runMooneyeTest(t, "acceptance/halt_ime1_timing.gb") }

// func TestMooneye27(t *testing.T) { runMooneyeTest(t, "acceptance/halt_ime1_timing2-GS.gb") }

// func TestMooneye28(t *testing.T) { runMooneyeTest(t, "acceptance/if_ie_registers.gb") }

func TestMooneye29(t *testing.T) { runMooneyeTest(t, "acceptance/instr/daa.gb") }

// func TestMooneye30(t *testing.T) { runMooneyeTest(t, "acceptance/interrupts/ie_push.gb") }

func TestMooneye31(t *testing.T) { runMooneyeTest(t, "acceptance/intr_timing.gb") }

func TestMooneye32(t *testing.T) { runMooneyeTest(t, "acceptance/jp_cc_timing.gb") }

func TestMooneye33(t *testing.T) { runMooneyeTest(t, "acceptance/jp_timing.gb") }

// func TestMooneye34(t *testing.T) { runMooneyeTest(t, "acceptance/ld_hl_sp_e_timing.gb") }

func TestMooneye35(t *testing.T) { runMooneyeTest(t, "acceptance/oam_dma/basic.gb") }

func TestMooneye36(t *testing.T) { runMooneyeTest(t, "acceptance/oam_dma/reg_read.gb") }

func TestMooneye37(t *testing.T) { runMooneyeTest(t, "acceptance/oam_dma/sources-dmgABCmgbS.gb") }

func TestMooneye38(t *testing.T) { runMooneyeTest(t, "acceptance/oam_dma_restart.gb") }

// func TestMooneye39(t *testing.T) { runMooneyeTest(t, "acceptance/oam_dma_start.gb") }

func TestMooneye40(t *testing.T) { runMooneyeTest(t, "acceptance/oam_dma_timing.gb") }

func TestMooneye41(t *testing.T) { runMooneyeTest(t, "acceptance/pop_timing.gb") }

// func TestMooneye42(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/hblank_ly_scx_timing-GS.gb") }

// func TestMooneye43(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/intr_1_2_timing-GS.gb") }

// func TestMooneye44(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/intr_2_0_timing.gb") }

// func TestMooneye45(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/intr_2_mode0_timing.gb") }

// func TestMooneye46(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/intr_2_mode0_timing_sprites.gb") }

// func TestMooneye47(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/intr_2_mode3_timing.gb") }

// func TestMooneye48(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/intr_2_oam_ok_timing.gb") }

// func TestMooneye49(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/lcdon_timing-dmgABCmgbS.gb") }

// func TestMooneye50(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/lcdon_write_timing-GS.gb") }

// func TestMooneye51(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/stat_irq_blocking.gb") }

// func TestMooneye52(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/stat_lyc_onoff.gb") }

// func TestMooneye53(t *testing.T) { runMooneyeTest(t, "acceptance/ppu/vblank_stat_intr-GS.gb") }

// func TestMooneye54(t *testing.T) { runMooneyeTest(t, "acceptance/push_timing.gb") }

// func TestMooneye55(t *testing.T) { runMooneyeTest(t, "acceptance/rapid_di_ei.gb") }

func TestMooneye56(t *testing.T) { runMooneyeTest(t, "acceptance/ret_cc_timing.gb") }

func TestMooneye57(t *testing.T) { runMooneyeTest(t, "acceptance/ret_timing.gb") }

// func TestMooneye58(t *testing.T) { runMooneyeTest(t, "acceptance/reti_intr_timing.gb") }

func TestMooneye59(t *testing.T) { runMooneyeTest(t, "acceptance/reti_timing.gb") }

// func TestMooneye60(t *testing.T) { runMooneyeTest(t, "acceptance/rst_timing.gb") }

// func TestMooneye61(t *testing.T) { runMooneyeTest(t, "acceptance/serial/boot_sclk_align-dmgABCmgb.gb") }

func TestMooneye62(t *testing.T) { runMooneyeTest(t, "acceptance/timer/div_write.gb") }

func TestMooneye63(t *testing.T) { runMooneyeTest(t, "acceptance/timer/rapid_toggle.gb") }

func TestMooneye64(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tim00.gb") }

func TestMooneye65(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tim00_div_trigger.gb") }

func TestMooneye66(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tim01.gb") }

func TestMooneye67(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tim01_div_trigger.gb") }

func TestMooneye68(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tim10.gb") }

func TestMooneye69(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tim10_div_trigger.gb") }

func TestMooneye70(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tim11.gb") }

func TestMooneye71(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tim11_div_trigger.gb") }

func TestMooneye72(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tima_reload.gb") }

func TestMooneye73(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tima_write_reloading.gb") }

func TestMooneye74(t *testing.T) { runMooneyeTest(t, "acceptance/timer/tma_write_reloading.gb") }

func TestMooneye75(t *testing.T) { runMooneyeTest(t, "emulator-only/mbc1/bits_ram_en.gb") }

// func TestMooneye76(t *testing.T) { runMooneyeTest(t, "emulator-only/mbc1/multicart_rom_8Mb.gb") }

func TestMooneye77(t *testing.T) { runMooneyeTest(t, "emulator-only/mbc1/ram_256Kb.gb") }

func TestMooneye78(t *testing.T) { runMooneyeTest(t, "emulator-only/mbc1/ram_64Kb.gb") }

func TestMooneye79(t *testing.T) { runMooneyeTest(t, "emulator-only/mbc1/rom_16Mb.gb") }

func TestMooneye80(t *testing.T) { runMooneyeTest(t, "emulator-only/mbc1/rom_1Mb.gb") }

func TestMooneye81(t *testing.T) { runMooneyeTest(t, "emulator-only/mbc1/rom_2Mb.gb") }

func TestMooneye82(t *testing.T) { runMooneyeTest(t, "emulator-only/mbc1/rom_4Mb.gb") }

func TestMooneye83(t *testing.T) { runMooneyeTest(t, "emulator-only/mbc1/rom_512Kb.gb") }

func TestMooneye84(t *testing.T) { runMooneyeTest(t, "emulator-only/mbc1/rom_8Mb.gb") }
