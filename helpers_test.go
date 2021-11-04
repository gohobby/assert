package assert

import "testing"

func Test_validateEqualArgs(t *testing.T) {
	if validateEqualArgs(func() {}, func() {}) == nil {
		t.Error("non-nil functions should error")
	}

	if validateEqualArgs(func() {}, func() {}) == nil {
		t.Error("non-nil functions should error")
	}

	if validateEqualArgs(nil, nil) != nil {
		t.Error("nil functions are equal")
	}
}
