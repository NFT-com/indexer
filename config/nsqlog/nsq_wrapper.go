package nsqlog

import (
	"github.com/rs/zerolog"
)

type NSQWrapper struct {
	log zerolog.Logger
}

func WrapForNSQ(log zerolog.Logger) NSQWrapper {

	w := NSQWrapper{
		log: log,
	}

	return w
}

func (w NSQWrapper) Output(calldepth int, msg string) error {
	w.log.Trace().Msg(msg)
	return nil
}
