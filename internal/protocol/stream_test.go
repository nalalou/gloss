package protocol

import (
	"strings"
	"testing"
)

func TestReadLinesLongLine(t *testing.T) {
	long := strings.Repeat("x", 500_000)
	var got []string
	err := ReadLines(strings.NewReader(long+"\r\n::ok after\nno newline"), func(l string) { got = append(got, l) })
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0] != long || got[1] != "::ok after" || got[2] != "no newline" {
		t.Errorf("got %d lines, first len %d", len(got), len(got[0]))
	}
}

func TestMasker(t *testing.T) {
	var m Masker
	if got := m.Apply("before hunter2"); got != "before hunter2" {
		t.Errorf("masked too early: %q", got)
	}
	m.Apply("::add-mask::hunter2")
	if got := m.Apply("pw=hunter2 again hunter2"); got != "pw=*** again ***" {
		t.Errorf("got %q", got)
	}
}
