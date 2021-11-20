.DEFINE EOL $0A
.SEGMENT "RESET"
.WORD $8000


.SEGMENT "CODE"
LDY #$00
JSR PrintMessageHello
LDY #$00
JSR ReadName
LDY #$00
JSR PrintMessageDear
NOP

ReleaseBuffer:
    LDA InputBuffer, Y
    STA $2000
    INY
    CMP EOL
    BNE ReleaseBuffer
    STA $2000
    RTS


ReadName:
    LDA $2000
    STA InputBuffer, Y
    INY
    CMP EOL
    BNE ReadName
    RTS

PrintMessageHello:
    LDA HelloMessage, Y
    INY
    STA $2000
    CMP EOL
    BNE PrintMessageHello
    RTS

PrintMessageDear:
    LDA Welcome, Y
    INY
    STA $2000
    CMP EOL
    BNE PrintMessageDear

    LDY #$00
    JSR ReleaseBuffer
    RTS

HelloMessage:
.BYTE "Hello, What is your name?", EOL
Welcome:
.BYTE "Welcome dear ", EOL
InputBuffer:
.RES 128
