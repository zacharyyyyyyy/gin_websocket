package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
	"time"
)

type Logger struct {
	*zap.Logger
}

var (
	loggers = make(map[string]*Logger)
	mu      sync.RWMutex
)

var (
	Api     = Get("api")
	Service = Get("service")
	Runtime = Get("runtime")
)

func Get(filename string) *Logger {
	mu.RLock()
	if l, ok := loggers[filename]; ok {
		mu.RUnlock()
		return l
	}
	mu.RUnlock()

	// 实例锁
	mu.Lock()
	defer mu.Unlock()
	if l, ok := loggers[filename]; ok {
		return l
	}
	encoder := getEncoder()
	writeSyncer := getLogWriter(fmt.Sprintf("%s/%s/", "Log", filename))
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	loggers[filename] = &Logger{zap.New(core, zap.AddCaller())}
	return loggers[filename]
}

func getLogWriter(path string) zapcore.WriteSyncer {
	fileName := path + "log"
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10, // 10M切割
		MaxAge:     5,  // 保留旧文件个数
		MaxBackups: 30, // 旧文件存活天数
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
func getEncoder() zapcore.Encoder {
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = logTimeFormat
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
func logTimeFormat(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 15:04:05]"))
}

func (l *Logger) Debug(msg string) {
	l.Logger.Debug(msg)
}

func (l *Logger) Info(msg string) {
	l.Logger.Info(msg)
}
func (l *Logger) Warn(msg string) {
	l.Logger.Warn(msg)
}
func (l *Logger) Error(msg string) {
	l.Logger.Error(msg)
}
