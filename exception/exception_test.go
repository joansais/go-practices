package exception

import (
    "testing"
    "errors"
)

var errValue = errors.New("sample error for testing purposes")

func TestCatch(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Catch: panic not recovered")
		}
	}()

	defer Catch(func (err error) {
		if err != errValue {
			t.Errorf("Catch: handler invoked with unexpected error value %q", err)
		}
	})

	Throw(errValue)
	t.Error("Throw: panic not called")
}

func TestCatchNonError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "xyz" {
				t.Errorf("Catch: re-panic with unexpected value %q", r)
			}
		}
	}()

	defer Catch(func (error) {
		t.Error("Catch: unexpected handler invocation")
	})

	panic("xyz")
}

func TestTry(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Try: panic not recovered")
		}
	}()

	err := Try(func() { Throw(errValue) })
	if err != errValue {
		t.Errorf("Try: unexpected error value %q", err)
	}

	otherErr := errors.New("foo")
	err = Try(func() { panic(otherErr) })
	if err != otherErr {
		t.Errorf("Try: unexpected error value %q", err)
	}

	err = Try(func() {})
	if err != nil {
		t.Errorf("Try: unexpected error returned %q", err)
	}
}

func TestTryNonError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "xyz" {
				t.Errorf("Try: re-panic with unexpected value %q", r)
			}
		}
	}()

	Try(func() { panic("xyz") })
	t.Error("Try: panic ignored")
}

func TestThrowIf(t *testing.T) {
	defer Catch(func (err error) {
		if err != errValue {
			t.Errorf("Catch: handler invoked with unexpected error value %q", err)
		}
	})

	ThrowIf(nil)
	ThrowIf(errValue)
	t.Error("ThrowIf: panic not called")
}

