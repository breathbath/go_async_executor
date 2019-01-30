package logger

import (
	"fmt"
	"log"
)

var IsVerbose = true

func Log(msg string, args ...interface{}) {
	if !IsVerbose {
		return
	}

	log.Printf(msg, args...)
}

func LogError(err error, msg string) {
	if msg != "" {
		log.Printf("%s: %v\n", msg, err)
	} else {
		log.Println(err)
	}
}

func LogErrorFromMsg(msg string, args ...interface{}) {
	LogError(fmt.Errorf(msg, args...), "")
}
