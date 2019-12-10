package hello

import "testing"

func TestHello(t *testing.T) {
	want := "HelloWorld"
	if got := HelloWorld(); got != want {
		t.Errorf("HelloWorld() = %q, want %q", got, want)
	}
}
