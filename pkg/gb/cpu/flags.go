package cpu

import "fmt"

func (cpu *CPU) zf() bool {
	return cpu.f&zFlag > 0
}

func (cpu *CPU) nf() bool {
	return cpu.f&nFlag > 0
}

func (cpu *CPU) hf() bool {
	return cpu.f&hFlag > 0
}

func (cpu *CPU) cf() bool {
	return cpu.f&cFlag > 0
}

func (cpu *CPU) setZf(value bool) {
	if value {
		cpu.f |= zFlag
	} else {
		cpu.f &^= zFlag
	}
}

func (cpu *CPU) setNf(value bool) {
	if value {
		cpu.f |= nFlag
	} else {
		cpu.f &^= nFlag
	}
}

func (cpu *CPU) setHf(value bool) {
	if value {
		cpu.f |= hFlag
	} else {
		cpu.f &^= hFlag
	}
}

func (cpu *CPU) setCf(value bool) {
	if value {
		cpu.f |= cFlag
	} else {
		cpu.f &^= cFlag
	}
}

func hc8(a, b uint8) bool {
	return a&0x0f+b&0x0f > 0x0f
}
func c8(a, b uint8) bool {
	return int(a)+int(b) > 0xff
}

func hc16(a, b uint16) bool {
	return a&0x0fff+b&0x0fff > 0x0fff
}

func c16(a, b uint16) bool {
	return int(a)+int(b) > 0xffff
}

func hc8Sub(a, b uint8) bool {
	return int(a)&0x0f-int(b)&0x0f < 0
}

func c8Sub(a, b uint8) bool {
	return int(a)-int(b) < 0
}

func flagMetadata(i uint, flags []string) string {
	if len(flags) == 0 {
		return "-"
	}
	return flags[i]
}

func validateFlag(label string, i uint, f1, f2 uint8, im metadata) {
	bit := uint8(0x80) >> i
	switch flagMetadata(i, im.Flags) {
	case "-":
		if f1&bit != f2&bit {
			panic(fmt.Sprintf("%s flag invalid! Should not change: before=0x%02x after=0x%02x metadata=%v", label, f1, f2, im))
		}
	case "0":
		if f2&bit != 0 {
			panic(fmt.Sprintf("%s flag invalid! Should be reset: flags=0x%02x metadata=%v", label, f2, im))
		}
	case "1":
		if f2&bit == 0 {
			panic(fmt.Sprintf("%s flag invalid! Should be set: flags=0x%02x metadata=%v", label, f2, im))
		}
	}
}

func validateFlags(f1, f2 uint8, im metadata) {
	validateFlag("Z", 0, f1, f2, im)
	validateFlag("N", 1, f1, f2, im)
	validateFlag("H", 2, f1, f2, im)
	validateFlag("C", 3, f1, f2, im)
}
