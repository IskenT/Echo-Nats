package zerolog

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type ZeroLogWrapper struct {
	log zerolog.Logger
}

func NewZeroLog(logWriter io.Writer) (*ZeroLogWrapper, error) {
	ioWriter := zerolog.ConsoleWriter{
		Out:        logWriter,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
	}

	lvl, err := zerolog.ParseLevel(zerolog.ErrorLevel.String())
	if err != nil {
		return nil, err
	}

	zeroLog := zerolog.New(ioWriter)
	zeroLog = zeroLog.Level(lvl).With().Timestamp().Logger()

	return &ZeroLogWrapper{
		log: zeroLog,
	}, nil
}

func (z *ZeroLogWrapper) Warn(kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Warn().Msg(msg)
}

func (z *ZeroLogWrapper) WarnF(str string, kv ...interface{}) {
	msg := fmt.Sprintf(str, kv...)
	z.log.Warn().Msg(msg)
}

func (z *ZeroLogWrapper) Error(kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Error().Msg(msg)
}

func (z *ZeroLogWrapper) ErrorF(str string, kv ...interface{}) {
	msg := fmt.Sprintf(str, kv...)
	z.log.Error().Msg(msg)
}

func (z *ZeroLogWrapper) Debug(kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Debug().Msg(msg)
}

func (z *ZeroLogWrapper) DebugF(str string, kv ...interface{}) {
	msg := fmt.Sprintf(str, kv...)
	z.log.Debug().Msg(msg)
}

func (z *ZeroLogWrapper) Info(kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Info().Msg(msg)
}

func (z *ZeroLogWrapper) InfoF(str string, kv ...interface{}) {
	msg := fmt.Sprintf(str, kv...)
	z.log.Info().Msg(msg)
}
