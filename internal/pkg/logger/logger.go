package logger

type Logger interface {
	Info(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
}

type Field struct {
	Key   string
	Value any
}
