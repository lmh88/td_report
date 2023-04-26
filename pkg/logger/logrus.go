package logger

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"td_report/vars"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Fields logrus.Fields

var AppName string = "app"
var NameRand bool

// Define logrus alias
var (
	Logger            *LoggerV2
	Log               = Logger
	Tracef            func(format string, args ...interface{})
	Debugf            func(format string, args ...interface{})
	Infof             func(format string, args ...interface{})
	Info              func(args ...interface{})
	Warnf             func(format string, args ...interface{})
	Errorf            func(format string, args ...interface{})
	Write             func(args ...interface{})
	Fatalf            func(format string, args ...interface{})
	Panicf            func(format string, args ...interface{})
	Printf            func(format string, args ...interface{})
	Println           func(args ...interface{})
	Error             func(args ...interface{})
	SetOutput         func(output io.Writer)
	WithFields        func(fields Fields) *Entry
	WithField         func(key string, value interface{}) *Entry
	SetReportCaller   func(reportCaller bool)
	StandardLogger    func() *logrus.Logger
	ParseLevel        = logrus.ParseLevel
	WithContext       func(ctx context.Context) *Entry
	NewTraceIDContext func(ctx context.Context, traceID string) context.Context
	debug             bool
)

// Define logger level
const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func init() {
	rand.Seed(time.Now().Unix())
	InitLog(AppName)
}

// Init 测试备选
func Init(appName string, nameRand bool) {
	NameRand = nameRand
	if appName == "" {
		InitLog(AppName)
	} else {
		InitLog(appName)
	}

}

func InitLog(appName string) {
	Logger = NewLog(appName)
	Tracef = Logger.Tracef
	Debugf = Logger.Debugf
	Infof = Logger.Infof
	Info = Logger.Info
	Warnf = Logger.Warnf
	Errorf = Logger.Errorf
	Write = Logger.Error
	Fatalf = Logger.Fatalf
	Panicf = Logger.Panicf
	Printf = Logger.Printf
	Println = Logger.Println
	Error = Logger.Error
	SetOutput = Logger.SetOutput
	WithFields = Logger.WithFields
	WithField = Logger.WithField
	SetReportCaller = Logger.SetReportCaller
	StandardLogger = logrus.StandardLogger
	ParseLevel = logrus.ParseLevel
	WithContext = Logger.WithContext
	NewTraceIDContext = Logger.NewTraceIDContext
}

// LoggerV2 Log Logs
type LoggerV2 struct {
	*logrus.Logger
}

// Entry logrus.Entry alias
type Entry = logrus.Entry

//
// Hook logrus.Hook alias
type Hook = logrus.Hook

type Level = logrus.Level

func NewLog(appName string) *LoggerV2 {
	loggerV2 := &LoggerV2{
		logrus.New(),
	}
	// 默认是json格式和终端输出
	format := "json"
	switch format {
	case "text":
		loggerV2.Formatter = &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.000", FullTimestamp: true}
	case "json":
		loggerV2.Formatter = &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"}
	}
	// Logger.Out       = os.Stdout
	logPath := g.Cfg().GetString("logger.Path")
	debug = false //config.GetBoolOrDefault("log.debug", false)
	// 显示行号
	//loggerV2.SetReportCaller(config.GetBoolOrDefault("log.caller", false))

	loggerV2.SetFileOutWriter(logPath, appName, 7*24*time.Hour, 3*time.Hour)
	return loggerV2
}

// SetOutWriter  可以设置文件记录日志
func (l *LoggerV2) SetOutWriter(writer io.Writer) {
	l.Out = writer
}

func (l *LoggerV2) getFormat() logrus.Formatter {
	return l.Logger.Formatter
}

// IsExist util存在对应方法，避免循环引用
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}

	return true
}

// SetFileOutWriter 设置带文件并且切割的日志模式
func (l *LoggerV2) SetFileOutWriter(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	// pro模式不开控制台日志
	if !debug {
		l.SetOutput(ioutil.Discard)
	}

	if !IsExist(logPath) {
		_ = os.Mkdir(logPath, 0755)
	}

	if NameRand {
		logFileName = fmt.Sprintf("%s-pro-%s-%d", logFileName, time.Now().Format(vars.TimeLayout), rand.Intn(1000))
	} else {
		logFileName = fmt.Sprintf("%s-pro-%s", logFileName, time.Now().Format(vars.TimeLayout))
	}

	path := fmt.Sprintf("%s/%s.log", logPath, logFileName)
	writer, _ := rotatelogs.New(
		path+".%Y%m%d.log",
		//path+".%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)

	l.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
		},
		l.Formatter,
	))
}

func (l *LoggerV2) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *LoggerV2) InfoWithContext(context context.Context, args ...interface{}) {
	l.WithContext(context).Info(args...)
}

func (l *LoggerV2) InfoWithContextMap(context context.Context, val map[string]interface{}, args ...interface{}) {
	l.WithContext(context).WithFields(val).Info(args...)
}

func (l *LoggerV2) ErrorWithContext(context context.Context, args ...interface{}) {
	l.WithContext(context).Error(args...)
}

func (l *LoggerV2) ErrorWithContextMap(context context.Context, val map[string]interface{}, args ...interface{}) {
	l.WithContext(context).WithFields(val).Error(args...)
}

func (l *LoggerV2) Tracef(format string, args ...interface{}) {
	l.Logger.Tracef(format, args...)
}

func (l *LoggerV2) Debugf(format string, args ...interface{}) {
	l.Logger.Tracef(format, args...)
}

func (l *LoggerV2) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *LoggerV2) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

func (l *LoggerV2) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}
func (l *LoggerV2) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}
func (l *LoggerV2) Panicf(format string, args ...interface{}) {
	l.Logger.Panicf(format, args...)
}
func (l *LoggerV2) Printf(format string, args ...interface{}) {
	l.Logger.Printf(format, args...)
}

func (l *LoggerV2) WithFields(fields Fields) *Entry {
	return l.Logger.WithFields(logrus.Fields(fields))
}

func (l *LoggerV2) WithField(key string, value interface{}) *Entry {
	return l.Logger.WithField(key, value)
}

func (l *LoggerV2) SetOutput(output io.Writer) {
	l.Logger.SetOutput(output)
}

func (l *LoggerV2) SetReportCaller(reportCaller bool) {
	l.Logger.SetReportCaller(reportCaller)
}

// Error 为了兼容以前老的代码
func (l *LoggerV2) Write(flag string, err error, infos ...interface{}) {
	var data = make(Fields)
	data["flag"] = flag
	data["err"] = err
	data["info"] = infos
	if Logger.Level == logrus.FatalLevel {
		l.WithFields(data).Error(infos...)
	} else {
		l.WithFields(data).Info(infos...)
	}
}

func (l *LoggerV2) WriteMap(val map[string]interface{}, infos ...interface{}) {
	if Logger.Level == logrus.FatalLevel {
		l.WithFields(val).Error(infos...)
	} else {
		l.WithFields(val).Info(infos...)
	}
}

// SetLevel Set logger level
func (l *LoggerV2) SetLevel(level Level) {
	l.Level = level
}

// SetFormatter Set logger output format (json/text)
func (l *LoggerV2) SetFormatter(format string) {
	switch format {
	case "json":
		l.Formatter = new(logrus.JSONFormatter)
	case "text":
		l.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	default:
		l.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	}
}

// AddLogHook AddHook Add logger hook
func (l *LoggerV2) AddLogHook(hook Hook) {
	l.Logger.AddHook(hook)
}

// Define key
const (
	TraceIDKey = "traceId"
	TagKey     = "tag"
	StackKey   = "stack"
)

type (
	traceIDKey  struct{}
	userIDKey   struct{}
	userNameKey struct{}
	tagKey      struct{}
	stackKey    struct{}
)

func (*LoggerV2) NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func (*LoggerV2) FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (*LoggerV2) NewUserIDContext(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func (*LoggerV2) FromUserIDContext(ctx context.Context) uint64 {
	v := ctx.Value(userIDKey{})
	if v != nil {
		if s, ok := v.(uint64); ok {
			return s
		}
	}
	return 0
}

func (*LoggerV2) NewUserNameContext(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, userNameKey{}, userName)
}

func (*LoggerV2) FromUserNameContext(ctx context.Context) string {
	v := ctx.Value(userNameKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (*LoggerV2) NewTagContext(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, tagKey{}, tag)
}

func (*LoggerV2) FromTagContext(ctx context.Context) string {
	v := ctx.Value(tagKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

func (*LoggerV2) NewStackContext(ctx context.Context, stack error) context.Context {
	return context.WithValue(ctx, stackKey{}, stack)
}

func (l *LoggerV2) FromStackContext(ctx context.Context) error {
	v := ctx.Value(stackKey{})
	if v != nil {
		if s, ok := v.(error); ok {
			return s
		}
	}

	return nil
}

// WithContext Use context create entry
func (l *LoggerV2) WithContext(ctx context.Context) *Entry {
	fields := logrus.Fields{}

	if v := l.FromTraceIDContext(ctx); v != "" {
		fields[TraceIDKey] = v
	}

	return l.Logger.WithContext(ctx).WithFields(fields)
}
