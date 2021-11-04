package assert

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gohobby/assert/tablewriter"
)

// Equal asserts that two objects are equal.
//
//    assert.Equal(t, 123, 123)
//
// Function equality cannot be determined and will always fail.
func Equal(t testing.TB, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if err := validateEqualArgs(expected, actual); err != nil {
		return Fail(t, fmt.Sprintf("Invalid operation: %#v == %#v (%s)", expected, actual, err), nil, msgAndArgs...)
	}

	if reflect.DeepEqual(expected, actual) {
		return true
	}

	if expected == nil {
		return Fail(t, fmt.Sprintf("Expected nil, but got: %#v", actual), nil, msgAndArgs...)
	}

	return Fail(t, "Not equal", func(formatter *tablewriter.Table) {
		formatter.Writef("\nExpect:\t%+v\t(%T)", expected, expected)
		formatter.Writef("\nActual:\t%+v\t(%T)", actual, actual)
	}, msgAndArgs...)
}

// NotEqual asserts that the specified values are NOT equal.
//
//    assert.NotEqual(t, obj1, obj2)
//
// Function equality cannot be determined and will always fail.
func NotEqual(t testing.TB, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if err := validateEqualArgs(expected, actual); err != nil {
		return Fail(t, fmt.Sprintf("Invalid operation: %#v == %#v (%s)", expected, actual, err), nil, msgAndArgs...)
	}

	if !reflect.DeepEqual(expected, actual) {
		return true
	}

	if expected == nil {
		return Fail(t, "Expected value not to be nil.", nil, msgAndArgs...)
	}

	return Fail(t, fmt.Sprintf("Should not be: %#v", actual), nil, msgAndArgs...)
}

// True asserts that the specified value is true.
//
//    assert.True(t, myBool)
func True(t *testing.T, value bool, msgAndArgs ...interface{}) bool {
	if !value {
		t.Helper()
		return Fail(t, "Should be true", nil, msgAndArgs...)
	}

	return true
}

// False asserts that the specified value is false.
//
//    assert.False(t, myBool)
func False(t *testing.T, value bool, msgAndArgs ...interface{}) bool {
	if value {
		t.Helper()
		return Fail(t, "Should be false", nil, msgAndArgs...)
	}

	return true
}

// Nil asserts that the specified object is nil.
//
//    assert.Nil(t, err)
func Nil(t *testing.T, object interface{}, msgAndArgs ...interface{}) bool {
	if isNil(object) {
		return true
	}

	t.Helper()

	return Fail(t, fmt.Sprintf("Expected nil, but got: %#v", object), nil, msgAndArgs...)
}

// NotNil asserts that the specified object is not nil.
//
//    assert.NotNil(t, err)
func NotNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) bool {
	if !isNil(object) {
		return true
	}

	t.Helper()

	return Fail(t, "Expected value not to be nil.", nil, msgAndArgs...)
}

// Implements asserts that an object is implemented by the specified interface.
//
//    assert.Implements(t, (*MyInterface)(nil), new(MyObject))
func Implements(t *testing.T, interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	interfaceType := reflect.TypeOf(interfaceObject).Elem()

	if object == nil {
		return Fail(t, fmt.Sprintf("Cannot check if nil implements %v", interfaceType), nil, msgAndArgs...)
	}

	if !reflect.TypeOf(object).Implements(interfaceType) {
		return Fail(t, fmt.Sprintf("%T must implement %v", object, interfaceType), nil, msgAndArgs...)
	}

	return true
}

// Fail reports a failure
func Fail(
	t testing.TB,
	failureMessage string,
	callback func(formatter *tablewriter.Table),
	msgAndArgs ...interface{},
) bool {
	t.Helper()

	stackTrace := StackTrace(3) // StackTrace + Fail + public function

	table := tablewriter.New()
	table.WriteRow("Test:", t.Name())
	table.Writef("\nTrace:\t%s", strings.Join(stackTrace, "\n\t"))

	if len(failureMessage) != 0 {
		table.WriteRow("Error:", failureMessage)
	}

	message := messageFromMsgAndArgs(msgAndArgs...)
	if len(message) > 0 {
		table.WriteRow("Message:", message)
	}

	if callback != nil {
		callback(table)
	}

	t.Error(table)

	return false
}
