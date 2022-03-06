package vp8

type MV struct {
	Row int16
	Col int16
}
type int_mv struct {
	// union
	As_int uint32
	As_mv  MV
}
