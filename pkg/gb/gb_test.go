package gb

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	blarggDir = os.Getenv("GOPATH") + "/src/github.com/scottyw/tetromino/test/blargg/"
	timeout   = 10 * time.Second
)

func runBlarggTest(t *testing.T, filename string) {
	sbWriter := &bytes.Buffer{}
	opts := Options{
		RomFilename: blarggDir + filename,
		SBWriter:    sbWriter,
		// DebugCPU:    true,
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	gameboy := NewGameboy(opts)
	go func() {
		for {
			result := sbWriter.String()
			select {
			case <-ctx.Done():
				return
			default:
				if strings.Contains(result, "Failed") || strings.Contains(result, "Passed") {
					cancelFunc()
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

func TestBlarggCpuInstrs01(t *testing.T) {
	runBlarggTest(t, "cpu_instrs/01-special.gb")
}

func TestBlarggCpuInstrs02(t *testing.T) {
	runBlarggTest(t, "cpu_instrs/02-interrupts.gb")
}

// func TestBlarggCpuInstrs03(t *testing.T) {
// 	runBlarggTest(t, "cpu_instrs/03-op sp,hl.gb")
// }

func TestBlarggCpuInstrs04(t *testing.T) {
	runBlarggTest(t, "cpu_instrs/04-op r,imm.gb")
}

func TestBlarggCpuInstrs05(t *testing.T) {
	runBlarggTest(t, "cpu_instrs/05-op rp.gb")
}

func TestBlarggCpuInstrs06(t *testing.T) {
	runBlarggTest(t, "cpu_instrs/06-ld r,r.gb")
}

// func TestBlarggCpuInstrs07(t *testing.T) {
// 	runBlarggTest(t, "cpu_instrs/07-jr,jp,call,ret,rst.gb")
// }

// func TestBlarggCpuInstrs08(t *testing.T) {
// 	runBlarggTest(t, "cpu_instrs/08-misc instrs.gb")
// }

func TestBlarggCpuInstrs09(t *testing.T) {
	runBlarggTest(t, "cpu_instrs/09-op r,r.gb")
}

func TestBlarggCpuInstrs10(t *testing.T) {
	runBlarggTest(t, "cpu_instrs/10-bit ops.gb")
}

func TestBlarggCpuInstrs11(t *testing.T) {
	runBlarggTest(t, "cpu_instrs/11-op a,(hl).gb")
}
