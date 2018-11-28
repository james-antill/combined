package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/james-antill/combined"
)

func main() {
	ndef := flag.Bool("no-default", false, "don't look in the default dirs")
	udir := flag.String("user-dir", "", "extra user dir")
	sdir := flag.String("sys-dir", "", "extra system dir")
	suffix := flag.String("suffix", ".conf", "suffix requirement on files")

	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: combine.d <application-name> <outputfile>\n")
		os.Exit(1)
	}
	c := &combined.Combiner{}

	if !*ndef {
		c = combined.New(flag.Arg(0))
	}

	c.SetSuffix(*suffix)

	if *udir != "" {
		c.AddUsrDir(*udir)
	}
	if *sdir != "" {
		c.AddSysDir(*sdir)
	}

	_, err := c.Write(flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)

	}
}
