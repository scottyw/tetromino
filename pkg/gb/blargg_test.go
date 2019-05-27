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

func TestBlarggCPU(t *testing.T) {
	runBlarggTest(t, "cpu_instrs/cpu_instrs.gb")
}

func TestBlarggInstrTiming(t *testing.T) {
	runBlarggTest(t, "instr_timing/instr_timing.gb")
}

func TestBlarggHaltBug(t *testing.T) {
	runBlarggTest(t, "halt_bug.gb")
}

func TestBlarggMemoryTiming(t *testing.T) {
	runBlarggTest(t, "mem_timing/mem_timing.gb")
}

func TestBlarggMemoryTiming2(t *testing.T) {
	runBlarggTest(t, "mem_timing-2/rom_singles/01-read_timing.gb")
}
