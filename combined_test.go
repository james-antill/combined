package combined

import "testing"

func TestCmpBasename(t *testing.T) {
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
		{"___a", "_b", true},
		{"____a", "_b", true},
		{"_____a", "_b", true},

		{"a", "_b", false},
		{"_a", "_b", true},
		{"_a", "__b", false},
		{"_a", "___b", false},
		{"_a", "____b", false},
		{"_a", "_____b", false},
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

func TestFilterNameConf(t *testing.T) {
	data := []struct {
		knownFiles map[string]bool
		name       string
		res        bool
	}{
		{nil, "_", true},
		{nil, "._", true},
		{nil, "._.conf", true},
		{nil, "._abcd.conf", true},

		{nil, "abcd", true},
		{nil, "abcd.", true},
		{nil, "abcd.c", true},
		{nil, "abcd.co", true},
		{nil, "abcd.con", true},
		{nil, "abcd.confx", true},
		{nil, "abcd.conf", false},
		{map[string]bool{"x.conf": true}, "abcd.conf", false},
		{map[string]bool{"abcd.conf": true}, "abcd.conf", true},
		{map[string]bool{"x.conf": true, "abcd.conf": true}, "abcd.conf", true},

		{nil, "_abcd", true},
		{nil, "_abcd.conf", false},
		{map[string]bool{"abcd.conf": true}, "_abcd.conf", false},
		{map[string]bool{"_abcd.conf": true}, "_abcd.conf", true},
	}

	for n := range data {
		knownFiles := data[n].knownFiles
		name := data[n].name
		res := data[n].res

		if knownFiles == nil {
			knownFiles = make(map[string]bool)
		}

		if val := filterName(knownFiles, name, ".conf"); val != res {
			t.Errorf("data not equl: name=%v\n tst=<%v>\n got <%v>\n",
				name, res, val)
		}
	}
}

func TestFilterNameEmpty(t *testing.T) {
	data := []struct {
		knownFiles map[string]bool
		name       string
		res        bool
	}{
		{nil, "_", false},
		{nil, "._", true},
		{nil, "._.conf", true},
		{nil, "._abcd.conf", true},

		{nil, "abcd", false},
		{nil, "abcd.", false},
		{nil, "abcd.c", false},
		{nil, "abcd.co", false},
		{nil, "abcd.con", false},
		{nil, "abcd.confx", false},
		{nil, "abcd.conf", false},
		{map[string]bool{"x.conf": true}, "abcd.conf", false},
		{map[string]bool{"abcd.conf": true}, "abcd.conf", true},

		{nil, "_abcd", false},
		{nil, "_abcd.conf", false},
		{map[string]bool{"abcd.conf": true}, "_abcd.conf", false},
		{map[string]bool{"_abcd.conf": true}, "_abcd.conf", true},
	}

	for n := range data {
		knownFiles := data[n].knownFiles
		name := data[n].name
		res := data[n].res

		if knownFiles == nil {
			knownFiles = make(map[string]bool)
		}

		if val := filterName(knownFiles, name, ""); val != res {
			t.Errorf("data not equl: name=%v\n tst=<%v>\n got <%v>\n",
				name, res, val)
		}
	}
}
