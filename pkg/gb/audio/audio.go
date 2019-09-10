package audio

import (
	"log"

	"github.com/gordonklaus/portaudio"
)

// Audio stream
type Audio struct {
	stream  *portaudio.Stream
	ch1     channel1
	ch2     channel2
	ch3     channel3
	ch4     channel4
	control control
}

// NewAudio initializes our internal channel for audio data
func NewAudio(silent bool) *Audio {
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

	// Sound 1 is also flagged as on at start time
	audio.control.sound1on = true

	// Initial the OS audio if enabled
	if !silent {
		portaudio.Initialize()
		err := audio.Start()
		if err != nil {
			log.Fatal(err)
		}
	}

	return &audio
}

// EndMachineCycle emulates the audio hardware at the end of a machine cycle
func (a *Audio) EndMachineCycle() {

	// mix audio from all channels here if enabled

}
