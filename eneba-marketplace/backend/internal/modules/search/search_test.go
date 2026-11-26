package search

import "testing"

func TestPattern(t *testing.T) {
	if got := Pattern("  GTA_5% "); got != `%gta\_5\%%` {
		t.Fatal(got)
	}
}
