package audio

import "fmt"

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
	ticks        int
}

func (w *wave) trigger() {

	if w.enabled && w.sampleTimer < 4 {
		fmt.Printf("BEFORE : %v - %v - %v - %02x - %+v\n", w.enabled, w.sampleTimer, w.ticks, w.lastAccessed, w.waveram)
		switch {

		case w.lastAccessed < 4:

			w.waveram[0] = w.waveram[w.lastAccessed]

		case w.lastAccessed < 8:

			w.waveram[0] = w.waveram[4]
			w.waveram[1] = w.waveram[5]
			w.waveram[2] = w.waveram[6]
			w.waveram[3] = w.waveram[7]

		case w.lastAccessed < 12:

			w.waveram[0] = w.waveram[8]
			w.waveram[1] = w.waveram[9]
			w.waveram[2] = w.waveram[10]
			w.waveram[3] = w.waveram[11]

		default:

			w.waveram[0] = w.waveram[12]
			w.waveram[1] = w.waveram[13]
			w.waveram[2] = w.waveram[14]
			w.waveram[3] = w.waveram[15]

		}
		fmt.Printf("AFTER : %v - %v - %v - %02x - %+v\n", w.enabled, w.sampleTimer, w.ticks, w.lastAccessed, w.waveram)

	} else {
		fmt.Printf("TRIGGER : %v - %v- %v - %02x\n", w.enabled, w.sampleTimer, w.ticks, w.lastAccessed)
	}

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
		w.ticks = 0
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
	w.ticks++
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
