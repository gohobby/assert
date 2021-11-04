package assert

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
)

// isNil checks if a specified object is nil or not, without Failing.
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)

	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.Slice:
		if value.IsNil() {
			return true
		}
	}

	return false
}

func isFunction(arg interface{}) bool {
	if arg == nil {
		return false
	}

	return reflect.TypeOf(arg).Kind() == reflect.Func
}

// validateEqualArgs checks if the supplied arguments can be safely used in the
// Equal/NotEqual functions.
func validateEqualArgs(expected, actual interface{}) error {
	if expected == nil && actual == nil {
		return nil
	}

	if isFunction(expected) || isFunction(actual) {
		return errors.New("cannot take func type as argument")
	}

	return nil
}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}

	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}

		return fmt.Sprintf("%+v", msg)
	}

	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}

	return ""
}

// StackTrace
func StackTrace(skip int) []string {
	var callers []string

	re := regexp.MustCompile(`._test\.go$`)

	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// Do not print the stack line if it is not in a test file.
		matched := re.MatchString(file)
		if !matched {
			continue
		}

		callers = append(callers, fmt.Sprintf("%s:%d", filepath.Base(file), line))
	}

	return callers
}
