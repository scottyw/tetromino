package interrupts

// Interrupts captures the current state of interrupts
type Interrupts struct {
	// FIXME these shouldn't exported really
	IE byte
	IF byte
}

// New Interrupts
func New() *Interrupts {
	return &Interrupts{
		IE: 0x00,
		IF: 0x01,
	}
}

// WriteIE handles writes to register IE
func (i *Interrupts) WriteIE(value uint8) {
	// fmt.Printf("> IE - 0x%02x\n", value)
	i.IE = value
}

// ReadIE handles reads from register IE
func (i *Interrupts) ReadIE() uint8 {
	ier := i.IE
	// fmt.Printf("< IE - 0x%02x\n", ier)
	return ier
}

// WriteIF handles writes to register IF
func (i *Interrupts) WriteIF(value uint8) {
	// fmt.Printf("> IF - 0x%02x\n", value)
	i.IF = value
}

// ReadIF handles reads from register IF
func (i *Interrupts) ReadIF() uint8 {
	// Top 3 bits are always high
	ifr := i.IF | 0xe0
	// fmt.Printf("< IF - 0x%02x\n", ifr)
	return ifr
}
