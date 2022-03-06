package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	_ "github.com/mearaj/libvpx/internal/dsp"
	_ "github.com/mearaj/libvpx/internal/mem"
	_ "github.com/mearaj/libvpx/internal/ports"
	_ "github.com/mearaj/libvpx/internal/scale"
	_ "github.com/mearaj/libvpx/internal/util"
	"github.com/mearaj/libvpx/internal/vpx"
	_ "log"
	"unsafe"
)

const VP8_LAST_FRAME = int(vpx.VP8_LAST_FRAME)
const VP8_GOLD_FRAME = int(vpx.VP8_GOLD_FRAME)
const VP8_ALTR_FRAME = int(vpx.VP8_ALTR_FRAME)
const VPX_CODEC_CAP_POSTPROC = vpx.VPX_CODEC_CAP_POSTPROC
const VPX_CODEC_CAP_DECODER = vpx.VPX_CODEC_CAP_DECODER
const VPX_CODEC_USE_INPUT_FRAGMENTS = vpx.VPX_CODEC_USE_INPUT_FRAGMENTS
const VPX_CODEC_ERROR = int(vpx.VPX_CODEC_ERROR)
const VPX_CODEC_OK = int(vpx.VPX_CODEC_OK)
const VP8_DEBLOCK = int(vpx.VP8_DEBLOCK)
const VP8_DEMACROBLOCK = int(vpx.VP8_DEMACROBLOCK)
const VP8_SET_REFERENCE = int(vpx.VP8_SET_REFERENCE)
const VP8_COPY_REFERENCE = int(vpx.VP8_COPY_REFERENCE)
const VP8D_GET_LAST_REF_UPDATES = int(vpx.VP8D_GET_LAST_REF_UPDATES)
const VP8D_GET_FRAME_CORRUPTED = int(vpx.VP8D_GET_FRAME_CORRUPTED)
const VP8D_GET_LAST_REF_USED = int(vpx.VP8D_GET_LAST_REF_USED)
const VPXD_GET_LAST_QUANTIZER = int(vpx.VPXD_GET_LAST_QUANTIZER)
const VP8_SET_POSTPROC = int(vpx.VP8_SET_POSTPROC)
const VPXD_SET_DECRYPTOR = int(vpx.VPXD_SET_DECRYPTOR)
const VP8_MFQE = int(vpx.VP8_MFQE)
const VPX_CODEC_USE_ERROR_CONCEALMENT = vpx.VPX_CODEC_USE_ERROR_CONCEALMENT
const VPX_CODEC_CAP_ERROR_CONCEALMENT = vpx.VPX_CODEC_CAP_ERROR_CONCEALMENT
const VPX_CODEC_CAP_INPUT_FRAGMENTS = vpx.VPX_CODEC_CAP_INPUT_FRAGMENTS

// Todo: from file vp8/decoder/dboolhuff.h
const CHAR_BIT = 8

var sub_mv_ref_prob = []uint8{180, 162, 25}
var vp8_sub_mv_ref_prob2 = [][]uint8{{147, 136, 18}, {106, 145, 1}, {179, 121, 1}, {223, 1, 34}, {208, 1, 1}}
var vp8_mbsplits = []vp8_mbsplit{{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1}, {0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}, {0, 0, 1, 1, 0, 0, 1, 1, 2, 2, 3, 3, 2, 2, 3, 3}, {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}}
var vp8_mbsplit_count = []int{2, 2, 4, 16}
var vp8_mbsplit_probs = []uint8{110, 111, 150}
var vp8_bmode_tree = []int8{int8(-B_DC_PRED), 2, int8(-B_TM_PRED), 4, int8(-B_VE_PRED), 6, 8, 12, int8(-B_HE_PRED), 10, int8(-B_RD_PRED), int8(-B_VR_PRED), int8(-B_LD_PRED), 14, int8(-B_VL_PRED), 16, int8(-B_HD_PRED), int8(-B_HU_PRED)}
var vp8_ymode_tree = []int8{int8(-DC_PRED), 2, 4, 6, int8(-V_PRED), int8(-H_PRED), int8(-TM_PRED), int8(-B_PRED)}
var vp8_kf_ymode_tree = []int8{int8(-B_PRED), 2, 4, 6, int8(-DC_PRED), int8(-V_PRED), int8(-H_PRED), int8(-TM_PRED)}
var vp8_uv_mode_tree = []int8{int8(-DC_PRED), 2, int8(-V_PRED), 4, int8(-H_PRED), int8(-TM_PRED)}
var vp8_mbsplit_tree = []int8{-3, 2, -2, 4, 0, -1}
var vp8_mv_ref_tree = []int8{int8(-ZEROMV), 2, int8(-NEARESTMV), 4, int8(-NEARMV), 6, int8(-NEWMV), int8(-SPLITMV)}
var vp8_sub_mv_ref_tree = []int8{int8(-LEFT4X4), 2, int8(-ABOVE4X4), 4, int8(-ZERO4X4), int8(-NEW4X4)}
var vp8_small_mvtree = []int8{2, 8, 4, 6, 0, -1, -2, -3, 10, 12, -4, -5, -6, -7}

// Todo: vp8/common/coefupdateprobs.h
var default_coef_probs = [BLOCK_TYPES][COEF_BANDS][PREV_COEF_CONTEXTS][ENTROPY_NODES]uint8{{ /* Block Type ( 0 ) */
	{ /* Coeff Band ( 0 )*/
		{128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128},
		{128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128},
		{128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128}},
	{ /* Coeff Band ( 1 )*/
		{253, 136, 254, 255, 228, 219, 128, 128, 128, 128, 128},
		{189, 129, 242, 255, 227, 213, 255, 219, 128, 128, 128},
		{106, 126, 227, 252, 214, 209, 255, 255, 128, 128, 128}},
	{ /* Coeff Band ( 2 )*/
		{1, 98, 248, 255, 236, 226, 255, 255, 128, 128, 128},
		{181, 133, 238, 254, 221, 234, 255, 154, 128, 128, 128},
		{78, 134, 202, 247, 198, 180, 255, 219, 128, 128, 128}},
	{ /* Coeff Band ( 3 )*/
		{1, 185, 249, 255, 243, 255, 128, 128, 128, 128, 128},
		{184, 150, 247, 255, 236, 224, 128, 128, 128, 128, 128},
		{77, 110, 216, 255, 236, 230, 128, 128, 128, 128, 128}},
	{ /* Coeff Band ( 4 )*/
		{1, 101, 251, 255, 241, 255, 128, 128, 128, 128, 128},
		{170, 139, 241, 252, 236, 209, 255, 255, 128, 128, 128},
		{37, 116, 196, 243, 228, 255, 255, 255, 128, 128, 128}},
	{ /* Coeff Band ( 5 )*/
		{1, 204, 254, 255, 245, 255, 128, 128, 128, 128, 128},
		{207, 160, 250, 255, 238, 128, 128, 128, 128, 128, 128},
		{102, 103, 231, 255, 211, 171, 128, 128, 128, 128, 128}},
	{ /* Coeff Band ( 6 )*/
		{1, 152, 252, 255, 240, 255, 128, 128, 128, 128, 128},
		{177, 135, 243, 255, 234, 225, 128, 128, 128, 128, 128},
		{80, 129, 211, 255, 194, 224, 128, 128, 128, 128, 128}},
	{ /* Coeff Band ( 7 )*/
		{1, 1, 255, 128, 128, 128, 128, 128, 128, 128, 128},
		{246, 1, 255, 128, 128, 128, 128, 128, 128, 128, 128},
		{255, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128}}},
	{ /* Block Type ( 1 ) */
		{ /* Coeff Band ( 0 )*/
			{198, 35, 237, 223, 193, 187, 162, 160, 145, 155, 62},
			{131, 45, 198, 221, 172, 176, 220, 157, 252, 221, 1},
			{68, 47, 146, 208, 149, 167, 221, 162, 255, 223, 128}},
		{ /* Coeff Band ( 1 )*/
			{1, 149, 241, 255, 221, 224, 255, 255, 128, 128, 128},
			{184, 141, 234, 253, 222, 220, 255, 199, 128, 128, 128},
			{81, 99, 181, 242, 176, 190, 249, 202, 255, 255, 128}},
		{ /* Coeff Band ( 2 )*/
			{1, 129, 232, 253, 214, 197, 242, 196, 255, 255, 128},
			{99, 121, 210, 250, 201, 198, 255, 202, 128, 128, 128},
			{23, 91, 163, 242, 170, 187, 247, 210, 255, 255, 128}},
		{ /* Coeff Band ( 3 )*/
			{1, 200, 246, 255, 234, 255, 128, 128, 128, 128, 128},
			{109, 178, 241, 255, 231, 245, 255, 255, 128, 128, 128},
			{44, 130, 201, 253, 205, 192, 255, 255, 128, 128, 128}},
		{ /* Coeff Band ( 4 )*/
			{1, 132, 239, 251, 219, 209, 255, 165, 128, 128, 128},
			{94, 136, 225, 251, 218, 190, 255, 255, 128, 128, 128},
			{22, 100, 174, 245, 186, 161, 255, 199, 128, 128, 128}},
		{ /* Coeff Band ( 5 )*/
			{1, 182, 249, 255, 232, 235, 128, 128, 128, 128, 128},
			{124, 143, 241, 255, 227, 234, 128, 128, 128, 128, 128},
			{35, 77, 181, 251, 193, 211, 255, 205, 128, 128, 128}},
		{ /* Coeff Band ( 6 )*/
			{1, 157, 247, 255, 236, 231, 255, 255, 128, 128, 128},
			{121, 141, 235, 255, 225, 227, 255, 255, 128, 128, 128},
			{45, 99, 188, 251, 195, 217, 255, 224, 128, 128, 128}},
		{ /* Coeff Band ( 7 )*/
			{1, 1, 251, 255, 213, 255, 128, 128, 128, 128, 128},
			{203, 1, 248, 255, 255, 128, 128, 128, 128, 128, 128},
			{137, 1, 177, 255, 224, 255, 128, 128, 128, 128, 128}}},
	{ /* Block Type ( 2 ) */
		{ /* Coeff Band ( 0 )*/
			{253, 9, 248, 251, 207, 208, 255, 192, 128, 128, 128},
			{175, 13, 224, 243, 193, 185, 249, 198, 255, 255, 128},
			{73, 17, 171, 221, 161, 179, 236, 167, 255, 234, 128}},
		{ /* Coeff Band ( 1 )*/
			{1, 95, 247, 253, 212, 183, 255, 255, 128, 128, 128},
			{239, 90, 244, 250, 211, 209, 255, 255, 128, 128, 128},
			{155, 77, 195, 248, 188, 195, 255, 255, 128, 128, 128}},
		{ /* Coeff Band ( 2 )*/
			{1, 24, 239, 251, 218, 219, 255, 205, 128, 128, 128},
			{201, 51, 219, 255, 196, 186, 128, 128, 128, 128, 128},
			{69, 46, 190, 239, 201, 218, 255, 228, 128, 128, 128}},
		{ /* Coeff Band ( 3 )*/
			{1, 191, 251, 255, 255, 128, 128, 128, 128, 128, 128},
			{223, 165, 249, 255, 213, 255, 128, 128, 128, 128, 128},
			{141, 124, 248, 255, 255, 128, 128, 128, 128, 128, 128}},
		{ /* Coeff Band ( 4 )*/
			{1, 16, 248, 255, 255, 128, 128, 128, 128, 128, 128},
			{190, 36, 230, 255, 236, 255, 128, 128, 128, 128, 128},
			{149, 1, 255, 128, 128, 128, 128, 128, 128, 128, 128}},
		{ /* Coeff Band ( 5 )*/
			{1, 226, 255, 128, 128, 128, 128, 128, 128, 128, 128},
			{247, 192, 255, 128, 128, 128, 128, 128, 128, 128, 128},
			{240, 128, 255, 128, 128, 128, 128, 128, 128, 128, 128}},
		{ /* Coeff Band ( 6 )*/
			{1, 134, 252, 255, 255, 128, 128, 128, 128, 128, 128},
			{213, 62, 250, 255, 255, 128, 128, 128, 128, 128, 128},
			{55, 93, 255, 128, 128, 128, 128, 128, 128, 128, 128}},
		{ /* Coeff Band ( 7 )*/
			{128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128},
			{128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128},
			{128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128}}},
	{ /* Block Type ( 3 ) */
		{ /* Coeff Band ( 0 )*/
			{202, 24, 213, 235, 186, 191, 220, 160, 240, 175, 255},
			{126, 38, 182, 232, 169, 184, 228, 174, 255, 187, 128},
			{61, 46, 138, 219, 151, 178, 240, 170, 255, 216, 128}},
		{ /* Coeff Band ( 1 )*/
			{1, 112, 230, 250, 199, 191, 247, 159, 255, 255, 128},
			{166, 109, 228, 252, 211, 215, 255, 174, 128, 128, 128},
			{39, 77, 162, 232, 172, 180, 245, 178, 255, 255, 128}},
		{ /* Coeff Band ( 2 )*/
			{1, 52, 220, 246, 198, 199, 249, 220, 255, 255, 128},
			{124, 74, 191, 243, 183, 193, 250, 221, 255, 255, 128},
			{24, 71, 130, 219, 154, 170, 243, 182, 255, 255, 128}},
		{ /* Coeff Band ( 3 )*/
			{1, 182, 225, 249, 219, 240, 255, 224, 128, 128, 128},
			{149, 150, 226, 252, 216, 205, 255, 171, 128, 128, 128},
			{28, 108, 170, 242, 183, 194, 254, 223, 255, 255, 128}},
		{ /* Coeff Band ( 4 )*/
			{1, 81, 230, 252, 204, 203, 255, 192, 128, 128, 128},
			{123, 102, 209, 247, 188, 196, 255, 233, 128, 128, 128},
			{20, 95, 153, 243, 164, 173, 255, 203, 128, 128, 128}},
		{ /* Coeff Band ( 5 )*/
			{1, 222, 248, 255, 216, 213, 128, 128, 128, 128, 128},
			{168, 175, 246, 252, 235, 205, 255, 255, 128, 128, 128},
			{47, 116, 215, 255, 211, 212, 255, 255, 128, 128, 128}},
		{ /* Coeff Band ( 6 )*/
			{1, 121, 236, 253, 212, 214, 255, 255, 128, 128, 128},
			{141, 84, 213, 252, 201, 202, 255, 219, 128, 128, 128},
			{42, 80, 160, 240, 162, 185, 255, 205, 128, 128, 128}},
		{ /* Coeff Band ( 7 )*/
			{1, 1, 255, 128, 128, 128, 128, 128, 128, 128, 128},
			{244, 1, 255, 128, 128, 128, 128, 128, 128, 128, 128},
			{238, 1, 255, 128, 128, 128, 128, 128, 128, 128, 128}}},
}

// Todo: used in dec_decodemv.go defined in vp8_entropymodedata.h
var vp8_kf_ymode_prob = [VP8_YMODES - 1]uint8{145, 156, 163, 128}

var vp8_kf_uv_mode_prob = [VP8_UV_MODES - 1]uint8{142, 114, 183}

// Todo: used in dec_decodemv.go defined in vp8_entropymodedata.h
var vp8_kf_bmode_prob = [VP8_BINTRAMODES][VP8_BINTRAMODES][VP8_BINTRAMODES - 1]uint8{
	{{231, 120, 48, 89, 115, 113, 120, 152, 112},
		{152, 179, 64, 126, 170, 118, 46, 70, 95},
		{175, 69, 143, 80, 85, 82, 72, 155, 103},
		{56, 58, 10, 171, 218, 189, 17, 13, 152},
		{144, 71, 10, 38, 171, 213, 144, 34, 26},
		{114, 26, 17, 163, 44, 195, 21, 10, 173},
		{121, 24, 80, 195, 26, 62, 44, 64, 85},
		{170, 46, 55, 19, 136, 160, 33, 206, 71},
		{63, 20, 8, 114, 114, 208, 12, 9, 226},
		{81, 40, 11, 96, 182, 84, 29, 16, 36}},
	{{134, 183, 89, 137, 98, 101, 106, 165, 148},
		{72, 187, 100, 130, 157, 111, 32, 75, 80},
		{66, 102, 167, 99, 74, 62, 40, 234, 128},
		{41, 53, 9, 178, 241, 141, 26, 8, 107},
		{104, 79, 12, 27, 217, 255, 87, 17, 7},
		{74, 43, 26, 146, 73, 166, 49, 23, 157},
		{65, 38, 105, 160, 51, 52, 31, 115, 128},
		{87, 68, 71, 44, 114, 51, 15, 186, 23},
		{47, 41, 14, 110, 182, 183, 21, 17, 194},
		{66, 45, 25, 102, 197, 189, 23, 18, 22}},
	{{88, 88, 147, 150, 42, 46, 45, 196, 205},
		{43, 97, 183, 117, 85, 38, 35, 179, 61},
		{39, 53, 200, 87, 26, 21, 43, 232, 171},
		{56, 34, 51, 104, 114, 102, 29, 93, 77},
		{107, 54, 32, 26, 51, 1, 81, 43, 31},
		{39, 28, 85, 171, 58, 165, 90, 98, 64},
		{34, 22, 116, 206, 23, 34, 43, 166, 73},
		{68, 25, 106, 22, 64, 171, 36, 225, 114},
		{34, 19, 21, 102, 132, 188, 16, 76, 124},
		{62, 18, 78, 95, 85, 57, 50, 48, 51}},
	{{193, 101, 35, 159, 215, 111, 89, 46, 111},
		{60, 148, 31, 172, 219, 228, 21, 18, 111},
		{112, 113, 77, 85, 179, 255, 38, 120, 114},
		{40, 42, 1, 196, 245, 209, 10, 25, 109},
		{100, 80, 8, 43, 154, 1, 51, 26, 71},
		{88, 43, 29, 140, 166, 213, 37, 43, 154},
		{61, 63, 30, 155, 67, 45, 68, 1, 209},
		{142, 78, 78, 16, 255, 128, 34, 197, 171},
		{41, 40, 5, 102, 211, 183, 4, 1, 221},
		{51, 50, 17, 168, 209, 192, 23, 25, 82}},
	{{125, 98, 42, 88, 104, 85, 117, 175, 82},
		{95, 84, 53, 89, 128, 100, 113, 101, 45},
		{75, 79, 123, 47, 51, 128, 81, 171, 1},
		{57, 17, 5, 71, 102, 57, 53, 41, 49},
		{115, 21, 2, 10, 102, 255, 166, 23, 6},
		{38, 33, 13, 121, 57, 73, 26, 1, 85},
		{41, 10, 67, 138, 77, 110, 90, 47, 114},
		{101, 29, 16, 10, 85, 128, 101, 196, 26},
		{57, 18, 10, 102, 102, 213, 34, 20, 43},
		{117, 20, 15, 36, 163, 128, 68, 1, 26}},
	{{138, 31, 36, 171, 27, 166, 38, 44, 229},
		{67, 87, 58, 169, 82, 115, 26, 59, 179},
		{63, 59, 90, 180, 59, 166, 93, 73, 154},
		{40, 40, 21, 116, 143, 209, 34, 39, 175},
		{57, 46, 22, 24, 128, 1, 54, 17, 37},
		{47, 15, 16, 183, 34, 223, 49, 45, 183},
		{46, 17, 33, 183, 6, 98, 15, 32, 183},
		{65, 32, 73, 115, 28, 128, 23, 128, 205},
		{40, 3, 9, 115, 51, 192, 18, 6, 223},
		{87, 37, 9, 115, 59, 77, 64, 21, 47}},
	{{104, 55, 44, 218, 9, 54, 53, 130, 226},
		{64, 90, 70, 205, 40, 41, 23, 26, 57},
		{54, 57, 112, 184, 5, 41, 38, 166, 213},
		{30, 34, 26, 133, 152, 116, 10, 32, 134},
		{75, 32, 12, 51, 192, 255, 160, 43, 51},
		{39, 19, 53, 221, 26, 114, 32, 73, 255},
		{31, 9, 65, 234, 2, 15, 1, 118, 73},
		{88, 31, 35, 67, 102, 85, 55, 186, 85},
		{56, 21, 23, 111, 59, 205, 45, 37, 192},
		{55, 38, 70, 124, 73, 102, 1, 34, 98}},
	{{102, 61, 71, 37, 34, 53, 31, 243, 192},
		{69, 60, 71, 38, 73, 119, 28, 222, 37},
		{68, 45, 128, 34, 1, 47, 11, 245, 171},
		{62, 17, 19, 70, 146, 85, 55, 62, 70},
		{75, 15, 9, 9, 64, 255, 184, 119, 16},
		{37, 43, 37, 154, 100, 163, 85, 160, 1},
		{63, 9, 92, 136, 28, 64, 32, 201, 85},
		{86, 6, 28, 5, 64, 255, 25, 248, 1},
		{56, 8, 17, 132, 137, 255, 55, 116, 128},
		{58, 15, 20, 82, 135, 57, 26, 121, 40}},
	{{164, 50, 31, 137, 154, 133, 25, 35, 218},
		{51, 103, 44, 131, 131, 123, 31, 6, 158},
		{86, 40, 64, 135, 148, 224, 45, 183, 128},
		{22, 26, 17, 131, 240, 154, 14, 1, 209},
		{83, 12, 13, 54, 192, 255, 68, 47, 28},
		{45, 16, 21, 91, 64, 222, 7, 1, 197},
		{56, 21, 39, 155, 60, 138, 23, 102, 213},
		{85, 26, 85, 85, 128, 128, 32, 146, 171},
		{18, 11, 7, 63, 144, 171, 4, 4, 246},
		{35, 27, 10, 146, 174, 171, 12, 26, 128}},
	{{190, 80, 35, 99, 180, 80, 126, 54, 45},
		{85, 126, 47, 87, 176, 51, 41, 20, 32},
		{101, 75, 128, 139, 118, 146, 116, 128, 85},
		{56, 41, 15, 176, 236, 85, 37, 9, 62},
		{146, 36, 19, 30, 171, 255, 97, 27, 20},
		{71, 30, 17, 119, 118, 255, 17, 18, 138},
		{101, 38, 60, 138, 55, 70, 43, 26, 142},
		{138, 45, 61, 62, 219, 1, 81, 188, 64},
		{32, 41, 20, 117, 151, 142, 20, 21, 163},
		{112, 19, 12, 61, 195, 128, 48, 4, 24}},
}

// Todo: defined in dx_iface.go
var vp8CtfMaps = [9]vpx.FnMap{{Ctrl_id: int(vpx.VP8_SET_REFERENCE), Fn: func(ctx *vpx.CodecAlgPvt, ap libc.ArgList) vpx.CodecErr {
	return vp8_set_reference((*ctx).(*CodecAlgPvt), ap)
}}, {Ctrl_id: int(vpx.VP8_COPY_REFERENCE), Fn: func(ctx *vpx.CodecAlgPvt, ap libc.ArgList) vpx.CodecErr {
	return vp8_get_reference((*ctx).(*CodecAlgPvt), ap)
}}, {Ctrl_id: int(vpx.VP8_SET_POSTPROC), Fn: func(ctx *vpx.CodecAlgPvt, ap libc.ArgList) vpx.CodecErr {
	return vp8_set_postproc((*ctx).(*CodecAlgPvt), ap)
}}, {Ctrl_id: int(vpx.VP8D_GET_LAST_REF_UPDATES), Fn: func(ctx *vpx.CodecAlgPvt, ap libc.ArgList) vpx.CodecErr {
	return vp8_get_last_ref_updates((*ctx).(*CodecAlgPvt), ap)
}}, {Ctrl_id: int(vpx.VP8D_GET_FRAME_CORRUPTED), Fn: func(ctx *vpx.CodecAlgPvt, ap libc.ArgList) vpx.CodecErr {
	return vp8_get_frame_corrupted((*ctx).(*CodecAlgPvt), ap)
}}, {Ctrl_id: int(vpx.VP8D_GET_LAST_REF_USED), Fn: func(ctx *vpx.CodecAlgPvt, ap libc.ArgList) vpx.CodecErr {
	return vp8_get_last_ref_frame((*ctx).(*CodecAlgPvt), ap)
}}, {Ctrl_id: int(vpx.VPXD_GET_LAST_QUANTIZER), Fn: func(ctx *vpx.CodecAlgPvt, ap libc.ArgList) vpx.CodecErr {
	return vp8_get_quantizer((*ctx).(*CodecAlgPvt), ap)
}}, {Ctrl_id: int(vpx.VPXD_SET_DECRYPTOR), Fn: func(ctx *vpx.CodecAlgPvt, ap libc.ArgList) vpx.CodecErr {
	return vp8_set_decryptor((*ctx).(*CodecAlgPvt), ap)
}}, {Ctrl_id: -1, Fn: nil}}

func UpdateErrorState(ctx *CodecAlgPvt, error *vpx.InternalErrorInfo) vpx.CodecErr {
	var res vpx.CodecErr
	res = error.Error_code
	if res != 0 {
		if error.Has_detail != 0 {
			ctx.Base.Err_detail = &error.Detail[0]
		} else {
			ctx.Base.Err_detail = nil
		}
	}
	return res
}

var vp8DxAlgo = vpx.CodecIFace{
	Name:        libc.CString("WebM Project VP8 Decoder v1.11.0-100-g2da19ac03"),
	Abi_version: 5,
	Caps: vpx.CodecCaps(VPX_CODEC_CAP_DECODER |
		(func() int {
			return VPX_CODEC_CAP_POSTPROC
		}()) |
		(func() int {
			return VPX_CODEC_CAP_ERROR_CONCEALMENT
		}()) | VPX_CODEC_CAP_INPUT_FRAGMENTS,
	),
	Init: func(ctx *vpx.CodecCtx, data *vpx.CodecPvtEncMrCfg) vpx.CodecErr {
		return vp8_init(ctx, data)
	},
	Destroy: func(ctx *vpx.CodecAlgPvt) vpx.CodecErr {
		return vp8_destroy((*ctx).(*CodecAlgPvt))
	},
	Ctrl_maps: &vp8CtfMaps[0],
	Dec: vpx.CodecDecIFace{
		PeekSi: func(data *uint8, data_sz uint, si *vpx.CodecStreamInfo) vpx.CodecErr {
			return vp8_peek_si(data, data_sz, si)
		},
		GetSi: func(ctx *vpx.CodecAlgPvt, si *vpx.CodecStreamInfo) vpx.CodecErr {
			return vp8_get_si((*ctx).(*CodecAlgPvt), si)
		},
		Decode: func(ctx *vpx.CodecAlgPvt, data *uint8, data_sz uint, user_priv unsafe.Pointer, deadline int) vpx.CodecErr {
			return vp8_decode((*ctx).(*CodecAlgPvt), data, data_sz, user_priv, deadline)
		},
		GetFrame: func(ctx *vpx.CodecAlgPvt, iter *vpx.CodecIter) *vpx.Image {
			return vp8_get_frame((*ctx).(*CodecAlgPvt), iter)
		},
		SetFbFn: nil,
	},
	Enc: vpx.CodecEncIFace{},
}

// vp8/common/generic/systemdependent.c file
func get_cpu_count() int {
	var core_count = 16
	// Todo:
	// core_count = sysconf(_SC_NPROCESSORS_ONLN)
	//if core_count > 0 {
	//return core_count
	// }
	return core_count
}

// vp8/common/generic/systemdependent.c file
func vp8_machine_specific_config(ctx *VP8Common) {
	ctx.Processor_core_count = get_cpu_count()
	// Todo:
	//ctx.Cpu_caps = x86_simd_caps()
}

// Todo: vp8/common/blockd.h && vp8/common/entropymodedata.h
const VP8_YMODES = B_PRED + 1
const VP8_UV_MODES = TM_PRED + 1
const VP8_BINTRAMODES = B_HU_PRED + 1

var vp8_ymode_prob = [VP8_YMODES - 1]uint8{112, 86, 140, 37}
var vp8_uv_mode_prob = [VP8_UV_MODES - 1]uint8{162, 101, 204}
var vp8_bmode_prob = [VP8_BINTRAMODES - 1]uint8{120, 90, 79, 133, 87, 85, 80, 111, 151}

var vp8_coef_update_probs = [BLOCK_TYPES][COEF_BANDS][PREV_COEF_CONTEXTS][ENTROPY_NODES]uint8{
	{
		{
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{176, 246, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{223, 241, 252, 255, 255, 255, 255, 255, 255, 255, 255},
			{249, 253, 253, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 244, 252, 255, 255, 255, 255, 255, 255, 255, 255},
			{234, 254, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{253, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 246, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{239, 253, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{254, 255, 254, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 248, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{251, 255, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 253, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{251, 254, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{254, 255, 254, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 254, 253, 255, 254, 255, 255, 255, 255, 255, 255},
			{250, 255, 254, 255, 254, 255, 255, 255, 255, 255, 255},
			{254, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
	},
	{
		{
			{217, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{225, 252, 241, 253, 255, 255, 254, 255, 255, 255, 255},
			{234, 250, 241, 250, 253, 255, 253, 254, 255, 255, 255},
		},
		{
			{255, 254, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{223, 254, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{238, 253, 254, 254, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 248, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{249, 254, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 253, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{247, 254, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 253, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{252, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 254, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{253, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 254, 253, 255, 255, 255, 255, 255, 255, 255, 255},
			{250, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{254, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
	},
	{
		{
			{186, 251, 250, 255, 255, 255, 255, 255, 255, 255, 255},
			{234, 251, 244, 254, 255, 255, 255, 255, 255, 255, 255},
			{251, 251, 243, 253, 254, 255, 254, 255, 255, 255, 255},
		},
		{
			{255, 253, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{236, 253, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{251, 253, 253, 254, 254, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 254, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{254, 254, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 254, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{254, 254, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{254, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{254, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
	},
	{
		{
			{248, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{250, 254, 252, 254, 255, 255, 255, 255, 255, 255, 255},
			{248, 254, 249, 253, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 253, 253, 255, 255, 255, 255, 255, 255, 255, 255},
			{246, 253, 253, 255, 255, 255, 255, 255, 255, 255, 255},
			{252, 254, 251, 254, 254, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 254, 252, 255, 255, 255, 255, 255, 255, 255, 255},
			{248, 254, 253, 255, 255, 255, 255, 255, 255, 255, 255},
			{253, 255, 254, 254, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 251, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{245, 251, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{253, 253, 254, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 251, 253, 255, 255, 255, 255, 255, 255, 255, 255},
			{252, 253, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 254, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 252, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{249, 255, 254, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 254, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 255, 253, 255, 255, 255, 255, 255, 255, 255, 255},
			{250, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{254, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
	},
}

// Todo: used by dec_decode_frame.go defined in standard lib
func fclose(z *stdio.File) {

}

// Todo: used by com_idct_blk_mmx.go, refer to vp8/common/x86/idct_blk_mmx.c (maps to dequantize_mmx.asm)
func vp8_dequantize_b_impl_mmx(sq *int16, dq *int16, dqc *int16) {

}

// Todo: used by com_idct_blk_sse2.go, refer to vp8/common/x86/idct_blk_sse2.c (maps to idctllm_sse2.asm)
func vp8_idct_dequant_full_2x_sse2(q *int16, dq *int16, dst *uint8, stride int) {

}

// Todo: used by com_idct_blk_sse2.go, refer to vp8/common/x86/idct_blk_sse2.c (maps to idctllm_sse2.asm)
func vp8_idct_dequant_0_2x_sse2(i *int16, dq *int16, u *uint8, stride int) {

}

// Todo: used by com_loopfilter.go, refer to vp8/common/x86/loopfilter_x86.c (maps to loopfilter_sse2.asm)
func vp8_mbloop_filter_horizontal_edge_uv_sse2(ptr *uint8, stride int, mblim *uint8, lim *uint8, thr *uint8, ptr2 *uint8) {

}

// Todo: used by com_loopfilter.go, refer to vp8/common/x86/loopfilter_x86.c (maps to loopfilter_sse2.asm)
func vp8_mbloop_filter_horizontal_edge_sse2(ptr *uint8, stride int, mblim *uint8, lim *uint8, thr *uint8) {

}

// Todo: used by com_loopfilter.go, refer to vp8/common/x86/loopfilter_x86.c (maps to loopfilter_sse2.asm)
func vp8_mbloop_filter_vertical_edge_uv_sse2(ptr *uint8, stride int, mblim *uint8, lim *uint8, thr *uint8, ptr2 *uint8) {

}

// Todo: used by com_loopfilter.go, refer to vp8/common/x86/loopfilter_x86.c (maps to loopfilter_sse2.asm)
func vp8_mbloop_filter_vertical_edge_sse2(ptr *uint8, stride int, mblim *uint8, lim *uint8, thr *uint8) {

}

// Todo: used by com_loopfilter.go, refer to vp8/common/x86/loopfilter_x86.c (maps to loopfilter_sse2.asm)
func vp8_loop_filter_vertical_edge_uv_sse2(u *uint8, stride int, blim *uint8, lim *uint8, thr *uint8, u2 *uint8) {

}

// Todo: used by com_loopfilter.go, refer to vp8/common/x86/loopfilter_x86.c (maps to loopfilter_sse2.asm)
func vp8_loop_filter_horizontal_edge_uv_sse2(u *uint8, stride int, blim *uint8, lim *uint8, thr *uint8, u2 *uint8) {

}

// Todo: used by com_loopfilter.go, refer to vp8/common/x86/loopfilter_x86.c (maps to loopfilter_block_sse2_x86_64.asm)
func vp8_loop_filter_bv_y_sse2(ptr *uint8, stride int, blim *uint8, lim *uint8, thr *uint8, i int) {

}

// Todo: used by com_loopfilter.go, refer to vp8/common/x86/loopfilter_x86.c (maps to loopfilter_block_sse2_x86_64.asm)
func vp8_loop_filter_bh_y_sse2(ptr *uint8, stride int, blim *uint8, lim *uint8, thr *uint8, i int) {

}

// Todo: callback for com_rtcd.go
func SetupRtcdInternal() {

}
