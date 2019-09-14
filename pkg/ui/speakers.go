package ui

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// PortaudioSpeakers implements speakers using portaudio
type PortaudioSpeakers struct {
	stream *portaudio.Stream
}

// NewPortaudioSpeakers starts audio output using portaudio
func NewPortaudioSpeakers() (*PortaudioSpeakers, error) {
	portaudio.Initialize()
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	parameters := portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice)
	speakers := &PortaudioSpeakers{}
	stream, err := portaudio.OpenStream(parameters, speakers.Callback)
	if err != nil {
		return nil, err
	}
	speakers.stream = stream
	if err := stream.Start(); err != nil {
		return nil, err
	}
	return speakers, nil
}

// Cleanup  returns resources to the OS
func (s *PortaudioSpeakers) Cleanup() {
	defer portaudio.Terminate()
	err := s.stream.Close()
	if err != nil {
		fmt.Println(err)
	}
}

// PlayAudio takes PCM values and buffers them for play
func (s *PortaudioSpeakers) PlayAudio(value uint8) {

	// Audio data incoming from the emulator arrives here

}

// Callback from portaudio to consume the audio data written to the channel
func (s *PortaudioSpeakers) Callback(out []float32) {

	length := len(out)

	// Left is 0th, 2nd, 4th ... array elements
	for i := 0; i < length; i += 2 {
		// do something here
		out[i] = 0
	}

	// Right  is 1st, 3rd, 5th ... array elements
	for i := 1; i < length; i += 2 {
		// do something here
		out[i] = 0
	}

}
