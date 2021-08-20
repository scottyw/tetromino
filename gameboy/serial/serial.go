package serial

import (
	"fmt"
	"io"
)

// Serial captures the current state of serial bus
type Serial struct {
	sc     byte
	writer io.Writer
}

// New Serial
func New(writer io.Writer) *Serial {
	return &Serial{
		writer: writer,
	}
}

// WriteSB handles writes to register SB
func (s *Serial) WriteSB(value uint8) {
	// fmt.Printf("> SB - 0x%02x\n", value)
	if s.writer == nil {
		return
	}
	_, err := s.writer.Write([]byte{value})
	if err != nil {
		panic(fmt.Sprintf("Write to SB failed: %v", err))
	}
}

// ReadSB handles reads from register SB
func (s *Serial) ReadSB() uint8 {
	// fmt.Printf("< SB - 0xff\n")
	return 0xff // FIXME
}

// WriteSC handles writes to register SC
func (s *Serial) WriteSC(value uint8) {
	// fmt.Printf("> SC - 0x%02x\n", value)
	// FIXME
}

// ReadSC handles reads from register SC
func (s *Serial) ReadSC() uint8 {
	// fmt.Printf("< SC - 0xff\n")
	return 0xff // FIXME
}
