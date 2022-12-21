package cpu

func adcM(cpu *CPU) {
	adc(cpu.mapper.Read(cpu.hl()))
}

func adcA(cpu *CPU) { adc(cpu.a) }

func adcB(cpu *CPU) { adc(cpu.b) }

func adcC(cpu *CPU) { adc(cpu.c) }

func adcD(cpu *CPU) { adc(cpu.d) }

func adcE(cpu *CPU) { adc(cpu.e) }

func adcH(cpu *CPU) { adc(cpu.h) }

func adcL(cpu *CPU) { adc(cpu.l) }

func adcU(cpu *CPU) { adc(cpu.u8a) }

func adc(u8 uint8) {
	a := cpu.a
	cpu.a += u8
	hf := hc8(a, u8)
	cf := c8(a, u8)
	if cpu.cf() {
		cpu.a++
		hf = hf || cpu.a&0x0f == 0
		cf = cf || cpu.a == 0
	}
	// [Z 0 H C]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(hf)
	cpu.setCf(cf)
}

func addM(cpu *CPU) {
	add(cpu.mapper.Read(cpu.hl()))
}

func addA(cpu *CPU) { add(cpu.a) }

func addB(cpu *CPU) { add(cpu.b) }

func addC(cpu *CPU) { add(cpu.c) }

func addD(cpu *CPU) { add(cpu.d) }

func addE(cpu *CPU) { add(cpu.e) }

func addH(cpu *CPU) { add(cpu.h) }

func addL(cpu *CPU) { add(cpu.l) }

func addU(cpu *CPU) { add(cpu.u8a) }

func add(u8 uint8) {
	a := cpu.a
	cpu.a += u8
	// [Z 0 H C]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(hc8(a, u8))
	cpu.setCf(c8(a, u8))
}

func addHLBC(cpu *CPU) { addHL(cpu.bc()) }

func addHLDE(cpu *CPU) { addHL(cpu.de()) }

func addHLHL(cpu *CPU) { addHL(cpu.hl()) }

func addHLSP(cpu *CPU) { addHL(cpu.sp) }

func addHL(u16 uint16) {
	hl := cpu.hl()
	new := hl + u16
	cpu.h = uint8(new >> 8)
	cpu.l = uint8(new)
	// [- 0 H C]
	cpu.setNf(false)
	cpu.setHf(hc16(hl, u16))
	cpu.setCf(c16(hl, u16))
}

func addSP(cpu *CPU) {
	i8 := int8(cpu.u8a)
	sp := cpu.sp
	cpu.sp = uint16(int(cpu.sp) + int(i8))
	// [0 0 H C]
	cpu.setZf(false)
	cpu.setNf(false)
	if i8 >= 0 {
		cpu.setHf(int(sp)&0x0f+int(i8)&0x0f > 0x0f)
		cpu.setCf(int(sp)&0xff+int(i8) > 0xff)
	} else {
		cpu.setHf(int(sp)&0x0f >= int(cpu.sp)&0x0f)
		cpu.setCf(int(sp)&0xff >= int(cpu.sp)&0xff)
	}
}

func andM(cpu *CPU) {
	and(cpu.mapper.Read(cpu.hl()))
}

func andA(cpu *CPU) { and(cpu.a) }

func andB(cpu *CPU) { and(cpu.b) }

func andC(cpu *CPU) { and(cpu.c) }

func andD(cpu *CPU) { and(cpu.d) }

func andE(cpu *CPU) { and(cpu.e) }

func andH(cpu *CPU) { and(cpu.h) }

func andL(cpu *CPU) { and(cpu.l) }

func andU(cpu *CPU) { and(cpu.u8a) }

func and(u8 uint8) {
	cpu.a &= u8
	// [Z 0 1 0]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(true)
	cpu.setCf(false)
}

func bitM(pos uint8) func() {
	return func() {
		u8 := cpu.mapper.Read(cpu.hl())
		bit(pos, &u8)()
	}
}

func bit(pos uint8, r8 *uint8) func() {
	return func() {
		zero := *r8&bits[pos] == 0
		// [Z 0 1 -]
		cpu.setZf(zero)
		cpu.setNf(false)
		cpu.setHf(true)
	}
}

func call(cpu *CPU) {
	// Store the old PC to write to memory in the next steps
	cpu.m8a = uint8(cpu.pc & 0xff)
	cpu.m8b = uint8(cpu.pc >> 8)
	// Update the PC
	cpu.pc = cpu.u16()
}

func ccf(cpu *CPU) {
	// [- 0 0 C]
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(!cpu.cf())
}

func cpM(cpu *CPU) {
	cp(cpu.mapper.Read(cpu.hl()))
}

func cpA(cpu *CPU) { cp(cpu.a) }

func cpB(cpu *CPU) { cp(cpu.b) }

func cpC(cpu *CPU) { cp(cpu.c) }

func cpD(cpu *CPU) { cp(cpu.d) }

func cpE(cpu *CPU) { cp(cpu.e) }

func cpH(cpu *CPU) { cp(cpu.h) }

func cpL(cpu *CPU) { cp(cpu.l) }

func cpU(cpu *CPU) { cp(cpu.u8a) }

func cp(u8 uint8) {
	// [Z 1 H C]
	cpu.setZf(cpu.a == u8)
	cpu.setNf(true)
	cpu.setHf(hc8Sub(cpu.a, u8))
	cpu.setCf(c8Sub(cpu.a, u8))
}

func cpl(cpu *CPU) {
	cpu.a = ^cpu.a
	// [- 1 1 -]
	cpu.setNf(true)
	cpu.setHf(true)
}

func daa(cpu *CPU) {
	if cpu.nf() {
		if cpu.hf() {
			cpu.a -= 0x06
		}
		if cpu.cf() {
			cpu.a -= 0x60
		}
	} else {
		a := cpu.a
		if cpu.hf() || cpu.a&0x0f > 0x09 {
			cpu.a += 0x06
		}
		if cpu.cf() || a&0xf0 > 0x90 || cpu.a&0xf0 > 0x90 {
			cpu.a += 0x60
			cpu.setCf(true)
		}
	}
	// [Z - 0 C]
	cpu.setZf(cpu.a == 0)
	cpu.setHf(false)
}

func decM(cpu *CPU) {
	dec(&cpu.m8a)
	cpu.mapper.Write(cpu.hl(), cpu.m8a)
}

func decA(cpu *CPU) { dec(&cpu.a) }

func decB(cpu *CPU) { dec(&cpu.b) }

func decC(cpu *CPU) { dec(&cpu.c) }

func decD(cpu *CPU) { dec(&cpu.d) }

func decE(cpu *CPU) { dec(&cpu.e) }

func decH(cpu *CPU) { dec(&cpu.h) }

func decL(cpu *CPU) { dec(&cpu.l) }

func dec(r8 *uint8) {
	old := *r8
	*r8--
	// [Z 1 H -]
	cpu.setZf(*r8 == 0)
	cpu.setNf(true)
	cpu.setHf(hc8Sub(old, 1))
}

func decBC(cpu *CPU) {
	cpu.oam.TriggerWriteCorruption(cpu.bc())
	dec16(&cpu.b, &cpu.c)
}

func decDE(cpu *CPU) {
	cpu.oam.TriggerWriteCorruption(cpu.de())
	dec16(&cpu.d, &cpu.e)
}

func decHL(cpu *CPU) {
	cpu.oam.TriggerWriteCorruption(cpu.hl())
	dec16(&cpu.h, &cpu.l)
}

func dec16(msb, lsb *uint8) {
	new := uint16(*msb)<<8 + uint16(*lsb) - 1
	*msb = uint8(new >> 8)
	*lsb = uint8(new)
}

func decSP(cpu *CPU) {
	cpu.oam.TriggerWriteCorruption(cpu.sp)
	cpu.sp--
}

func di(cpu *CPU) {
	cpu.interrupts.Disable()
}

func ei(cpu *CPU) {
	cpu.interrupts.Enable()
}

func halt(cpu *CPU) {
	if cpu.interrupts.Enabled() {
		cpu.halted = true
	} else {
		if !cpu.interrupts.Pending() {
			cpu.halted = true
		} else {
			cpu.haltbug = true
		}
	}
}

func incM(cpu *CPU) {
	inc(&cpu.m8a)
	cpu.mapper.Write(cpu.hl(), cpu.m8a)
}

func incA(cpu *CPU) { inc(&cpu.a) }

func incB(cpu *CPU) { inc(&cpu.b) }

func incC(cpu *CPU) { inc(&cpu.c) }

func incD(cpu *CPU) { inc(&cpu.d) }

func incE(cpu *CPU) { inc(&cpu.e) }

func incH(cpu *CPU) { inc(&cpu.h) }

func incL(cpu *CPU) { inc(&cpu.l) }

func inc(r8 *uint8) {
	old := *r8
	*r8++
	// [Z 0 H -]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(hc8(old, 1))
}

func incBC(cpu *CPU) {
	cpu.oam.TriggerWriteCorruption(cpu.bc())
	inc16(&cpu.b, &cpu.c)
}

func incDE(cpu *CPU) {
	cpu.oam.TriggerWriteCorruption(cpu.de())
	inc16(&cpu.d, &cpu.e)
}

func incHL(cpu *CPU) {
	cpu.oam.TriggerWriteCorruption(cpu.hl())
	inc16(&cpu.h, &cpu.l)
}

func inc16(msb, lsb *uint8) {
	new := uint16(*msb)<<8 + uint16(*lsb) + 1
	*msb = uint8(new >> 8)
	*lsb = uint8(new)
}

func incSP(cpu *CPU) {
	cpu.oam.TriggerWriteCorruption(cpu.sp)
	cpu.sp++
}

func jp(cpu *CPU) {
	cpu.pc = cpu.u16()
}

func jpHL(cpu *CPU) {
	cpu.pc = cpu.hl()
}

func jr(cpu *CPU) {
	i8 := int8(cpu.u8a)
	cpu.pc = uint16(int16(cpu.pc) + int16(i8))
}

func ldAB(cpu *CPU) { cpu.a = cpu.b }

func ldAC(cpu *CPU) { cpu.a = cpu.c }

func ldAD(cpu *CPU) { cpu.a = cpu.d }

func ldAE(cpu *CPU) { cpu.a = cpu.e }

func ldAH(cpu *CPU) { cpu.a = cpu.h }

func ldAL(cpu *CPU) { cpu.a = cpu.l }

func ldAU(cpu *CPU) { cpu.a = cpu.u8a }

func ldBA(cpu *CPU) { cpu.b = cpu.a }

func ldBC(cpu *CPU) { cpu.b = cpu.c }

func ldBD(cpu *CPU) { cpu.b = cpu.d }

func ldBE(cpu *CPU) { cpu.b = cpu.e }

func ldBH(cpu *CPU) { cpu.b = cpu.h }

func ldBL(cpu *CPU) { cpu.b = cpu.l }

func ldBU(cpu *CPU) { cpu.b = cpu.u8a }

func ldCA(cpu *CPU) { cpu.c = cpu.a }

func ldCB(cpu *CPU) { cpu.c = cpu.b }

func ldCD(cpu *CPU) { cpu.c = cpu.d }

func ldCE(cpu *CPU) { cpu.c = cpu.e }

func ldCH(cpu *CPU) { cpu.c = cpu.h }

func ldCL(cpu *CPU) { cpu.c = cpu.l }

func ldCU(cpu *CPU) { cpu.c = cpu.u8a }

func ldDA(cpu *CPU) { cpu.d = cpu.a }

func ldDB(cpu *CPU) { cpu.d = cpu.b }

func ldDC(cpu *CPU) { cpu.d = cpu.c }

func ldDE(cpu *CPU) { cpu.d = cpu.e }

func ldDH(cpu *CPU) { cpu.d = cpu.h }

func ldDL(cpu *CPU) { cpu.d = cpu.l }

func ldDU(cpu *CPU) { cpu.d = cpu.u8a }

func ldEA(cpu *CPU) { cpu.e = cpu.a }

func ldEB(cpu *CPU) { cpu.e = cpu.b }

func ldEC(cpu *CPU) { cpu.e = cpu.c }

func ldED(cpu *CPU) { cpu.e = cpu.d }

func ldEH(cpu *CPU) { cpu.e = cpu.h }

func ldEL(cpu *CPU) { cpu.e = cpu.l }

func ldEU(cpu *CPU) { cpu.e = cpu.u8a }

func ldHA(cpu *CPU) { cpu.h = cpu.a }

func ldHB(cpu *CPU) { cpu.h = cpu.b }

func ldHC(cpu *CPU) { cpu.h = cpu.c }

func ldHD(cpu *CPU) { cpu.h = cpu.d }

func ldHE(cpu *CPU) { cpu.h = cpu.e }

func ldHL(cpu *CPU) { cpu.h = cpu.l }

func ldHU(cpu *CPU) { cpu.h = cpu.u8a }

func ldLA(cpu *CPU) { cpu.l = cpu.a }

func ldLB(cpu *CPU) { cpu.l = cpu.b }

func ldLC(cpu *CPU) { cpu.l = cpu.c }

func ldLD(cpu *CPU) { cpu.l = cpu.d }

func ldLE(cpu *CPU) { cpu.l = cpu.e }

func ldLH(cpu *CPU) { cpu.l = cpu.h }

func ldLU(cpu *CPU) { cpu.l = cpu.u8a }

func ldHLU8(cpu *CPU) {
	cpu.mapper.Write(cpu.hl(), cpu.u8a)
}

func ldBCU16(cpu *CPU) {
	u16 := cpu.u16()
	cpu.b = uint8(u16 >> 8)
	cpu.c = uint8(u16)
}

func ldDEU16(cpu *CPU) {
	u16 := cpu.u16()
	cpu.d = uint8(u16 >> 8)
	cpu.e = uint8(u16)
}

func ldHLU16(cpu *CPU) {
	u16 := cpu.u16()
	cpu.h = uint8(u16 >> 8)
	cpu.l = uint8(u16)
}

func ldSPU16(cpu *CPU) {
	cpu.sp = cpu.u16()
}

func ldSPHL(cpu *CPU) {
	cpu.sp = cpu.hl()
}

func ldHLSP(cpu *CPU) {
	i8 := int8(cpu.u8a)
	new := int(int(cpu.sp) + int(i8))
	cpu.h = uint8(new >> 8)
	cpu.l = uint8(new)
	// [0 0 H C]
	cpu.setZf(false)
	cpu.setNf(false)
	if i8 >= 0 {
		cpu.setHf(int(cpu.sp)&0x0f+int(i8)&0x0f > 0x0f)
		cpu.setCf(int(cpu.sp)&0xff+int(i8) > 0xff)
	} else {
		cpu.setHf(int(cpu.sp)&0x0f >= new&0x0f)
		cpu.setCf(int(cpu.sp)&0xff >= new&0xff)
	}
}

func ldBCA(cpu *CPU) {
	cpu.mapper.Write(cpu.bc(), cpu.a)
}

func ldABC(cpu *CPU) {
	cpu.a = cpu.mapper.Read(cpu.bc())
}

func ldDEA(cpu *CPU) {
	cpu.mapper.Write(cpu.de(), cpu.a)
}

func ldADE(cpu *CPU) {
	cpu.a = cpu.mapper.Read(cpu.de())
}

func ldHLDA(cpu *CPU) {
	cpu.mapper.Write(cpu.hl(), cpu.a)
	decHL(cpu)
}

func ldAHLD(cpu *CPU) {
	cpu.a = cpu.mapper.Read(cpu.hl())
	decHL(cpu)
}

func ldHLIA(cpu *CPU) {
	cpu.mapper.Write(cpu.hl(), cpu.a)
	incHL(cpu)
}

func ldAHLI(cpu *CPU) {
	cpu.a = cpu.mapper.Read(cpu.hl())
	incHL(cpu)
}

func ldHLA(cpu *CPU) { cpu.mapper.Write(cpu.hl(), cpu.a) }

func ldHLB(cpu *CPU) { cpu.mapper.Write(cpu.hl(), cpu.b) }

func ldHLC(cpu *CPU) { cpu.mapper.Write(cpu.hl(), cpu.c) }

func ldHLD(cpu *CPU) { cpu.mapper.Write(cpu.hl(), cpu.d) }

func ldHLE(cpu *CPU) { cpu.mapper.Write(cpu.hl(), cpu.e) }

func ldHLH(cpu *CPU) { cpu.mapper.Write(cpu.hl(), cpu.h) }

func ldHLL(cpu *CPU) { cpu.mapper.Write(cpu.hl(), cpu.l) }

func ldAHL(cpu *CPU) { cpu.a = cpu.mapper.Read(cpu.hl()) }

func ldBHL(cpu *CPU) { cpu.b = cpu.mapper.Read(cpu.hl()) }

func ldCHL(cpu *CPU) { cpu.c = cpu.mapper.Read(cpu.hl()) }

func ldDHL(cpu *CPU) { cpu.d = cpu.mapper.Read(cpu.hl()) }

func ldEHL(cpu *CPU) { cpu.e = cpu.mapper.Read(cpu.hl()) }

func ldHHL(cpu *CPU) { cpu.h = cpu.mapper.Read(cpu.hl()) }

func ldLHL(cpu *CPU) { cpu.l = cpu.mapper.Read(cpu.hl()) }

func ldMHL(cpu *CPU) {
	cpu.m8a = cpu.mapper.Read(cpu.hl())
}

func ldACX(cpu *CPU) {
	a16 := uint16(0xff00 + uint16(cpu.c))
	cpu.a = cpu.mapper.Read(a16)
}

func ldCXA(cpu *CPU) {
	a16 := uint16(0xff00 + uint16(cpu.c))
	cpu.mapper.Write(a16, cpu.a)
}

func ldAUX(cpu *CPU) {
	a16 := uint16(0xff00 + uint16(cpu.u8a))
	cpu.a = cpu.mapper.Read(a16)
}

func ldUXA(cpu *CPU) {
	a16 := uint16(0xff00 + uint16(cpu.u8a))
	cpu.mapper.Write(a16, cpu.a)
}

func ldAUX16(cpu *CPU) {
	cpu.a = cpu.mapper.Read(cpu.u16())
}

func ldUX16A(cpu *CPU) {
	cpu.mapper.Write(cpu.u16(), cpu.a)
}

func writeLowSP(cpu *CPU) {
	cpu.mapper.Write(cpu.u16(), uint8(cpu.sp))
}

func writeHighSP(cpu *CPU) {
	cpu.mapper.Write(cpu.u16()+1, uint8(cpu.sp>>8))
}

func orM(cpu *CPU) {
	or(cpu.mapper.Read(cpu.hl()))
}

func orA(cpu *CPU) { or(cpu.a) }

func orB(cpu *CPU) { or(cpu.b) }

func orC(cpu *CPU) { or(cpu.c) }

func orD(cpu *CPU) { or(cpu.d) }

func orE(cpu *CPU) { or(cpu.e) }

func orH(cpu *CPU) { or(cpu.h) }

func orL(cpu *CPU) { or(cpu.l) }

func orU(cpu *CPU) { or(cpu.u8a) }

func or(u8 uint8) {
	cpu.a |= u8
	// [Z 0 0 0]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(false)
}

func pop(r8 *uint8) func() {
	return func() {
		*r8 = cpu.mapper.Read(cpu.sp)
		incSP(cpu)
	}
}

func popF(cpu *CPU) {
	// Lower nibble is always zero no matter what data was written
	cpu.f = cpu.mapper.Read(cpu.sp) & 0xf0
	cpu.oam.TriggerWriteCorruption(cpu.sp)
	cpu.sp++
}

func push(r8 *uint8) func(*CPU) {
	return func(cpu *CPU) {
		decSP(cpu)
		cpu.mapper.Write(cpu.sp, *r8)
	}
}

func resM(pos uint8) func(*CPU) {
	return func(cpu *CPU) {
		res(pos, &cpu.m8a)
		cpu.mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func res(pos uint8, r8 *uint8) func(*CPU) {
	return func(cpu *CPU) {
		*r8 &^= bits[pos]
	}
}

func ret(cpu *CPU) {
	cpu.pc = uint16(cpu.m8b)<<8 | uint16(cpu.m8a)
}

func reti(cpu *CPU) {
	ret(cpu)
	ei(cpu)
}

func rlM(cpu *CPU) {
	rl(&cpu.m8a)
	cpu.mapper.Write(cpu.hl(), cpu.m8a)
}

func rlA(cpu *CPU) { rl(&cpu.a) }

func rlB(cpu *CPU) { rl(&cpu.b) }

func rlC(cpu *CPU) { rl(&cpu.c) }

func rlD(cpu *CPU) { rl(&cpu.d) }

func rlE(cpu *CPU) { rl(&cpu.e) }

func rlH(cpu *CPU) { rl(&cpu.h) }

func rlL(cpu *CPU) { rl(&cpu.l) }

func rl(r8 *uint8) {
	cf := (*r8 & 0x80) > 0
	*r8 <<= 1
	if cpu.cf() {
		*r8 |= 0x01
	}
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func rla(cpu *CPU) {
	rl(&cpu.a)
	cpu.f &^= zFlag
	// [0 0 0 C]
	cpu.setZf(false)
}

func rlcM(cpu *CPU) {
	rlc(&cpu.m8a)
	cpu.mapper.Write(cpu.hl(), cpu.m8a)
}

func rlcA(cpu *CPU) { rlc(&cpu.a) }

func rlcB(cpu *CPU) { rlc(&cpu.b) }

func rlcC(cpu *CPU) { rlc(&cpu.c) }

func rlcD(cpu *CPU) { rlc(&cpu.d) }

func rlcE(cpu *CPU) { rlc(&cpu.e) }

func rlcH(cpu *CPU) { rlc(&cpu.h) }

func rlcL(cpu *CPU) { rlc(&cpu.l) }

func rlc(r8 *uint8) {
	cf := (*r8 & 0x80) > 0
	*r8 <<= 1
	if cf {
		*r8 |= 0x01
	}
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func rlca(cpu *CPU) {
	rlc(&cpu.a)
	cpu.f &^= zFlag
	// [0 0 0 C]
	cpu.setZf(false)
}

func rrM(cpu *CPU) {
	rr(&cpu.m8a)
	cpu.mapper.Write(cpu.hl(), cpu.m8a)
}

func rrA(cpu *CPU) { rr(&cpu.a) }

func rrB(cpu *CPU) { rr(&cpu.b) }

func rrC(cpu *CPU) { rr(&cpu.c) }

func rrD(cpu *CPU) { rr(&cpu.d) }

func rrE(cpu *CPU) { rr(&cpu.e) }

func rrH(cpu *CPU) { rr(&cpu.h) }

func rrL(cpu *CPU) { rr(&cpu.l) }

func rr(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	*r8 >>= 1
	if cpu.cf() {
		*r8 |= 0x80
	}
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func rra(cpu *CPU) {
	rr(&cpu.a)
	// [0 0 0 C]
	cpu.setZf(false)
}

func rrcM(cpu *CPU) {
	rrc(&cpu.m8a)
	cpu.mapper.Write(cpu.hl(), cpu.m8a)
}

func rrcA(cpu *CPU) { rrc(&cpu.a) }

func rrcB(cpu *CPU) { rrc(&cpu.b) }

func rrcC(cpu *CPU) { rrc(&cpu.c) }

func rrcD(cpu *CPU) { rrc(&cpu.d) }

func rrcE(cpu *CPU) { rrc(&cpu.e) }

func rrcH(cpu *CPU) { rrc(&cpu.h) }

func rrcL(cpu *CPU) { rrc(&cpu.l) }

func rrc(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	*r8 >>= 1
	if cf {
		*r8 |= 0x80
	}
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func rrca(cpu *CPU) {
	rrc(&cpu.a)
	// [0 0 0 C]
	cpu.setZf(false)
}

func rst(a16 uint16) func(*CPU) {
	return func(cpu *CPU) {
		// Store the old PC to write to memory in the next steps
		cpu.m8a = uint8(cpu.pc & 0xff)
		cpu.m8b = uint8(cpu.pc >> 8)
		// Update the PC
		cpu.pc = a16
	}
}

func setM(pos uint8) func(*CPU) {
	return func(cpu *CPU) {
		set(pos, &cpu.m8a)()
		cpu.mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func set(pos uint8, r8 *uint8) func() {
	return func() {
		*r8 |= bits[pos]
	}
}

func slaM(cpu *CPU) {
	sla(&cpu.m8a)
	cpu.mapper.Write(cpu.hl(), cpu.m8a)
}

func slaA(cpu *CPU) { sla(&cpu.a) }

func slaB(cpu *CPU) { sla(&cpu.b) }

func slaC(cpu *CPU) { sla(&cpu.c) }

func slaD(cpu *CPU) { sla(&cpu.d) }

func slaE(cpu *CPU) { sla(&cpu.e) }

func slaH(cpu *CPU) { sla(&cpu.h) }

func slaL(cpu *CPU) { sla(&cpu.l) }

func sla(r8 *uint8) {
	cf := (*r8 & 0x80) > 0
	*r8 <<= 1
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func sraM(cpu *CPU) {
	sra(&cpu.m8a)
	cpu.mapper.Write(cpu.hl(), cpu.m8a)
}

func sraA(cpu *CPU) { sra(&cpu.a) }

func sraB(cpu *CPU) { sra(&cpu.b) }

func sraC(cpu *CPU) { sra(&cpu.c) }

func sraD(cpu *CPU) { sra(&cpu.d) }

func sraE(cpu *CPU) { sra(&cpu.e) }

func sraH(cpu *CPU) { sra(&cpu.h) }

func sraL(cpu *CPU) { sra(&cpu.l) }

func sra(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	bit7 := (*r8 & 0x80) > 0
	*r8 >>= 1
	if bit7 {
		*r8 |= 0x80
	}
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func srlM(cpu *CPU) {
	srl(&cpu.m8a)
	cpu.mapper.Write(cpu.hl(), cpu.m8a)
}

func srlA(cpu *CPU) { srl(&cpu.a) }

func srlB(cpu *CPU) { srl(&cpu.b) }

func srlC(cpu *CPU) { srl(&cpu.c) }

func srlD(cpu *CPU) { srl(&cpu.d) }

func srlE(cpu *CPU) { srl(&cpu.e) }

func srlH(cpu *CPU) { srl(&cpu.h) }

func srlL(cpu *CPU) { srl(&cpu.l) }

func srl(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	*r8 >>= 1
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func swapM(cpu *CPU) {
	swap(&cpu.m8a)
	cpu.mapper.Write(cpu.hl(), cpu.m8a)
}

func swapA(cpu *CPU) { swap(&cpu.a) }

func swapB(cpu *CPU) { swap(&cpu.b) }

func swapC(cpu *CPU) { swap(&cpu.c) }

func swapD(cpu *CPU) { swap(&cpu.d) }

func swapE(cpu *CPU) { swap(&cpu.e) }

func swapH(cpu *CPU) { swap(&cpu.h) }

func swapL(cpu *CPU) { swap(&cpu.l) }

func swap(r8 *uint8) {
	u8 := *r8
	*r8 = u8<<4 | u8>>4
	// [Z 0 0 0]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(false)
}

func scf(cpu *CPU) {
	// [- 0 0 1]
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(true)
}

func stop(cpu *CPU) {
	cpu.stopped = true
}

func sbcM(cpu *CPU) {
	sbc(cpu.mapper.Read(cpu.hl()))
}

func sbcA(cpu *CPU) { sbc(cpu.a) }

func sbcB(cpu *CPU) { sbc(cpu.b) }

func sbcC(cpu *CPU) { sbc(cpu.c) }

func sbcD(cpu *CPU) { sbc(cpu.d) }

func sbcE(cpu *CPU) { sbc(cpu.e) }

func sbcH(cpu *CPU) { sbc(cpu.h) }

func sbcL(cpu *CPU) { sbc(cpu.l) }

func sbcU(cpu *CPU) { sbc(cpu.u8a) }

func sbc(u8 uint8) {
	a := cpu.a
	cpu.a -= u8
	hf := hc8Sub(a, u8)
	cf := c8Sub(a, u8)
	if cpu.cf() {
		cpu.a--
		hf = hf || cpu.a&0x0f == 0x0f
		cf = cf || cpu.a == 0xff
	}
	// [Z 1 H C]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(true)
	cpu.setHf(hf)
	cpu.setCf(cf)

}

func subM(cpu *CPU) {
	sub(cpu.mapper.Read(cpu.hl()))
}

func subA(cpu *CPU) { sub(cpu.a) }

func subB(cpu *CPU) { sub(cpu.b) }

func subC(cpu *CPU) { sub(cpu.c) }

func subD(cpu *CPU) { sub(cpu.d) }

func subE(cpu *CPU) { sub(cpu.e) }

func subH(cpu *CPU) { sub(cpu.h) }

func subL(cpu *CPU) { sub(cpu.l) }

func subU(cpu *CPU) { sub(cpu.u8a) }

func sub(u8 uint8) {
	a := cpu.a
	cpu.a -= u8
	// [Z 1 H C]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(true)
	cpu.setHf(hc8Sub(a, u8))
	cpu.setCf(c8Sub(a, u8))
}

func xorM(cpu *CPU) {
	xor(cpu.mapper.Read(cpu.hl()))
}

func xorA(cpu *CPU) { xor(cpu.a) }

func xorB(cpu *CPU) { xor(cpu.b) }

func xorC(cpu *CPU) { xor(cpu.c) }

func xorD(cpu *CPU) { xor(cpu.d) }

func xorE(cpu *CPU) { xor(cpu.e) }

func xorH(cpu *CPU) { xor(cpu.h) }

func xorL(cpu *CPU) { xor(cpu.l) }

func xorU(cpu *CPU) { xor(cpu.u8a) }

func xor(u8 uint8) {
	cpu.a = cpu.a ^ u8
	// [Z 0 0 0]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(false)
}
