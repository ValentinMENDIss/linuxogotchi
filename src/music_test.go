package main

import "testing"

func testInit(t *testing.T) {
	t.Run("Initializing Music Player", func(t *testing.T) {
		got := Init()
		want := "Hello"
		assertCorrectMessage(t, got, want)
	})
}

func assertCorrectMessage(t testing.TB, got, want string) {
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
