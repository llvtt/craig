package utils

import "errors"

func PanicOnErr(args ...interface{}) {
	if len(args) < 1 {
		return
	}
	if err, ok := args[len(args)-1].(error); ok {
		panic(err)
	}
}

func makeError(msg string) error {
	return errors.New(msg)
}

func WrapError(msg string, err error) error {
	return errors.New(msg + " --> Caused by: "+err.Error())
}
