package vp8

type VP8D_CONFIG struct {
	Width             int
	Height            int
	Version           int
	Postprocess       int
	Max_threads       int
	Error_concealment int
}
type VP8D_SETTING int

const VP8D_OK VP8D_SETTING = 0
