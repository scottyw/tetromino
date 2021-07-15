package audio

type sweep struct {
	sweepPeriod   uint8
	sweepIncrease bool
	sweepShift    uint8

	// Internal state
	sweepEnabled    bool
	sweepTimer      uint8
	shadowFrequency uint16
}

type square struct {
	duty             uint8
	length           uint8
	initialVolume    uint8
	envelopeIncrease bool
	envelopeSweep    uint8
	frequency        uint16
	lengthEnable     bool
	*sweep

	// Internal state
	enabled       bool
	dacEnabled    bool
	dutyIndex     uint8
	volume        uint8
	timer         uint16
	envelopeTimer uint8
}

func (s *square) trigger() {

	// Channel is enabled (see length counter).
	s.enabled = true

	// If length counter is zero, it is set to 64 (256 for wave channel).
	if s.length == 0 {
		s.length = 64
	}

	// Frequency timer is reloaded with period.
	s.timer = (2048 - s.frequency) * 4

	// Volume envelope timer is reloaded with period.
	s.envelopeTimer = s.envelopeSweep
	if s.envelopeTimer == 0 {
		s.envelopeTimer = 8
	}

	// Channel volume is reloaded from NRx2.
	s.volume = s.initialVolume

	if s.sweep != nil {

		// Square 1's frequency is copied to the shadow register.
		s.shadowFrequency = s.frequency

		// The sweep timer is reloaded.
		s.sweepTimer = s.sweepPeriod
		if s.sweepTimer == 0 {
			s.sweepTimer = 8
		}

		// The internal enabled flag is set if either the sweep period or shift are non-zero, cleared otherwise.
		s.sweepEnabled = s.sweepPeriod > 0 || s.sweepShift > 0

		// If the sweep shift is non-zero, frequency calculation and the overflow check are performed immediately.
		if s.sweepShift > 0 {
			s.calculateFrequency()
		}

	}

	// Note that if the channel's DAC is off, after the above actions occur the channel will be immediately disabled again.
	if !s.dacEnabled {
		s.enabled = false
	}

}

func (s *square) tickTimer() {
	if s.timer == 0 {
		s.timer = (2048 - s.frequency) * 4
		s.dutyIndex++
		if s.dutyIndex >= 8 {
			s.dutyIndex = 0
		}
	}
	s.timer--
}

func (s *square) tickLength() {
	if !s.lengthEnable {
		return
	}
	if s.length > 0 {
		s.length--
		if s.length == 0 {
			s.enabled = false
		}
	}
}

func (s *square) tickVolumeEnvelope() {
	if s.envelopeSweep == 0 {
		return
	}
	if s.envelopeTimer == 0 {
		if s.envelopeIncrease {
			if s.volume < 15 {
				s.volume++
				s.envelopeTimer = s.envelopeSweep
			}
		} else {
			if s.volume > 0 {
				s.volume--
				s.envelopeTimer = s.envelopeSweep
			}
		}
	}
	s.envelopeTimer--
}

func (s *square) tickSweep() {
	if !s.sweepEnabled {
		return
	}
	s.sweepTimer--
	if s.sweepTimer == 0 {
		s.sweepTimer = s.sweepPeriod
		if s.sweepTimer == 0 {
			s.sweepTimer = 8
		} else {
			newFrequency := s.calculateFrequency()
			if newFrequency < 2048 && s.sweepShift > 0 {
				s.frequency = newFrequency
				s.shadowFrequency = newFrequency
				s.calculateFrequency()
			}
		}
	}
}

func (s *square) calculateFrequency() uint16 {
	newFrequency := s.shadowFrequency
	newFrequency >>= s.sweepShift
	if !s.sweepIncrease {
		newFrequency = -newFrequency
	}
	newFrequency = s.shadowFrequency + newFrequency
	if newFrequency > 2047 {
		s.enabled = false
	}
	return newFrequency
}

func (s *square) takeSample() float32 {

	if !s.enabled || !s.dacEnabled {
		return 0
	}

	wave := waveduty[s.duty][s.dutyIndex]

	wave *= float32(s.volume) / 8

	return wave

}

// Wave Duty:
//   00: 12.5% ( _-------_-------_------- )
//   01: 25%   ( __------__------__------ )
//   10: 50%   ( ____----____----____---- ) (normal)
//   11: 75%   ( ______--______--______-- )
var waveduty = [][]float32{
	{0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
	{0, 0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
	{0, 0, 0, 0, 1.0, 1.0, 1.0, 1.0},
	{0, 0, 0, 0, 0, 0, 1.0, 1.0},
}
