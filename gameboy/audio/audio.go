package audio

const (
	frameSeqPeriod = 4194304 / 512   // 512Hz
	samplerPeriod  = 4194304 / 44100 // 44100 Hz
)

// Audio stream
type Audio struct {
	l             chan float32
	r             chan float32
	ch1           *square
	ch2           *square
	ch3           *wave
	ch4           *noise
	control       *control
	ticks         uint64
	frameSeqTicks uint64
}

// NewAudio initializes our internal channel for audio data
func New(l, r chan float32) *Audio {
	audio := Audio{
		l:   l,
		r:   r,
		ch1: &square{sweep: &sweep{}},
		ch2: &square{},
		ch3: &wave{
			waveram: [16]uint8{0x84, 0x40, 0x43, 0xAA, 0x2D, 0x78, 0x92, 0x3C, 0x60, 0x59, 0x59, 0xB0, 0x34, 0xB8, 0x2E, 0xDA},
		},
		ch4:     &noise{},
		control: &control{},
		ticks:   1,
	}

	// Set default values for the NR registers
	audio.WriteNR10(0x80)
	audio.WriteNR11(0xbf)
	audio.WriteNR12(0xf3)
	audio.WriteNR13(0xff)
	audio.WriteNR14(0xbf)
	audio.WriteNR21(0x3f)
	audio.WriteNR23(0xff)
	audio.WriteNR24(0xbf)
	audio.WriteNR30(0x7f)
	audio.WriteNR31(0xff)
	audio.WriteNR32(0x9f)
	audio.WriteNR33(0xff)
	audio.WriteNR34(0xbf)
	audio.WriteNR41(0xff)
	audio.WriteNR44(0xbf)
	audio.WriteNR50(0x77)
	audio.WriteNR51(0xf3)
	audio.WriteNR52(0xf1)

	return &audio
}

// EndMachineCycle emulates the audio hardware at the end of a machine cycle
func (a *Audio) EndMachineCycle() {
	// Each machine cycle is four clock cycles
	a.tickClock()
	a.tickClock()
	a.tickClock()
	a.tickClock()
	a.ch1.triggered = false
	a.ch2.triggered = false
	a.ch3.triggered = false
	a.ch4.triggered = false
}

func (a *Audio) tickClock() {
	if a.ticks > 4194304 {
		a.ticks = 1
		a.frameSeqTicks = 0
	}

	// Tick every clock cycle
	a.tickTimer()

	// Tick the frame sequencer at 512 Hz
	if a.ticks%frameSeqPeriod == 0 {
		a.tickFrameSequencer()
		if a.frameSeqTicks >= 512 {
			a.frameSeqTicks = 0
		}
	}

	// Tick this function at 44100 Hz
	if a.ticks%samplerPeriod == 0 {
		a.tickSampler()
	}

	a.ticks++

}

func (a *Audio) tickTimer() {
	if !a.ch1.triggered {
		a.ch1.tickTimer()
	}
	if !a.ch2.triggered {
		a.ch2.tickTimer()
	}
	if !a.ch3.triggered {
		a.ch3.tickTimer()
	}
	if !a.ch4.triggered {
		a.ch4.tickTimer()
	}
}

func (a *Audio) tickFrameSequencer() {

	// Step   Length Ctr  Vol Env     Sweep
	// ---------------------------------------
	// 0      Clock       -           -
	// 1      -           -           -
	// 2      Clock       -           Clock
	// 3      -           -           -
	// 4      Clock       -           -
	// 5      -           -           -
	// 6      Clock       -           Clock
	// 7      -           Clock       -
	// ---------------------------------------
	// Rate   256 Hz      64 Hz       128 Hz

	if a.frameSeqTicks%2 == 0 {
		a.ch1.tickLength()
		a.ch2.tickLength()
		a.ch3.tickLength()
		a.ch4.tickLength()
	}

	if (a.frameSeqTicks-7)%8 == 0 {
		a.ch1.tickVolumeEnvelope()
		a.ch2.tickVolumeEnvelope()
		a.ch4.tickVolumeEnvelope()
	}

	if (a.frameSeqTicks-2)%4 == 0 {
		a.ch1.tickSweep()
	}

	a.frameSeqTicks++

}

func (a *Audio) tickSampler() {
	a.takeSample()
}
