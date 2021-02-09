package main

import (
	"os"
	"testing"
	"time"
)

func setDisplay(val string) func() {
	old := os.Getenv("DISPLAY")
	os.Setenv("DISPLAY", val)
	return func() {
		os.Setenv("DISPLAY", old)
	}
}

func TestGetCurrentColorTemp(t *testing.T) {
	defer setDisplay("testdisplay" + time.Now().String())()

	temp, err := getCurrentColorTemp()
	if err != nil {
		t.Fatal(err)
	}
	if temp != *dayTemp {
		t.Fatalf("expected default day temp")
	}

	for _, want := range []int{3000, 6500} {
		if err := saveCurrentColorTemp(want); err != nil {
			t.Fatal(err)
		}
		temp, err := getCurrentColorTemp()
		if err != nil {
			t.Fatal(err)
		}
		if temp != want {
			t.Fatalf("expected %d temp", want)
		}
	}
}

func TestInterpolate(t *testing.T) {
	totalTime = 100 * time.Millisecond

	for _, want := range []int{3000, 6500} {
		if err := interpolateColorTemp(want); err != nil {
			t.Fatal(err)
		}
		temp, err := getCurrentColorTemp()
		if err != nil {
			t.Fatal(err)
		}
		if temp != want {
			t.Fatalf("expected %d temp", want)
		}
	}
}
