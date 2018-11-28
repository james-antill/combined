package combined

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"

	roc "github.com/james-antill/rename-on-close"
)

// Combiner combines file given in the correct order/precedence/etc.
type Combiner struct {
	udirs  []string
	sdirs  []string
	suffix string
	mode   os.FileMode
}

// AddUsrDir adds a user dir.
func (c *Combiner) AddUsrDir(dname string) {
	c.udirs = append(c.udirs, dname)
}

// AddSysDir adds a system dir.
func (c *Combiner) AddSysDir(dname string) {
	c.sdirs = append(c.sdirs, dname)
}

// New creates a normal combiner, given an application name
func New(name string) *Combiner {
	c := &Combiner{}
	c.suffix = ".conf"
	c.AddSysDir("/usr/lib/" + name + ".d")
	c.AddUsrDir("/etc/" + name + ".d")

	return c
}

type cName struct {
	dname string
	fname string
}

func getCNames(knownFiles map[string]bool, suffix, dname string) ([]cName, error) {
	cfiles := []cName{}
	files, err := ioutil.ReadDir(dname)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		if suffix != "" && !strings.HasSuffix(f.Name(), suffix) {
			continue
		}
		if knownFiles[f.Name()] {
			continue
		}
		knownFiles[f.Name()] = true

		cfiles = append(cfiles, cName{dname, f.Name()})
	}
	return cfiles, nil
}

func countUnderbar(name string) int {
	ret := 0
	for i := range name {
		if name[i] != '_' {
			break
		}
		ret++
	}

	return ret
}

func cmpBasename(ifname, jfname string) bool {
	pi := countUnderbar(ifname)
	pj := countUnderbar(jfname)
	if pi != pj {
		return pi > pj
	}

	return ifname < jfname
}

// Files returns an ordered list of files to combine
func (c *Combiner) Files() ([]string, error) {
	knownFiles := make(map[string]bool)

	// Get list of files from user dirs, from most prefered to least.
	cfiles := []cName{}
	end := len(c.udirs) - 1
	for i := range c.udirs {
		dname := c.udirs[end-i]
		ncfiles, err := getCNames(knownFiles, c.suffix, dname)
		if err != nil {
			return nil, err
		}
		cfiles = append(cfiles, ncfiles...)
	}

	// Get list of files from system dirs, from most prefered to least.
	// Skips any files overridden by users, ie. above
	end = len(c.sdirs) - 1
	for i := range c.sdirs {
		dname := c.sdirs[end-i]
		ncfiles, err := getCNames(knownFiles, c.suffix, dname)
		if err != nil {
			return nil, err
		}
		cfiles = append(cfiles, ncfiles...)
	}

	// Sort based on basename of files
	sort.Slice(cfiles, func(i, j int) bool {
		return cmpBasename(cfiles[i].fname, cfiles[j].fname)
	})

	// Convert internal cname struct into simple file paths
	ret := []string{}
	for _, cname := range cfiles {
		pname := cname.dname + "/" + cname.fname
		ret = append(ret, pname)
	}

	return ret, nil
}

// Data returns a combined data set
func (c *Combiner) Data() ([]byte, error) {
	ret := []byte{}

	files, err := c.Files()
	if err != nil {
		return nil, err
	}

	for _, fname := range files {
		r, err := os.Open(fname)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		ret = append(ret, b...)
	}

	return ret, nil
}

// Write combines all the input files to an output file
func (c *Combiner) Write(fname string) (int, error) {

	d, err := c.Data()
	if err != nil {
		return 0, err
	}
	f, err := roc.Create(fname)
	defer f.Close() // clean up

	if _, err := f.Write(d); err != nil {
		return 0, err
	}
	if err := f.CloseRename(); err != nil {
		return 0, err
	}

	return len(d), nil
}
