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
	tac         uint8
	tima        uint8
	tma         uint8
	lastEdgeSet bool
	timaWrite   bool
	tmaWrite    bool
	overflow    bool
	endCycleA   uint16
	endCycleB   uint16
}

// NewTimer creates an initialized timer
func New() *Timer {
	return &Timer{
		counter: 0xabcc,
	}
}

// EndMachineCycle updates the timer after a machine cycle
func (t *Timer) EndMachineCycle() bool {
	t.counter += 4
	if t.overflow && t.counter == t.endCycleA {
		t.endCycleA = 0xffff
		if !t.timaWrite {
			t.tima = t.tma
		}
	}
	t.timaWrite = false
	if t.overflow && t.counter == t.endCycleB {
		t.endCycleB = 0xffff
		t.overflow = false
		if t.tmaWrite {
			t.tima = t.tma
		}
	}
	t.tmaWrite = false
	// Check for a falling edge
	var interrupt bool
	enableBitSet := t.tac&0x04 > 0
	counterBitSet := t.counter&counterBitMasks[t.tac&0x03] > 0
	edgeSet := enableBitSet && counterBitSet
	if t.lastEdgeSet && !edgeSet {
		t.tima++
		// Check for overflow
		if t.tima == 0 {
			t.overflow = true
			t.endCycleA = t.counter + 4
			t.endCycleB = t.counter + 8
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
func (t *Timer) ReadDIV() uint8 {
	return uint8(t.counter >> 8)
}

// TAC returns the value of the TAC register
func (t *Timer) ReadTAC() uint8 {
	// First 5 bits are always high
	return t.tac | 0xf8
}

// TIMA returns the value of the TIMA register
func (t *Timer) ReadTIMA() uint8 {
	return t.tima
}

// TMA returns the value of the TMA register
func (t *Timer) ReadTMA() uint8 {
	return t.tma
}

// WriteDIV is called on writes to the DIV register
func (t *Timer) WriteDIV(value uint8) {
	t.Reset()
}

// WriteTAC is called on writes to the TAC register
func (t *Timer) WriteTAC(value uint8) {
	t.tac = value
}

// WriteTIMA is called on writes to the TIMA register
func (t *Timer) WriteTIMA(value uint8) {
	if t.counter != t.endCycleB-4 {
		t.tima = value
		t.timaWrite = true
	}
}

// WriteTMA is called on writes to the TMA register
func (t *Timer) WriteTMA(value uint8) {
	t.tma = value
	t.tmaWrite = true
	if t.counter == t.endCycleB {
		t.tima = value
	}
}
