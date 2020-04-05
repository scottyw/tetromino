package audio

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

func (a *Audio) takeSample() {

	if !a.control.on || a.l == nil || a.r == nil {
		return
	}

	var wave1, wave2 float32

	// channel 1
	if a.control.ch1Enable {
		wave1 = a.ch1.takeSample()
	}

	// channel 2
	if a.control.ch2Enable {
		wave1 = a.ch2.takeSample()
	}

	// channel 3
	if a.control.ch3Enable {
		panic("ch3")
	}

	// channel 4
	if a.control.ch4Enable {
		panic("ch4")
	}

	// Mix channels
	masterVolume := float32(0.6)
	left := float32(0)
	if a.control.ch1Left {
		left += wave1
	}
	if a.control.ch2Left {
		left += wave2
	}
	left /= 2
	left *= float32(a.control.volumeLeft) / 8 * masterVolume
	a.l <- left

	right := float32(0)
	if a.control.ch1Right {
		right += wave1
	}
	if a.control.ch2Right {
		right += wave2
	}
	right /= 2
	right *= float32(a.control.volumeRight) / 8 * masterVolume
	a.r <- right

}

func (s *square) tickTimer() {
	s.timer--
	if s.timer == 0 {
		s.timerOut = waveduty[s.duty][s.dutyIndex]
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

}

func (s *square) tickSweep() {

}

func (s *square) takeSample() float32 {

	if !s.enabled {
		return 0
	}

	wave1 := s.timerOut * float32(s.initialVolume) / 8

	return wave1

}
