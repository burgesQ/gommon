// assert package hold common function wrapping the
// testing.T object to perfome various assertion.
// The assert pacakge need no init, and  outputed are colored.
// To disable the colored output, set the `GOTESTNOCOLOR` env var to any value
// other than nothing
package assert

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/logrusorgru/aurora"
)

const (
	// Assertion varible type to be perfome
	String VariableType = iota
	Int

	_success = "\t\t\t\t ✓"
	_fail    = "\t\t\t\t ✗"

	_flagType = "STRING = 1 - INT = 2"
)

var (
	// fetch env value
	disableColorVal, disableColor = os.LookupEnv("GOTESTNOCOLOR")
	colorEnable                   = (!disableColor || (disableColor && disableColorVal == ""))
	au                            = aurora.NewAurora(colorEnable)
	_notEqual                     = Sprintf("\n\t[%s] :\t> %%v <\t\n\t[%s] :\t> %%v <\t\n",
		au.Bold(Red("✗")), au.Bold(Green("✓")))
)

type VariableType int

func assert(t *testing.T, method func() bool, context string, args ...interface{}) {
	t.Helper()

	if method() {
		t.Log(au.Green(_success))
		return
	}

	t.Log(au.Bold(au.Red(_fail)))

	switch {
	case len(args) == 2 && context == _notEqual:
		t.Errorf(context, au.Bold(Yellow(args[0])), au.Bold(Green(args[1])))
	case len(args) > 0:
		t.Errorf(context, args...)
	default:
		t.Errorf(context)
	}
}

// SimpleEqualContext assert that have and want are equal and output the
// formated context args as error message
func SimpleEqualContext(t *testing.T, have, want interface{},
	context string, args ...interface{}) {
	t.Helper()
	assert(t, func() bool { return have == want }, context, args...)
}

// SimpleNotEqualContext assert that have and want are not equal and output the
// formated context args as error message
func SimpleNotEqualContext(t *testing.T, have, want interface{},
	context string, args ...interface{}) {
	t.Helper()
	assert(t, func() bool { return have != want }, context, args...)
}

// SimpleEqual assert that have and want are equal
func SimpleEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	SimpleEqualContext(t, have, want, _notEqual, have, want)
}

// SimpleNotEqual assert that have and want are not equal
func SimpleNotEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	SimpleNotEqualContext(t, have, want, _notEqual, have, want)
}

// Equal run an assertion that the argument are equal
func Equal(t *testing.T, have, want interface{}) {
	t.Helper()
	SimpleEqual(t, have, want)
}

// EqualContext run a assertion that the arguments are equal
// on fail the error message context is yield
func EqualContext(t *testing.T, have, want interface{},
	context string, args ...interface{}) {
	t.Helper()
	SimpleEqualContext(t, have, want, context, args...)
}

// NotEqual run an assertion that the argument are not equal
func NotEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	SimpleNotEqual(t, have, want)
}

// NotNil run an assertion that the argument is not nil
func NotNil(t *testing.T, have interface{}) {
	t.Helper()
	SimpleNotEqual(t, have, nil)
}

// Nil run an ion that the argument is nil
func Nil(t *testing.T, have interface{}) {
	t.Helper()
	SimpleEqual(t, have, nil)
}

// True run an assertion that the bool argument is true
func True(t *testing.T, have bool) {
	t.Helper()
	SimpleEqual(t, have, true)
}

// TrueContext run an assertion that the bool argument is true
func TrueContext(t *testing.T, have bool, fmt string, args ...interface{}) {
	t.Helper()
	SimpleEqualContext(t, have, true, fmt, args...)
}

// False run an assertion that the bool argument is false
func False(t *testing.T, have bool) {
	t.Helper()
	SimpleEqual(t, have, false)
}

// FalseContext run an assertion that the bool argument is false with the
// provied context
func FalseContext(t *testing.T, have bool, context string, args ...interface{}) {
	t.Helper()
	SimpleEqualContext(t, have, false, context, args...)
}

// StringEqual run an assertion that the two string arguments are equal
func StringEqual(t *testing.T, have, want string) {
	t.Helper()
	SimpleEqual(t, have, want)
}

// StringNotEqual run an assertion that the two string arguments are not equal
func StringNotEqual(t *testing.T, have, want string) {
	t.Helper()
	SimpleNotEqual(t, have, want)
}

// IntEqual run an assertion that the two string arguments are equal
func IntEqual(t *testing.T, have, want int) {
	t.Helper()
	SimpleEqual(t, have, want)
}

// UInt16Equal run an assertion that the two string arguments are equal
func UInt16Equal(t *testing.T, have, want uint16) {
	t.Helper()
	SimpleEqual(t, have, want)
}

// IntNotEqual run an assertion that the two string arguments are not equal
func IntNotEqual(t *testing.T, have, want int) {
	t.Helper()
	SimpleNotEqual(t, have, want)
}

// MapOfInterfaceEqual
func MapOfInterfaceEqual(t *testing.T, have, want map[string]interface{}) {
	t.Helper()
	TrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}

// MapOfStringEqual
func MapOfStringEqual(t *testing.T, have, want map[string]string) {
	t.Helper()
	TrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}

// Is assert the type of have
func Is(t *testing.T, have interface{}, what VariableType) {
	t.Helper()
	var (
		fmt = ""
		ok  = false
	)
	switch what {
	case String:
		_, ok = have.(string)
		fmt = "%s is not a string (%T)"
	case Int:
		_, ok = have.(int)
		fmt = "%s is not a int (%T)"
	default:
		t.Errorf("type assertion not supported (%d)[%s]", what, _flagType)
	}
	TrueContext(t, ok, fmt, have, have)
}

// IsString assert that have is a string
func IsString(t *testing.T, have interface{}) {
	t.Helper()
	Is(t, have, String)
}

// IsInt assert that have in as int
func IsInt(t *testing.T, have interface{}) {
	t.Helper()
	Is(t, have, Int)
}

// Lower assert that have is lower than want
// 0 - 0 -- false
// 1 - 0 -- false
// 0 - 1 -- true
func Lower(t *testing.T, have, want int) {
	t.Helper()
	assert(t, func() bool { return have < want }, "%d is greater than %d ", have, want)
}

// EuqalOrGreater assert that have is lower than want
// 0 - 0 -- true
// 1 - 0 -- true
// 0 - 1 -- false
func EqualOrGreater(t *testing.T, have, want int) {
	t.Helper()
	assert(t, func() bool { return have >= want }, "%d is strictly lower than %d ", have, want)
}

// SliceByteEqual assert that the have and want slice of byte are equal
func SliceByteEqual(t *testing.T, have, want []byte) {
	t.Helper()
	TrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}

// SliceByteU16Equal assert that the have and want slice of byte are equal
func SliceU16Equal(t *testing.T, have, want []uint16) {
	t.Helper()
	TrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}

// SliceEqual assert that the have and want slice of byte are equal
func SliceEqual(t *testing.T, have, want []interface{}) {
	t.Helper()
	TrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}

// SliceOfStringMap assert that the slice have the same values
func SliceOfStringMap(t *testing.T, have, want []map[string]string) {
	t.Helper()
	TrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}

// ErrorIs assert that an error belong to an error type
func ErrorIs(t *testing.T, e, what error) {
	t.Helper()

	TrueContext(t, errors.Is(e, what), "errors %T isn't %T", e, what)
}
