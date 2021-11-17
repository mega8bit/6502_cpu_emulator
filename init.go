package main

func init() {
	cpu = new(CPU)
	cpu.PC = 0x0
	cpu.A = 0x0
	cpu.X = 0x0
	cpu.Y = 0x0
	cpu.Ram = make([]byte, 0x0800)

	setAsmOpcodes()

	console := new(Console)
	console.ReadBuf = make([]byte, 1024)
	cpu.console = *console
}
