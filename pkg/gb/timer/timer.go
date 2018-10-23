package timer

var counterBitMasks = []uint16{
	uint16(1) << 9,
	uint16(1) << 3,
	uint16(1) << 5,
	uint16(1) << 7,
}

// Timer stores the state of the internal timer
type Timer struct {
	counter     uint16
	overflow    uint16
	tac         uint8
	tima        uint8
	tma         uint8
	lastEdgeSet bool
	timaWrite   bool
}

// NewTimer creates an initialized timer
func NewTimer() *Timer {
	return &Timer{
		counter: 0xabcc,
	}
}

// MTick updates the timer after a machine cycle
func (t *Timer) MTick() bool {
	t.counter += 4
	var interrupt bool
	if !t.timaWrite && t.counter == t.overflow {
		t.tima = t.tma
	}
	t.timaWrite = false
	// Check for a falling edge
	enableBitSet := t.tac&0x04 > 0
	counterBitSet := t.counter&counterBitMasks[t.tac&0x03] > 0
	edgeSet := enableBitSet && counterBitSet
	if t.lastEdgeSet && !edgeSet {
		t.tima++
		// Check for overflow
		if t.tima == 0 {
			t.overflow = t.counter + 4
			interrupt = true
		}
	}
	t.lastEdgeSet = edgeSet
	return interrupt
}

// Reset the counter to zero, used when a value is written to DIV
func (t *Timer) Reset() {
	t.counter = 0
}

// DIV returns the value of the DIV register
func (t *Timer) DIV() uint8 {
	return uint8(t.counter >> 8)
}

// TAC returns the value of the TAC register
func (t *Timer) TAC() uint8 {
	return t.tac
}

// TIMA returns the value of the TIMA register
func (t *Timer) TIMA() uint8 {
	return t.tima
}

// TMA returns the value of the TMA register
func (t *Timer) TMA() uint8 {
	return t.tma
}

// WriteTAC returns the value of the TAC register
func (t *Timer) WriteTAC(value uint8) {
	t.tac = value
}

// WriteTIMA returns the value of the TIMA register
func (t *Timer) WriteTIMA(value uint8) {
	if t.counter != t.overflow {
		t.tima = value
		t.timaWrite = true
	}
}

// WriteTMA returns the value of the TMA register
func (t *Timer) WriteTMA(value uint8) {
	t.tma = value
	if t.counter == t.overflow {
		t.tima = value
	}
}
