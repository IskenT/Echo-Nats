package logger

type Logger interface {
	Warn(kv ...interface{})
	Error(kv ...interface{})
	Debug(kv ...interface{})
	Info(kv ...interface{})
	WarnF(str string, kv ...interface{})
	ErrorF(str string, kv ...interface{})
	DebugF(str string, kv ...interface{})
	InfoF(str string, kv ...interface{})
}
