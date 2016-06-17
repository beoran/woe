
/* A package that provides advanced logging facilities. */
package monolog

import (
    "os"
    "fmt"
    "time"
    "runtime"
    "path/filepath"
    "unicode"
    "strings"
)

  
type Logger interface {
    Log(level string, file string, line int, format string, args...interface{})
    Close()
}


func GetCallerName(depth int) {
    pc := make([]uintptr, depth+1)
    runtime.Callers(depth, pc)
    for i, v := range pc {
        fun := runtime.FuncForPC(v)
        if fun != nil { 
            fmt.Printf("GetCallerName %d %s\n", i, fun.Name()) 
        } else {
              fmt.Printf("GetCallerName %d nil\n", i) 
        }
    }
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

type Logbook struct {
    loggers          [] Logger
    levels  map[string] bool
}


func NewLog() * Logbook {
    loggers :=  make([] Logger, 32)
    levels  :=  make(map[string] bool)
    return &Logbook{loggers, levels}
}

func (me * Logbook) AddLogger(logger Logger) {
    me.loggers = append(me.loggers, logger)
}

func (me * Logbook) EnableLevel(level string) {
    me.levels[level] = true
}

func (me * Logbook) DisableLevel(level string) {
    me.levels[level] = false
}

func enableDisableSplitter(c rune) (bool) {
        ok :=  (!unicode.IsLetter(c))
        ok = ok && (!unicode.IsNumber(c))
        ok = ok && (c != '_')
        return ok
}


func (me * Logbook) EnableLevels(list string) {    
    to_enable := strings.FieldsFunc(list, enableDisableSplitter)
    for _, level := range to_enable {
        me.EnableLevel(level)
    }
}

    
func (me * Logbook) DisableLevels(list string) {    
    to_disable := strings.FieldsFunc(list, enableDisableSplitter)
    for _, level := range to_disable {
        me.DisableLevel(level)
    }
}


func (me * Logbook) LogVa(name string, file string, line int, format string, args...interface{}) {
    enabled, ok := me.levels[name]
    
    if (!ok) || (!enabled) {
        return
    }
    
    for _ , logger := range me.loggers {
        if (logger != nil) {
            logger.Log(name, file, line, format, args...)
        }
    }
}

func (me * Logbook) Close() {    
    for index , logger := range me.loggers {
        if logger != nil {
            logger.Close()
            me.loggers[index] = nil
        }
    }
}

var DefaultLog * Logbook

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

func EnableLevels(list string) {    
  DefaultLog.EnableLevels(list)
}

func DisableLevels(list string) {    
  DefaultLog.DisableLevels(list)
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

func Log(name string, format string, args ...interface{}) {
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

func エラ(err error) {
    WriteLog(2, "ERROR", "%s", err.Error())
}

func WriteError(err error) {
    WriteLog(2, "ERROR", "%s", err.Error())
}

