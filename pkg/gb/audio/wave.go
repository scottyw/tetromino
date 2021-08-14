package audio

type wave struct {
	length       uint16
	outputLevel  uint8
	frequency    uint16
	lengthEnable bool
	waveram      [16]uint8

	// Internal state
	enabled      bool
	dacEnabled   bool
	timer        uint16
	outputShift  uint8
	position     uint8
	lastAccessed uint8
	sampleBuffer uint8
	sampleTimer  uint8
	triggered    bool
}

func (w *wave) trigger() {

	w.triggered = true

	// Channel is enabled (see length counter).
	w.enabled = true

	// If length counter is zero, it is set to 64 (256 for wave channel).
	if w.length == 0 {
		w.length = 256
	}

	// Frequency timer is reloaded with period.
	w.timer = (2048 - w.frequency) * 2

	// Channel volume is reloaded from NRx2.
	if w.outputLevel == 0 {
		w.outputShift = 4
	} else {
		w.outputShift = w.outputLevel - 1
	}

	// Wave channel's position is set to 0 but sample buffer is NOT refilled.
	w.position = 0

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !w.dacEnabled {
		w.enabled = false
	}

}

func (w *wave) tickTimer() {
	if !w.enabled {
		return
	}
	if w.timer == 0 {
		w.timer = (2048 - w.frequency) * 2
		w.position++
		if w.position >= 32 {
			w.position = 0
		}
		w.lastAccessed = w.position / 2
		if w.position%2 == 0 {
			w.sampleBuffer = w.waveram[w.lastAccessed] >> 4
		} else {
			w.sampleBuffer = w.waveram[w.lastAccessed] & 0x0f
		}
		w.sampleTimer = 0
	}
	w.timer--
	w.sampleTimer++
}

func (w *wave) tickLength() {
	if !w.lengthEnable {
		return
	}
	if w.length > 0 {
		w.length--
		if w.length == 0 {
			w.enabled = false
		}
	}
}

func (w *wave) takeSample() float32 {
	if !w.enabled {
		return 0
	}
	return float32(w.sampleBuffer>>w.outputShift) / 15
}
