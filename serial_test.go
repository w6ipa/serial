package serial

import (
	"testing"
)

func TestOpen(t *testing.T) {
	c := Config{
		Name:     "/dev/tty.Bluetooth-Incoming-Port",
		Baud:     9600,
		StopBits: 2,
	}
	if _, err := OpenPort(&c); err != nil {
		t.Fatal(err)
	}
}
