package vp8

import (
	"github.com/mearaj/libvpx/internal/scale"
	"github.com/mearaj/libvpx/internal/util"
	"github.com/mearaj/libvpx/internal/vpx"
	"unsafe"
)

const MAX_FB_MT_DEC = 32

type DECODETHREAD_DATA struct {
	Ithread int
	Ptr1    unsafe.Pointer
	Ptr2    unsafe.Pointer
}
type MB_ROW_DEC struct {
	Mbd MacroBlockd
}
type FRAGMENT_DATA struct {
	Enabled int
	Count   uint
	Ptrs    [9]*uint8
	Sizes   [9]uint
}
type frame_buffers struct {
	Pbi [32]*VP8D_COMP
}
type VP8D_COMP struct {
	Mb                              MacroBlockd
	Dec_fb_ref                      [4]*scale.Yv12BufferConfig
	Common                          VP8Common
	Mbc                             [9]vp8_reader
	Oxcf                            VP8D_CONFIG
	Fragments                       FRAGMENT_DATA
	B_multithreaded_rd              util.VpxAtomicInt
	Max_threads                     int
	Current_mb_col_main             int
	Decoding_thread_count           uint
	Allocated_decoding_thread_count int
	Mt_baseline_filter_level        [4]int
	Sync_range                      int
	Mt_current_mb_col               *util.VpxAtomicInt
	Mt_yabove_row                   **uint8
	Mt_uabove_row                   **uint8
	Mt_vabove_row                   **uint8
	Mt_yleft_col                    **uint8
	Mt_uleft_col                    **uint8
	Mt_vleft_col                    **uint8
	Mb_row_di                       *MB_ROW_DEC
	De_thread_data                  *DECODETHREAD_DATA
	H_decoding_thread               *pthread_t
	H_event_start_decoding          *sem_t
	H_event_end_decoding            sem_t
	Last_time_stamp                 int64
	Ready_for_new_data              int
	Prob_intra                      uint8
	Prob_last                       uint8
	Prob_gf                         uint8
	Prob_skip_false                 uint8
	Ec_enabled                      int
	Ec_active                       int
	Decoded_key_frame               int
	Independent_partitions          int
	Frame_corrupt_residual          int
	Decrypt_cb                      vpx.DecryptCb
	Decrypt_state                   unsafe.Pointer
	Restart_threads                 int
}
