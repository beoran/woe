
/* A package that provides advanced logging facilities. */
package monolog

import (
    "os"
    "fmt"
    "time"
    "runtime"
    "path/filepath"
)

  
type Logger interface {
    Log(level string, file string, line int, format string, args...interface{})
    Close()
}


func (me * FileLogger) WriteLog(depth int, level string, format string, args ... interface{}) {
    _ ,file, line, ok := runtime.Caller(depth)
    
    if !ok {
        file = "unknown"
        line = 0
    }
    
    me.Log(level, file, line, format, args...)
}

func (me * FileLogger) NamedLog(name string, format string, args ... interface{}) {
    me.WriteLog(2, name, format, args...)
}

func (me * FileLogger) Info(format string, args ... interface{}) {
    me.WriteLog(2, "INFO", format, args...)
}

func (me * FileLogger) Warning(format string, args ... interface{}) {
    me.WriteLog(2, "WARNING", format, args...)
}

func (me * FileLogger) Error(format string, args ... interface{}) {
    me.WriteLog(2, "ERROR", format, args...)
}

func (me * FileLogger) Fatal(format string, args ... interface{}) {
    me.WriteLog(2, "FATAL", format, args...)
}

func (me * FileLogger) Debug(format string, args ... interface{}) {
    me.WriteLog(2, "DEBUG", format, args...)
}

  
type FileLogger struct {
    filename      string
    file        * os.File
}

func (me * FileLogger) Close() {
    if (me.file != nil) { 
        me.file.Close()
    }
    me.file = nil
}

func (me * FileLogger) Log(level string, file string, line int, format string, args...interface{}) {
    fileshort := filepath.Base(file)
    now := time.Now().Format(time.RFC3339)
    fmt.Fprintf(me.file, "%s: %s: %s: %d: ", now, level, fileshort, line)
    if args != nil && len(args) > 0 { 
        fmt.Fprintf(me.file, format, args...)
    } else {
        fmt.Fprint(me.file, format)
    }
    fmt.Fprint(me.file, "\n")
}


func NewFileLogger(filename string) (logger Logger, err error) {
    file, err := os.OpenFile(filename, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0660) 
    
    if err != nil { 
        return nil, err
    }
    
    return &FileLogger{filename, file}, nil
}    

func NewStderrLogger() (logger Logger, err error) {    
    return &FileLogger{"/dev/stderr", os.Stderr}, nil
}


func NewStdoutLogger() (logger Logger, err error) {    
    return &FileLogger{"/dev/stderr", os.Stdout}, nil
}

type Log struct {
    loggers          [] Logger
    levels  map[string] bool
}


func NewLog() * Log {
    loggers :=  make([] Logger, 32)
    levels  :=  make(map[string] bool)
    return &Log{loggers, levels}
}

func (me * Log) AddLogger(logger Logger) {
    me.loggers = append(me.loggers, logger)
}

func (me * Log) EnableLevel(level string) {
    me.levels[level] = true
}

func (me * Log) DisableLevel(level string) {
    me.levels[level] = false
}

func (me * Log) LogVa(name string, file string, line int, format string, args...interface{}) {
    _, ok := me.levels[name]
    
    if !ok {
        return
    }
    
    for _ , logger := range me.loggers {
        if (logger != nil) {
            logger.Log(name, file, line, format, args...)
        }
    }
}

func (me * Log) Close() {    
    for index , logger := range me.loggers {
        if logger != nil {
            logger.Close()
            me.loggers[index] = nil
        }
    }
}

var DefaultLog * Log

func init() {
    DefaultLog = NewLog()
    // runtime.SetFinalizer(DefaultLog, DefaultLog.Close)
}

func EnableLevel(level string) {
    DefaultLog.EnableLevel(level)
}

func DisableLevel(level string) {
    DefaultLog.DisableLevel(level)
}

func AddLogger(logger Logger, err error) {
    if err == nil {
        DefaultLog.AddLogger(logger)
    }
}


func Setup(name string, stderr bool, stdout bool) {    
    if name != "" {
        AddLogger(NewFileLogger(name))
    }
    
    if stderr { 
        AddLogger(NewStderrLogger())
    }
    
    if stdout { 
        AddLogger(NewStdoutLogger())
    }
    
    EnableLevel("INFO")
    EnableLevel("WARNING")
    EnableLevel("ERROR")
    EnableLevel("FATAL")    
}     

func Close() {
    DefaultLog.Close()
}


func WriteLog(depth int, name string, format string, args ... interface{}) {
    _ ,file, line, ok := runtime.Caller(depth)
    if !ok {
        file = "unknown"
        line = 0
    }    
    DefaultLog.LogVa(name, file, line, format, args...)
}

func NamedLog(name string, format string, args ...interface{}) {
    WriteLog(2, name, format, args)
}

func Info(format string, args ...interface{}) {
    WriteLog(2, "INFO", format, args...)
}

func Warning(format string, args ...interface{}) {
    WriteLog(2, "WARNING", format, args...)
}

func Error(format string, args ...interface{}) {
    WriteLog(2, "ERROR", format, args...)
}

func Fatal(format string, args ...interface{}) {
    WriteLog(2, "FATAL", format, args...)
}

func Debug(format string, args ...interface{}) {
    WriteLog(2, "DEBUG", format, args...)
}



