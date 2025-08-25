package lib

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var onceLog sync.Once
var instanceLog *Log

type TypePrint int

const (
	LG_Silent  TypePrint = 0
	LG_Print   TypePrint = 1
	LG_PrintLn TypePrint = 2
)

type NivelLog int

const (
	N_NENHUM   NivelLog = 0
	N_DEBUG    NivelLog = 1
	N_INFO     NivelLog = 2
	N_WARNING  NivelLog = 3
	N_ERROR    NivelLog = 4
	N_CRITICAL NivelLog = 5
)

type Log struct {
	AppName   string
	TypePrint TypePrint
	LogNivel  NivelLog
	Log       map[string]*LogEmp
	mutex     sync.Mutex
}
type LogEmp struct {
	File *os.File
	Hora string
}

func NewLog() *Log {
	onceLog.Do(
		func() {
			instanceLog = &Log{
				AppName:   "Init",
				TypePrint: LG_Silent,
				LogNivel:  N_ERROR,
				Log:       make(map[string]*LogEmp),
			}
		})
	return instanceLog
}

func NewLogApp(app string) *LogEmp {
	timestamp := time.Now().Format("02-01-2006-15")
	hora := time.Now().Format("15")
	file, _ := os.OpenFile(fmt.Sprintf("log_%s_%v.log", app, timestamp), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	return &LogEmp{
		File: file,
		Hora: hora,
	}
}

func (lg *Log) Print(msg ...interface{}) {
	switch lg.TypePrint {
	case LG_Silent:
	case LG_Print:
		fmt.Print(msg)
	case LG_PrintLn:
		fmt.Println(msg)
	default:
		fmt.Print(msg)
	}
}
func (lg *Log) WriteLog(app string, setTime bool, msg ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	hora := time.Now().Format("15")
	lg.mutex.Lock()
	if _, ok := lg.Log[app]; !ok {
		lg.Log[app] = NewLogApp(lg.AppName)
	}

	if hora != lg.Log[app].Hora {
		lg.Log[app].File.Close()
		lg.Log[app] = NewLogApp(lg.AppName)
	}
	msg_print := ""
	if setTime {
		msg_print = fmt.Sprint(timestamp, msg)
	} else {
		msg_print = fmt.Sprint(msg)
	}

	lg.Print(msg_print)
	lg.Log[app].File.WriteString(msg_print + "\n")
	lg.mutex.Unlock()
}
func (lg *Log) Debug(app string, v ...interface{}) {
	if lg.LogNivel > N_NENHUM && lg.LogNivel <= N_CRITICAL {
		lg.WriteLog(app, true, "Debug:", v)
	}
}
func (lg *Log) Info(app string, v ...interface{}) {
	if lg.LogNivel >= N_DEBUG && lg.LogNivel <= N_CRITICAL {
		lg.WriteLog(app, true, "Info:", v)
	}
}
func (lg *Log) Screen(app string, v ...interface{}) {

	fmt.Println(v)
}
func (lg *Log) Error(app string, v ...interface{}) {
	if lg.LogNivel >= N_WARNING && lg.LogNivel <= N_CRITICAL {
		lg.WriteLog(app, true, "Error:", v)
	}
}
func (lg *Log) Warning(app string, v ...interface{}) {
	if lg.LogNivel >= N_INFO && lg.LogNivel <= N_CRITICAL {
		lg.WriteLog(app, true, "Warning:", v)
	}
}
func (lg *Log) Critical(app string, v ...interface{}) {
	if lg.LogNivel >= N_ERROR && lg.LogNivel <= N_CRITICAL {
		lg.WriteLog(app, true, "Critical:", v)
	}
}
func (lg *Log) Fatal(app string, v ...interface{}) {
	lg.WriteLog(app, true, "Fatal:", v)
	os.Exit(1)
}
