package utils

import (
	"log"
	"reflect"
	"runtime"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func LogPrintln(funcName string, v string) {
	log.Printf("[LOG - %s]    %s\n", funcName, v)
}

func LogError(funcName string, v ...interface{}) {
	log.Printf("[ERR - %s]    %s\n", funcName, v)
}
