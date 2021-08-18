package timer

import (
	"runtime"
	"testing"
)

func mticks(timer *Timer, mticks int) bool {
	var interrupt bool
	for i := 0; i < mticks; i++ {
		interrupt = interrupt || timer.EndMachineCycle()
	}
	return interrupt
}

func assertCounter(t *testing.T, timer *Timer, counter uint16) {
	if timer.counter != counter {
		_, line, file, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: Wrong counter: 0x%02x", line, file, timer.counter)
	}
}

func assertDiv(t *testing.T, timer *Timer, div uint8) {
	if timer.DIV() != div {
		_, line, file, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: Wrong DIV: 0x%02x", line, file, timer.DIV())
	}
}

func assertTac(t *testing.T, timer *Timer, tac uint8) {
	if timer.TAC() != tac {
		_, line, file, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: Wrong TAC: 0x%02x", line, file, timer.TAC())
	}
}

func assertTima(t *testing.T, timer *Timer, tima uint8) {
	if timer.TIMA() != tima {
		_, line, file, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: Wrong TIMA: 0x%02x", line, file, timer.TIMA())
	}
}

func TestEndMachineCycle(t *testing.T) {
	timer := NewTimer()
	assertCounter(t, timer, 0xabcc)
	timer.EndMachineCycle()
	assertCounter(t, timer, 0xabd0)
	timer.EndMachineCycle()
	assertCounter(t, timer, 0xabd4)
}

func TestDIV(t *testing.T) {
	timer := NewTimer()
	assertDiv(t, timer, 0xab)
	mticks(timer, 12)
	assertDiv(t, timer, 0xab)
	mticks(timer, 1)
	assertDiv(t, timer, 0xac)
	mticks(timer, 63)
	assertDiv(t, timer, 0xac)
	mticks(timer, 1)
	assertDiv(t, timer, 0xad)
}

func TestReset(t *testing.T) {
	timer := NewTimer()
	assertCounter(t, timer, 0xabcc)
	timer.Reset()
	assertCounter(t, timer, 0x0000)
}

func TestDisabledTAC(t *testing.T) {
	timer := NewTimer()
	for _, tac := range []uint8{0x00, 0x01, 0x02, 0x03} {
		timer.WriteTAC(tac)
		assertTac(t, timer, tac)
		assertTima(t, timer, 0x0000)
		mticks(timer, 2000)
		assertTima(t, timer, 0x0000)
	}
}

func testTIMA(t *testing.T, tac uint8, mticksPerIncrement int) {
	timer := NewTimer()
	timer.WriteTAC(tac)
	timer.Reset()
	assertTima(t, timer, 0x00)
	mticks(timer, mticksPerIncrement-1)
	assertTima(t, timer, 0x00)
	mticks(timer, 1)
	assertTima(t, timer, 0x01)
	mticks(timer, mticksPerIncrement-1)
	assertTima(t, timer, 0x01)
	mticks(timer, 1)
	assertTima(t, timer, 0x02)
	mticks(timer, mticksPerIncrement-1)
	assertTima(t, timer, 0x02)
	mticks(timer, 1)
	assertTima(t, timer, 0x03)
}

func TestTIMA1024(t *testing.T) {
	testTIMA(t, 0x04, 1024/4)
}

func TestTIMA16(t *testing.T) {
	testTIMA(t, 0x05, 16/4)
}

func TestTIMA64(t *testing.T) {
	testTIMA(t, 0x06, 64/4)
}

func TestTIMA256(t *testing.T) {
	testTIMA(t, 0x07, 256/4)
}

func testTIMAOnReset(t *testing.T, tac uint8, mticksPerIncrement int) {
	timer := NewTimer()
	timer.WriteTAC(tac)
	// Reset after less than half the required mticks does not increment
	timer.Reset()
	assertTima(t, timer, 0x00)
	mticks(timer, mticksPerIncrement/2-1)
	timer.Reset()
	mticks(timer, 1)
	assertTima(t, timer, 0x00)
	// Reset after more than half the required mticks does increment
	timer.Reset()
	assertTima(t, timer, 0x00)
	mticks(timer, mticksPerIncrement/2+1)
	timer.Reset()
	mticks(timer, 1)
	assertTima(t, timer, 0x01)
}

func TestTIMAOnReset1024(t *testing.T) {
	testTIMAOnReset(t, 0x04, 1024/4)
}

func TestTIMAOnReset16(t *testing.T) {
	testTIMAOnReset(t, 0x05, 16/4)
}

func TestTIMAOnReset64(t *testing.T) {
	testTIMAOnReset(t, 0x06, 64/4)
}

func TestTIMAOnReset256(t *testing.T) {
	testTIMAOnReset(t, 0x07, 256/4)
}

func testTIMAOnFrequentReset(t *testing.T, tac uint8) {
	timer := NewTimer()
	timer.WriteTAC(tac)
	timer.Reset()
	for i := 0; i < 2000; i++ {
		assertTima(t, timer, 0x00)
		timer.Reset()
		mticks(timer, 1)
		assertTima(t, timer, 0x00)
	}
}

func TestTIMAOnFrequentReset1024(t *testing.T) {
	testTIMAOnFrequentReset(t, 0x04)
}

func TestTIMAOnFrequentReset16(t *testing.T) {
	testTIMAOnFrequentReset(t, 0x05)
}

func TestTIMAOnFrequentReset64(t *testing.T) {
	testTIMAOnFrequentReset(t, 0x06)
}

func TestTIMAOnFrequentReset256(t *testing.T) {
	testTIMAOnFrequentReset(t, 0x07)
}

func testTIMAOnDisable(t *testing.T, tac uint8, mticksPerIncrement int) {
	// Disable after less than half the required mticks does not increment
	timer := NewTimer()
	timer.WriteTAC(tac)
	timer.Reset()
	mticks(timer, mticksPerIncrement*3)
	assertTima(t, timer, 0x03)
	mticks(timer, mticksPerIncrement/2-1)
	assertTima(t, timer, 0x03)
	timer.WriteTAC(tac & 0x03)
	mticks(timer, 1)
	assertTima(t, timer, 0x03)
	// Disable after more than half the required mticks does increment
	timer = NewTimer()
	timer.WriteTAC(tac)
	timer.Reset()
	mticks(timer, mticksPerIncrement*3)
	assertTima(t, timer, 0x03)
	mticks(timer, mticksPerIncrement/2+1)
	assertTima(t, timer, 0x03)
	timer.WriteTAC(tac & 0x03)
	mticks(timer, 1)
	assertTima(t, timer, 0x04)
}

func TestTIMAOnDisable1024(t *testing.T) {
	testTIMAOnDisable(t, 0x04, 1024/4)
}

func TestTIMAOnDisable16(t *testing.T) {
	testTIMAOnDisable(t, 0x05, 16/4)
}

func TestTIMAOnDisable64(t *testing.T) {
	testTIMAOnDisable(t, 0x06, 64/4)
}

func TestTIMAOnDisable256(t *testing.T) {
	testTIMAOnDisable(t, 0x07, 256/4)
}

func TestTIMAOnTACChange(t *testing.T) {
	// Setup
	timer := NewTimer()
	timer.WriteTAC(0x06)
	timer.Reset()
	mticks(timer, 0x1110/4-8)
	assertTima(t, timer, 0x43)
	assertCounter(t, timer, 0x10f0)
	// Now change TAC and expect a TIMA increment
	timer.WriteTAC(0x05)
	mticks(timer, 1)
	assertTima(t, timer, 0x44)
	assertCounter(t, timer, 0x10f4)
}

func TestTIMAReload(t *testing.T) {
	// Setup
	timer := NewTimer()
	timer.WriteTAC(0x05)
	timer.counter = 0
	timer.tima = 0xff
	timer.tma = 0x23
	// Execute to one tick before rollover
	interrupt := mticks(timer, 3)
	assertTima(t, timer, 0xff)
	if interrupt {
		t.Errorf("Timer interrupt should not have occurred")
	}
	// One more tick should rollover but not set TIMA correctly
	interrupt = timer.EndMachineCycle()
	// Cycle A starts
	assertTima(t, timer, 0x00)
	if !interrupt {
		t.Errorf("Timer interrupt should have occurred")
	}
	// Cycle A ends and another tick sets TIMA correctly
	interrupt = timer.EndMachineCycle()
	// Cycle B starts
	assertTima(t, timer, 0x23)
	if interrupt {
		t.Errorf("Timer interrupt should not have occurred")
	}
}

func TestTIMAReloadWithWriteA(t *testing.T) {
	// Setup
	timer := NewTimer()
	timer.WriteTAC(0x05)
	timer.counter = 0
	timer.tima = 0xff
	timer.tma = 0x23
	// Execute to one tick before rollover
	interrupt := mticks(timer, 3)
	assertTima(t, timer, 0xff)
	if interrupt {
		t.Errorf("Timer interrupt should not have occurred")
	}
	// One more tick should rollover but not set TIMA correctly
	interrupt = timer.EndMachineCycle()
	// Cycle A starts
	assertTima(t, timer, 0x00)
	if !interrupt {
		t.Errorf("Timer interrupt should have occurred")
	}
	// Write TIMA during cycle A
	timer.WriteTIMA(0x57)
	// Cycle A ends and another tick after the write should retain the written value
	interrupt = timer.EndMachineCycle()
	assertTima(t, timer, 0x57)
	if interrupt {
		t.Errorf("Timer interrupt should not have occurred")
	}
}

func TestTIMAReloadWithWriteB(t *testing.T) {
	// Setup
	timer := NewTimer()
	timer.WriteTAC(0x05)
	timer.counter = 0
	timer.tima = 0xff
	timer.tma = 0x23
	// Execute to one tick before rollover
	interrupt := mticks(timer, 3)
	assertTima(t, timer, 0xff)
	if interrupt {
		t.Errorf("Timer interrupt should not have occurred")
	}
	// One more tick should rollover but not set TIMA correctly
	interrupt = timer.EndMachineCycle()
	// Cycle A starts
	assertTima(t, timer, 0x00)
	if !interrupt {
		t.Errorf("Timer interrupt should have occurred")
	}
	// Cycle A ends and another tick sets TIMA correctly
	interrupt = timer.EndMachineCycle()
	// Cycle B starts
	assertTima(t, timer, 0x23)
	if interrupt {
		t.Errorf("Timer interrupt should not have occurred")
	}
	// Write TIMA during cycle B
	timer.WriteTIMA(0x57)
	// Cycle B ends and the write has been ignored
	interrupt = timer.EndMachineCycle()
	assertTima(t, timer, 0x23)
	if interrupt {
		t.Errorf("Timer interrupt should not have occurred")
	}
}

func TestTIMAReloadWithTMAWrite(t *testing.T) {
	// Setup
	timer := NewTimer()
	timer.WriteTAC(0x05)
	timer.counter = 0
	timer.tima = 0xff
	timer.tma = 0x23
	// Execute to one tick before rollover
	interrupt := mticks(timer, 3)
	assertTima(t, timer, 0xff)
	if interrupt {
		t.Errorf("Timer interrupt should not have occurred")
	}
	// One more tick should rollover but not set TIMA correctly
	interrupt = timer.EndMachineCycle()
	// Cycle A starts
	assertTima(t, timer, 0x00)
	if !interrupt {
		t.Errorf("Timer interrupt should have occurred")
	}
	// Cycle A ends and another tick sets TIMA correctly
	interrupt = timer.EndMachineCycle()
	// Cycle B starts
	assertTima(t, timer, 0x23)
	if interrupt {
		t.Errorf("Timer interrupt should not have occurred")
	}
	// Write TMA during cycle B
	timer.WriteTMA(0x57)
	// Cycle B ends and TIMA has been updated
	interrupt = timer.EndMachineCycle()
	assertTima(t, timer, 0x57)
	if interrupt {
		t.Errorf("Timer interrupt should not have occurred")
	}
}
