MEMORY
{
    ROM:
        start = $8000, size = $fffc - $8000, fill = yes;
    RESET:
        start = $fffc, size = $2;
}
SEGMENTS
{
        RESET: load = RESET, type = ro;
        CODE: load = ROM, type = ro;
}