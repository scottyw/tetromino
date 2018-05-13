package lcd

import "testing"

func assert(t *testing.T, expected, actual uint16) {
	if expected != actual {
		t.Errorf("Expected %04x, actual %04x", expected, actual)
	}
}
func TestLowTileAbsoluteAddress(t *testing.T) {
	assert(t, 0x8000, lowTileAbsoluteAddress(0))
	assert(t, 0x82f0, lowTileAbsoluteAddress(47))
	assert(t, 0x8800, lowTileAbsoluteAddress(128))
	assert(t, 0x8ff0, lowTileAbsoluteAddress(255))
}

func TestHighTileAbsoluteAddress(t *testing.T) {
	assert(t, 0x9000, highTileAbsoluteAddress(0))
	assert(t, 0x97f0, highTileAbsoluteAddress(127))
	assert(t, 0x8800, highTileAbsoluteAddress(-128))
	assert(t, 0x8d10, highTileAbsoluteAddress(-47))
	assert(t, 0x92f0, highTileAbsoluteAddress(47))
}
