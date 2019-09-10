package audio

import (
	"github.com/gordonklaus/portaudio"
)

// Start the portaudio stream
func (a *Audio) Start() error {
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return err
	}
	parameters := portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice)
	stream, err := portaudio.OpenStream(parameters, a.Callback)
	if err != nil {
		return err
	}
	if err := stream.Start(); err != nil {
		return err
	}
	a.stream = stream
	// a.sampleRate = parameters.SampleRate
	// a.outputChannels = parameters.Output.Channels
	return nil
}

// Stop the portaudio stream
func (a *Audio) Stop() error {
	return a.stream.Close()
}

// Callback from portaudio to consume the audio data written to the channel
func (a *Audio) Callback(out []float32) {

	length := len(out)

	// Left is 0th, 2nd, 4th ... array elements
	for i := 0; i < length; i += 2 {
		// do something here
	}

	// Right  is 1st, 3rd, 5th ... array elements
	for i := 1; i < length; i += 2 {
		// do something here

	}

}
