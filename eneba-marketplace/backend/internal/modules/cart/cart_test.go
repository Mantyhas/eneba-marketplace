package cart

import "testing"

func TestValidate(t *testing.T) {
	for _, c := range []struct {
		items []Item
		valid bool
	}{{nil, false}, {[]Item{{1, 1}}, true}, {[]Item{{1, 0}}, false}, {[]Item{{1, 11}}, false}, {[]Item{{1, 1}, {1, 2}}, false}, {[]Item{{1, 1}, {2, 2}}, true}} {
		if (Validate(c.items) == nil) != c.valid {
			t.Fatalf("unexpected result: %+v", c)
		}
	}
}
