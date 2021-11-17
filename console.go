package main

import "fmt"

type Console struct {
	ReadBuf []byte
}

func (c *Console) Read() byte {
	if len(c.ReadBuf) == 0 {
		n, _ := fmt.Scanln(c.ReadBuf)
		c.ReadBuf[n+1] = '\n'
	}
	result := c.ReadBuf[0]
	c.ReadBuf = c.ReadBuf[1:]
	return result
}

func (c Console) Write(value byte) {
	fmt.Printf("%c", rune(value))
}
