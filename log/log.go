package log

/*
Author: chenwenjiang
Email:chenwenjiang@oneniceapp.com
Create date: 2015-04-29
*/

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

type LogConf struct {
	Module       string `json:"module"`
	Level        string `json:"level"`
	LogDir       string `json:"log_dir"`
	SuffixFormat string `json:"suffix_format"`
	Color        bool   `json:"color,omitempty"`
	Stdout       bool   `json:"stdout,omitempty"`
	MaxLogCount  int    `json:"max_log_count"`
}

const (
	debugLevel = iota
	infoLevel
	noticeLevel
	warningLevel
	errorLevel
)

const (
	nocolor = 0
	red     = 30 + iota
	green
	yellow
	blue
	purple
	cyan
)

var (
	logPrefix = map[int]string{
		debugLevel:   "DEBUG",
		infoLevel:    "INFO",
		noticeLevel:  "NOTICE",
		warningLevel: "WARNING",
		errorLevel:   "ERROR",
	}
	logColor = map[int]int{
		debugLevel:   cyan,
		infoLevel:    nocolor,
		noticeLevel:  green,
		warningLevel: yellow,
		errorLevel:   red,
	}

	LogStr2Int = map[string]int{
		"DEBUG":   debugLevel,
		"INFO":    infoLevel,
		"NOTICE":  noticeLevel,
		"WARNING": warningLevel,
		"ERROR":   errorLevel,
	}
	Levels = map[int]bool{}
	mtime  string

	logger *Logger

	hostName = ""
	timeNow  time.Time
)

type FileStat struct {
	name   string
	fmtime int64
}
type FSTS []FileStat

func (fs FSTS) Len() int {
	return len(fs)
}
func (fs FSTS) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}
func (fs FSTS) Less(i, j int) bool {
	return fs[i].fmtime > fs[j].fmtime
}

type Logger struct {
	sync.Mutex
	isConsole      bool
	isColorfull    bool
	reserveCounter int
	timeFormat     string
	fileName       string
	fileWriter     io.WriteCloser
}

// 滚动切割文件
func rollingLogFile(toFileName string, logger *Logger) {
	logger.Lock()
	defer logger.Unlock()

	logger.fileWriter.Close()
	logger.fileWriter = nil
	err := os.Rename(logger.fileName, toFileName)
	if err != nil {
		panic(err)
	}
	fileWriter, err := os.OpenFile(logger.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logger.fileWriter = fileWriter
}

// log split checker
func logSplitChecker(logger *Logger) {
	ticker := time.NewTicker(1 * time.Second) // 文件切割
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if mtime == "" {
				mtime = getCurrentTime().Format(logger.timeFormat)
				continue
			}
			currentTime := getCurrentTime().Format(logger.timeFormat)
			if currentTime != mtime {
				toFileName := fmt.Sprintf("%s.%s", logger.fileName, mtime)
				rollingLogFile(toFileName, logger)
				mtime = currentTime
			}
		}
	}
}

func logCounterChecker(logger *Logger) {
	ticker := time.NewTicker(24 * time.Hour) // 一天检查一次日志个数
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			dirname := filepath.Dir(logger.fileName)
			basename := filepath.Base(logger.fileName)
			logLists := make(FSTS, 0)
			filepath.Walk(fmt.Sprintf("%s/", dirname), func(path string, f os.FileInfo, err error) error {
				if strings.HasPrefix(f.Name(), fmt.Sprintf("%s.", basename)) {
					fs := FileStat{
						name:   f.Name(),
						fmtime: f.ModTime().Unix(),
					}
					logLists = append(logLists, fs)
				}
				return nil
			})
			sort.Sort(logLists)

			if len(logLists) > logger.reserveCounter {
				removes := logLists[logger.reserveCounter:]
				for _, fname := range removes {
					rmname := fmt.Sprintf("%s/%s", dirname, fname.name)
					os.Remove(rmname)
				}
			}
		}
	}
}

// 设置日志级别
func setLevels(levels []string) {
	Levels = make(map[int]bool)

	if len(levels) > 0 {
		flag := false
		for _, v := range levels {
			if v == "NOSET" {
				// 表示所有的log都大于
				flag = true
			}
		}
		if !flag {
			for _, v := range levels {
				if sv, ok := LogStr2Int[v]; ok {
					Levels[sv] = true
				}
			}
		}
	}
}

func NewLogger(config LogConf) *Logger {
	logLevel := config.Level
	logDir := config.LogDir
	logFile := fmt.Sprintf("%s.log", config.Module)
	logReserveCounter := config.MaxLogCount
	logSuffix := config.SuffixFormat
	var logConsole int
	if config.Stdout {
		logConsole = 1
	}
	var logColorfull int
	if config.Color {
		logColorfull = 1
	}
	return newLoggerWithArg(logLevel, logDir, logFile, logReserveCounter, logSuffix, logConsole, logColorfull)
}

func InitLogger(config LogConf) {
	file := fmt.Sprintf("%s.log", config.Module)
	var stdout int
	if config.Stdout {
		stdout = 1
	}
	var color int
	if config.Color {
		color = 1
	}
	logger = newLoggerWithArg(config.Level, config.LogDir, file, config.MaxLogCount, config.SuffixFormat, stdout, color)
}

func initLooger(level string, dir string, file string, reserve int, suffix string, console int, color int) {
	logger = newLoggerWithArg(level, dir, file, reserve, suffix, console, color)
}

// func InitLoggerWithConfigFile(loggerFile string) {
// 	logger = NewLogger(loggerFile)
// }

func newLoggerWithArg(level string, dir string, file string, reserve int, suffix string, console int, color int) *Logger {
	if _, err := os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, 0777); err != nil {
			panic(err)
		}
	}
	logFilePath := fmt.Sprintf("%s/%s", dir, file)
	fileWriter := getFileWriter(logFilePath)
	setLevels(strings.Split(level, ","))

	boolConsole := false
	if console == 1 {
		boolConsole = true
	}
	boolColorfull := false
	if color == 1 {
		boolColorfull = true
	}

	loggerHandle := &Logger{
		isConsole:      boolConsole,
		isColorfull:    boolColorfull,
		reserveCounter: reserve,
		timeFormat:     suffix,
		fileName:       logFilePath,
		fileWriter:     fileWriter,
	}
	timeNow = time.Now()
	go tick()
	go logSplitChecker(loggerHandle)
	go logCounterChecker(loggerHandle)
	return loggerHandle
}

// 创建默认的logger log的配置在pwd
// func NewDefaultLogger() *Logger {
// 	file, _ := os.Getwd()
// 	logFile := fmt.Sprintf("%s/%s", file, "golog.cfg")
// 	return NewLogger(logFile)
// }

// 新建console的日志handle
func NewConsoleLogger() *Logger {
	return nil
}

// 创建文件
func getFileWriter(fileName string) io.WriteCloser {
	fileWriter, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return fileWriter
}

func getDetail() (string, int) {
	// 获取调用函数的数据
	_, file, line, _ := runtime.Caller(3)
	filename := path.Base(file)
	return filename, line
}

func isMasterNetCard(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 && ip4[1] == 10
	}
	return false

}
func getHostIpAddr() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if isMasterNetCard(ipnet.IP) {
				return addr.String()
			}
		}
	}
	return ""
}

func getHostName() string {
	if hostName == "" {
		hostname, _ := os.Hostname()
		hostname = strings.Split(hostname, ".")[0]
		return hostname
	}
	return hostName
}

func (l *Logger) write(level int, format string, content ...interface{}) {
	if _, ok := Levels[level]; !ok {
		return
	}
	filename, line := getDetail() // 获取文件的信息
	hostname := getHostName()
	now := getCurrentTime()
	var s string
	if format == "" {
		s = renderColor(fmt.Sprintf("[%s] [%s] [%s:%d] [%s] %s\n", now.Format("2006/01/02 15:04:05"), logPrefix[level], filename, line, hostname, fmt.Sprint(content...)), logColor[level], l.isColorfull)
	} else {
		s = renderColor(fmt.Sprintf("[%s] [%s] [%s:%d] [%s] %s\n", now.Format("2006/01/02 15:04:05"), logPrefix[level], filename, line, hostname, fmt.Sprintf(format, content...)), logColor[level], l.isColorfull)
	}

	l.Lock()
	defer l.Unlock()
	l.fileWriter.Write([]byte(s))
	if l.isConsole {
		fmt.Print(s)
	}
}

func (l *Logger) Info(content ...interface{}) {
	l.write(infoLevel, "", content...)
}

func (l *Logger) Infof(format string, content ...interface{}) {
	l.write(infoLevel, format, content...)
}

func (l *Logger) Warning(content ...interface{}) {
	l.write(warningLevel, "", content...)
}

func (l *Logger) Warningf(format string, content ...interface{}) {
	l.write(warningLevel, format, content...)
}

func (l *Logger) Notice(content ...interface{}) {
	l.write(noticeLevel, "", content...)
}

func (l *Logger) Noticef(format string, content ...interface{}) {
	l.write(noticeLevel, format, content...)
}

func (l *Logger) Debug(content ...interface{}) {
	l.write(debugLevel, "", content...)
}

func (l *Logger) Debugf(format string, content ...interface{}) {
	l.write(debugLevel, format, content...)
}

func (l *Logger) Error(content ...interface{}) {
	l.write(errorLevel, "", content...)
}

func (l *Logger) Errorf(format string, content ...interface{}) {
	l.write(errorLevel, format, content...)
}

func renderColor(s string, color int, isColorfull bool) string {
	if isColorfull {
		return fmt.Sprintf("\033[%dm%s\033[0m", color, s)
	} else {
		return s
	}
}

func Info(content ...interface{}) {
	logger.write(infoLevel, "", content...)
}

func Infof(format string, content ...interface{}) {
	logger.write(infoLevel, format, content...)
}

func Warning(content ...interface{}) {
	logger.write(warningLevel, "", content...)
}

func Warningf(format string, content ...interface{}) {
	logger.write(warningLevel, format, content...)
}

func Notice(content ...interface{}) {
	logger.write(noticeLevel, "", content...)
}

func Noticef(format string, content ...interface{}) {
	logger.write(noticeLevel, format, content...)
}

func Debug(content ...interface{}) {
	logger.write(debugLevel, "", content...)
}

func Debugf(format string, content ...interface{}) {
	logger.write(debugLevel, format, content...)
}

func Error(content ...interface{}) {
	logger.write(errorLevel, "", content...)
}

func Errorf(format string, content ...interface{}) {
	logger.write(errorLevel, format, content...)
}

func getCurrentTime() time.Time {
	return timeNow
}

func tick() {
	for {
		select {
		case <-time.Tick(time.Second):
			timeNow = time.Now()
		}
	}
}
