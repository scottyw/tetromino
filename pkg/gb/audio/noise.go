package audio

type noise struct {
	length           uint8
	initialVolume    uint8
	envelopeIncrease bool
	envelopeSweep    uint8
	shift            uint8
	step             uint8
	ratio            uint8
	lengthEnable     bool
}

func (n *noise) trigger() {

}
