package audio

import (
	"math"
)

const (
	frameSeqTicks = 4194304 / 512      // 512Hz
	samplerPeriod = 95.108934240362812 // 44100 Hz
)

// Audio stream
type Audio struct {
	l             chan float32
	r             chan float32
	ch1           square
	ch2           square
	ch3           wave
	ch4           noise
	control       control
	waveram       [16]uint8
	ticks         uint64
	frameSeqTicks uint64
	samplerTicks  float64
}

// NewAudio initializes our internal channel for audio data
func NewAudio() *Audio {
	audio := Audio{
		waveram: [16]uint8{0x84, 0x40, 0x43, 0xAA, 0x2D, 0x78, 0x92, 0x3C, 0x60, 0x59, 0x59, 0xB0, 0x34, 0xB8, 0x2E, 0xDA},
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

	// Ch 1 supports sweep, ch 2 does not
	audio.ch1.supportSweep = true
	audio.ch2.supportSweep = false

	return &audio
}

// Speakers abstracts over a real-world implementation of the Gameboy speakers
type Speakers interface {
	Left() chan float32
	Right() chan float32
}

// RegisterSpeakers associates real-world audio output with the audio subsystem
func (a *Audio) RegisterSpeakers(speakers Speakers) {
	a.l = speakers.Left()
	a.r = speakers.Right()
}

// EndMachineCycle emulates the audio hardware at the end of a machine cycle
func (a *Audio) EndMachineCycle() {
	// Each machine cycle is four clock cycles
	a.tickClock()
	a.tickClock()
	a.tickClock()
	a.tickClock()
}

func (a *Audio) tickClock() {
	if a.ticks >= 4194304 {
		a.ticks = 0
		a.ticks = 0
		a.samplerTicks = 0
	}

	// Tick every clock cycle
	a.tick()

	// Tick the frame sequencer at 512 Hz
	if a.ticks%frameSeqTicks == 0 {
		a.tickFrameSequencer()
		a.frameSeqTicks++
		if a.frameSeqTicks >= 512 {
			a.frameSeqTicks = 0
		}
	}

	// Tick this function at 44100 Hz
	if a.ticks == uint64(math.Round(a.samplerTicks*samplerPeriod)) {
		a.tickSampler()
		a.samplerTicks++
	}

	a.ticks++
}

func (a *Audio) tick() {
	a.ch1.tickTimer()
	a.ch2.tickTimer()
	a.ch3.tickTimer()
	a.ch4.tickTimer()
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

}

func (a *Audio) tickSampler() {
	a.takeSample()
}
