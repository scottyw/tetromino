package audio

type noise struct {
	length           uint8
	initialVolume    uint8
	envelopeIncrease bool
	envelopeSweep    uint8
	shift            uint8
	lfsrWidth        uint8
	divisor          uint8
	lengthEnable     bool

	// Internal state
	enabled       bool
	dacEnabled    bool
	volume        uint8
	timer         uint16
	envelopeTimer uint8
	lfsr          uint16
}

func (n *noise) trigger() {

	// Channel is enabled (see length counter).
	n.enabled = true

	// If length counter is zero, it is set to 64 (256 for wave channel).
	if n.length == 0 {
		n.length = 64
	}

	// Frequency timer is reloaded with period.
	n.timer = 8 << n.divisor

	// Volume envelope timer is reloaded with period.
	n.envelopeTimer = n.envelopeSweep

	// Channel volume is reloaded from NRx2.
	n.volume = n.initialVolume

	// Noise channel's LFSR bits are all set to 1.
	n.lfsr = 0xffff

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !n.dacEnabled {
		n.enabled = false
	}

}

func (n *noise) tickTimer() {
	if n.timer == 0 {
		n.timer = 8 << n.divisor
	}
	n.timer--
}

func (n *noise) tickLength() {
	if n.length == 0 {
		n.enabled = false
	} else {
		n.length--
	}
}

func (n *noise) tickVolumeEnvelope() {

}

func (n *noise) takeSample() float32 {

	if !n.enabled || !n.dacEnabled {
		return 0
	}

	wave := float32(0)

	return wave

}
