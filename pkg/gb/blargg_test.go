package gb

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func runBlarggTest(t *testing.T, filename string) {
	sbWriter := &bytes.Buffer{}
	opts := Options{
		RomFilename: filename,
		SBWriter:    sbWriter,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	gameboy := NewGameboy(opts)
	var result string
	go func() {
		// The main test ROMs write to serial but we need to directly check RAM for singles
		checkRAM := strings.Contains(filename, "rom_singles") ||
			strings.Contains(filename, "halt_bug") ||
			strings.Contains(filename, "mem_timing-2")
		for {
			if checkRAM {
				result = string(gameboy.memory.CartRAM()[0][:])
			} else {
				result = sbWriter.String()
			}
			select {
			case <-ctx.Done():
				return
			default:
				if strings.Contains(result, "Failed") || strings.Contains(result, "Passed") {
					cancel()
				}
			}
			// Check every 100ms for a result
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	gameboy.Run(ctx)
	<-ctx.Done()
	screenshotFilename := fmt.Sprintf(
		"testresults/%s.png",
		strings.Replace(
			strings.TrimPrefix(filename, "testdata/blargg/"),
			"/", "_", -1),
	)
	gameboy.lcd.Screenshot(screenshotFilename)
	if !strings.Contains(result, "Passed") {
		t.Errorf("--------\n%s--------\n%s--------\n", filename, result)
		// fmt.Printf("| :boom: fail | %s | [pic](pkg/gb/%s) |\n", filename, screenshotFilename)
	} else {
		// fmt.Printf("| :green_heart: pass | %s | [pic](pkg/gb/%s) |\n", filename, screenshotFilename)
	}
}

func TestBlarggCPUInstrs(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/cpu_instrs/cpu_instrs.gb")
}

func TestBlarggDMGSound(t *testing.T) {

	// runBlarggTest(t, "testdata/blargg/dmg_sound/dmg_sound.gb")

	for _, filename := range []string{
		// Test individual sound ROMs until they are all passing
		"testdata/blargg/dmg_sound/rom_singles/01-registers.gb",
		"testdata/blargg/dmg_sound/rom_singles/02-len ctr.gb",
		"testdata/blargg/dmg_sound/rom_singles/03-trigger.gb",
		"testdata/blargg/dmg_sound/rom_singles/04-sweep.gb",
		"testdata/blargg/dmg_sound/rom_singles/05-sweep details.gb",
		"testdata/blargg/dmg_sound/rom_singles/06-overflow on trigger.gb",
		"testdata/blargg/dmg_sound/rom_singles/07-len sweep period sync.gb",
		"testdata/blargg/dmg_sound/rom_singles/08-len ctr during power.gb",
		"testdata/blargg/dmg_sound/rom_singles/09-wave read while on.gb",
		"testdata/blargg/dmg_sound/rom_singles/10-wave trigger while on.gb",
		"testdata/blargg/dmg_sound/rom_singles/11-regs after power.gb",
		// "testdata/blargg/dmg_sound/rom_singles/12-wave write while on.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runBlarggTest(t, filename)
		})
	}

}

func TestBlarggHaltBug(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/halt_bug.gb")
}

func TestBlarggInstrTiming(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/instr_timing/instr_timing.gb")
}

func TestBlarggMemTiming(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/mem_timing/mem_timing.gb")
}

func TestBlarggMemTiming2(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/mem_timing-2/mem_timing.gb")
}

func TestBlarggOAMBug(t *testing.T) {

	// runBlarggTest(t, "testdata/blargg/oam_bug/oam_bug.gb")

	for _, filename := range []string{
		// Test individual OAM ROMs until they are all passing
		// "testdata/blargg/oam_bug/rom_singles/1-lcd_sync.gb",
		// "testdata/blargg/oam_bug/rom_singles/2-causes.gb",
		// "testdata/blargg/oam_bug/rom_singles/3-non_causes.gb",
		// "testdata/blargg/oam_bug/rom_singles/4-scanline_timing.gb",
		// "testdata/blargg/oam_bug/rom_singles/5-timing_bug.gb",
		// "testdata/blargg/oam_bug/rom_singles/6-timing_no_bug.gb",
		// "testdata/blargg/oam_bug/rom_singles/7-timing_effect.gb",
		// "testdata/blargg/oam_bug/rom_singles/8-instr_effect.gb",
	} {
		t.Run(filename, func(t *testing.T) {
			runBlarggTest(t, filename)
		})
	}

}
