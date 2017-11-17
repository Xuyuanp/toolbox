package log

import (
	"flag"
	"fmt"
	"log"
)

var verbose bool

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
}

// Copy from log lib
var (
	SetOutput = log.SetOutput
	SetFlags  = log.SetFlags
	SetPrefix = log.SetPrefix
	Prefix    = log.Prefix
	Flags     = log.Flags
	Output    = log.Output

	Print   = log.Print
	Printf  = log.Printf
	Println = log.Println

	Fatal   = log.Fatal
	Fatalf  = log.Fatalf
	Fatalln = log.Fatalln

	Panic   = log.Panic
	Panicf  = log.Panicf
	Panicln = log.Panicln
)

// Copy from log lib
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// Verbose calls Output to print to the standard logger only if verbose is true.
// Arguments are handled in the manner of fmt.Print.
func Verbose(args ...interface{}) {
	if !verbose {
		return
	}
	Output(2, fmt.Sprint(args...))
}

// Verbosef calls Output to print to the standard logger only if verbose is true.
// Arguments are handled in the manner of fmt.Printf.
func Verbosef(format string, args ...interface{}) {
	if !verbose {
		return
	}
	Output(2, fmt.Sprintf(format, args...))
}

// Verboseln calls Output to print to the standard logger only if verbose is true.
// Arguments are handled in the manner of fmt.Println.
func Verboseln(args ...interface{}) {
	if !verbose {
		return
	}
	Output(2, fmt.Sprintln(args...))
}
