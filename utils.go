package scalaimports

import (
	"fmt"
	"os"
	"runtime"
)

func debug(v ...interface{}) {
	if !Verbose {
		return
	}
	fmt.Println(v...)
}

func debugf(f string, v ...interface{}) {
	if !Verbose {
		return
	}
	fmt.Printf(f, v...)
}

func callerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return "[" + runtime.FuncForPC(pc).Name() + "]"
}

func check(err error) {
	if err != nil {
		fmt.Println(callerName(), err)
		os.Exit(1)
	}
}
