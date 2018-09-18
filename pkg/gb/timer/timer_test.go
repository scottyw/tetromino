package timer

import (
	"runtime"
	"testing"
)

func ticks(timer *Timer, ticks int) {
	for i := 0; i < ticks; i++ {
		timer.Tick()
	}
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

func TestTick(t *testing.T) {
	timer := NewTimer()
	assertCounter(t, timer, 0xabcc)
	timer.Tick()
	assertCounter(t, timer, 0xabd0)
	timer.Tick()
	assertCounter(t, timer, 0xabd4)
}

func TestDIV(t *testing.T) {
	timer := NewTimer()
	assertDiv(t, timer, 0xab)
	ticks(timer, 12)
	assertDiv(t, timer, 0xab)
	ticks(timer, 1)
	assertDiv(t, timer, 0xac)
	ticks(timer, 63)
	assertDiv(t, timer, 0xac)
	ticks(timer, 1)
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
		ticks(timer, 2000)
		assertTima(t, timer, 0x0000)
	}
}

func testTIMA(t *testing.T, tac uint8, ticksPerIncrement int) {
	timer := NewTimer()
	timer.WriteTAC(tac)
	timer.Reset()
	assertTima(t, timer, 0x00)
	ticks(timer, ticksPerIncrement-1)
	assertTima(t, timer, 0x00)
	ticks(timer, 1)
	assertTima(t, timer, 0x01)
	ticks(timer, ticksPerIncrement-1)
	assertTima(t, timer, 0x01)
	ticks(timer, 1)
	assertTima(t, timer, 0x02)
	ticks(timer, ticksPerIncrement-1)
	assertTima(t, timer, 0x02)
	ticks(timer, 1)
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

func testTIMAOnReset(t *testing.T, tac uint8, ticksPerIncrement int) {
	timer := NewTimer()
	timer.WriteTAC(tac)
	// Reset after less than half the required ticks does not increment
	timer.Reset()
	assertTima(t, timer, 0x00)
	ticks(timer, ticksPerIncrement/2-1)
	timer.Reset()
	ticks(timer, 1)
	assertTima(t, timer, 0x00)
	// Reset after more than half the required ticks does increment
	timer.Reset()
	assertTima(t, timer, 0x00)
	ticks(timer, ticksPerIncrement/2+1)
	timer.Reset()
	ticks(timer, 1)
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
		ticks(timer, 1)
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
