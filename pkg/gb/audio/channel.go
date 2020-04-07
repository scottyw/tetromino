package audio

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
