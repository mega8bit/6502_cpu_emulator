.SEGMENT "RESET"
.WORD $8000


.SEGMENT "CODE"
LDA Symbol

DrawRectangle:
    LDX #$00
    DrawUp:
        STA $2000
        INX
        CPX LineLen
        BNE DrawUp
    STA $2000
    LDA EOL
    STA $2000

    LDA Symbol
    LDY #$0A
    LDX #$00
    DrawLeftRight:
        CPY #$00
        BEQ Bottom
        CPX #$00
        BEQ DrawLeftRight_ColLeft
        CPX LineLen
        BEQ DrawLeftRight_ColRight
        JMP DrawLeftRight_Space

        DrawLeftRight_ColLeft:
            LDA Symbol
            STA $2000
            INX
            JMP DrawLeftRight

        DrawLeftRight_ColRight:
            LDA Symbol
            STA $2000
            LDA EOL
            STA $2000
            LDX #$00
            DEY
            JMP DrawLeftRight

        DrawLeftRight_Space:
            LDA Space
            STA $2000
            INX
            JMP DrawLeftRight


    Bottom:
        LDA Symbol
        LDX #$00
    DrawBottom:
        STA $2000
        INX
        CPX LineLen
        BNE DrawBottom

    STA $2000
    LDA EOL
    STA $2000
    STA $2000
    
NOP

LineLen:
.BYTE $8
Symbol:
.BYTE "#"
EOL:
.BYTE $0A
Space:
.BYTE " "