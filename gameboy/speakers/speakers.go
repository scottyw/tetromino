package speakers

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// PortaudioSpeakers implements speakers using portaudio
type PortaudioSpeakers struct {
	stream *portaudio.Stream
	l      chan float32
	r      chan float32
}

// NewPortaudioSpeakers starts audio output using portaudio
func NewPortaudioSpeakers() (*PortaudioSpeakers, error) {
	portaudio.Initialize()
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	parameters := portaudio.LowLatencyParameters(nil, host.DefaultOutputDevice)
	speakers := &PortaudioSpeakers{
		l: make(chan float32, 200),
		r: make(chan float32, 200),
	}
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

// Cleanup returns resources to the OS
func (s *PortaudioSpeakers) Cleanup() {
	defer portaudio.Terminate()
	close(s.l)
	close(s.r)
	err := s.stream.Close()
	if err != nil {
		fmt.Println(err)
	}
}

// Left returns the channel that feeds the left speaker
func (s *PortaudioSpeakers) Left() chan float32 {
	return s.l
}

// Right returns the channel that feeds the right speaker
func (s *PortaudioSpeakers) Right() chan float32 {
	return s.r
}

// Callback from portaudio to consume the audio data written to the channel
func (s *PortaudioSpeakers) Callback(out []float32) {

	// Low latency callback every 1.44216ms approx i.e. 693.4 times per second approx
	// Array size is always 126 i.e. 88200 elements per second

	// High latency callback every 11.581337ms approx i.e. 86.3 times per second approx
	// Array size is always 1022 i.e. 88200 elements per second

	length := len(out)

	// Left is 0th, 2nd, 4th ... array elements
	// Right  is 1st, 3rd, 5th ... array elements
	for i := 0; i < length; i += 2 {
		out[i] = <-s.l
		out[i+1] = <-s.r
	}

}
