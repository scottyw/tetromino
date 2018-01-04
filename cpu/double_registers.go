package cpu

type r16 struct {
	r1, r2 *uint8
}

// Get the 16-bit register value
func (r r16) Get() uint16 {
	return uint16(*r.r1)<<8 + uint16(*r.r2)
}

// Set the 16-bit register value
func (r r16) Set(val uint16) {
	*r.r1 = uint8(val >> 8)
	*r.r2 = uint8(val)
}

func newRegister16(r1, r2 *uint8) register16 {
	return r16{r1: r1, r2: r2}
}
