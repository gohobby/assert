package assert

import (
	"fmt"
	"testing"
)

// AssertionTesterInterface defines an interface to be used for testing assertion methods
type AssertionTesterInterface interface {
	TestMethod()
}

// AssertionTesterConformingObject is an object that conforms to the AssertionTesterInterface interface
type AssertionTesterConformingObject struct{}

func (a *AssertionTesterConformingObject) TestMethod() {}

// AssertionTesterNonConformingObject is an object that does not conform to the AssertionTesterInterface interface
type AssertionTesterNonConformingObject struct{}

func TestEqual(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	type myType string

	var m map[string]interface{}

	cases := []struct {
		expected interface{}
		actual   interface{}
		result   bool
		remark   string
	}{
		{"Hello World", "Hello World", true, ""},
		{123, 123, true, ""},
		{123.5, 123.5, true, ""},
		{[]byte("Hello World"), []byte("Hello World"), true, ""},
		{nil, nil, true, ""},
		{int32(123), int32(123), true, ""},
		{uint64(123), uint64(123), true, ""},
		{float32(123), float32(123), true, ""},
		{myType("1"), myType("1"), true, ""},
		{&struct{}{}, &struct{}{}, true, "pointer equality is based on equality of underlying value"},

		// Not expected to be equal
		{m["bar"], "something", false, ""},
		{float32(1), float64(1), false, ""},
		{float32(1), float32(2), false, ""},
		{nil, float32(1), false, ""},
		{float32(1), nil, false, ""},
		{myType("1"), myType("2"), false, ""},
		{func() int { return 23 }, func() int { return 24 }, false, ""},

		// A case that might be confusing, especially with numeric literals
		{10, uint(10), false, ""},
	}

	for _, c := range cases {
		c := c
		t.Run(fmt.Sprintf("Equal(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := Equal(mockT, c.expected, c.actual)

			if res != c.result {
				t.Errorf("Equal(%#v, %#v) should return %#v: %s", c.expected, c.actual, c.result, c.remark)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	cases := []struct {
		expected interface{}
		actual   interface{}
		result   bool
	}{
		// cases that are expected not to match
		{"Hello World", "Hello World!", true},
		{123, 1234, true},
		{123.5, 123.55, true},
		{[]byte("Hello World"), []byte("Hello World!"), true},
		{nil, new(AssertionTesterConformingObject), true},

		// cases that are expected to match
		{nil, nil, false},
		{"Hello World", "Hello World", false},
		{123, 123, false},
		{123.5, 123.5, false},
		{[]byte("Hello World"), []byte("Hello World"), false},
		{new(AssertionTesterConformingObject), new(AssertionTesterConformingObject), false},
		{&struct{}{}, &struct{}{}, false},
		{func() int { return 23 }, func() int { return 24 }, false},
		// A case that might be confusing, especially with numeric literals
		{int(10), uint(10), true},
	}

	for _, c := range cases {
		c := c
		t.Run(fmt.Sprintf("NotEqual(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := NotEqual(mockT, c.expected, c.actual)

			if res != c.result {
				t.Errorf("NotEqual(%#v, %#v) should return %#v", c.expected, c.actual, c.result)
			}
		})
	}
}

func TestTrue(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	if !True(mockT, true) {
		t.Error("True should return true")
	}

	if True(mockT, false) {
		t.Error("True should return false")
	}
}

func TestFalse(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	if !False(mockT, false) {
		t.Error("False should return true")
	}

	if False(mockT, true) {
		t.Error("False should return false")
	}
}

func TestNil(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	if !Nil(mockT, nil) {
		t.Error("Nil should return true: object is nil")
	}

	if !Nil(mockT, (*struct{})(nil)) {
		t.Error("Nil should return true: object is (*struct{})(nil)")
	}

	if Nil(mockT, new(AssertionTesterConformingObject)) {
		t.Error("Nil should return false: object is not nil")
	}
}

func TestNotNil(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	if !NotNil(mockT, new(AssertionTesterConformingObject)) {
		t.Error("NotNil should return true: object is not nil")
	}

	if NotNil(mockT, nil) {
		t.Error("NotNil should return false: object is nil")
	}

	if NotNil(mockT, (*struct{})(nil)) {
		t.Error("NotNil should return false: object is (*struct{})(nil)")
	}
}

func TestImplements(t *testing.T) {
	mockT := new(testing.T)

	if !Implements(mockT, (*AssertionTesterInterface)(nil), new(AssertionTesterConformingObject)) {
		t.Error("Implements method should return true: AssertionTesterConformingObject implements AssertionTesterInterface")
	}

	if Implements(mockT, (*AssertionTesterInterface)(nil), new(AssertionTesterNonConformingObject)) {
		t.Error("Implements method should return false: AssertionTesterNonConformingObject does not implements AssertionTesterInterface")
	}

	if Implements(mockT, (*AssertionTesterInterface)(nil), nil) {
		t.Error("Implements method should return false: nil does not implement AssertionTesterInterface")
	}
}
