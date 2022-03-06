# libvpx in golang #

### Story
* The idea is to transpile libvpx c code into libvpx go code inside [internal](/internal) package and once the job is done,
then expose cleaner interfaces (closer to go standard and other necessary standards) to the outside world.
* [libvpxsrc](/libvpxsrc) folder was created by cloning `git clone https://chromium.googlesource.com/webm/libvpx libvpxsrc`<br>
* [libvpxbuild](/libvpxbuild) folder was created to keep compiled output of [libvpxsrc](/libvpxsrc).<br>
`cd libvpxbuild `<br>
`../libvpxsrc/configure`<br>
`make`
* The header files generated from compilation i.e.
[vp8_rtcd.h](/libvpxbuild/vp8_rtcd.h),
[vp9_rtcd.h](/libvpxbuild/vp9_rtcd.h),
[vpx_config.h](/libvpxbuild/vpx_config.h),
[vpx_dsp_rtcd.h](/libvpxbuild/vpx_dsp_rtcd.h),
[vpx_scale_rtcd.h](/libvpxbuild/vpx_scale_rtcd.h) and 
[vpx_version.h](/libvpxbuild/vpx_version.h)
were copied into [libvpxsrc](/libvpxsrc) folder.
* The [include](/include) and [includes](/includes) folder were created at the root directory of this project.
They are copied into [libvpxsrc](/libvpxsrc) directory by [compile.sh](/compile.sh) script,
so that [cxgo](https://github.com/gotranspile/cxgo) transpiler auto identifies it.
* The following is the mapping of src directories to transpiled directories<br>
[libvpxsrc](/libvpxsrc)                        => [internal](/internal)<br>
[libvpxsrc/vpx_scale](libvpxsrc/vpx_scale)     => [internal/scale](internal/scale)<br>
[libvpxsrc/vp9](libvpxsrc/vp9)                 => [internal/vp9](internal/vp9)<br>
[libvpxsrc/vp8](libvpxsrc/vp8)                 => [internal/vp8](internal/vp8)<br>
[libvpxsrc/vpx](libvpxsrc/vpx)                 => [internal/vpx](internal/vpx)<br>
[libvpxsrc/vpx_mem](libvpxsrc/vpx_mem)         => [internal/mem](internal/mem)<br>
[libvpxsrc/vpx_dsp](libvpxsrc/vpx_dsp)         => [internal/dsp](internal/dsp)<br>
[libvpxsrc/vpx_util](libvpxsrc/vpx_util)       => [internal/util](internal/util)<br>
[libvpxsrc/vpx_ports](libvpxsrc/vpx_ports)     => [internal/ports](internal/ports)<br>
In the above folders, there exists a corresponding cxgo's yml file and one hack.go file which may or may not exist.<br>
For example, [internal/ports/ports.yml](internal/ports/ports.yml) and [internal/ports/hack.go](internal/ports/hack.go).
* The transpilation process starts by running [compile.sh](compile.sh) script that does some heck sort of work.<br>
It deletes all *.go files except hack.go file and runs cxgo's yml file for transpilation and new go files are generated.
* Due to my current lack of knowledge in this area, the implementation is not so clean and there's much room for improvement.
Ultimate goal is that with very little or no intervention from human side, the code should be able to transpile libvpx
c code into libvpx go code inside the internal package and then expose cleaner interface to the outside world.


* dsp - Digital Signal Processor
* rtc - Real Time Communication


### Directory structure of [libvpxsrc](/libvpxsrc) 
```
third_party
third_party/googletest
third_party/googletest/src
third_party/googletest/src/include
third_party/googletest/src/include/gtest
third_party/googletest/src/include/gtest/internal
third_party/googletest/src/include/gtest/internal/custom
third_party/googletest/src/src
third_party/libwebm
third_party/libwebm/mkvparser
third_party/libwebm/common
third_party/libwebm/mkvmuxer
third_party/x86inc
third_party/libyuv
third_party/libyuv/source
third_party/libyuv/include
third_party/libyuv/include/libyuv
examples
includes
includes/bits
includes/bits/types
build
build/make
vpx_dsp
vpx_dsp/ppc
vpx_dsp/mips
vpx_dsp/loongarch
vpx_dsp/arm
vpx_dsp/x86
vpx_mem
vpx_mem/include
tools
tools/3D-Reconstruction
tools/3D-Reconstruction/genY4M
tools/3D-Reconstruction/sketch_3D_reconstruction
tools/3D-Reconstruction/MotionEST
tools/non_greedy_mv
build_debug
build_debug/non_greedy_mv_test_files
include
include/bits
include/bits/types
vp9
vp9/common
vp9/common/ppc
vp9/common/mips
vp9/common/mips/dspr2
vp9/common/mips/msa
vp9/common/arm
vp9/common/arm/neon
vp9/common/x86
vp9/decoder
vp9/encoder
vp9/encoder/ppc
vp9/encoder/mips
vp9/encoder/mips/msa
vp9/encoder/arm
vp9/encoder/arm/neon
vp9/encoder/x86
vpx_scale
vpx_scale/mips
vpx_scale/mips/dspr2
vpx_scale/generic
vpx_ports
test
test/android
vp8
vp8/common
vp8/common/mips
vp8/common/mips/dspr2
vp8/common/mips/msa
vp8/common/mips/mmi
vp8/common/loongarch
vp8/common/generic
vp8/common/arm
vp8/common/arm/neon
vp8/common/x86
vp8/decoder
vp8/encoder
vp8/encoder/mips
vp8/encoder/mips/msa
vp8/encoder/mips/mmi
vp8/encoder/arm
vp8/encoder/arm/neon
vp8/encoder/x86
vpx
vpx/internal
vpx/src
vpx_util
```

### The following files maps to libvpx package on Arch Linux
```
/usr/
/usr/bin/
/usr/bin/vpxdec
/usr/bin/vpxenc
/usr/include/
/usr/include/vpx/
/usr/include/vpx/vp8.h
/usr/include/vpx/vp8cx.h
/usr/include/vpx/vp8dx.h
/usr/include/vpx/vpx_codec.h
/usr/include/vpx/vpx_decoder.h
/usr/include/vpx/vpx_encoder.h
/usr/include/vpx/vpx_ext_ratectrl.h
/usr/include/vpx/vpx_frame_buffer.h
/usr/include/vpx/vpx_image.h
/usr/include/vpx/vpx_integer.h
/usr/lib/
/usr/lib/libvpx.so
/usr/lib/libvpx.so.7
/usr/lib/libvpx.so.7.0
/usr/lib/libvpx.so.7.0.0
/usr/lib/pkgconfig/
/usr/lib/pkgconfig/vpx.pc
/usr/share/
/usr/share/licenses/
/usr/share/licenses/libvpx/
/usr/share/licenses/libvpx/LICENSE
```

### Issues
* typedef of void in struct_FILE.h at line 43
  CurrentSolution: replace typedef void _IO_lock_t; with typedef void* _IO_lock_t;