package util

import (
	"fmt"
	"log"
	"runtime/debug"
)

// ErrCheck - saves a bunch of typing, prints error if it exists
//            and provides a traceback as well
func ErrCheck(err error) {
	if err != nil {
		fmt.Printf("error = %v\n", err)
		debug.PrintStack()
		log.Fatal(err)
	}
}

// LogAndPrintError encapsulates logging and printing an error.
// Note that the error is printed only if the environment is NOT production.
func LogAndPrintError(funcname string, err error) {
	errmsg := fmt.Sprintf("%s: err = %v\n", funcname, err)
	Ulog(errmsg)
	fmt.Println(errmsg)
}
