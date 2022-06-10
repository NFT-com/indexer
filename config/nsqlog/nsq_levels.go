package nsqlog

import (
	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog"
)

func ToNSQLevel(lvl zerolog.Level) nsq.LogLevel {

	switch lvl {

	case zerolog.TraceLevel:
		return nsq.LogLevelDebug

	case zerolog.DebugLevel:
		return nsq.LogLevelDebug

	case zerolog.InfoLevel:
		return nsq.LogLevelInfo

	case zerolog.WarnLevel:
		return nsq.LogLevelWarning

	case zerolog.ErrorLevel:
		return nsq.LogLevelError

	default:
		return nsq.LogLevelError
	}
}
