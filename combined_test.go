package combined

import "testing"

func TestBasename(t *testing.T) {
	data := []struct {
		i       string
		j       string
		resLess bool
	}{
		{"_", "_", false},
		{"a", "a", false},
		{"__a", "__a", false},

		{"a", "b", true},
		{"_a", "b", true},
		{"_a", "_b", true},
		{"__a", "_b", true},

		{"a", "_b", false},
		{"_a", "_b", true},
		{"_a", "__b", false},
	}

	for n := range data {
		i := data[n].i
		j := data[n].j
		resLess := data[n].resLess

		if val := cmpBasename(i, j); val != resLess {
			t.Errorf("data not equl: %v %v\n tst=<%v>\n got <%v>\n",
				i, j, resLess, val)
		}
	}

}
