package audio

type wave struct {
	length             uint16
	initialOutputLevel uint8
	frequency          uint16
	lengthEnable       bool

	// Internal state
	enabled     bool
	dacEnabled  bool
	timer       uint16
	outputLevel uint8
	position    uint8
}

func (w *wave) trigger() {

	// Channel is enabled (see length counter).
	w.enabled = true

	// If length counter is zero, it is set to 64 (256 for wave channel).
	if w.length == 0 {
		w.length = 256
	}

	// Frequency timer is reloaded with period.
	w.timer = (2048 - w.frequency) * 2

	// Channel volume is reloaded from NRx2.
	w.outputLevel = w.initialOutputLevel

	// Wave channel's position is set to 0 but sample buffer is NOT refilled.
	w.position = 0

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !w.dacEnabled {
		w.enabled = false
	}

}

func (w *wave) tickTimer() {
	if w.timer == 0 {
		w.timer = (2048 - w.frequency) * 4
	}
	w.timer--
}

func (w *wave) tickLength() {
	if !w.lengthEnable {
		return
	}
	if w.length > 0 {
		w.length--
	}
	if w.length == 0 {
		w.enabled = false
	}
}

func (w *wave) takeSample() float32 {

	if !w.enabled {
		return 0
	}

	wave := float32(0)

	return wave

}
