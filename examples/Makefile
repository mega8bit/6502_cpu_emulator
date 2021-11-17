hello_world:
	ca65 --cpu 6502 hello_world.asm
	ld65 -C linker.ld --obj hello_world.o -o rom.bin

clean:
	rm -rf *.o
	rm -rf *.bin
	rm -rf *.out