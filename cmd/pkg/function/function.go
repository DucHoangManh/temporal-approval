package function

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(f any) string {
	functionDescription := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	functionTruncatedName := strings.Split(functionDescription, ".")
	if len(functionTruncatedName) == 2 {
		return functionTruncatedName[1]
	}
	return ""
}
