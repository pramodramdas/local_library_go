package utils

import (
	"fmt"
)

// type UtilMock struct {
// }

// func (u UtilMock) SendErrorResponse(res http.ResponseWriter, msg string) {
// 	fmt.Println("sdsd")
// }

// //return true if zero value
// func (u UtilMock) IsZeroValue(x interface{}) bool {
// 	fmt.Println("yyyyyyyyyyyyy")
// 	return true
// }

type CustomLogInterfaceMock interface {
	WriteDevLogs()
	WriteLog()
}

type CustomLogMock struct {
}

func WriteLogMock(Level string, Msg string, LogPayload []byte, ExtraInfo []byte, Err error) {
	fmt.Println("ddyyyyydddd")
}

func (c *CustomLogMock) WriteDevLogs() {
	fmt.Println("dddddd")
}

func init() {
	NewCustomLog = WriteLogMock
}
