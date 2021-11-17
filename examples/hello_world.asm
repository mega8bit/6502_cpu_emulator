.SEGMENT "RESET"
.WORD $8000


.SEGMENT "CODE"
LDY $00
LDA Message, Y

Loop:
STA $2000
INY
LDA Message, Y
CMP $0A
BNE Loop

STA $2000
NOP

Message:
.BYTE "Hello World!", $0A