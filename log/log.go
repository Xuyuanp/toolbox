package log

import (
	"flag"
	"fmt"
	"log"
	"os"
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

// Output calls log.Output by calldepth inc.
func Output(calldepth int, s string) {
	log.Output(calldepth+1, s)
}

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

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	Output(2, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Output(2, fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	Output(2, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	Output(2, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	Output(2, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	Output(2, s)
	panic(s)
}
