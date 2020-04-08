package assert

import (
	"reflect"
	"testing"
)

const (
	STRING = 1
	INT    = 2

	_notEqual = "assertion failed\ngot :\t>[\t%v\t]<\nwant :\t>[\t%v\t]<"
	_flagType = "STRING = 1 - INT = 2"
)

func assert(t *testing.T, method func() bool,
	context string, args ...interface{}) {
	t.Helper()
	if !method() {
		if len(args) > 0 {
			t.Errorf(context, args...)
		} else {
			t.Errorf(context)
		}
	}
}

func SimpleEqualContext(t *testing.T, have, want interface{},
	context string, args ...interface{}) {
	t.Helper()
	assert(t, func() bool { return have == want }, context, args...)
}

func SimpleNotEqualContext(t *testing.T, have, want interface{},
	context string, args ...interface{}) {
	t.Helper()
	assert(t, func() bool { return have != want }, context, args...)
}

func SimpleEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	SimpleEqualContext(t, have, want, _notEqual, have, want)
}

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
func Is(t *testing.T, have interface{}, what int) {
	t.Helper()
	var (
		fmt = ""
		ok  = false
	)
	switch what {
	case STRING:
		_, ok = have.(string)
		fmt = "%s is not a string (%T)"
	case INT:
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
	Is(t, have, STRING)
}

// IsInt assert that have in as int
func IsInt(t *testing.T, have interface{}) {
	t.Helper()
	Is(t, have, INT)
}

// 0 - 0 -- false
// 1 - 0 -- false
// 0 - 1 -- true
func Lower(t *testing.T, have, want int) {
	t.Helper()
	assert(t, func() bool { return have < want }, "%d is greater than %d ", have, want)
}

// 0 - 0 -- true
// 1 - 0 -- true
// 0 - 1 -- false
func EqualOrGreater(t *testing.T, have, want int) {
	t.Helper()
	assert(t, func() bool { return have >= want }, "%d is strictly lower than %d ", have, want)
}

// SliceByteEqual
func SliceByteEqual(t *testing.T, have, want []byte) {
	t.Helper()
	TrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}

// SliceEqual
func SliceEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	TrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}
