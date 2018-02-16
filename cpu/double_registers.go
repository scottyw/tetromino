package cpu

type r16 struct {
	msb, lsb *uint8
}

// Get the 16-bit register value
func (r r16) Get() uint16 {
	return uint16(*r.msb)<<8 + uint16(*r.lsb)
}

// Get the most significant byte
func (r r16) GetMsb() uint8 {
	return *r.msb
}

// Get the least significant byte
func (r r16) GetLsb() uint8 {
	return *r.lsb
}

// Set the 16-bit register value
func (r r16) Set(val uint16) {
	*r.msb = uint8(val >> 8)
	*r.lsb = uint8(val)
}

// Set the most significant byte
func (r r16) SetMsb(val uint8) {
	*r.msb = val
}

// Set the least significant byte
func (r r16) SetLsb(val uint8) {
	*r.lsb = val
}

func newRegister16(msb, lsb *uint8) register16 {
	return r16{msb: msb, lsb: lsb}
}
