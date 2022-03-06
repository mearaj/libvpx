package internal

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unsafe"
)

type Arg struct {
	Argv      **byte
	Name      *byte
	Val       *byte
	Argv_step uint
	Def       *ArgDef
}
type arg_enum_list struct {
	Name *byte
	Val  int
}
type ArgDef struct {
	Short_name *byte
	Long_name  *byte
	Has_val    int
	Desc       *byte
	Enums      *arg_enum_list
}
type ArgDefT ArgDef

func arg_init(argv **byte) Arg {
	var a Arg
	a.Argv = argv
	a.Argv_step = 1
	a.Name = nil
	a.Val = nil
	a.Def = nil
	return a
}
func ArgMatch(arg_ *Arg, def *ArgDef, argv **byte) int {
	var Arg Arg
	if *(**byte)(unsafe.Add(unsafe.Pointer(argv), unsafe.Sizeof((*byte)(nil))*0)) == nil || *(*byte)(unsafe.Add(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(argv), unsafe.Sizeof((*byte)(nil))*0))), 0)) != byte('-') {
		return 0
	}
	Arg = arg_init(argv)
	if def.Short_name != nil && libc.StrLen(*(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*0))) == libc.StrLen(def.Short_name)+1 && libc.StrCmp((*byte)(unsafe.Add(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*0))), 1)), def.Short_name) == 0 {
		Arg.Name = (*byte)(unsafe.Add(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*0))), 1))
		if def.Has_val != 0 {
			Arg.Val = *(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*1))
		} else {
			Arg.Val = nil
		}
		if def.Has_val != 0 {
			Arg.Argv_step = 2
		} else {
			Arg.Argv_step = 1
		}
	} else if def.Long_name != nil {
		var name_len uint64 = uint64(libc.StrLen(def.Long_name))
		if libc.StrLen(*(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*0))) >= int(name_len+2) && *(*byte)(unsafe.Add(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*0))), 1)) == byte('-') && libc.StrNCmp((*byte)(unsafe.Add(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*0))), 2)), def.Long_name, int(name_len)) == 0 && (*(*byte)(unsafe.Add(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*0))), name_len+2)) == byte('=') || *(*byte)(unsafe.Add(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*0))), name_len+2)) == byte('\x00')) {
			Arg.Name = (*byte)(unsafe.Add(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*0))), 2))
			if *(*byte)(unsafe.Add(unsafe.Pointer(Arg.Name), name_len)) == byte('=') {
				Arg.Val = (*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(Arg.Name), name_len))), 1))
			} else {
				Arg.Val = nil
			}
			Arg.Argv_step = 1
		}
	}
	if Arg.Name != nil && Arg.Val == nil && def.Has_val != 0 {
		Die(libc.CString("Error: option %s requires argument.\n"), Arg.Name)
	}
	if Arg.Name != nil && Arg.Val != nil && def.Has_val == 0 {
		Die(libc.CString("Error: option %s requires no argument.\n"), Arg.Name)
	}
	if Arg.Name != nil && (Arg.Val != nil || def.Has_val == 0) {
		Arg.Def = def
		*arg_ = Arg
		return 1
	}
	return 0
}
func arg_next(Arg *Arg) *byte {
	if *(**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*0)) != nil {
		Arg.Argv = (**byte)(unsafe.Add(unsafe.Pointer(Arg.Argv), unsafe.Sizeof((*byte)(nil))*uintptr(Arg.Argv_step)))
	}
	return *Arg.Argv
}
func ArgvDup(argc int, argv **byte) **byte {
	var new_argv **byte = (**byte)(libc.Malloc((argc + 1) * int(unsafe.Sizeof((*byte)(nil)))))
	libc.MemCpy(unsafe.Pointer(new_argv), unsafe.Pointer(argv), argc*int(unsafe.Sizeof((*byte)(nil))))
	*(**byte)(unsafe.Add(unsafe.Pointer(new_argv), unsafe.Sizeof((*byte)(nil))*uintptr(argc))) = nil
	return new_argv
}
func ArgShowUsage(fp *stdio.File, defs **ArgDef) {
	var option_text [40]byte = [40]byte{}
	for ; *defs != nil; defs = (**ArgDef)(unsafe.Add(unsafe.Pointer(defs), unsafe.Sizeof((*ArgDef)(nil))*1)) {
		var (
			def       *ArgDef = *defs
			short_val *byte
		)
		if def.Has_val != 0 {
			short_val = libc.CString(" <arg>")
		} else {
			short_val = libc.CString("")
		}
		var long_val *byte
		if def.Has_val != 0 {
			long_val = libc.CString("=<arg>")
		} else {
			long_val = libc.CString("")
		}
		if def.Short_name != nil && def.Long_name != nil {
			var comma *byte
			if def.Has_val != 0 {
				comma = libc.CString(",")
			} else {
				comma = libc.CString(",      ")
			}
			stdio.Snprintf(&option_text[0], 37, "-%s%s%s --%s%6s", def.Short_name, short_val, comma, def.Long_name, long_val)
		} else if def.Short_name != nil {
			stdio.Snprintf(&option_text[0], 37, "-%s%s", def.Short_name, short_val)
		} else if def.Long_name != nil {
			stdio.Snprintf(&option_text[0], 37, "          --%s%s", def.Long_name, long_val)
		}
		stdio.Fprintf(fp, "  %-37s\t%s\n", &option_text[0], def.Desc)
		if def.Enums != nil {
			var listptr *arg_enum_list
			stdio.Fprintf(fp, "  %-37s\t  ", "")
			for listptr = def.Enums; listptr.Name != nil; listptr = (*arg_enum_list)(unsafe.Add(unsafe.Pointer(listptr), unsafe.Sizeof(arg_enum_list{})*1)) {
				stdio.Fprintf(fp, "%s%s", listptr.Name, func() string {
					if (*(*arg_enum_list)(unsafe.Add(unsafe.Pointer(listptr), unsafe.Sizeof(arg_enum_list{})*1))).Name != nil {
						return ", "
					}
					return "\n"
				}())
			}
		}
	}
}
func ArgParseUint(Arg *Arg) uint {
	var (
		rawval uint32
		endptr *byte
	)
	rawval = uint32(strtoul(Arg.Val, &endptr, 10))
	if *(*byte)(unsafe.Add(unsafe.Pointer(Arg.Val), 0)) != byte('\x00') && *(*byte)(unsafe.Add(unsafe.Pointer(endptr), 0)) == byte('\x00') {
		if int(rawval) <= math.MaxUint32 {
			return uint(rawval)
		}
		Die(libc.CString("Option %s: Value %ld out of range for unsigned int\n"), Arg.Name, rawval)
	}
	Die(libc.CString("Option %s: Invalid character '%c'\n"), Arg.Name, *endptr)
	return 0
}
func arg_parse_int(Arg *Arg) int {
	var (
		rawval int32
		endptr *byte
	)
	rawval = int32(strtol(Arg.Val, &endptr, 10))
	if *(*byte)(unsafe.Add(unsafe.Pointer(Arg.Val), 0)) != byte('\x00') && *(*byte)(unsafe.Add(unsafe.Pointer(endptr), 0)) == byte('\x00') {
		if int(rawval) >= math.MinInt32 && int(rawval) <= math.MaxInt32 {
			return int(rawval)
		}
		Die(libc.CString("Option %s: Value %ld out of range for signed int\n"), Arg.Name, rawval)
	}
	Die(libc.CString("Option %s: Invalid character '%c'\n"), Arg.Name, *endptr)
	return 0
}

type vpx_rational struct {
	Num int
	Den int
}

func arg_parse_rational(Arg *Arg) vpx_rational {
	var (
		rawval int
		endptr *byte
		rat    vpx_rational
	)
	rawval = strtol(Arg.Val, &endptr, 10)
	if *(*byte)(unsafe.Add(unsafe.Pointer(Arg.Val), 0)) != byte('\x00') && *(*byte)(unsafe.Add(unsafe.Pointer(endptr), 0)) == byte('/') {
		if rawval >= math.MinInt32 && rawval <= math.MaxInt32 {
			rat.Num = rawval
		} else {
			Die(libc.CString("Option %s: Value %ld out of range for signed int\n"), Arg.Name, rawval)
		}
	} else {
		Die(libc.CString("Option %s: Expected / at '%c'\n"), Arg.Name, *endptr)
	}
	rawval = strtol((*byte)(unsafe.Add(unsafe.Pointer(endptr), 1)), &endptr, 10)
	if *(*byte)(unsafe.Add(unsafe.Pointer(Arg.Val), 0)) != byte('\x00') && *(*byte)(unsafe.Add(unsafe.Pointer(endptr), 0)) == byte('\x00') {
		if rawval >= math.MinInt32 && rawval <= math.MaxInt32 {
			rat.Den = rawval
		} else {
			Die(libc.CString("Option %s: Value %ld out of range for signed int\n"), Arg.Name, rawval)
		}
	} else {
		Die(libc.CString("Option %s: Invalid character '%c'\n"), Arg.Name, *endptr)
	}
	return rat
}
func arg_parse_enum(Arg *Arg) int {
	var (
		listptr *arg_enum_list
		rawval  int
		endptr  *byte
	)
	rawval = strtol(Arg.Val, &endptr, 10)
	if *(*byte)(unsafe.Add(unsafe.Pointer(Arg.Val), 0)) != byte('\x00') && *(*byte)(unsafe.Add(unsafe.Pointer(endptr), 0)) == byte('\x00') {
		for listptr = Arg.Def.Enums; listptr.Name != nil; listptr = (*arg_enum_list)(unsafe.Add(unsafe.Pointer(listptr), unsafe.Sizeof(arg_enum_list{})*1)) {
			if listptr.Val == rawval {
				return rawval
			}
		}
	}
	for listptr = Arg.Def.Enums; listptr.Name != nil; listptr = (*arg_enum_list)(unsafe.Add(unsafe.Pointer(listptr), unsafe.Sizeof(arg_enum_list{})*1)) {
		if libc.StrCmp(Arg.Val, listptr.Name) == 0 {
			return listptr.Val
		}
	}
	Die(libc.CString("Option %s: Invalid value '%s'\n"), Arg.Name, Arg.Val)
	return 0
}
func arg_parse_enum_or_int(Arg *Arg) int {
	if Arg.Def.Enums != nil {
		return arg_parse_enum(Arg)
	}
	return arg_parse_int(Arg)
}
