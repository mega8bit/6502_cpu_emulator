package main

type CPU struct {
	PC uint16
	P  byte
	S  byte
	A  byte
	X  byte
	Y  byte

	Ram     []byte
	Rom     []byte
	console Console
}

func (c *CPU) Read(address uint16) byte {
	if address < 0x2000 {
		return c.Ram[address]
	}

	if address == 0x2000 {
		return c.console.Read()
	}

	return c.Rom[address%0x8000]
}

func (c *CPU) Write(address uint16, value byte) {
	if address < 0x2000 {
		c.Ram[address] = value
		return
	}

	if address == 0x2000 {
		c.console.Write(value)
		return
	}
}

func (c *CPU) Reset() {
	c.PC = read16(0xFFFC)
}

func (c *CPU) SetConsole(console Console) {
	c.console = console
}

func (c *CPU) Step() bool {
	opcodeNum := cpu.Read(cpu.PC)
	if opcodeNum == 0xEA {
		return true
	}

	op := opcodes[opcodeNum]

	if op.Instruction != nil {
		op.Instruction(op.GetAddress())
	}

	return false
}
