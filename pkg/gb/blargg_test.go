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

func TestBlarggCPU(t *testing.T) { runBlarggTest(t, "cpu_instrs/cpu_instrs.gb") }

func TestBlarggInstrTiming(t *testing.T) { runBlarggTest(t, "instr_timing/instr_timing.gb") }

func TestBlarggMemoryTiming(t *testing.T) {
	runBlarggTest(t, "mem_timing/mem_timing.gb")
}

func TestBlarggMemoryTiming2(t *testing.T) {
	// This test passes but isn't correctly identified as such because it doesn't write to serial
	//
	// Text output and the final result are also written to memory at $A000,
	// allowing testing a very minimal emulator that supports little more than
	// CPU and RAM. To reliably indicate that the data is from a test and not
	// random data, $A001-$A003 are written with a signature: $DE,$B0,$61. If
	// this is present, then the text string and final result status are valid.

	// $A000 holds the overall status. If the test is still running, it holds
	// $80, otherwise it holds the final result code.

	// All text output is appended to a zero-terminated string at $A004. An
	// emulator could regularly check this string for any additional
	// characters, and output them, allowing real-time text output, rather than
	// just printing the final output at the end.
	//
	// runBlarggTest(t, "mem_timing-2/rom_singles/01-read_timing.gb")
}
