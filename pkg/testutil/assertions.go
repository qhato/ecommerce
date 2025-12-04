package testutil

import (
	"reflect"
	"testing"
	"time"
)

// AssertEqual checks if two values are equal
func AssertEqual(t *testing.T, got, want interface{}, msg string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s: got %v, want %v", msg, got, want)
	}
}

// AssertNotEqual checks if two values are not equal
func AssertNotEqual(t *testing.T, got, want interface{}, msg string) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		t.Errorf("%s: got %v, want different value", msg, got)
	}
}

// AssertNil checks if value is nil
func AssertNil(t *testing.T, got interface{}, msg string) {
	t.Helper()
	if !isNil(got) {
		t.Errorf("%s: got %v, want nil", msg, got)
	}
}

// AssertNotNil checks if value is not nil
func AssertNotNil(t *testing.T, got interface{}, msg string) {
	t.Helper()
	if isNil(got) {
		t.Errorf("%s: got nil, want non-nil value", msg)
	}
}

// AssertNoError checks if error is nil
func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: unexpected error: %v", msg, err)
	}
}

// AssertError checks if error is not nil
func AssertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: expected error, got nil", msg)
	}
}

// AssertErrorContains checks if error contains specific text
func AssertErrorContains(t *testing.T, err error, want string, msg string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: expected error containing %q, got nil", msg, want)
		return
	}
	if !contains(err.Error(), want) {
		t.Errorf("%s: error %q does not contain %q", msg, err.Error(), want)
	}
}

// AssertTrue checks if condition is true
func AssertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("%s: got false, want true", msg)
	}
}

// AssertFalse checks if condition is false
func AssertFalse(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Errorf("%s: got true, want false", msg)
	}
}

// AssertTimeAlmostEqual checks if two times are within delta
func AssertTimeAlmostEqual(t *testing.T, got, want time.Time, delta time.Duration, msg string) {
	t.Helper()
	diff := got.Sub(want)
	if diff < 0 {
		diff = -diff
	}
	if diff > delta {
		t.Errorf("%s: times differ by %v, want within %v", msg, diff, delta)
	}
}

// AssertLen checks if slice/map has expected length
func AssertLen(t *testing.T, got interface{}, want int, msg string) {
	t.Helper()
	v := reflect.ValueOf(got)
	if v.Len() != want {
		t.Errorf("%s: got length %d, want %d", msg, v.Len(), want)
	}
}

// Helper functions

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	k := v.Kind()
	return (k == reflect.Chan || k == reflect.Func || k == reflect.Interface ||
		k == reflect.Map || k == reflect.Ptr || k == reflect.Slice) && v.IsNil()
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
