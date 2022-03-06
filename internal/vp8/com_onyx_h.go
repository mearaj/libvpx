package vp8

import (
	"github.com/mearaj/libvpx/internal/vpx"
	"log"
)

const (
	NORMAL    int = 0
	FOURFIVE  int = 1
	THREEFIVE int = 2
	ONETWO    int = 3
)
const (
	USAGE_LOCAL_FILE_PLAYBACK int = 0
	USAGE_STREAM_FROM_SERVER  int = 1
	USAGE_CONSTRAINED_QUALITY int = 2
	USAGE_CONSTANT_QUALITY    int = 3
)
const (
	MODE_REALTIME        int = 0
	MODE_GOODQUALITY     int = 1
	MODE_BESTQUALITY     int = 2
	MODE_FIRSTPASS       int = 3
	MODE_SECONDPASS      int = 4
	MODE_SECONDPASS_BEST int = 5
)
const (
	FRAMEFLAGS_KEY    int = 1
	FRAMEFLAGS_GOLDEN int = 2
	FRAMEFLAGS_ALTREF int = 4
)

func Scale2Ratio(mode int, hr *int, hs *int) {
	switch mode {
	case NORMAL:
		*hr = 1
		*hs = 1
	case FOURFIVE:
		*hr = 4
		*hs = 5
	case THREEFIVE:
		*hr = 3
		*hs = 5
	case ONETWO:
		*hr = 1
		*hs = 2
	default:
		*hr = 1
		*hs = 1
		if false {
		} else {
			// Todo:
			log.Fatal("error")

		}
	}
}

type VP8_CONFIG struct {
	Version                     int
	Width                       int
	Height                      int
	Timebase                    vpx.Rational
	Target_bandwidth            uint
	Noise_sensitivity           int
	Sharpness                   int
	Cpu_used                    int
	Rc_max_intra_bitrate_pct    uint
	Gf_cbr_boost_pct            uint
	Screen_content_mode         uint
	Mode                        int
	Auto_key                    int
	Key_freq                    int
	Allow_lag                   int
	Lag_in_frames               int
	End_usage                   int
	Under_shoot_pct             int
	Over_shoot_pct              int
	Starting_buffer_level       int64
	Optimal_buffer_level        int64
	Maximum_buffer_size         int64
	Starting_buffer_level_in_ms int64
	Optimal_buffer_level_in_ms  int64
	Maximum_buffer_size_in_ms   int64
	Fixed_q                     int
	Worst_allowed_q             int
	Best_allowed_q              int
	Cq_level                    int
	Allow_spatial_resampling    int
	Resample_down_water_mark    int
	Resample_up_water_mark      int
	Allow_df                    int
	Drop_frames_water_mark      int
	Two_pass_vbrbias            int
	Two_pass_vbrmin_section     int
	Two_pass_vbrmax_section     int
	Play_alternate              int
	Alt_freq                    int
	Alt_q                       int
	Key_q                       int
	Gold_q                      int
	Multi_threaded              int
	Token_partitions            int
	Encode_breakout             int
	Error_resilient_mode        uint
	Arnr_max_frames             int
	Arnr_strength               int
	Arnr_type                   int
	Two_pass_stats_in           vpx.FixedBuf
	Output_pkt_list             *vpx.CodecPktList
	Tuning                      vpx.Vp8eTuning
	Number_of_layers            uint
	Target_bitrate              [16]uint
	Rate_decimator              [16]uint
	Periodicity                 uint
	Layer_id                    [16]uint
}
