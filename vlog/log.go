package vlog

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var verbose bool

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose mode")
}

// Copy from log lib
var (
	SetOutput = log.SetOutput
	SetFlags  = log.SetFlags
	SetPrefix = log.SetPrefix
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

func Output(calldepth int, s string) {
	log.Output(calldepth+1, s)
}

func output(ok bool, calldepth int, s string) {
	if !ok {
		return
	}
	log.Output(calldepth+1, s)
}

func Verbose(args ...interface{}) {
	output(verbose, 2, fmt.Sprint(args...))
}

func Verbosef(format string, args ...interface{}) {
	output(verbose, 2, fmt.Sprintf(format, args...))
}

func Verboseln(args ...interface{}) {
	output(verbose, 2, fmt.Sprintln(args...))
}

func Print(args ...interface{}) {
	output(true, 2, fmt.Sprint(args...))
}

func Printf(format string, args ...interface{}) {
	output(true, 2, fmt.Sprintf(format, args...))
}

func Println(args ...interface{}) {
	output(true, 2, fmt.Sprintln(args...))
}

func Panic(args ...interface{}) {
	s := fmt.Sprint(args...)
	output(true, 2, s)
	panic(s)
}

func Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	output(true, 2, s)
	panic(s)
}

func Panicln(args ...interface{}) {
	s := fmt.Sprintln(args...)
	output(true, 2, s)
	panic(s)
}

func Fatal(args ...interface{}) {
	output(true, 2, fmt.Sprint(args...))
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	output(true, 2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Fatalln(args ...interface{}) {
	output(true, 2, fmt.Sprintln(args...))
	os.Exit(1)
}
