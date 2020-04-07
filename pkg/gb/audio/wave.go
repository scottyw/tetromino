package audio

type wave struct {
	enable       bool
	length       uint8
	outputLevel  uint8
	frequency    uint16
	lengthEnable bool
}

func (w *wave) trigger() {

}
