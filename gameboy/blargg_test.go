package gameboy

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func runBlarggTest(t *testing.T, filename string, checkRAM bool) {
	serialWriter := &bytes.Buffer{}
	config := Config{
		RomFilename:        filename,
		DisableVideoOutput: true,
		DisableAudioOutput: true,
		SerialWriter:       serialWriter,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	gameboy := New(config)
	var result string
	go func() {
		for {
			// Some ROMs write results to serial but for some we need to directly check RAM
			if checkRAM {
				result = string(gameboy.memory.CartRAM()[0][:])
			} else {
				result = serialWriter.String()
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
		t.Errorf("\n--------\n%s\n--------\n%s\n--------\n", filename, result)
		// fmt.Printf("| :boom: fail | %s | [pic](pkg/gb/%s) |\n", filename, screenshotFilename)
	} else {
		// fmt.Printf("| :green_heart: pass | %s | [pic](pkg/gb/%s) |\n", filename, screenshotFilename)
	}
}

func TestBlarggCPUInstrs(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/cpu_instrs/cpu_instrs.gb", false)
}

func TestBlarggDMGSound(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/dmg_sound/dmg_sound.gb", true)
}

func TestBlarggHaltBug(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/halt_bug.gb", true)
}

func TestBlarggInstrTiming(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/instr_timing/instr_timing.gb", false)
}

func TestBlarggMemTiming(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/mem_timing/mem_timing.gb", false)
}

func TestBlarggMemTiming2(t *testing.T) {
	runBlarggTest(t, "testdata/blargg/mem_timing-2/mem_timing.gb", true)
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
			runBlarggTest(t, filename, false)
		})
	}

}
