module github.com/mearaj/libvpx

go 1.17

replace (
	github.com/mearaj/libvpx => ./
	github.com/mearaj/libvpx/dsp => ./internal/dsp
	github.com/mearaj/libvpx/internal => ./internal
	github.com/mearaj/libvpx/mem => ./internal/mem
	github.com/mearaj/libvpx/ports => ./internal/ports
	github.com/mearaj/libvpx/scale => ./internal/scale
	github.com/mearaj/libvpx/util => ./internal/util
	github.com/mearaj/libvpx/vp8 => ./internal/vp8
	github.com/mearaj/libvpx/vp9 => ./internal/vp9
	github.com/mearaj/libvpx/vpx => ./internal/vpx
)

require github.com/gotranspile/cxgo v0.3.3-0.20220306160746-757678937454

require maze.io/x/math32 v0.0.0-20181106113604-c78ed91899f1 // indirect
