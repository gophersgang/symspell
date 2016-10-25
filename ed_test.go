package symspell

import "testing"

func TestEditMatch(t *testing.T) {
	t.Skip()
	if x := EditMatch("abcdefg", "abcedefg", 1); !x {
		t.Fail()
	}
}
