all:
	GOOS=linux go build -o 6502em_linux
	GOOS=windows go build -o 6502em.exe
	GOOS=darwin go build -o 6502em_macos

linux:
	GOOS=linux go build -o 6502em_linux
windows:
	GOOS=windows go build -o 6502em.exe
macos:
	GOOS=darwin go build -o 6502em_macos


clean:
	rm -rf 6502em_linux
	rm -rf 6502em.exe
	rm -rf 6502em_macos