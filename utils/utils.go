package utils

func PanicOnErr(args ...interface{}) {
	if len(args) < 1 {
		return
	}
	if err, ok := args[len(args)-1].(error); ok {
		panic(err)
	}
}
