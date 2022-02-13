package memory

type rtc struct {

	// RTC registers
	s     uint8
	m     uint8
	h     uint8
	d     uint16
	carry bool
	halt  bool

	// Latched RTC registers
	ls     uint8
	lm     uint8
	lh     uint8
	ld     uint16
	lcarry bool
	lhalt  bool

	// Internal state
	ticks int
	low   bool
}

func newRTC() *rtc {
	return &rtc{}
}

func (r *rtc) tick() {
	if r.halt {
		return
	}
	r.ticks++
	if r.ticks == 1048576 {
		r.ticks = 0
		r.increment()
	}
}

func (r *rtc) increment() {

	// Seconds
	r.s++
	if r.s == 60 {
		r.s = 0

		// Minutes
		r.m++
		if r.m == 60 {
			r.m = 0

			// Hours
			r.h++
			if r.h == 24 {
				r.h = 0

				// Days
				r.d++
				if r.d == 512 {
					r.d = 0
					r.carry = true
				}
			}
		}
	}

	// Make sure we're not using more bits than we should
	r.s &= 0x3f
	r.m &= 0x3f
	r.h &= 0x1f
	r.d &= 0x01ff

}

func (r *rtc) latchLow() {
	r.low = true
}

func (r *rtc) latchHigh() {
	if r.low {
		r.ls = r.s
		r.lm = r.m
		r.lh = r.h
		r.ld = r.d
		r.lcarry = r.carry
		r.lhalt = r.halt
	}
	r.low = false
}

func (r *rtc) read(ramBank uint8) uint8 {
	switch ramBank {
	case 0x08:
		return r.ls & 0x3f
	case 0x09:
		return r.lm & 0x3f
	case 0x0a:
		return r.lh & 0x1f
	case 0x0b:
		return uint8(r.ld)
	case 0x0c:
		control := uint8(r.ld>>8) & 0x01
		if r.lcarry {
			control += 0x80
		}
		if r.lhalt {
			control += 0x40
		}
		return control
	default:
		panic("invalid RTC register")
	}
}

func (r *rtc) write(ramBank uint8, value uint8) {
	switch ramBank {
	case 0x08:
		r.s = value & 0x3f
		r.ticks = 0
	case 0x09:
		r.m = value & 0x3f
	case 0x0a:
		r.h = value & 0x1f
	case 0x0b:
		r.d = (r.d & 0x0100) | uint16(value)
	case 0x0c:
		r.carry = (value >> 7) > 0
		r.halt = (value>>6)&0x01 > 0
		r.d = (uint16(value)&0x01)<<8 | (r.d & 0xff)
	}
}
