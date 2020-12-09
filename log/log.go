package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 通用DLTag常量定义
const (
	DLTagUndefind      = "_undef"
	DLTagMySQLFailed   = "_com_mysql_failure"
	DLTagRedisFailed   = "_com_redis_failure"
	DLTagMySQLSuccess  = "_com_mysql_success"
	DLTagRedisSuccess  = "_com_redis_success"
	DLTagThriftFailed  = "_com_thrift_failure"
	DLTagThriftSuccess = "_com_thrift_success"
	DLTagHTTPSuccess   = "_com_http_success"
	DLTagHTTPFailed    = "_com_http_failure"
	DLTagTCPFailed     = "_com_tcp_failure"
	DLTagRequestIn     = "_com_request_in"
	DLTagRequestOut    = "_com_request_out"
)

const (
	_dlTag          = "dltag"
	_traceID        = "traceid"
	_spanID         = "spanid"
	_childSpanID    = "cspanid"
	_dlTagBizPrefix = "_com_"
	_dlTagBizUndef  = "_com_undef"
)

// Trace 链路日志结构体
type Trace struct {
	TraceID     string
	SpandID     string
	Caller      string
	SrcMethod   string
	HintContent string
	HintCode    int64
}

// TraceContext 链路日志上下文信息
type TraceContext struct {
	Trace
	CSpanID string
}

// Logger 全局变量
type Logger struct {
	L *zap.Logger
}

// Log 全局日志实例
var Log *Logger

// GetEncoder 编码器
func GetEncoder() zapcore.Encoder {

	// 配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger", // 名字是什么
		CallerKey:     "caller", // 调用者的名字
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		//EncodeTime: zapcore.EpochTimeEncoder, // 默认的时间编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder, // 修改之后的时间编码器
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeCaller: myCollerEncoder,
	}

	// p配置 JSON 编码器
	// return zapcore.NewJSONEncoder(encoderConfig)

	// 配置 Console 编码器
	return zapcore.NewConsoleEncoder(encoderConfig)

}

// 自定义调用函数解析器
func myCollerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// TODO: consider using a byte-oriented API to save an allocation.
	pc, file, line, ok := runtime.Caller(6)
	// caller := "undefined"
	// if ok {
	// 	// code = path.Base(file) + ":" + strconv.Itoa(line)
	// 	caller = fmt.Sprintf("%s:%d", path.Base(file), line)
	// }
	caller.PC = pc
	caller.File = file
	caller.Line = line
	caller.Defined = ok
	// mycaller := EntryCaller{
	// 	Defined ok,
	// 	PC      pc
	// 	File    file
	// 	Line    line,
	// }
	enc.AppendString(caller.TrimmedPath())
}

// GetLogWriter 分割归档
func GetLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// Debug debug 级别日志
func (l *Logger) Debug(msg string, trace *TraceContext, dltag string, m map[string]interface{}) {
	// m[_dlTag] = CheckDLTag(dltag)
	// m[_traceID] = trace.TraceID
	// m[_childSpanID] = trace.CSpanID
	// m[_spanID] = trace.SpandID
	// l.L.Debug(parseParams(m))
	l.L.Debug(msg, zap.String(_traceID, trace.TraceID),
		zap.String(_spanID, trace.SpandID),
		zap.String(_childSpanID, trace.CSpanID),
		zap.String("msg", parseParams(m)))
}

// Info Info级别日志
func (l *Logger) Info(msg string, trace *TraceContext, dltag string, m map[string]interface{}) {
	// dltag = CheckDLTag(dltag)
	// m[_traceID] = trace.TraceID
	// m[_childSpanID] = trace.CSpanID
	// m[_spanID] = trace.SpandID
	// traceMsg := traceFormat(trace, dltag)
	// l.L.Info(msg, zap.String(traceMsg, fmt.Sprintf("%v", m["message"])))

	// eg: 2020-12-09T14:45:54.100+0800	info	log/log.go:107	/test	{"traceid=000000005fd072a2a1b0eba1658221|spanid=9f786d5978629a0f|cspanid=": "_undef|error=text string|balabala=xxxx|message=todo sth"}
	// l.L.Info(msg, zap.String(traceMsg, parseParams(m)))

	// eg: 2020-12-09T14:52:44.440+0800	info	log/log.go:111	/test	{"traceid": "000000005fd0743c70d618e9658221", "spanid": "9f786bc778629a0f", "cspanid": "", "msg": "_undef|message=todo sth|error=text string|balabala=xxxx"}
	// source code, file and line num
	// pc, file, line, ok := runtime.Caller(1)
	// caller := "undefined"
	// if ok {
	// 	// code = path.Base(file) + ":" + strconv.Itoa(line)
	// 	caller = fmt.Sprintf("%s:%d", path.Base(file), line)
	// }

	// l.L.Core().Write(zapcore.Entry{
	// 	Level:  zapcore.InfoLevel,
	// 	Caller: zapcore.NewEntryCaller(pc, file, line, ok)}, []zapcore.Field{zap.String("callertest", "测试caller")})

	// m["caller"] = caller

	l.L.Info(msg, zap.String(_traceID, trace.TraceID),
		zap.String(_spanID, trace.SpandID),
		zap.String(_childSpanID, trace.CSpanID),
		zap.String("msg", parseParams(m)))
}

// Warn warn 级别日志
func (l *Logger) Warn(msg string, trace *TraceContext, dltag string, m map[string]interface{}) {
	l.L.Warn(msg, zap.String(_traceID, trace.TraceID),
		zap.String(_spanID, trace.SpandID),
		zap.String(_childSpanID, trace.CSpanID),
		zap.String("msg", parseParams(m)))
}

// Error error 级别日志
func (l *Logger) Error(msg string, trace *TraceContext, dltag string, m map[string]interface{}) {
	// m[_dlTag] = CheckDLTag(dltag)
	// m[_traceID] = trace.TraceID
	// m[_childSpanID] = trace.CSpanID
	// m[_spanID] = trace.SpandID
	// l.L.Error(parseParams(m))
	// pc, file, line, ok := runtime.Caller(1)
	// caller := "undefined"
	// if ok {
	// 	// code = path.Base(file) + ":" + strconv.Itoa(line)
	// 	caller = fmt.Sprintf("%s:%d", path.Base(file), line)
	// }
	// l.L.Core().Write(zapcore.Entry{
	// 	// Time: zapcore.Con
	// 	Level:  zapcore.ErrorLevel,
	// 	Caller: zapcore.NewEntryCaller(pc, file, line, ok)}, []zapcore.Field{zap.String("callertest", "测试caller")})

	// m["caller"] = caller

	l.L.Error(msg, zap.String(_traceID, trace.TraceID),
		zap.String(_spanID, trace.SpandID),
		zap.String(_childSpanID, trace.CSpanID),
		zap.String("msg", parseParams(m)))
}

// DPanic DPanic 级别日志
func (l *Logger) DPanic(msg string, trace *TraceContext, dltag string, m map[string]interface{}) {
	l.L.DPanic(msg, zap.String(_traceID, trace.TraceID),
		zap.String(_spanID, trace.SpandID),
		zap.String(_childSpanID, trace.CSpanID),
		zap.String("msg", parseParams(m)))
}

// Panic Panic 级别日志
func (l *Logger) Panic(msg string, trace *TraceContext, dltag string, m map[string]interface{}) {
	l.L.Panic(msg, zap.String(_traceID, trace.TraceID),
		zap.String(_spanID, trace.SpandID),
		zap.String(_childSpanID, trace.CSpanID),
		zap.String("msg", parseParams(m)))
}

// Fatal Fatal 级别日志
func (l *Logger) Fatal(msg string, trace *TraceContext, dltag string, m map[string]interface{}) {
	l.L.Fatal(msg, zap.String(_traceID, trace.TraceID),
		zap.String(_spanID, trace.SpandID),
		zap.String(_childSpanID, trace.CSpanID),
		zap.String("msg", parseParams(m)))
}

// func Trace() {
// }

// CheckDLTag 检验 dltag 合法性
func CheckDLTag(dltag string) string {
	if strings.HasPrefix(dltag, _dlTagBizPrefix) {
		return dltag
	}

	if strings.HasPrefix(dltag, "_com_") {
		return dltag
	}

	if dltag == DLTagUndefind {
		return dltag
	}
	return dltag
}

func traceFormat(trace *TraceContext, dltag string) string {

	return fmt.Sprintf("traceid=%s|spanid=%s|cspanid=%s",
		trace.TraceID, trace.SpandID, trace.CSpanID)
}

// map 格式化为 string
func parseParams(m map[string]interface{}) string {
	var dltag string = "_undef"
	if _dltag, _have := m["dltag"]; _have {
		if _val, _ok := _dltag.(string); _ok {
			dltag = _val
		}
	}

	for _key, _val := range m {
		if _key == "dltag" {
			continue
		}
		dltag = dltag + "|" + fmt.Sprintf("%v=%+v", _key, _val)
	}
	dltag = strings.Trim(fmt.Sprintf("%q", dltag), "\"")
	return dltag
}
