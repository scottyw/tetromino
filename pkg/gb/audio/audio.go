package audio

// Speakers abstracts over a real-world implementation of the Gameboy speakers
type Speakers interface {
	Left() chan float32
	Right() chan float32
}

// Audio stream
type Audio struct {
	waitCycles uint32
	sample     uint32
	l          chan float32
	r          chan float32
	ch1        channel1
	ch2        channel2
	ch3        channel3
	ch4        channel4
	control    control
}

// NewAudio initializes our internal channel for audio data
func NewAudio() *Audio {
	audio := Audio{}

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

	// Ch1 is flagged as enabled at start time
	audio.control.ch1Enable = true

	return &audio
}

// EndMachineCycle emulates the audio hardware at the end of a machine cycle
func (a *Audio) EndMachineCycle() {

	// Audio is 44.1KHz which means writing every 24 cycles roughly
	// This messes up timing accuracy because once every 24 cycles is slightly too infrequently so the whole Gameboy slows down
	a.waitCycles++
	if a.waitCycles <= 24 {
		return
	}

	a.generateSample()

	a.sample++
	a.waitCycles = 0
}

// RegisterSpeakers associates real-world audio output with the audio subsystem
func (a *Audio) RegisterSpeakers(speakers Speakers) {
	a.l = speakers.Left()
	a.r = speakers.Right()
}
