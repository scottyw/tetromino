package audio

func (a *Audio) takeSample() {

	if !a.control.on || a.l == nil || a.r == nil {
		return
	}

	var wave1, wave2, wave3, wave4 float32

	// channel 1
	if a.control.ch1Enable {
		wave1 = a.ch1.takeSample()
	}

	// channel 2
	if a.control.ch2Enable {
		wave2 = a.ch2.takeSample()
	}

	// channel 3
	if a.control.ch3Enable {
		wave3 = a.ch3.takeSample()
	}

	// channel 4
	if a.control.ch4Enable {
		wave4 = a.ch4.takeSample()
	}

	// Hardcode master volume for now
	masterVolume := float32(0.6)

	// Mix left channel
	left := float32(0)
	if a.control.ch1Left {
		left += wave1
	}
	if a.control.ch2Left {
		left += wave2
	}
	if a.control.ch3Left {
		left += wave3
	}
	if a.control.ch4Left {
		left += wave4
	}
	left /= 4
	left *= float32(a.control.volumeLeft) / 8 * masterVolume
	a.l <- left

	// Mix right channel
	right := float32(0)
	if a.control.ch1Right {
		right += wave1
	}
	if a.control.ch2Right {
		right += wave2
	}
	if a.control.ch3Right {
		right += wave3
	}
	if a.control.ch4Right {
		right += wave4
	}
	right /= 4
	right *= float32(a.control.volumeRight) / 8 * masterVolume
	a.r <- right

}
