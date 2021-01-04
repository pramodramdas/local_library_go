package utils

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// type utilF UtilMembers

// var utilFunc *utilF

// func init() {
// 	utilFunc = &utilF{UtilInterface: Util{}}
// }

func TestLog() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("hello world")
	log.Info().Str("foo", "bar").Send()

	// log.Info().Str("foo", "bar").
	// 	Msgf("Cannot start %s", service)

	empJson := `{
        "id" : 11,
        "name" : "Irshad",
        "department" : "IT",
        "designation" : "Product Manager"
	}`

	var result map[string]interface{}
	json.Unmarshal([]byte(empJson), &result)

	log.Info().RawJSON("json", []byte(empJson)).Send()
	// c := CustomLog{
	// 	Level:      "fatal",
	// 	Msg:        "abc",
	// 	LogPayload: []byte(empJson),
	// 	ExtraInfo:  []byte(empJson),
	// 	Err:        errors.New("A repo man spends his life getting into tense situations"),
	// }
	//CustomLogFactoryStruct.NewCustomLog("info", "", make([]byte, 0), make([]byte, 0), nil).WriteLog()
	//fmt.Println(Log.logger)
	NewCustomLog("info", "", make([]byte, 0), make([]byte, 0), nil)
}

type CustomLog struct {
	Level      string
	Msg        string
	LogPayload []byte
	ExtraInfo  []byte
	Err        error
}

type CustomLogInterface interface {
	WriteDevLogs()
	WriteLog()
}

var NewCustomLog LogFunc

type LogFunc func(Level string, Msg string, LogPayload []byte, ExtraInfo []byte, Err error)

func init() {
	NewCustomLog = WriteLog
}

func WriteLog(Level string, Msg string, LogPayload []byte, ExtraInfo []byte, Err error) {
	c := CustomLog{Level, Msg, LogPayload, ExtraInfo, Err}
	if os.Getenv("GO_ENV") == "production" { //production
		c.WriteDevLogs()
	} else { //development
		c.WriteDevLogs()
	}
}

func (c *CustomLog) WriteDevLogs() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var logRef *zerolog.Event
	if c.Level == "info" {
		logRef = log.Info()
	} else if c.Level == "warn" {
		logRef = log.Info()
	} else if c.Level == "error" {
		logRef = log.Info()
	} else if c.Level == "fatal" {
		logRef = log.Fatal()
	}

	if UtilStruct.IsZeroValue(c.Msg) == false {
		logRef.Str("message", c.Msg)
	}
	if UtilStruct.IsZeroValue(c.LogPayload) == false {
		logRef.RawJSON("LogPayload", c.LogPayload)
	}
	if UtilStruct.IsZeroValue(c.ExtraInfo) == false {
		logRef.RawJSON("ExtraInfo", c.ExtraInfo)
	}
	if c.Err != nil {
		logRef.Err(c.Err)
	}

	logRef.Send()
}
