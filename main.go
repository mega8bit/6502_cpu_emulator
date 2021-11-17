package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong input parameters for emulator")
		return
	}

	programPath := os.Args[1]
	fileData, err := ioutil.ReadFile(programPath)

	if err != nil {
		fmt.Println("Cannot read program data", err)
		return
	}

	cpu.Rom = fileData
	cpu.Reset()

	var isExit bool
	for !isExit {
		isExit = cpu.Step()
	}

	fmt.Println("Program has been executed")
}
