package sdklog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/openimsdk/tools/errs"
	rotatelogs "github.com/openimsdk/tools/log/file-rotatelogs"
	"github.com/openimsdk/tools/utils/stringutil"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Foreground colors.
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

var (
	_levelToColor = map[zapcore.Level]Color{
		zapcore.DebugLevel:  White,
		zapcore.InfoLevel:   Blue,
		zapcore.WarnLevel:   Yellow,
		zapcore.ErrorLevel:  Red,
		zapcore.DPanicLevel: Red,
		zapcore.PanicLevel:  Red,
		zapcore.FatalLevel:  Red,
	}
	_unknownLevelColor = make(map[zapcore.Level]string, len(_levelToColor))

	_levelToLowercaseColorString = make(map[zapcore.Level]string, len(_levelToColor))
	_levelToCapitalColorString   = make(map[zapcore.Level]string, len(_levelToColor))
)

func init() {
	for level, color := range _levelToColor {
		_levelToLowercaseColorString[level] = color.Add(level.String())
		_levelToCapitalColorString[level] = color.Add(level.CapitalString())
	}
}

// Color represents a text color.
type Color uint8

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

//-----------------------//

const sdkCallDepth = 2

var (
	pkgLogger   Logger
	sp          = string(filepath.Separator)
	logLevelMap = map[int]zapcore.Level{
		6: zapcore.DebugLevel,
		5: zapcore.DebugLevel,
		4: zapcore.InfoLevel,
		3: zapcore.WarnLevel,
		2: zapcore.ErrorLevel,
		1: zapcore.FatalLevel,
		0: zapcore.PanicLevel,
	}
)

const hoursPerDay = 24

type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...any)
	Info(ctx context.Context, msg string, keysAndValues ...any)
	Warn(ctx context.Context, msg string, err error, keysAndValues ...any)
	Error(ctx context.Context, msg string, err error, keysAndValues ...any)
	WithValues(keysAndValues ...any) Logger
	WithName(name string) Logger
	WithCallDepth(depth int) Logger
	ToZap() *zap.SugaredLogger
}

type SDKLogger struct {
	zapLogger        *zap.SugaredLogger
	level            zapcore.Level
	sdkType          string
	platformID       string
	moduleName       string
	moduleVersion    string
	loggerPrefix     string
	loggerPrefixName string
	rotationTime     time.Duration
	isSimplify       bool
}

func InitSDKLogger(
	loggerPrefixName, moduleName string, sdkType, platformID string,
	logLevel int,
	isStdout bool,
	isJson bool,
	logLocation string,
	rotateCount uint,
	rotationTime uint,
	moduleVersion string,
	isSimplify bool,
) error {
	l, err := NewSDKLogger(loggerPrefixName, moduleName, sdkType, platformID, logLevel, isStdout, isJson, logLocation, rotateCount, rotationTime, moduleVersion, isSimplify)
	if err != nil {
		return err
	}

	pkgLogger = l.WithCallDepth(sdkCallDepth)
	return nil
}

func NewSDKLogger(
	loggerPrefixName, moduleName string, sdkType, platformID string,
	logLevel int,
	isStdout bool,
	isJson bool,
	logLocation string,
	rotateCount uint,
	rotationTime uint,
	moduleVersion string,
	isSimplify bool,
) (*SDKLogger, error) {
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(logLevelMap[logLevel]),
		DisableStacktrace: true,
	}
	if isJson {
		zapConfig.Encoding = "json"
	} else {
		zapConfig.Encoding = "console"
	}
	zl := &SDKLogger{
		level:         logLevelMap[logLevel],
		sdkType:       sdkType,
		platformID:    platformID,
		moduleName:    moduleName,
		loggerPrefix:  loggerPrefixName,
		rotationTime:  time.Duration(rotationTime) * time.Hour,
		moduleVersion: moduleVersion,
		isSimplify:    isSimplify,
	}
	opts, err := zl.cores(isStdout, isJson, logLocation, rotateCount)
	if err != nil {
		return nil, err
	}
	l, err := zapConfig.Build(opts)
	if err != nil {
		return nil, err
	}
	zl.zapLogger = l.Sugar()
	return zl, nil
}

func (l *SDKLogger) cores(isStdout bool, isJson bool, logLocation string, rotateCount uint) (zap.Option, error) {
	c := zap.NewProductionEncoderConfig()
	c.EncodeTime = l.timeEncoder
	c.EncodeLevel = l.capitalColorLevelEncoder
	c.EncodeDuration = zapcore.SecondsDurationEncoder
	c.MessageKey = "msg"
	c.LevelKey = "level"
	c.TimeKey = "time"
	c.CallerKey = "caller"
	c.NameKey = "logger"

	var fileEncoder zapcore.Encoder
	if isJson {
		c.EncodeLevel = zapcore.CapitalLevelEncoder
		fileEncoder = zapcore.NewJSONEncoder(c)
		fileEncoder.AddInt("PID", os.Getpid())
		fileEncoder.AddString("sdkType", l.sdkType)
		fileEncoder.AddString("platformID", l.platformID)
		fileEncoder.AddString("version", l.moduleVersion)
	} else {
		c.EncodeLevel = zapcore.CapitalLevelEncoder
		c.EncodeCaller = l.customCallerEncoder
		fileEncoder = zapcore.NewConsoleEncoder(c)
	}

	writer, err := l.getWriter(logLocation, rotateCount)
	if err != nil {
		return nil, err
	}

	var cores []zapcore.Core
	if logLocation != "" {
		cores = []zapcore.Core{
			zapcore.NewCore(fileEncoder, writer, zap.NewAtomicLevelAt(l.level)),
		}
	}
	if isStdout {
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.Lock(os.Stdout), zap.NewAtomicLevelAt(l.level)))
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	}), nil
}

func (l *SDKLogger) getWriter(logLocation string, rotateCount uint) (zapcore.WriteSyncer, error) {
	var path string
	if l.rotationTime%(time.Hour*time.Duration(hoursPerDay)) == 0 {
		path = logLocation + sp + l.loggerPrefixName + ".%Y-%m-%d"
	} else if l.rotationTime%time.Hour == 0 {
		path = logLocation + sp + l.loggerPrefixName + ".%Y-%m-%d_%H"
	} else {
		path = logLocation + sp + l.loggerPrefixName + ".%Y-%m-%d_%H_%M_%S"
	}
	logf, err := rotatelogs.New(path,
		rotatelogs.WithRotationCount(rotateCount),
		rotatelogs.WithRotationTime(l.rotationTime),
	)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(logf), nil
}

func (l *SDKLogger) timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := "2006-01-02 15:04:05.000"
	enc.AppendString(t.Format(layout))
}

func (l *SDKLogger) capitalColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s, ok := _levelToCapitalColorString[level]
	if !ok {
		s = _unknownLevelColor[zapcore.ErrorLevel]
	}
	pid := stringutil.FormatString(fmt.Sprintf("[PID:%d]", os.Getpid()), 15, true)
	color := _levelToColor[level]
	enc.AppendString(s)
	enc.AppendString(color.Add(pid))
	if l.moduleName != "" {
		moduleName := stringutil.FormatString(l.moduleName, 25, true)
		enc.AppendString(color.Add(moduleName))
	}
	if l.moduleVersion != "" {
		moduleVersion := stringutil.FormatString(fmt.Sprintf("[version:%s]", l.moduleVersion), 17, true)
		enc.AppendString(moduleVersion)
	}
}

func (l *SDKLogger) customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	fixedLength := 50
	trimmedPath := caller.TrimmedPath()
	trimmedPath = "[" + trimmedPath + "]"
	s := stringutil.FormatString(trimmedPath, fixedLength, true)
	enc.AppendString(s)
}

func (l *SDKLogger) ToZap() *zap.SugaredLogger {
	return l.zapLogger
}

func SDKLog(ctx context.Context, logLevel int, path string, line string, msg, err string, keysAndValues map[string]string) {
	var kvSlice []any
	for k, v := range keysAndValues {
		kvSlice = append(kvSlice, k, v)
	}
	switch logLevel {
	case 6:
		SDKDebug(ctx, path, line, msg, kvSlice...)
	case 4:
		SDKInfo(ctx, path, line, msg, kvSlice...)
	case 3:
		SDKWarn(ctx, path, line, msg, errs.New(err), kvSlice...)
	case 2:
		SDKError(ctx, path, line, msg, errs.New(err), kvSlice...)
	}
}

func (l *SDKLogger) Debug(ctx context.Context, msg string, keysAndValues ...any) {
	keysAndValues = l.kvAppend(ctx, keysAndValues)
	l.zapLogger.Debugw(msg, keysAndValues...)
}

func (l *SDKLogger) Info(ctx context.Context, msg string, keysAndValues ...any) {
	keysAndValues = l.kvAppend(ctx, keysAndValues)
	l.zapLogger.Infow(msg, keysAndValues...)
}

func (l *SDKLogger) Warn(ctx context.Context, msg string, err error, keysAndValues ...any) {
	if err != nil {
		keysAndValues = append(keysAndValues, "error", err.Error())
	}
	keysAndValues = l.kvAppend(ctx, keysAndValues)
	l.zapLogger.Warnw(msg, keysAndValues...)
}

func (l *SDKLogger) Error(ctx context.Context, msg string, err error, keysAndValues ...any) {
	if err != nil {
		keysAndValues = append(keysAndValues, "error", err.Error())
	}
	keysAndValues = l.kvAppend(ctx, keysAndValues)
	l.zapLogger.Errorw(msg, keysAndValues...)
}

func (l *SDKLogger) WithValues(keysAndValues ...any) Logger {
	dup := *l
	dup.zapLogger = l.zapLogger.With(keysAndValues...)
	return &dup
}

func (l *SDKLogger) WithName(name string) Logger {
	dup := *l
	dup.zapLogger = l.zapLogger.Named(name)
	return &dup
}

func (l *SDKLogger) WithCallDepth(depth int) Logger {
	dup := *l
	dup.zapLogger = l.zapLogger.WithOptions(zap.AddCallerSkip(depth))
	return &dup
}

func (l *SDKLogger) kvAppend(ctx context.Context, keysAndValues []any) []any {
	if ctx == nil {
		return keysAndValues
	}
	keysAndValues = append([]any{"sdkType", l.sdkType, "platformID", l.platformID}, keysAndValues...)
	keysAndValues = append(keysAndValues, "moduleName", l.moduleName, "version", l.moduleVersion)
	return keysAndValues
}

func SDKDebug(ctx context.Context, path string, line string, msg string, keysAndValues ...any) {
	if pkgLogger == nil {
		return
	}
	pkgLogger.Debug(ctx, msg, keysAndValues...)
}

func SDKInfo(ctx context.Context, path string, line string, msg string, keysAndValues ...any) {
	if pkgLogger == nil {
		return
	}
	pkgLogger.Info(ctx, msg, keysAndValues...)
}

func SDKWarn(ctx context.Context, path string, line string, msg string, err error, keysAndValues ...any) {
	if pkgLogger == nil {
		return
	}
	pkgLogger.Warn(ctx, msg, err, keysAndValues...)
}

func SDKError(ctx context.Context, path string, line string, msg string, err error, keysAndValues ...any) {
	if pkgLogger == nil {
		return
	}
	pkgLogger.Error(ctx, msg, err, keysAndValues...)
}

// func CInfo(ctx context.Context, msg string, keysAndValues ...any) {
// 	if osStdout == nil {
// 		return
// 	}
// 	osStdout.Info(ctx, msg, keysAndValues...)
// }
