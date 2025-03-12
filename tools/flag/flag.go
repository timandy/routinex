package flag

import (
	"fmt"
	"path/filepath"
	"reflect"
	"time"
	"unsafe"

	"github.com/timandy/pflag"
)

func PrintUsage(flagSet *pflag.FlagSet) {
	usage := fmt.Sprintf("Usage: %s [OPTIONS] [NEXT_TOOLEXEC] [ARGS] [COMMAND] [ARGS]...\n\nOptions:\n%s", flagSet.Name(), flagSet.FlagUsages())
	_, _ = fmt.Fprint(flagSet.Output(), usage)
}

func ParseStruct(structPtr any, execName string, args []string) *pflag.FlagSet {
	flagSet := createStructFlagSet(structPtr, execName)
	_ = flagSet.Parse(args)
	return flagSet
}

func createStructFlagSet(structPtr any, execName string) *pflag.FlagSet {
	shortName := filepath.Base(execName)
	elem := reflect.ValueOf(structPtr).Elem()
	typ := elem.Type()
	flagSet := pflag.NewFlagSet(shortName, pflag.ContinueOnError)
	flagSet.AllowMultCharsShorthand = true
	flagSet.ParseErrorsWhitelist.UnknownFlags = true
	flagSet.SetInterspersed(false)
	flagSet.Usage = func() {}
	for i := 0; i < typ.NumField(); i++ {
		field := elem.Field(i)
		tag := typ.Field(i).Tag
		name := getStructTagValue(tag, "name")
		shorthand := getStructTagValue(tag, "shorthand")
		usage := getStructTagValue(tag, "usage")
		if name == "" && shorthand == "" {
			continue
		}
		defineFiledFlag(flagSet, field, name, shorthand, usage)
	}
	return flagSet
}

func getStructTagValue(tag reflect.StructTag, key string) string {
	if value, ok := tag.Lookup(key); ok {
		return value
	}
	return ""
}

func defineFiledFlag(flagSet *pflag.FlagSet, field reflect.Value, name string, shorthand string, usage string) {
	valuePtr := unsafe.Pointer(field.UnsafeAddr())
	value := field.Interface()
	switch raw := value.(type) {
	case bool:
		flagSet.BoolVarP((*bool)(valuePtr), name, shorthand, raw, usage)
	case []bool:
		flagSet.BoolSliceVarP((*[]bool)(valuePtr), name, shorthand, raw, usage)
	case []byte:
		flagSet.BytesHexVarP((*[]byte)(valuePtr), name, shorthand, raw, usage)
	case time.Duration:
		flagSet.DurationVarP((*time.Duration)(valuePtr), name, shorthand, raw, usage)
	case []time.Duration:
		flagSet.DurationSliceVarP((*[]time.Duration)(valuePtr), name, shorthand, raw, usage)
	case float32:
		flagSet.Float32VarP((*float32)(valuePtr), name, shorthand, raw, usage)
	case float64:
		flagSet.Float64VarP((*float64)(valuePtr), name, shorthand, raw, usage)
	case int:
		flagSet.IntVarP((*int)(valuePtr), name, shorthand, raw, usage)
	case []int:
		flagSet.IntSliceVarP((*[]int)(valuePtr), name, shorthand, raw, usage)
	case int8:
		flagSet.Int8VarP((*int8)(valuePtr), name, shorthand, raw, usage)
	case int16:
		flagSet.Int16VarP((*int16)(valuePtr), name, shorthand, raw, usage)
	case int32:
		flagSet.Int32VarP((*int32)(valuePtr), name, shorthand, raw, usage)
	case []int32:
		flagSet.Int32SliceVarP((*[]int32)(valuePtr), name, shorthand, raw, usage)
	case int64:
		flagSet.Int64VarP((*int64)(valuePtr), name, shorthand, raw, usage)
	case []int64:
		flagSet.Int64SliceVarP((*[]int64)(valuePtr), name, shorthand, raw, usage)
	case string:
		flagSet.StringVarP((*string)(valuePtr), name, shorthand, raw, usage)
	case []string:
		flagSet.StringSliceVarP((*[]string)(valuePtr), name, shorthand, raw, usage)
	case uint:
		flagSet.UintVarP((*uint)(valuePtr), name, shorthand, raw, usage)
	case []uint:
		flagSet.UintSliceVarP((*[]uint)(valuePtr), name, shorthand, raw, usage)
	case uint8:
		flagSet.Uint8VarP((*uint8)(valuePtr), name, shorthand, raw, usage)
	case uint16:
		flagSet.Uint16VarP((*uint16)(valuePtr), name, shorthand, raw, usage)
	case uint32:
		flagSet.Uint32VarP((*uint32)(valuePtr), name, shorthand, raw, usage)
	case uint64:
		flagSet.Uint64VarP((*uint64)(valuePtr), name, shorthand, raw, usage)
	}
}
