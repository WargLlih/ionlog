package recordhistory

import "testing"

func TestLogOnce(t *testing.T) {
	t.Run("First Log", func(t *testing.T) {
		r := NewRecordHistory()
		if !LogOnce(r, "pkg", "function", "file", 1, "msg") {
			t.Errorf("LogOnce() failed")
		}
	})

	t.Run("Two Logs Check", func(t *testing.T) {
		r := NewRecordHistory()

		if !LogOnce(r, "pkg", "function", "file", 1, "msg") {
			t.Errorf("LogOnce() failed")
		}

		LogOnce(r, "pkg", "function", "file", 1, "msg")
		if LogOnce(r, "pkg", "function", "file", 1, "msg") {
			t.Errorf("LogOnce() failed")
		}
	})
}

func TestLogOnChange(t *testing.T) {
	t.Run("On Change Check", func(t *testing.T) {
		r := NewRecordHistory()

		if !LogOnChange(r, "pkg", "function", "file", 1, "msg") {
			t.Errorf("LogOnChange() failed")
		}

		if LogOnChange(r, "pkg", "function", "file", 1, "msg") {
			t.Errorf("LogOnChange() failed")
		}

		if !LogOnChange(r, "pkg", "function", "file", 1, "msg 2") {
			t.Errorf("LogOnChange() failed")
		}
	})
}
