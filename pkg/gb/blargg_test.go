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
