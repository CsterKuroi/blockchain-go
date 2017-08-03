package log

import (
	"fmt"
	"testing"
)

func Test_log(t *testing.T) {
	var logType ErrorType
	logType = API_ERROR
	var logErrorMsg string
	logErrorMsg = "Token is error!"
	Error(fmt.Sprintf("[%s][%s]", logType, logErrorMsg))

	logType = DEBUG_NO_ERROR
	logErrorMsg = "debug no error!"
	Debug(fmt.Sprintf("[%s][%s]", logType, logErrorMsg))
}

type TestStruct struct {
	a   int
	msg string
}

func (t TestStruct) String() string {
	return fmt.Sprintf("TestStruct is { a = %d , msg = %s }", t.a, t.msg)
}

func (t TestStruct) func1() {
	Debug(t)
}

func (t *TestStruct) func2() {
	Debug(t)
}

func Test_log1(t *testing.T) {
	test_struct := TestStruct{a: 100, msg: "hello world!"}
	Debug("--------------------------func1--------------------------")
	test_struct.func1()
	Debug("--------------------------func1--------------------------")
	test_struct.func2()

	t.Log("ok")
}
