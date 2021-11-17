package main

// status register layout
// 7 6 5 4 3 2 1 0
// N V   B D I Z C

type (
	Opcode struct {
		GetAddress  func() *uint16
		Instruction func(*uint16)
		Title       string
	}
)

const (
	FlagC = iota
	FlagZ
	FlagI
	FlagD
	FlagB
	_
	FlagV
	FlagN
)

func getFlag(flagNum uint8) uint8 {
	return (cpu.P >> flagNum) & 0x01
}

func setFlag(flagNum uint8, value uint8) {
	if value == 1 {
		cpu.P = cpu.P | (value << flagNum)
		return
	}

	cpu.P = cpu.P & ^(1 << flagNum)
}

func accumulator() *uint16 {
	cpu.PC++
	return nil
}

func immediate() *uint16 {
	cpu.PC++
	p := cpu.PC
	cpu.PC++
	return &p
}

func zeroPage() *uint16 {
	cpu.PC++
	address := uint16(cpu.Read(cpu.PC))
	cpu.PC++
	return &address
}

func zeroPageX() *uint16 {
	cpu.PC++
	address := uint16(cpu.Read(cpu.PC))
	cpu.PC++
	address += uint16(cpu.X)
	address &= 0xff
	return &address
}

func zeroPageY() *uint16 {
	cpu.PC++
	address := uint16(cpu.Read(cpu.PC))
	cpu.PC++
	address += uint16(cpu.Y)
	address &= 0xff
	return &address
}

func relative() *uint16 {
	cpu.PC++
	offset := uint16(cpu.Read(cpu.PC))
	return &offset
}

func absolute() *uint16 {
	cpu.PC++
	address := read16(cpu.PC)
	cpu.PC += 2
	return &address
}

func absoluteX() *uint16 {
	cpu.PC++
	address := read16(cpu.PC)
	cpu.PC += 2
	return &address
}

func absoluteY() *uint16 {
	cpu.PC++
	address := read16(cpu.PC)
	cpu.PC += 2
	address += uint16(cpu.Y)
	return &address
}

func indirect() *uint16 {
	cpu.PC++
	pointer := read16(cpu.PC)
	cpu.PC++
	address := read16bug(pointer)
	return &address
}

func indirectX() *uint16 {
	cpu.PC++
	pointer := uint16(cpu.Read(cpu.PC))
	cpu.PC++
	address := read16bug(pointer)
	address += uint16(cpu.X)
	return &address
}

func indirectY() *uint16 {
	cpu.PC++
	pointer := uint16(cpu.Read(cpu.PC))
	cpu.PC++
	address := read16bug(pointer)
	address += uint16(cpu.Y)
	return &address
}

func implied() *uint16 {
	cpu.PC++
	return nil
}

func read16(address uint16) uint16 {
	low := uint16(cpu.Read(address))
	high := uint16(cpu.Read(address + 1))
	return high<<8 | low
}

func read16bug(address uint16) uint16 {
	a := address
	b := (a & 0xFF00) | uint16(byte(a)+1)
	low := cpu.Read(a)
	high := cpu.Read(b)
	return uint16(high)<<8 | uint16(low)
}

func adc(address *uint16) {
	a := cpu.A
	b := cpu.Read(*address)
	c := getFlag(FlagC)
	cpu.A = a + b + c

	if cpu.A == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.A>>7)

	if uint16(a)+uint16(b)+uint16(c) > 0xFF {
		setFlag(FlagC, 1)
	} else {
		setFlag(FlagC, 0)
	}

	if (a^b)&0x80 == 0 && (a^cpu.A)&0x80 != 0 {
		setFlag(FlagV, 1)
	} else {
		setFlag(FlagV, 0)
	}
}

func and(address *uint16) {
	value := cpu.Read(*address)
	cpu.A &= value

	if cpu.A == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.A>>7)
}

func asl(address *uint16) {
	var value *uint8

	if address == nil {
		value = &cpu.A
	} else {
		v := cpu.Read(*address)
		value = &v
	}

	setFlag(FlagC, (*value>>7)&1)
	*value = *value << 1

	if address != nil {
		cpu.Write(*address, *value)
	}

	if *value == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, *value>>7)

}

func bcc(address *uint16) {
	if getFlag(FlagC) == 0 {
		cpu.PC++
		cpu.PC += *address

		if *address >= 0x80 {
			cpu.PC -= 0x100
		}
		return
	}
	cpu.PC++
}

func bcs(address *uint16) {
	if getFlag(FlagC) == 1 {
		cpu.PC++
		cpu.PC += *address

		if *address >= 0x80 {
			cpu.PC -= 0x100
		}

		return
	}
	cpu.PC++
}

func beq(address *uint16) {
	if getFlag(FlagZ) == 1 {
		cpu.PC++
		cpu.PC += *address

		if *address >= 0x80 {
			cpu.PC -= 0x100
		}

		return
	}
	cpu.PC++
}

func bit(address *uint16) {
	value := cpu.Read(*address)
	setFlag(FlagN, value>>7)
	setFlag(FlagV, (value>>6)&0x1)
	if value&cpu.A == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}
}

func bmi(address *uint16) {
	if getFlag(FlagN) != 0 {
		cpu.PC++
		cpu.PC += *address

		if *address >= 0x80 {
			cpu.PC -= 0x100
		}

		return
	}
	cpu.PC++
}

func bne(address *uint16) {
	if getFlag(FlagZ) == 0 {
		cpu.PC++
		cpu.PC += *address

		if *address >= 0x80 {
			cpu.PC -= 0x100
		}
		return
	}
	cpu.PC++
}

func bpl(address *uint16) {
	if getFlag(FlagN) == 0 {
		cpu.PC++
		cpu.PC += *address

		if *address >= 0x80 {
			cpu.PC -= 0x100
		}

		return
	}
	cpu.PC++
}

func bvc(address *uint16) {
	if getFlag(FlagV) == 0 {
		cpu.PC++
		cpu.PC += *address

		if *address >= 0x80 {
			cpu.PC -= 0x100
		}
		return
	}
	cpu.PC++
}

func bvs(address *uint16) {
	if getFlag(FlagV) == 1 {
		cpu.PC++
		cpu.PC += *address

		if *address >= 0x80 {
			cpu.PC -= 0x100
		}
		return
	}
	cpu.PC++
}

func clc(*uint16) {
	setFlag(FlagC, 0)
}

func cld(*uint16) {
	setFlag(FlagD, 0)
}

func cli(*uint16) {
	setFlag(FlagI, 0)
}

func clv(*uint16) {
	setFlag(FlagV, 0)
}

func cmp(address *uint16) {
	value := uint8(*address)

	if cpu.A >= value {
		setFlag(FlagC, 1)
	} else {
		setFlag(FlagC, 0)
	}

	if cpu.A == value {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	value = cpu.A - value
	setFlag(FlagN, value>>7)
}

func cpx(address *uint16) {
	value := cpu.Read(*address)
	if cpu.X >= value {
		setFlag(FlagC, 1)
	} else {
		setFlag(FlagC, 0)
	}

	value = cpu.X - value

	if value == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, value>>7)
}

func cpy(address *uint16) {
	value := cpu.Read(*address)
	if cpu.Y >= value {
		setFlag(FlagC, 1)
	} else {
		setFlag(FlagC, 0)
	}

	value = cpu.Y - value

	if value == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, value>>7)
}

func dec(address *uint16) {
	value := cpu.Read(*address) - 1
	cpu.Write(*address, value)

	if value == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, value>>7)
}

func dex(*uint16) {
	cpu.X--
	if cpu.X == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.X>>7)
}

func dey(*uint16) {
	cpu.Y--
	if cpu.Y == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.Y>>7)
}

func eor(address *uint16) {
	cpu.A = cpu.A ^ cpu.Read(*address)
	if cpu.A == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.A>>7)
}

func inc(address *uint16) {
	value := cpu.Read(*address) + 1
	cpu.Write(*address, value)
	if value == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, value>>7)
}

func inx(*uint16) {
	cpu.X = cpu.X + 1
	if cpu.X == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.X>>7)
}

func iny(*uint16) {
	cpu.Y = cpu.Y + 1
	if cpu.Y == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.Y>>7)
}

func jmp(address *uint16) {
	cpu.PC = *address
}

func jsr(address *uint16) {
	cpu.PC--
	pushStack16(cpu.PC)
	cpu.PC = *address
}

func lda(address *uint16) {
	cpu.A = cpu.Read(*address)
	if cpu.A == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.A>>7)
}

func ldx(address *uint16) {
	cpu.X = cpu.Read(*address)
	if cpu.X == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.X>>7)
}

func ldy(address *uint16) {
	cpu.Y = cpu.Read(*address)
	if cpu.Y == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.Y>>7)
}

func lsr(address *uint16) {
	var value *uint8

	if address == nil {
		value = &cpu.A
	} else {
		v := cpu.Read(*address)
		value = &v
	}

	setFlag(FlagC, *value&0x01)
	*value = *value >> 1

	if address != nil {
		cpu.Write(*address, *value)
	}

	if *value == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, *value>>7)
}

func nop(*uint16) {
	// do nothing
}

func ora(address *uint16) {
	value := cpu.Read(*address)
	cpu.A = cpu.A | value

	if cpu.A == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.A>>7)
}

func pha(*uint16) {
	pushStack(cpu.A)
}

func php(*uint16) {
	pushStack(cpu.P | 0x10)
}

func pla(*uint16) {
	cpu.A = pullStack()
	if cpu.A == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.A>>7)
}

func plp(*uint16) {
	cpu.P = pullStack()
	cpu.P &= 0xef
	cpu.P |= 0x20
}

func rol(address *uint16) {
	var cFlag = getFlag(FlagC)

	var value *uint8
	if address == nil {
		value = &cpu.A
	} else {
		v := cpu.Read(*address)
		value = &v
	}

	setFlag(FlagC, (*value>>7)&0x1)

	*value = (*value << 1) | cFlag

	if address != nil {
		cpu.Write(*address, *value)
	}

	if *value == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, *value>>7)
}

func ror(address *uint16) {
	var cFlag = getFlag(FlagC)

	var value *uint8
	if address == nil {
		value = &cpu.A
	} else {
		v := cpu.Read(*address)
		value = &v
	}

	setFlag(FlagC, *value&0x01)
	*value = (*value >> 1) | (cFlag << 7)

	if address != nil {
		cpu.Write(*address, *value)
	}

	if *value == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, *value>>7)
}

func rts(*uint16) {
	cpu.PC = pullStack16()
	cpu.PC++
}

func sbc(address *uint16) {
	a := cpu.A
	b := cpu.Read(*address)
	c := getFlag(FlagC)
	cpu.A = a - b - (1 - c)

	if cpu.A == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.A>>7)

	if int(a)-int(b)-int(1-c) >= 0 {
		setFlag(FlagC, 1)
	} else {
		setFlag(FlagC, 0)
	}
	if (a^b)&0x80 != 0 && (a^cpu.A)&0x80 != 0 {
		setFlag(FlagV, 1)
	} else {
		setFlag(FlagV, 0)
	}
}

func sec(*uint16) {
	setFlag(FlagC, 1)
}

func sed(*uint16) {
	setFlag(FlagD, 1)
}

func sei(*uint16) {
	setFlag(FlagI, 1)
}

func sta(address *uint16) {
	cpu.Write(*address, cpu.A)
}

func stx(address *uint16) {
	cpu.Write(*address, cpu.X)
}

func sty(address *uint16) {
	cpu.Write(*address, cpu.Y)
}

func tax(*uint16) {
	cpu.X = cpu.A
	if cpu.X == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.X>>7)
}

func tay(*uint16) {
	cpu.Y = cpu.A
	if cpu.Y == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.Y>>7)
}

func tsx(*uint16) {
	cpu.X = cpu.S
	if cpu.X == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.X>>7)
}

func txa(*uint16) {
	cpu.A = cpu.X
	if cpu.A == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.A>>7)
}

func txs(*uint16) {
	cpu.S = cpu.X
}

func tya(*uint16) {
	cpu.A = cpu.Y
	if cpu.A == 0 {
		setFlag(FlagZ, 1)
	} else {
		setFlag(FlagZ, 0)
	}

	setFlag(FlagN, cpu.A>>7)
}

func brk(*uint16) {
	pushStack16(cpu.PC)
	php(nil)
	sei(nil)
	cpu.PC = read16(0xFFFE)
}

func rti(*uint16) {
	cpu.P = pullStack()
	cpu.P &= 0xef
	cpu.P |= 0x20
	cpu.PC = pullStack16()
}

func pushStack(value byte) {
	cpu.Write(0x100|uint16(cpu.S), value)
	cpu.S--
}

func pullStack() byte {
	cpu.S++
	return cpu.Read(0x100 | uint16(cpu.S))
}

func pushStack16(value uint16) {
	pushStack(byte(value >> 8))
	pushStack(byte(value & 0xFF))
}

func pullStack16() uint16 {
	low := uint16(pullStack())
	hight := uint16(pullStack())
	return hight<<8 | low
}

func setAsmOpcodes() {
	cpu.PC = 0x8000

	opcodes[0x00].GetAddress = implied
	opcodes[0x00].Instruction = brk
	opcodes[0x00].Title = "BRK (implied)"

	opcodes[0x40].GetAddress = implied
	opcodes[0x40].Instruction = rti
	opcodes[0x40].Title = "RTI (implied)"

	opcodes[0x98].GetAddress = implied
	opcodes[0x98].Instruction = tya
	opcodes[0x98].Title = "TYA (implied)"

	opcodes[0x9A].GetAddress = implied
	opcodes[0x9A].Instruction = txs
	opcodes[0x9A].Title = "TXS (implied)"

	opcodes[0x8A].GetAddress = implied
	opcodes[0x8A].Instruction = txa
	opcodes[0x8A].Title = "TXA (implied)"

	opcodes[0xBA].GetAddress = implied
	opcodes[0xBA].Instruction = tsx
	opcodes[0xBA].Title = "TSX (implied)"

	opcodes[0xA8].GetAddress = implied
	opcodes[0xA8].Instruction = tay
	opcodes[0xA8].Title = "TAY (implied)"

	opcodes[0xAA].GetAddress = implied
	opcodes[0xAA].Instruction = tax
	opcodes[0xAA].Title = "TAX (implied)"

	opcodes[0x84].GetAddress = zeroPage
	opcodes[0x84].Instruction = sty
	opcodes[0x84].Title = "STY (zeroPage)"
	opcodes[0x94].GetAddress = zeroPageX
	opcodes[0x94].Instruction = sty
	opcodes[0x94].Title = "STY (zeroPageX)"
	opcodes[0x8C].GetAddress = absolute
	opcodes[0x8C].Instruction = sty
	opcodes[0x8C].Title = "STY (absolute)"

	opcodes[0x86].GetAddress = zeroPage
	opcodes[0x86].Instruction = stx
	opcodes[0x86].Title = "STX (zeroPage)"
	opcodes[0x96].GetAddress = zeroPageY
	opcodes[0x96].Instruction = stx
	opcodes[0x96].Title = "STX (zeroPageY)"
	opcodes[0x8E].GetAddress = absolute
	opcodes[0x8E].Instruction = stx
	opcodes[0x8E].Title = "STX (absolute)"

	opcodes[0x85].GetAddress = zeroPage
	opcodes[0x85].Instruction = sta
	opcodes[0x85].Title = "STA (zeroPage)"
	opcodes[0x95].GetAddress = zeroPageX
	opcodes[0x95].Instruction = sta
	opcodes[0x95].Title = "STA (zeroPageX)"
	opcodes[0x8D].GetAddress = absolute
	opcodes[0x8D].Instruction = sta
	opcodes[0x8D].Title = "STA (absolute)"
	opcodes[0x9D].GetAddress = absoluteX
	opcodes[0x9D].Instruction = sta
	opcodes[0x9D].Title = "STA (absoluteX)"
	opcodes[0x99].GetAddress = absoluteY
	opcodes[0x99].Instruction = sta
	opcodes[0x99].Title = "STA (absoluteY)"
	opcodes[0x81].GetAddress = indirectX
	opcodes[0x81].Instruction = sta
	opcodes[0x81].Title = "STA (indirectX)"
	opcodes[0x91].GetAddress = indirectY
	opcodes[0x91].Instruction = sta
	opcodes[0x91].Title = "STA (indirectY)"

	opcodes[0x78].GetAddress = implied
	opcodes[0x78].Instruction = sei
	opcodes[0x78].Title = "SEI (implied)"

	opcodes[0xF8].GetAddress = implied
	opcodes[0xF8].Instruction = sed
	opcodes[0xF8].Title = "SED (implied)"

	opcodes[0x38].GetAddress = implied
	opcodes[0x38].Instruction = sec
	opcodes[0x38].Title = "SEC (implied)"

	opcodes[0xE9].GetAddress = immediate
	opcodes[0xE9].Instruction = sbc
	opcodes[0xE9].Title = "SBC (immediate)"
	opcodes[0xE5].GetAddress = zeroPage
	opcodes[0xE5].Instruction = sbc
	opcodes[0xE5].Title = "SBC (zeroPage)"
	opcodes[0xF5].GetAddress = zeroPageX
	opcodes[0xF5].Instruction = sbc
	opcodes[0xF5].Title = "SBC (zeroPageX)"
	opcodes[0xED].GetAddress = absolute
	opcodes[0xED].Instruction = sbc
	opcodes[0xED].Title = "SBC (absolute)"
	opcodes[0xFD].GetAddress = absoluteX
	opcodes[0xFD].Instruction = sbc
	opcodes[0xFD].Title = "SBC (absoluteX)"
	opcodes[0xF9].GetAddress = absoluteY
	opcodes[0xF9].Instruction = sbc
	opcodes[0xF9].Title = "SBC (absoluteY)"
	opcodes[0xE1].GetAddress = indirectX
	opcodes[0xE1].Instruction = sbc
	opcodes[0xE1].Title = "SBC (indirectX)"
	opcodes[0xF1].GetAddress = indirectY
	opcodes[0xF1].Instruction = sbc
	opcodes[0xF1].Title = "SBC (indirectY)"

	opcodes[0x60].GetAddress = implied
	opcodes[0x60].Instruction = rts //x
	opcodes[0x60].Title = "RTS (implied)"

	opcodes[0x6A].GetAddress = accumulator
	opcodes[0x6A].Instruction = ror
	opcodes[0x6A].Title = "ROR (accumulator)"
	opcodes[0x66].GetAddress = zeroPage
	opcodes[0x66].Instruction = ror
	opcodes[0x66].Title = "ROR (zeroPage)"
	opcodes[0x76].GetAddress = zeroPageX
	opcodes[0x76].Instruction = ror
	opcodes[0x76].Title = "ROR (zeroPageX)"
	opcodes[0x6E].GetAddress = absolute
	opcodes[0x6E].Instruction = ror
	opcodes[0x6E].Title = "ROR (absolute)"
	opcodes[0x7E].GetAddress = absoluteX
	opcodes[0x7E].Instruction = ror
	opcodes[0x7E].Title = "ROR (absoluteX)"

	opcodes[0x2A].GetAddress = accumulator
	opcodes[0x2A].Instruction = rol
	opcodes[0x2A].Title = "ROL (accumulator)"
	opcodes[0x26].GetAddress = zeroPage
	opcodes[0x26].Instruction = rol
	opcodes[0x26].Title = "ROL (zeroPage)"
	opcodes[0x36].GetAddress = zeroPageX
	opcodes[0x36].Instruction = rol
	opcodes[0x36].Title = "ROL (zeroPageX)"
	opcodes[0x2E].GetAddress = absolute
	opcodes[0x2E].Instruction = rol
	opcodes[0x2E].Title = "ROL (absolute)"
	opcodes[0x3E].GetAddress = absoluteX
	opcodes[0x3E].Instruction = rol
	opcodes[0x3E].Title = "ROL (absoluteX)"

	opcodes[0x28].GetAddress = implied
	opcodes[0x28].Instruction = plp
	opcodes[0x28].Title = "PLP (implied)"

	opcodes[0x68].GetAddress = implied
	opcodes[0x68].Instruction = pla
	opcodes[0x68].Title = "PLA (implied)"

	opcodes[0x08].GetAddress = implied
	opcodes[0x08].Instruction = php
	opcodes[0x08].Title = "PHP (implied)"

	opcodes[0x48].GetAddress = implied
	opcodes[0x48].Instruction = pha
	opcodes[0x48].Title = "PHA (implied)"

	opcodes[0x09].GetAddress = immediate
	opcodes[0x09].Instruction = ora
	opcodes[0x09].Title = "ORA (immediate)"
	opcodes[0x05].GetAddress = zeroPage
	opcodes[0x05].Instruction = ora
	opcodes[0x05].Title = "ORA (zeroPage)"
	opcodes[0x15].GetAddress = zeroPageX
	opcodes[0x15].Instruction = ora
	opcodes[0x15].Title = "ORA (zeroPageX)"
	opcodes[0x0D].GetAddress = absolute
	opcodes[0x0D].Instruction = ora
	opcodes[0x0D].Title = "ORA (absolute)"
	opcodes[0x1D].GetAddress = absoluteX
	opcodes[0x1D].Instruction = ora
	opcodes[0x1D].Title = "ORA (absoluteX)"
	opcodes[0x19].GetAddress = absoluteY
	opcodes[0x19].Instruction = ora
	opcodes[0x19].Title = "ORA (absoluteY)"
	opcodes[0x01].GetAddress = indirectX
	opcodes[0x01].Instruction = ora
	opcodes[0x01].Title = "ORA (indirectX)"
	opcodes[0x11].GetAddress = indirectY
	opcodes[0x11].Instruction = ora
	opcodes[0x11].Title = "ORA (indirectY)"

	opcodes[0xEA].GetAddress = implied
	opcodes[0xEA].Instruction = nop
	opcodes[0xEA].Title = "NOP (implied)"

	opcodes[0x4A].GetAddress = accumulator
	opcodes[0x4A].Instruction = lsr
	opcodes[0x4A].Title = "LSR (accumulator)"
	opcodes[0x46].GetAddress = zeroPage
	opcodes[0x46].Instruction = lsr
	opcodes[0x46].Title = "LSR (zeroPage)"
	opcodes[0x56].GetAddress = zeroPageX
	opcodes[0x56].Instruction = lsr
	opcodes[0x56].Title = "LSR (zeroPageX)"
	opcodes[0x4E].GetAddress = absolute
	opcodes[0x4E].Instruction = lsr
	opcodes[0x4E].Title = "LSR (absolute)"
	opcodes[0x5E].GetAddress = absoluteX
	opcodes[0x5E].Instruction = lsr
	opcodes[0x5E].Title = "LSR (absoluteX)"

	opcodes[0xA0].GetAddress = immediate
	opcodes[0xA0].Instruction = ldy //x
	opcodes[0xA0].Title = "LDY (immediate)"
	opcodes[0xA4].GetAddress = zeroPage
	opcodes[0xA4].Instruction = ldy
	opcodes[0xA4].Title = "LDY (zeroPage)"
	opcodes[0xB4].GetAddress = zeroPageX
	opcodes[0xB4].Instruction = ldy
	opcodes[0xB4].Title = "LDY (zeroPageX)"
	opcodes[0xAC].GetAddress = absolute
	opcodes[0xAC].Instruction = ldy
	opcodes[0xAC].Title = "LDY (absolute)"
	opcodes[0xBC].GetAddress = absoluteX
	opcodes[0xBC].Instruction = ldy
	opcodes[0xBC].Title = "LDY (absoluteX)"

	opcodes[0xA2].GetAddress = immediate
	opcodes[0xA2].Instruction = ldx //x
	opcodes[0xA2].Title = "LDX (immediate)"
	opcodes[0xA6].GetAddress = zeroPage
	opcodes[0xA6].Instruction = ldx
	opcodes[0xA6].Title = "LDX (zeroPage)"
	opcodes[0xB6].GetAddress = zeroPageY
	opcodes[0xB6].Instruction = ldx
	opcodes[0xB6].Title = "LDX (zeroPageY)"
	opcodes[0xAE].GetAddress = absolute
	opcodes[0xAE].Instruction = ldx
	opcodes[0xAE].Title = "LDX (absolute)"
	opcodes[0xBE].GetAddress = absoluteY
	opcodes[0xBE].Instruction = ldx
	opcodes[0xBE].Title = "LDX (absoluteY)"

	opcodes[0xA9].GetAddress = immediate
	opcodes[0xA9].Instruction = lda //x
	opcodes[0xA9].Title = "LDA (immediate)"
	opcodes[0xA5].GetAddress = zeroPage
	opcodes[0xA5].Instruction = lda
	opcodes[0xA5].Title = "LDA (zeroPage)"
	opcodes[0xB5].GetAddress = zeroPageX
	opcodes[0xB5].Instruction = lda
	opcodes[0xB5].Title = "LDA (zeroPageX)"
	opcodes[0xAD].GetAddress = absolute
	opcodes[0xAD].Instruction = lda
	opcodes[0xAD].Title = "LDA (absolute)"
	opcodes[0xBD].GetAddress = absoluteX
	opcodes[0xBD].Instruction = lda
	opcodes[0xBD].Title = "LDA (absoluteX)"
	opcodes[0xB9].GetAddress = absoluteY
	opcodes[0xB9].Instruction = lda
	opcodes[0xB9].Title = "LDA (absoluteY)"
	opcodes[0xA1].GetAddress = indirectX
	opcodes[0xA1].Instruction = lda
	opcodes[0xA1].Title = "LDA (indirectX)"
	opcodes[0xB1].GetAddress = indirectY
	opcodes[0xB1].Instruction = lda
	opcodes[0xB1].Title = "LDA (indirectY)"

	opcodes[0x20].GetAddress = absolute
	opcodes[0x20].Instruction = jsr //x
	opcodes[0x20].Title = "JSR (absolute)"

	opcodes[0x4C].GetAddress = absolute
	opcodes[0x4C].Instruction = jmp //x
	opcodes[0x4C].Title = "JMP (absolute)"
	opcodes[0x6C].GetAddress = indirect
	opcodes[0x6C].Instruction = jmp
	opcodes[0x6C].Title = "JMP (indirect)"

	opcodes[0xC8].GetAddress = implied
	opcodes[0xC8].Instruction = iny //x
	opcodes[0xC8].Title = "INY (implied)"

	opcodes[0xE8].GetAddress = implied
	opcodes[0xE8].Instruction = inx //x
	opcodes[0xE8].Title = "INX (implied)"

	opcodes[0xE6].GetAddress = zeroPage
	opcodes[0xE6].Instruction = inc //x
	opcodes[0xE6].Title = "INC (zeroPage)"
	opcodes[0xF6].GetAddress = zeroPageX
	opcodes[0xF6].Instruction = inc
	opcodes[0xF6].Title = "INC (zeroPageX)"
	opcodes[0xEE].GetAddress = absolute
	opcodes[0xEE].Instruction = inc
	opcodes[0xEE].Title = "INC (absolute)"
	opcodes[0xFE].GetAddress = absoluteX
	opcodes[0xFE].Instruction = inc
	opcodes[0xFE].Title = "INC (absoluteX)"

	opcodes[0x49].GetAddress = immediate
	opcodes[0x49].Instruction = eor
	opcodes[0x49].Title = "EOR (immediate)"
	opcodes[0x45].GetAddress = zeroPage
	opcodes[0x45].Instruction = eor
	opcodes[0x45].Title = "EOR (zeroPage)"
	opcodes[0x55].GetAddress = zeroPageX
	opcodes[0x55].Instruction = eor
	opcodes[0x55].Title = "EOR (zeroPageX)"
	opcodes[0x4D].GetAddress = absolute
	opcodes[0x4D].Instruction = eor
	opcodes[0x4D].Title = "EOR (absolute)"
	opcodes[0x5D].GetAddress = absoluteX
	opcodes[0x5D].Instruction = eor
	opcodes[0x5D].Title = "EOR (absoluteX)"
	opcodes[0x59].GetAddress = absoluteY
	opcodes[0x59].Instruction = eor
	opcodes[0x59].Title = "EOR (absoluteY)"
	opcodes[0x41].GetAddress = indirectX
	opcodes[0x41].Instruction = eor
	opcodes[0x41].Title = "EOR (indirectX)"
	opcodes[0x51].GetAddress = indirectY
	opcodes[0x51].Instruction = eor
	opcodes[0x51].Title = "EOR (indirectY)"

	opcodes[0x88].GetAddress = implied
	opcodes[0x88].Instruction = dey //x
	opcodes[0x88].Title = "DEY (implied)"

	opcodes[0xCA].GetAddress = implied
	opcodes[0xCA].Instruction = dex //x
	opcodes[0xCA].Title = "DEX (implied)"

	opcodes[0xC6].GetAddress = zeroPage
	opcodes[0xC6].Instruction = dec
	opcodes[0xC6].Title = "DEC (zeroPage)"
	opcodes[0xD6].GetAddress = zeroPageX
	opcodes[0xD6].Instruction = dec
	opcodes[0xD6].Title = "DEC (zeroPageX)"
	opcodes[0xCE].GetAddress = absolute
	opcodes[0xCE].Instruction = dec
	opcodes[0xCE].Title = "DEC (absolute)"
	opcodes[0xDE].GetAddress = absoluteX
	opcodes[0xDE].Instruction = dec
	opcodes[0xDE].Title = "DEC absoluteX)"

	opcodes[0xC0].GetAddress = immediate
	opcodes[0xC0].Instruction = cpy
	opcodes[0xC0].Title = "CPY (immediate)"
	opcodes[0xC4].GetAddress = zeroPage
	opcodes[0xC4].Instruction = cpy
	opcodes[0xC4].Title = "CPY (zeroPage)"
	opcodes[0xCC].GetAddress = absolute
	opcodes[0xCC].Instruction = cpy
	opcodes[0xCC].Title = "CPY (absolute)"

	opcodes[0xE0].GetAddress = immediate
	opcodes[0xE0].Instruction = cpx
	opcodes[0xE0].Title = "CPX (immediate)"
	opcodes[0xE4].GetAddress = zeroPage
	opcodes[0xE4].Instruction = cpx
	opcodes[0xE4].Title = "CPX (zeroPage)"
	opcodes[0xEC].GetAddress = absolute
	opcodes[0xEC].Instruction = cpx
	opcodes[0xEC].Title = "CPX (absolute)"

	opcodes[0xC9].GetAddress = immediate
	opcodes[0xC9].Instruction = cmp
	opcodes[0xC9].Title = "CMP (immediate)"
	opcodes[0xC5].GetAddress = zeroPage
	opcodes[0xC5].Instruction = cmp
	opcodes[0xC5].Title = "CMP (zeroPage)"
	opcodes[0xD5].GetAddress = zeroPageX
	opcodes[0xD5].Instruction = cmp
	opcodes[0xD5].Title = "CMP (zeroPageX)"
	opcodes[0xCD].GetAddress = absolute
	opcodes[0xCD].Instruction = cmp
	opcodes[0xCD].Title = "CMP (absolute)"
	opcodes[0xDD].GetAddress = absoluteX
	opcodes[0xDD].Instruction = cmp
	opcodes[0xDD].Title = "CMP (absoluteX)"
	opcodes[0xD9].GetAddress = absoluteY
	opcodes[0xD9].Instruction = cmp
	opcodes[0xD9].Title = "CMP (absoluteY)"
	opcodes[0xC1].GetAddress = indirectX
	opcodes[0xC1].Instruction = cmp
	opcodes[0xC1].Title = "CMP (indirectX)"
	opcodes[0xD1].GetAddress = indirectY
	opcodes[0xD1].Instruction = cmp
	opcodes[0xD1].Title = "CMP (indirectY)"

	opcodes[0xB8].GetAddress = implied
	opcodes[0xB8].Instruction = clv
	opcodes[0xB8].Title = "CLV (implied)"

	opcodes[0x58].GetAddress = implied
	opcodes[0x58].Instruction = cli
	opcodes[0x58].Title = "CLI (implied)"

	opcodes[0xD8].GetAddress = implied
	opcodes[0xD8].Instruction = cld //x
	opcodes[0xD8].Title = "CLD (implied)"

	opcodes[0x18].GetAddress = implied
	opcodes[0x18].Instruction = clc
	opcodes[0x18].Title = "CLC (implied)"

	opcodes[0x70].GetAddress = relative
	opcodes[0x70].Instruction = bvs
	opcodes[0x70].Title = "BVS (relative)"

	opcodes[0x50].GetAddress = relative
	opcodes[0x50].Instruction = bvc
	opcodes[0x50].Title = "BVC (relative)"

	opcodes[0x10].GetAddress = relative
	opcodes[0x10].Instruction = bpl //x
	opcodes[0x10].Title = "BPL (relative)"

	opcodes[0xD0].GetAddress = relative
	opcodes[0xD0].Instruction = bne //x
	opcodes[0xD0].Title = "BNE (relative)"

	opcodes[0x30].GetAddress = relative
	opcodes[0x30].Instruction = bmi //x
	opcodes[0x30].Title = "BMI (relative)"

	opcodes[0x24].GetAddress = zeroPage
	opcodes[0x24].Instruction = bit
	opcodes[0x24].Title = "BIT (zeroPage)"
	opcodes[0x2C].GetAddress = absolute
	opcodes[0x2C].Instruction = bit //x
	opcodes[0x2C].Title = "BIT (absolute)"

	opcodes[0xF0].GetAddress = relative
	opcodes[0xF0].Instruction = beq
	opcodes[0xF0].Title = "BEQ (relative)"

	opcodes[0xB0].GetAddress = relative
	opcodes[0xB0].Instruction = bcs
	opcodes[0xB0].Title = "BCS (relative)"

	opcodes[0x90].GetAddress = relative
	opcodes[0x90].Instruction = bcc
	opcodes[0x90].Title = "BCC (relative)"

	opcodes[0x0A].GetAddress = accumulator
	opcodes[0x0A].Instruction = asl
	opcodes[0x0A].Title = "ASL (accumulator)"
	opcodes[0x06].GetAddress = zeroPage
	opcodes[0x06].Instruction = asl
	opcodes[0x06].Title = "ASL (zeroPage)"
	opcodes[0x16].GetAddress = zeroPageX
	opcodes[0x16].Instruction = asl
	opcodes[0x16].Title = "ASL (zeroPageX)"
	opcodes[0x0E].GetAddress = absolute
	opcodes[0x0E].Instruction = asl
	opcodes[0x0E].Title = "ASL (absolute)"
	opcodes[0x1E].GetAddress = absoluteX
	opcodes[0x1E].Instruction = asl
	opcodes[0x1E].Title = "ASL (absoluteX)"

	opcodes[0x29].GetAddress = immediate
	opcodes[0x29].Instruction = and
	opcodes[0x29].Title = "AND (immediate)"
	opcodes[0x25].GetAddress = zeroPage
	opcodes[0x25].Instruction = and
	opcodes[0x25].Title = "AND (zeroPage)"
	opcodes[0x35].GetAddress = zeroPageX
	opcodes[0x35].Instruction = and
	opcodes[0x35].Title = "AND (zeroPageX)"
	opcodes[0x2D].GetAddress = absolute
	opcodes[0x2D].Instruction = and
	opcodes[0x2D].Title = "AND (absolute)"
	opcodes[0x3D].GetAddress = absoluteX
	opcodes[0x3D].Instruction = and
	opcodes[0x3D].Title = "AND (absoluteX)"
	opcodes[0x39].GetAddress = absoluteY
	opcodes[0x39].Instruction = and
	opcodes[0x39].Title = "AND (absoluteY)"
	opcodes[0x21].GetAddress = indirectX
	opcodes[0x21].Instruction = and
	opcodes[0x21].Title = "AND (indirectY)"
	opcodes[0x31].GetAddress = indirectY
	opcodes[0x31].Instruction = and
	opcodes[0x31].Title = "AND (indirectY)"

	opcodes[0x69].GetAddress = immediate
	opcodes[0x69].Instruction = adc
	opcodes[0x69].Title = "ADC (immediate)"
	opcodes[0x65].GetAddress = zeroPage
	opcodes[0x65].Instruction = adc
	opcodes[0x65].Title = "ADC (zeroPage)"
	opcodes[0x75].GetAddress = zeroPageX
	opcodes[0x75].Instruction = adc
	opcodes[0x75].Title = "ADC (zeroPageX)"
	opcodes[0x6D].GetAddress = absolute
	opcodes[0x6D].Instruction = adc
	opcodes[0x6D].Title = "ADC (absolute)"
	opcodes[0x7D].GetAddress = absoluteX
	opcodes[0x7D].Instruction = adc
	opcodes[0x7D].Title = "ADC (absoluteX)"
	opcodes[0x79].GetAddress = absoluteY
	opcodes[0x79].Instruction = adc
	opcodes[0x79].Title = "ADC (absoluteY)"
	opcodes[0x61].GetAddress = indirectX
	opcodes[0x61].Instruction = adc
	opcodes[0x61].Title = "ADC (indirectX)"
	opcodes[0x71].GetAddress = indirectY
	opcodes[0x71].Instruction = adc
	opcodes[0x71].Title = "ADC (indirectY)"
}
