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
			select {
			case <-ctx.Done():
				return
			default:
				if strings.Contains(result, "Failed") || strings.Contains(result, "Passed") {
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
	if !strings.Contains(result, "Passed") {
		t.Errorf(result)
	}
}

func TestBlargg00(t *testing.T) { runBlarggTest(t, "cpu_instrs.gb") }

// func TestBlargg12(t *testing.T) { runBlarggTest(t, "dmg_sound-2.gb") }

// func TestBlargg13(t *testing.T) { runBlarggTest(t, "dmg_sound-2/01-registers.gb") }

// func TestBlargg14(t *testing.T) { runBlarggTest(t, "dmg_sound-2/02-len ctr.gb") }

// func TestBlargg15(t *testing.T) { runBlarggTest(t, "dmg_sound-2/03-trigger.gb") }

// func TestBlargg16(t *testing.T) { runBlarggTest(t, "dmg_sound-2/04-sweep.gb") }

// func TestBlargg17(t *testing.T) { runBlarggTest(t, "dmg_sound-2/05-sweep details.gb") }

// func TestBlargg18(t *testing.T) { runBlarggTest(t, "dmg_sound-2/06-overflow on trigger.gb") }

// func TestBlargg19(t *testing.T) { runBlarggTest(t, "dmg_sound-2/07-len sweep period sync.gb") }

// func TestBlargg20(t *testing.T) { runBlarggTest(t, "dmg_sound-2/08-len ctr during power.gb") }

// func TestBlargg21(t *testing.T) { runBlarggTest(t, "dmg_sound-2/09-wave read while on.gb") }

// func TestBlargg22(t *testing.T) { runBlarggTest(t, "dmg_sound-2/10-wave trigger while on.gb") }

// func TestBlargg23(t *testing.T) { runBlarggTest(t, "dmg_sound-2/11-regs after power.gb") }

// func TestBlargg24(t *testing.T) { runBlarggTest(t, "dmg_sound-2/12-wave write while on.gb") }

// func TestBlargg25(t *testing.T) { runBlarggTest(t, "halt_bug.gb") }

func TestBlargg26(t *testing.T) { runBlarggTest(t, "instr_timing.gb") }

// func TestBlargg27(t *testing.T) { runBlarggTest(t, "interrupt_time.gb") }

// func TestBlargg28(t *testing.T) { runBlarggTest(t, "mem_timing-2.gb") }

// func TestBlargg29(t *testing.T) { runBlarggTest(t, "mem_timing-2/01-read_timing.gb") }

// func TestBlargg30(t *testing.T) { runBlarggTest(t, "mem_timing-2/02-write_timing.gb") }

// func TestBlargg31(t *testing.T) { runBlarggTest(t, "mem_timing-2/03-modify_timing.gb") }

// func TestBlargg32(t *testing.T) { runBlarggTest(t, "oam_bug-2.gb") }

// func TestBlargg33(t *testing.T) { runBlarggTest(t, "oam_bug-2/1-lcd_sync.gb") }

// func TestBlargg34(t *testing.T) { runBlarggTest(t, "oam_bug-2/2-causes.gb") }

// func TestBlargg35(t *testing.T) { runBlarggTest(t, "oam_bug-2/3-non_causes.gb") }

// func TestBlargg36(t *testing.T) { runBlarggTest(t, "oam_bug-2/4-scanline_timing.gb") }

// func TestBlargg37(t *testing.T) { runBlarggTest(t, "oam_bug-2/5-timing_bug.gb") }

// func TestBlargg38(t *testing.T) { runBlarggTest(t, "oam_bug-2/6-timing_no_bug.gb") }

// func TestBlargg39(t *testing.T) { runBlarggTest(t, "oam_bug-2/7-timing_effect.gb") }

// func TestBlargg40(t *testing.T) { runBlarggTest(t, "oam_bug-2/8-instr_effect.gb") }
