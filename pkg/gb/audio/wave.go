package audio

type wave struct {
	enable       bool
	length       uint8
	outputLevel  uint8
	frequency    uint16
	lengthEnable bool
}

func (w *wave) trigger() {

	// Channel is enabled (see length counter).
	// If length counter is zero, it is set to 64 (256 for wave channel).
	// Frequency timer is reloaded with period.
	// Volume envelope timer is reloaded with period.
	// Channel volume is reloaded from NRx2.
	// Noise channel's LFSR bits are all set to 1.
	// Wave channel's position is set to 0 but sample buffer is NOT refilled.
	// Square 1's sweep does several things (see frequency sweep).

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.

}
