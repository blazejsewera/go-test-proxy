package log

import (
	"os"

	"github.com/blazejsewera/go-test-proxy/colorfmt"
)

var instance = colorfmt.New(false, os.Stdout, os.Stderr)

func SetFmt(newInstance *colorfmt.Fmt) {
	instance = newInstance
}

func Printf(format string, v ...any) {
	instance.Cerrprintf(colorfmt.Faint, colorfmt.Base, format, v...)
}

func Fatalf(format string, v ...any) {
	instance.Cerrprintf(colorfmt.Normal, colorfmt.BrightRed, format, v...)
	os.Exit(1)
}

func Fatalln(v any) {
	Fatalf("%v\n", v)
}
