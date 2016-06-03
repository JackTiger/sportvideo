package youlog

import (
    "log"
	"runtime"
    "strconv"
    "os"
    "os/exec"
    "path/filepath"
	"strings"
    "time"
)

var appStart = false
var logfilepath = ""
var debugLog * log.Logger

//SetLogFileName SetLogFileName
func SetLogFileName(filename string) {

    appPath := appName()
    appdir := filepath.Dir(appPath)
    logdir := appdir + string(filepath.Separator) + "log"
    os.MkdirAll(logdir, os.ModePerm)  

    now := time.Now()
    strNow := now.Format("2006-01-02")
    logfilepath := logdir + string(filepath.Separator) + filename + "." + strNow + ".log"
    logFile, _ := os.OpenFile(logfilepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm|os.ModeTemporary)
    if logFile != nil {
        debugLog = log.New(logFile, "",log.LstdFlags)
    } else {
        debugLog = nil
    }
}

//Info Info
func Info(strContent string) {
	
    funcName, _, line, ok := runtime.Caller(1)
    var message = ""
    if ok {
        message = "[Info]" + "Function: " + runtime.FuncForPC(funcName).Name() + ", Line: " + strconv.Itoa(line) + ", Message: " + strContent
    } else {
        message = "[Info]" + " Message: " + strContent
    }
    
    print(message)
}

//Debug Debug
func Debug(strContent string) {
	
    funcName, file, line, ok := runtime.Caller(1)
    var message = ""
    if ok {
        message = "[Debug]" + " File: " + file + ", Function: " + runtime.FuncForPC(funcName).Name() + ", Line: " + strconv.Itoa(line) + ", Message: " + strContent
    } else {
        message = "[Debug]" + " Message: " + strContent
    }
    
    print(message)
}

//Warnning Warnning
func Warnning(strContent string) {
	
    funcName, file, line, ok := runtime.Caller(1)
    var message = ""
    if ok {
        message = "[Warnning]" + " File: " + file + ", Function: " + runtime.FuncForPC(funcName).Name() + ", Line: " + strconv.Itoa(line) + ", Message: " + strContent
    } else {
        message = "[Warnning]" + " Message: " + strContent
    }
    
    print(message)
}

//Trace Trace
func Trace(strContent string) {
	
    funcName, file, line, ok := runtime.Caller(1)
    var message = ""
    if ok {
        message = "[Trace]" + " File: " + file + ", Function: " + runtime.FuncForPC(funcName).Name() + ", Line: " + strconv.Itoa(line) + ", Message: " + strContent
    } else {
        message = "[Trace]" + " Message: " + strContent
    }
    
    print(message)
}

//Error Error
func Error(strContent string) {
	
    funcName, file, line, ok := runtime.Caller(1)
    var message = ""
    if ok {
        message = "[Error]" + " File: " + file + ", Function: " + runtime.FuncForPC(funcName).Name() + ", Line: " + strconv.Itoa(line) + ", Message: " + strContent
    } else {
        message = "[Error]" + " Message: " + strContent
    }
    
    fatal(message)
}

//Fatal Fatal
func Fatal(strContent string) {
	
    funcName, file, line, ok := runtime.Caller(1)
    var message = ""
    if ok {
        message = "[Fatal]" + " File: " + file + ", Function: " + runtime.FuncForPC(funcName).Name() + ", Line: " + strconv.Itoa(line) + ", Message: " + strContent
    } else {
        message = "[Fatal]" + " Message: " + strContent
    }
    
    fatal(message)
}

func print(message string) {
    
    log.Println(message)
    exportToFile(message)
}

func fatal(message string) {
    
    log.Fatalln(message)    
    exportToFile(message)
}

func exportToFile(message string) {
    
    if debugLog != nil {
        debugLog.Println(message)
    }
}

func appName() string {
    
    file, _ := exec.LookPath(os.Args[0])
    path, _ := filepath.Abs(file)
    return path
}

func logfileName() string {
    
    appPath := appName()
    index := strings.LastIndex(appPath, ".")
    filename := substr(appPath, 0, index)
    if index == -1 {
        filename = appPath
    }

    return filename
}

func substr(str string, start, length int) string {
    rs := []rune(str)
    rl := len(rs)
    end := 0

    if start < 0 {
        start = rl - 1 + start
    }
    end = start + length

    if start > end {
        start, end = end, start
    }

    if start < 0 {
        start = 0
    }
    if start > rl {
        start = rl
    }
    if end < 0 {
        end = 0
    }
    if end > rl {
        end = rl
    }

    return string(rs[start:end])
}