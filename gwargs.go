package gwargs

import (
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/*
Config struct reserved for future use
*/
type Config struct{}

/*
Parse takes a pointer to a struct s and attempts to parse
os.Args into that struct.

Checks will prevent any float and integer under or overflow
from occurring and will result in an error.

The config argument is reserved for future use.

Currently, the only supported types are string, bool, and
signed and unsigned ints and floats.
*/
func Parse(s any, config *Config) error {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("expected pointer, received %s", v.Kind())
	}
	if v.IsNil() {
		return errors.New("s cannot be nil")
	}
	if v.Elem().Kind() != reflect.Struct {
		return errors.New("s must be a pointer to a struct")
	}

	deRef := v.Elem()
	args := mapArgSlice(os.Args[1:])

	for i := 0; i < deRef.NumField(); i++ {
		field := t.Elem().Field(i)
		val := deRef.Field(i)

		switch val.Kind() {

		case reflect.String:
			val.SetString(args[field.Name])

		case reflect.Bool:
			if value, ok := args[field.Name]; ok {
				val.SetBool(strings.EqualFold(value, "true"))
			}

		case
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:

			n, err := strconv.ParseInt(args[field.Name], 10, 64)
			if err != nil {
				return err
			}
			if ok := checkIntOverflow(n, val.Kind()); !ok {
				return fmt.Errorf("overflow detected, cannot fit %v into %s", n, val.Kind())
			}
			val.SetInt(n)

		case
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:

			if strings.Contains(args[field.Name], "-") {
				return fmt.Errorf(
					"underflow detected, cannot fit '%v' into '%s'",
					args[field.Name],
					val.Kind(),
				)
			}

			n, err := strconv.ParseUint(args[field.Name], 10, 64)
			if err != nil {
				return err
			}

			if ok := checkUIntOverflow(n, val.Kind()); !ok {
				return fmt.Errorf("overflow detected, cannot fit %v into %s", n, val.Kind())
			}
			val.SetUint(n)

		case
			reflect.Float32,
			reflect.Float64:

			n, err := strconv.ParseFloat(args[field.Name], 64)
			if err != nil {
				return err
			}
			if ok := checkFloatOverflow(n, val.Kind()); !ok {
				return fmt.Errorf("overflow detected, cannot fit %v into %s", n, val.Kind())
			}
			val.SetFloat(n)

		default:
			return fmt.Errorf("unsupported type '%s' in field '%s'", val.Kind(), field.Name)
		}
	}

	return nil
}

/*
Parses a unix-style slice of arguments into a map

Arguments starting with two dashes '--' are treated as named
arguments and will split on '=' if present OR take the next
arg in the slice.

Arguments starting with a single dash '-' are treated as
boolean flags and split on empty space, e.g. -lahR results
in a map entry for l, a, h, and R.
*/
func mapArgSlice(args []string) map[string]string {
	res := map[string]string{}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Named args separated by " " or "="
		if strings.HasPrefix(arg, "--") {
			arg = strings.TrimPrefix(arg, "--")

			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg, "=", 2)
				res[parts[0]] = parts[1]
				continue
			}

			next := args[i+1]
			if strings.HasPrefix(next, "-") {
				res[arg] = ""
				continue
			}

			res[arg] = next

			i++
			continue
		}

		// Single character flags
		if strings.HasPrefix(arg, "-") {
			for _, flag := range strings.Split(arg[1:], "") {
				res[flag] = ""
			}
			continue
		}

		// Junk
		res[arg] = ""
	}

	return res
}

// Reports whether an unknown int will overflow
func checkIntOverflow(n int64, t reflect.Kind) (ok bool) {
	switch t {
	case reflect.Int8:
		return !(n > math.MaxInt8 || n < math.MinInt8)

	case reflect.Int16:
		return !(n > math.MaxInt16 || n < math.MinInt16)

	case reflect.Int32:
		return !(n > math.MaxInt32 || n < math.MinInt32)

	case
		reflect.Int,
		reflect.Int64:

		return !(n > math.MaxInt64 || n < math.MinInt64)

	default:
		return false
	}
}

// Reports whether an unsigned integer will overflow
func checkUIntOverflow(n uint64, k reflect.Kind) (ok bool) {
	switch k {
	case reflect.Uint8:
		return !(n > math.MaxUint8 || n < 0)

	case reflect.Uint16:
		return !(n > math.MaxUint16 || n < 0)

	case reflect.Uint32:
		return !(n > math.MaxUint32 || n < 0)

	case reflect.Uint, reflect.Uint64:
		return !(n > math.MaxUint64 || n < 0)

	default:
		return false
	}
}

// Reports whether an unknown float will overflow
func checkFloatOverflow(n float64, k reflect.Kind) (ok bool) {
	switch k {
	case reflect.Float32:
		return !(n > math.MaxFloat32 || n < math.SmallestNonzeroFloat32)

	case reflect.Float64:
		return !(n > math.MaxFloat64 || n < math.SmallestNonzeroFloat64)

	default:
		return false
	}
}
