package dsp

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func od_bin_fdct8x8(y *int16, ystride int, x *int16, xstride int) {
	var (
		i int
		j int
	)
	_ = xstride
	vpx_fdct8x8(x, y, ystride)
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			*((*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(int16(0))*uintptr(ystride*i)))), unsafe.Sizeof(int16(0))*uintptr(j)))) = int16((int(*((*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(int16(0))*uintptr(ystride*i)))), unsafe.Sizeof(int16(0))*uintptr(j))))) + 4) >> 3)
		}
	}
}

var csf_y [8][8]float64 = [8][8]float64{{1.6193873005, 2.2901594831, 2.08509755623, 1.48366094411, 1.00227514334, 0.678296995242, 0.466224900598, 0.3265091542}, {2.2901594831, 1.94321815382, 2.04793073064, 1.68731108984, 1.2305666963, 0.868920337363, 0.61280991668, 0.436405793551}, {2.08509755623, 2.04793073064, 1.34329019223, 1.09205635862, 0.875748795257, 0.670882927016, 0.501731932449, 0.372504254596}, {1.48366094411, 1.68731108984, 1.09205635862, 0.772819797575, 0.605636379554, 0.48309405692, 0.380429446972, 0.295774038565}, {1.00227514334, 1.2305666963, 0.875748795257, 0.605636379554, 0.448996256676, 0.352889268808, 0.283006984131, 0.226951348204}, {0.678296995242, 0.868920337363, 0.670882927016, 0.48309405692, 0.352889268808, 0.27032073436, 0.215017739696, 0.17408067321}, {0.466224900598, 0.61280991668, 0.501731932449, 0.380429446972, 0.283006984131, 0.215017739696, 0.168869545842, 0.136153931001}, {0.3265091542, 0.436405793551, 0.372504254596, 0.295774038565, 0.226951348204, 0.17408067321, 0.136153931001, 0.109083846276}}
var csf_cb420 [8][8]float64 = [8][8]float64{{1.91113096927, 2.46074210438, 1.18284184739, 1.14982565193, 1.05017074788, 0.898018824055, 0.74725392039, 0.615105596242}, {2.46074210438, 1.58529308355, 1.21363250036, 1.38190029285, 1.33100189972, 1.17428548929, 0.996404342439, 0.830890433625}, {1.18284184739, 1.21363250036, 0.978712413627, 1.02624506078, 1.03145147362, 0.960060382087, 0.849823426169, 0.731221236837}, {1.14982565193, 1.38190029285, 1.02624506078, 0.861317501629, 0.801821139099, 0.751437590932, 0.685398513368, 0.608694761374}, {1.05017074788, 1.33100189972, 1.03145147362, 0.801821139099, 0.676555426187, 0.605503172737, 0.55002013668, 0.495804539034}, {0.898018824055, 1.17428548929, 0.960060382087, 0.751437590932, 0.605503172737, 0.514674450957, 0.454353482512, 0.407050308965}, {0.74725392039, 0.996404342439, 0.849823426169, 0.685398513368, 0.55002013668, 0.454353482512, 0.389234902883, 0.342353999733}, {0.615105596242, 0.830890433625, 0.731221236837, 0.608694761374, 0.495804539034, 0.407050308965, 0.342353999733, 0.295530605237}}
var csf_cr420 [8][8]float64 = [8][8]float64{{2.03871978502, 2.62502345193, 1.26180942886, 1.11019789803, 1.01397751469, 0.867069376285, 0.721500455585, 0.593906509971}, {2.62502345193, 1.69112867013, 1.17180569821, 1.3342742857, 1.28513006198, 1.13381474809, 0.962064122248, 0.802254508198}, {1.26180942886, 1.17180569821, 0.944981930573, 0.990876405848, 0.995903384143, 0.926972725286, 0.820534991409, 0.706020324706}, {1.11019789803, 1.3342742857, 0.990876405848, 0.831632933426, 0.77418706195, 0.725539939514, 0.661776842059, 0.587716619023}, {1.01397751469, 1.28513006198, 0.995903384143, 0.77418706195, 0.653238524286, 0.584635025748, 0.531064164893, 0.478717061273}, {0.867069376285, 1.13381474809, 0.926972725286, 0.725539939514, 0.584635025748, 0.496936637883, 0.438694579826, 0.393021669543}, {0.721500455585, 0.962064122248, 0.820534991409, 0.661776842059, 0.531064164893, 0.438694579826, 0.375820256136, 0.330555063063}, {0.593906509971, 0.802254508198, 0.706020324706, 0.587716619023, 0.478717061273, 0.393021669543, 0.330555063063, 0.285345396658}}

func convert_score_db(_score float64, _weight float64, bit_depth int) float64 {
	var pix_max int16 = math.MaxUint8
	if _score*_weight >= 0.0 {
	} else {
		__assert_fail(libc.CString("_score * _weight >= 0.0"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	if bit_depth == 10 {
		pix_max = 1023
	} else if bit_depth == 12 {
		pix_max = 4095
	}
	if _weight*_score < float64(int(pix_max)*int(pix_max))*1e-10 {
		return MAX_PSNR
	}
	return (log10(float64(int(pix_max)*int(pix_max))) - log10(_weight*_score)) * 10
}
func calc_psnrhvs(src *uint8, _systride int, dst *uint8, _dystride int, _par float64, _w int, _h int, _step int, _csf [8][8]float64, bit_depth uint32, _shift uint32) float64 {
	var (
		ret        float64
		_src8      *uint8  = (*uint8)(unsafe.Pointer(src))
		_dst8      *uint8  = (*uint8)(unsafe.Pointer(dst))
		_src16     *uint16 = ((*uint16)(unsafe.Pointer(uintptr((uint64(uintptr(unsafe.Pointer(src)))) << 1))))
		_dst16     *uint16 = ((*uint16)(unsafe.Pointer(uintptr((uint64(uintptr(unsafe.Pointer(dst)))) << 1))))
		dct_s      [64]int16
		dct_d      [64]int16
		dct_s_coef [64]int16
		dct_d_coef [64]int16
		mask       [8][8]float64
		pixels     int
		x          int
		y          int
	)
	_ = _par
	ret = float64(func() int {
		pixels = 0
		return pixels
	}())
	for x = 0; x < 8; x++ {
		for y = 0; y < 8; y++ {
			mask[x][y] = (_csf[x][y] / _csf[1][0]) * (_csf[x][y] / _csf[1][0])
		}
	}
	for y = 0; y < _h-7; y += _step {
		for x = 0; x < _w-7; x += _step {
			var (
				i       int
				j       int
				s_means [4]float64
				d_means [4]float64
				s_vars  [4]float64
				d_vars  [4]float64
				s_gmean float64 = 0
				d_gmean float64 = 0
				s_gvar  float64 = 0
				d_gvar  float64 = 0
				s_mask  float64 = 0
				d_mask  float64 = 0
			)
			for i = 0; i < 4; i++ {
				s_means[i] = func() float64 {
					p := &d_means[i]
					d_means[i] = func() float64 {
						p := &s_vars[i]
						s_vars[i] = func() float64 {
							p := &d_vars[i]
							d_vars[i] = 0
							return *p
						}()
						return *p
					}()
					return *p
				}()
			}
			for i = 0; i < 8; i++ {
				for j = 0; j < 8; j++ {
					var sub int = ((i & 12) >> 2) + ((j & 12) >> 1)
					if bit_depth == 8 && _shift == 0 {
						dct_s[i*8+j] = int16(*(*uint8)(unsafe.Add(unsafe.Pointer(_src8), (y+i)*_systride+(j+x))))
						dct_d[i*8+j] = int16(*(*uint8)(unsafe.Add(unsafe.Pointer(_dst8), (y+i)*_dystride+(j+x))))
					} else if bit_depth == 10 || bit_depth == 12 {
						dct_s[i*8+j] = int16(uint16(uint32(*(*uint16)(unsafe.Add(unsafe.Pointer(_src16), unsafe.Sizeof(uint16(0))*uintptr((y+i)*_systride+(j+x))))) >> _shift))
						dct_d[i*8+j] = int16(uint16(uint32(*(*uint16)(unsafe.Add(unsafe.Pointer(_dst16), unsafe.Sizeof(uint16(0))*uintptr((y+i)*_dystride+(j+x))))) >> _shift))
					}
					s_gmean += float64(dct_s[i*8+j])
					d_gmean += float64(dct_d[i*8+j])
					s_means[sub] += float64(dct_s[i*8+j])
					d_means[sub] += float64(dct_d[i*8+j])
				}
			}
			s_gmean /= 64.0
			d_gmean /= 64.0
			for i = 0; i < 4; i++ {
				s_means[i] /= 16.0
			}
			for i = 0; i < 4; i++ {
				d_means[i] /= 16.0
			}
			for i = 0; i < 8; i++ {
				for j = 0; j < 8; j++ {
					var sub int = ((i & 12) >> 2) + ((j & 12) >> 1)
					s_gvar += (float64(dct_s[i*8+j]) - s_gmean) * (float64(dct_s[i*8+j]) - s_gmean)
					d_gvar += (float64(dct_d[i*8+j]) - d_gmean) * (float64(dct_d[i*8+j]) - d_gmean)
					s_vars[sub] += (float64(dct_s[i*8+j]) - s_means[sub]) * (float64(dct_s[i*8+j]) - s_means[sub])
					d_vars[sub] += (float64(dct_d[i*8+j]) - d_means[sub]) * (float64(dct_d[i*8+j]) - d_means[sub])
				}
			}
			s_gvar *= 1 / 63.0 * 64
			d_gvar *= 1 / 63.0 * 64
			for i = 0; i < 4; i++ {
				s_vars[i] *= 1 / 15.0 * 16
			}
			for i = 0; i < 4; i++ {
				d_vars[i] *= 1 / 15.0 * 16
			}
			if s_gvar > 0 {
				s_gvar = (s_vars[0] + s_vars[1] + s_vars[2] + s_vars[3]) / s_gvar
			}
			if d_gvar > 0 {
				d_gvar = (d_vars[0] + d_vars[1] + d_vars[2] + d_vars[3]) / d_gvar
			}
			if bit_depth == 8 {
				od_bin_fdct8x8(&dct_s_coef[0], 8, &dct_s[0], 8)
				od_bin_fdct8x8(&dct_d_coef[0], 8, &dct_d[0], 8)
			}
			for i = 0; i < 8; i++ {
				for j = int(libc.BoolToInt(i == 0)); j < 8; j++ {
					s_mask += float64(int(dct_s_coef[i*8+j])*int(dct_s_coef[i*8+j])) * mask[i][j]
				}
			}
			for i = 0; i < 8; i++ {
				for j = int(libc.BoolToInt(i == 0)); j < 8; j++ {
					d_mask += float64(int(dct_d_coef[i*8+j])*int(dct_d_coef[i*8+j])) * mask[i][j]
				}
			}
			s_mask = sqrt(s_mask*s_gvar) / 32.0
			d_mask = sqrt(d_mask*d_gvar) / 32.0
			if d_mask > s_mask {
				s_mask = d_mask
			}
			for i = 0; i < 8; i++ {
				for j = 0; j < 8; j++ {
					var err float64
					err = fabs(float64(int(dct_s_coef[i*8+j]) - int(dct_d_coef[i*8+j])))
					if i != 0 || j != 0 {
						if err < s_mask/mask[i][j] {
							err = 0
						} else {
							err = err - s_mask/mask[i][j]
						}
					}
					ret += (err * _csf[i][j]) * (err * _csf[i][j])
					pixels++
				}
			}
		}
	}
	if pixels <= 0 {
		return 0
	}
	ret /= float64(pixels)
	return ret
}
func vpx_psnrhvs(src *YV12_BUFFER_CONFIG, dest *YV12_BUFFER_CONFIG, y_psnrhvs *float64, u_psnrhvs *float64, v_psnrhvs *float64, bd uint32, in_bd uint32) float64 {
	var (
		psnrhvs  float64
		par      float64 = 1.0
		step     int     = 7
		bd_shift uint32  = 0
	)
	vpx_clear_system_state()
	if bd == 8 || bd == 10 || bd == 12 {
	} else {
		__assert_fail(libc.CString("bd == 8 || bd == 10 || bd == 12"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	if bd >= in_bd {
	} else {
		__assert_fail(libc.CString("bd >= in_bd"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	bd_shift = bd - in_bd
	*y_psnrhvs = calc_psnrhvs((*uint8)(unsafe.Pointer(src.Y_buffer)), src.Y_stride, (*uint8)(unsafe.Pointer(dest.Y_buffer)), dest.Y_stride, par, src.Y_crop_width, src.Y_crop_height, step, csf_y, bd, bd_shift)
	*u_psnrhvs = calc_psnrhvs((*uint8)(unsafe.Pointer(src.U_buffer)), src.Uv_stride, (*uint8)(unsafe.Pointer(dest.U_buffer)), dest.Uv_stride, par, src.Uv_crop_width, src.Uv_crop_height, step, csf_cb420, bd, bd_shift)
	*v_psnrhvs = calc_psnrhvs((*uint8)(unsafe.Pointer(src.V_buffer)), src.Uv_stride, (*uint8)(unsafe.Pointer(dest.V_buffer)), dest.Uv_stride, par, src.Uv_crop_width, src.Uv_crop_height, step, csf_cr420, bd, bd_shift)
	psnrhvs = (*y_psnrhvs)*0.8 + ((*u_psnrhvs)+(*v_psnrhvs))*0.1
	return convert_score_db(psnrhvs, 1.0, int(in_bd))
}
