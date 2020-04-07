package audio

type square struct {
	sweepTime        uint8
	sweepIncrease    bool
	sweepShift       uint8
	duty             uint8
	length           uint8
	initialVolume    uint8
	envelopeIncrease bool
	envelopeSweep    uint8
	frequency        uint16
	lengthEnable     bool

	// Internal state
	enabled       bool
	dutyIndex     uint8
	timer         uint16
	sweepTimer    uint8
	envelopeTimer uint8
	freqShadow    uint16
}

func (s *square) trigger() {

	// a.ch1.timer = (2048 - a.ch1.frequency) * 4

}

func (s *square) tickTimer() {
	s.timer--
	if s.timer == 0 {
		s.timer = (2048 - s.frequency) * 4
		s.dutyIndex++
		if s.dutyIndex >= 8 {
			s.dutyIndex = 0
		}
	}
}

func (s *square) tickLength() {
	if !s.lengthEnable {
		s.enabled = true
		return
	}
	if s.length == 0 {
		s.enabled = false
	} else {
		s.enabled = true
		s.length--
	}
}

func (s *square) tickVolumeEnvelope() {
	if s.envelopeSweep == 0 {
		return
	}
	s.envelopeTimer++
	if s.envelopeTimer >= s.envelopeSweep {
		if s.envelopeIncrease {
			if s.initialVolume < 15 {
				s.initialVolume++
				s.envelopeTimer = 0
			}
		} else {
			if s.initialVolume > 0 {
				s.initialVolume--
				s.envelopeTimer = 0
			}
		}
	}
}

func (s *square) tickSweep() {

}

func (s *square) takeSample() float32 {

	if !s.enabled {
		return 0
	}

	wave := waveduty[s.duty][s.dutyIndex]

	wave *= float32(s.initialVolume) / 8

	return wave

}

// Wave Duty:
//   00: 12.5% ( _-------_-------_------- )
//   01: 25%   ( __------__------__------ )
//   10: 50%   ( ____----____----____---- ) (normal)
//   11: 75%   ( ______--______--______-- )
var waveduty = [][]float32{
	[]float32{0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
	[]float32{0, 0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
	[]float32{0, 0, 0, 0, 1.0, 1.0, 1.0, 1.0},
	[]float32{0, 0, 0, 0, 0, 0, 1.0, 1.0},
}
