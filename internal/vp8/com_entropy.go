package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const ZERO_TOKEN = 0
const ONE_TOKEN = 1
const TWO_TOKEN = 2
const THREE_TOKEN = 3
const FOUR_TOKEN = 4
const DCT_VAL_CATEGORY1 = 5
const DCT_VAL_CATEGORY2 = 6
const DCT_VAL_CATEGORY3 = 7
const DCT_VAL_CATEGORY4 = 8
const DCT_VAL_CATEGORY5 = 9
const DCT_VAL_CATEGORY6 = 10
const DCT_EOB_TOKEN = 11
const MAX_ENTROPY_TOKENS = 12
const ENTROPY_NODES = 11
const PROB_UPDATE_BASELINE_COST = 7
const MAX_PROB uint8 = math.MaxUint8
const DCT_MAX_VALUE = 2048
const BLOCK_TYPES = 4
const COEF_BANDS = 8
const PREV_COEF_CONTEXTS = 3

type vp8_extra_bit_struct struct {
	Tree     vp8_tree_p
	Prob     *uint8
	Len      int
	Base_val int
}

var vp8_norm [256]uint8 = [256]uint8{0, 7, 6, 6, 5, 5, 5, 5, 4, 4, 4, 4, 4, 4, 4, 4, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var vp8_coef_bands [16]uint8 = [16]uint8{0, 1, 2, 3, 6, 4, 5, 6, 6, 6, 6, 6, 6, 6, 6, 7}
var vp8_prev_token_class [12]uint8 = [12]uint8{0, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0}
var vp8_default_zig_zag1d [16]int = [16]int{0, 1, 4, 8, 5, 2, 3, 6, 9, 12, 13, 10, 7, 11, 14, 15}
var vp8_default_inv_zig_zag [16]int16 = [16]int16{1, 2, 6, 7, 3, 5, 8, 13, 4, 9, 12, 14, 10, 11, 15, 16}
var vp8_default_zig_zag_mask [16]int16 = [16]int16{1, 2, 32, 64, 4, 16, 128, 4096, 8, 256, 2048, 8192, 512, 1024, 0x4000, math.MinInt16}
var vp8_mb_feature_data_bits [2]int = [2]int{7, 6}
var vp8_coef_tree [22]int8 = [22]int8{-DCT_EOB_TOKEN, 2, -ZERO_TOKEN, 4, -ONE_TOKEN, 6, 8, 12, -TWO_TOKEN, 10, -THREE_TOKEN, -FOUR_TOKEN, 14, 16, -DCT_VAL_CATEGORY1, -DCT_VAL_CATEGORY2, 18, 20, -DCT_VAL_CATEGORY3, -DCT_VAL_CATEGORY4, -DCT_VAL_CATEGORY5, -DCT_VAL_CATEGORY6}
var vp8_coef_encodings [12]vp8_token = [12]vp8_token{{Value: 2, Len: 2}, {Value: 6, Len: 3}, {Value: 28, Len: 5}, {Value: 58, Len: 6}, {Value: 59, Len: 6}, {Value: 60, Len: 6}, {Value: 61, Len: 6}, {Value: 124, Len: 7}, {Value: 125, Len: 7}, {Value: 126, Len: 7}, {Value: math.MaxInt8, Len: 7}, {Value: 0, Len: 1}}
var Pcat1 [1]uint8 = [1]uint8{159}
var Pcat2 [2]uint8 = [2]uint8{165, 145}
var Pcat3 [3]uint8 = [3]uint8{173, 148, 140}
var Pcat4 [4]uint8 = [4]uint8{176, 155, 140, 135}
var Pcat5 [5]uint8 = [5]uint8{180, 157, 141, 134, 130}
var Pcat6 [11]uint8 = [11]uint8{254, 254, 243, 230, 196, 177, 153, 140, 133, 130, 129}
var cat1 [2]int8 = [2]int8{}
var cat2 [4]int8 = [4]int8{2, 2, 0, 0}
var cat3 [6]int8 = [6]int8{2, 2, 4, 4, 0, 0}
var cat4 [8]int8 = [8]int8{2, 2, 4, 4, 6, 6, 0, 0}
var cat5 [10]int8 = [10]int8{2, 2, 4, 4, 6, 6, 8, 8, 0, 0}
var cat6 [22]int8 = [22]int8{2, 2, 4, 4, 6, 6, 8, 8, 10, 10, 12, 12, 14, 14, 16, 16, 18, 18, 20, 20, 0, 0}
var vp8_extra_bits [12]vp8_extra_bit_struct = [12]vp8_extra_bit_struct{{}, {Tree: nil, Prob: nil, Len: 0, Base_val: 1}, {Tree: nil, Prob: nil, Len: 0, Base_val: 2}, {Tree: nil, Prob: nil, Len: 0, Base_val: 3}, {Tree: nil, Prob: nil, Len: 0, Base_val: 4}, {Tree: &cat1[0], Prob: &Pcat1[0], Len: 1, Base_val: 5}, {Tree: &cat2[0], Prob: &Pcat2[0], Len: 2, Base_val: 7}, {Tree: &cat3[0], Prob: &Pcat3[0], Len: 3, Base_val: 11}, {Tree: &cat4[0], Prob: &Pcat4[0], Len: 4, Base_val: 19}, {Tree: &cat5[0], Prob: &Pcat5[0], Len: 5, Base_val: 35}, {Tree: &cat6[0], Prob: &Pcat6[0], Len: 11, Base_val: 67}, {}}

func vp8_default_coef_probs(pc *VP8Common) {
	libc.MemCpy(unsafe.Pointer(&pc.Fc.Coef_probs[0][0][0][0]), unsafe.Pointer(&default_coef_probs[0][0][0][0]), int(unsafe.Sizeof([4][8][3][11]uint8{})))
}
