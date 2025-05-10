package util

import (
	"fmt"
	"reflect"

	"github.com/rogpeppe/go-internal/diff"
	"github.com/sanity-io/litter"
)

var LitterConfig = litter.Options{
	HidePrivateFields: true,
	HideZeroValues:    true,
}

func Format[T any](args ...T) string {
	if len(args) == 1 {
		return LitterConfig.Sdump(args[0])
	}
	a := make([]any, len(args))
	for i, arg := range args {
		a[i] = arg
	}
	return LitterConfig.Sdump(a...)
}

func NamedDiff(xname, yname string, x, y any) string {
	return string(diff.Diff(
		xname,
		[]byte(fmt.Sprint(x)),
		yname,
		[]byte(fmt.Sprint(y)),
	))
}

func Diff(x, y any) string {
	return string(diff.Diff(
		reflect.TypeOf(x).String(),
		[]byte(fmt.Sprint(x)),
		reflect.TypeOf(y).String(),
		[]byte(fmt.Sprint(y)),
	))
}
