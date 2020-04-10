package audio

type wave struct {
	length       uint8
	outputLevel  uint8
	frequency    uint16
	lengthEnable bool

	// Internal state
	enabled bool
	timer   uint16
}

func (w *wave) trigger(dacEnabled bool) {

	// Channel is enabled (see length counter).
	// If length counter is zero, it is set to 64 (256 for wave channel).
	// Frequency timer is reloaded with period.
	// Volume envelope timer is reloaded with period.
	// Channel volume is reloaded from NRx2.
	// Noise channel's LFSR bits are all set to 1.
	// Wave channel's position is set to 0 but sample buffer is NOT refilled.
	// Square 1's sweep does several things (see frequency sweep).

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.

	// Channel is enabled (see length counter).
	w.enabled = true

	// If length counter is zero, it is set to 64 (256 for wave channel).
	if w.length == 0 {
		w.length = 64
	}

	// Frequency timer is reloaded with period.
	w.timer = (2048 - w.frequency) * 4

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !dacEnabled {
		w.enabled = false
	}

}

func (w *wave) tickTimer() {
	if w.timer == 0 {
		w.timer = (2048 - w.frequency) * 4
	}
	w.timer--
}

func (w *wave) takeSample() float32 {

	if !w.enabled {
		return 0
	}

	wave := float32(0)

	return wave

}
