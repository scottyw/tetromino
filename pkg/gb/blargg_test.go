package gb

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func runBlarggTest(t *testing.T, filename string) {
	sbWriter := &bytes.Buffer{}
	opts := Options{
		RomFilename: "testdata/blargg/" + filename,
		SBWriter:    sbWriter,
		Fast:        true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	gameboy := NewGameboy(opts, cancel)
	go func() {
		for {
			result := sbWriter.String()
			ram := string(gameboy.memory.CartRAM()[0][:])
			select {
			case <-ctx.Done():
				return
			default:
				if strings.Contains(result, "Failed") ||
					strings.Contains(result, "Passed") ||
					strings.Contains(ram, "Failed") ||
					strings.Contains(ram, "Passed") {
					cancel()
				}
			}
			// Check every 100ms for a result
			time.Sleep(100 * time.Millisecond)
		}
	}()
	gameboy.Run(ctx, nil)
	<-ctx.Done()
	result := sbWriter.String()
	ram := string(gameboy.memory.CartRAM()[0][:])
	if !strings.Contains(result, "Passed") &&
		!strings.Contains(ram, "Passed") {
		t.Errorf(result)
	}
}

func TestBlarggCPUInstrs(t *testing.T) {
	runBlarggTest(t, "cpu_instrs/cpu_instrs.gb")
}

// func TestBlarggDMGSound(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/dmg_sound.gb")
// }

// func TestBlarggDMGSound01(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/01-registers.gb")
// }

// func TestBlarggDMGSound02(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/02-len ctr.gb")
// }

// func TestBlarggDMGSound03(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/03-trigger.gb")
// }

// func TestBlarggDMGSound04(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/04-sweep.gb")
// }

// func TestBlarggDMGSound05(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/05-sweep details.gb")
// }

// func TestBlarggDMGSound06(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/06-overflow on trigger.gb")
// }

// func TestBlarggDMGSound07(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/07-len sweep period sync.gb")
// }

// func TestBlarggDMGSound08(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/08-len ctr during power.gb")
// }

// func TestBlarggDMGSound09(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/09-wave read while on.gb")
// }

// func TestBlarggDMGSound10(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/10-wave trigger while on.gb")
// }

// func TestBlarggDMGSound11(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/11-regs after power.gb")
// }

// func TestBlarggDMGSound12(t *testing.T) {
// 	runBlarggTest(t, "dmg_sound/rom_singles/12-wave write while on.gb")
// }

func TestBlarggHaltBug(t *testing.T) {
	runBlarggTest(t, "halt_bug.gb")
}

func TestBlarggInstrTiming(t *testing.T) {
	runBlarggTest(t, "instr_timing/instr_timing.gb")
}

func TestBlarggMemTiming(t *testing.T) {
	runBlarggTest(t, "mem_timing/mem_timing.gb")
}

func TestBlarggMemTiming2(t *testing.T) {
	runBlarggTest(t, "mem_timing-2/mem_timing.gb")
}

// func TestBlarggOAMBug(t *testing.T) {
// 	runBlarggTest(t, "oam_bug/oam_bug.gb")
// }

// func TestBlarggOAMBug01(t *testing.T) {
// 	runBlarggTest(t, "oam_bug/rom_singles/1-lcd_sync.gb")
// }

// func TestBlarggOAMBug02(t *testing.T) {
// 	runBlarggTest(t, "oam_bug/rom_singles/2-causes.gb")
// }

// func TestBlarggOAMBug03(t *testing.T) {
// 	runBlarggTest(t, "oam_bug/rom_singles/3-non_causes.gb")
// }

// func TestBlarggOAMBug04(t *testing.T) {
// 	runBlarggTest(t, "oam_bug/rom_singles/4-scanline_timing.gb")
// }

// func TestBlarggOAMBug05(t *testing.T) {
// 	runBlarggTest(t, "oam_bug/rom_singles/5-timing_bug.gb")
// }

// func TestBlarggOAMBug06(t *testing.T) {
// 	runBlarggTest(t, "oam_bug/rom_singles/6-timing_no_bug.gb")
// }

// func TestBlarggOAMBug07(t *testing.T) {
// 	runBlarggTest(t, "oam_bug/rom_singles/7-timing_effect.gb")
// }

// func TestBlarggOAMBug08(t *testing.T) {
// 	runBlarggTest(t, "oam_bug/rom_singles/8-instr_effect.gb")
// }
