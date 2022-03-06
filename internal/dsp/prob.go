package dsp

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const MAX_PROB uint8 = math.MaxUint8
const MODE_MV_COUNT_SAT = 20

type vpx_tree [0]int8

func get_prob(num uint, den uint) uint8 {
	libc.Assert(den != 0)
	{
		var (
			p            int = int((uint64(num)*256 + uint64(den>>1)) / uint64(den))
			clipped_prob int = p | (math.MaxUint8-p)>>23 | int(libc.BoolToInt(p == 0))
		)
		return uint8(int8(clipped_prob))
	}
}
func get_binary_prob(n0 uint, n1 uint) uint8 {
	var den uint = n0 + n1
	if den == 0 {
		return 128
	}
	return get_prob(n0, den)
}
func weighted_prob(prob1 int, prob2 int, factor int) uint8 {
	return uint8(int8(((prob1*(256-factor) + prob2*factor) + (1 << (8 - 1))) >> 8))
}
func merge_probs(pre_prob uint8, ct [2]uint, count_sat uint, max_update_factor uint) uint8 {
	var (
		prob  uint8 = get_binary_prob(ct[0], ct[1])
		count uint  = (func() uint {
			if (ct[0] + ct[1]) < count_sat {
				return ct[0] + ct[1]
			}
			return count_sat
		}())
		factor uint = max_update_factor * count / count_sat
	)
	return weighted_prob(int(pre_prob), int(prob), int(factor))
}

var count_to_update_factor [21]int = [21]int{0, 6, 12, 19, 25, 32, 38, 44, 51, 57, 64, 70, 76, 83, 89, 96, 102, 108, 115, 121, 128}

func mode_mv_merge_probs(pre_prob uint8, ct [2]uint) uint8 {
	var den uint = ct[0] + ct[1]
	if den == 0 {
		return pre_prob
	} else {
		var (
			count uint = (func() uint {
				if den < MODE_MV_COUNT_SAT {
					return den
				}
				return MODE_MV_COUNT_SAT
			}())
			factor uint  = uint(count_to_update_factor[count])
			prob   uint8 = get_prob(ct[0], den)
		)
		return weighted_prob(int(pre_prob), int(prob), int(factor))
	}
}

var vpx_norm [256]uint8 = [256]uint8{0, 7, 6, 6, 5, 5, 5, 5, 4, 4, 4, 4, 4, 4, 4, 4, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func tree_merge_probs_impl(i uint, tree *int8, pre_probs *uint8, counts *uint, probs *uint8) uint {
	var (
		l          int = int(*(*int8)(unsafe.Add(unsafe.Pointer(tree), i)))
		left_count uint
	)
	if l <= 0 {
		left_count = *(*uint)(unsafe.Add(unsafe.Pointer(counts), -int(unsafe.Sizeof(uint(0))*uintptr(l))))
	} else {
		left_count = tree_merge_probs_impl(uint(l), tree, pre_probs, counts, probs)
	}
	var r int = int(*(*int8)(unsafe.Add(unsafe.Pointer(tree), i+1)))
	var right_count uint
	if r <= 0 {
		right_count = *(*uint)(unsafe.Add(unsafe.Pointer(counts), -int(unsafe.Sizeof(uint(0))*uintptr(r))))
	} else {
		right_count = tree_merge_probs_impl(uint(r), tree, pre_probs, counts, probs)
	}
	var ct [2]uint = [2]uint{left_count, right_count}
	*(*uint8)(unsafe.Add(unsafe.Pointer(probs), i>>1)) = mode_mv_merge_probs(*(*uint8)(unsafe.Add(unsafe.Pointer(pre_probs), i>>1)), ct)
	return left_count + right_count
}
func vpx_tree_merge_probs(tree *int8, pre_probs *uint8, counts *uint, probs *uint8) {
	tree_merge_probs_impl(0, tree, pre_probs, counts, probs)
}
