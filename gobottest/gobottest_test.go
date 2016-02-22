package gobottest

import "testing"

func TestAssert(t *testing.T) {
	err := ""
	errFunc = func(t *testing.T, message string) {
		err = message
	}

	Assert(t, 1, 1)
	if err != "" {
		t.Errorf("Assert failed: 1 should equal 1")
	}

	Assert(t, 1, 2)
	if err == "" {
		t.Errorf("Assert failed: 1 should not equal 2")
	}
}

func TestRefute(t *testing.T) {
	err := ""
	errFunc = func(t *testing.T, message string) {
		err = message
	}

	Refute(t, 1, 2)
	if err != "" {
		t.Errorf("Refute failed: 1 should not be 2")
	}

	Refute(t, 1, 1)
	if err == "" {
		t.Errorf("Refute failed: 1 should not be 1")
	}
}
