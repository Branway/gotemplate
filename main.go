// Go template command
package main

/*
Are there any examples of wanting more than one type exported from the
same package? Possibly for functional type utilities.

Could import multiple types from the same package and the builder
would do the right thing.

Path generation for generated files could do with work - args may have
spaces in, may have upper and lower case characters which will fold
together on Windows.

Detect dupliace template definitions so we don't write them multiple times

write some test

manage all the generated files - find them - delete stale ones, etg

Put comment in generated file, generated by gotemplate from xyz on date?

do replacements in comments too?
*/

import (
	"flag"
	"fmt"

	"log"
	"os"
	"path"
)

// Globals
var (
	// Flags
	verbose = flag.Bool("v", false, "Verbose - print lots of stuff")
)

// Logging function
var logf = log.Printf

// Log then fatal error
var fatalf = func(format string, args ...interface{}) {
	logf(format, args...)
	os.Exit(1)
}

// Log if -v set
func debugf(format string, args ...interface{}) {
	if *verbose {
		logf(format, args...)
	}
}

// usage prints the syntax and exists
func usage() {
	BaseName := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Syntax: %s [flags] package_name parameter\n\n"+
			"Flags:\n\n",
		BaseName)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fatalf("Need 2 arguments, package and parameters")
	}

	cwd, err := os.Getwd()
	if err != nil {
		fatalf("Couldn't get wd: %v", err)
	}

	t := newTemplate(cwd, args[0], args[1])
	t.instantiate()
}
