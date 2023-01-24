package logger
import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct {
	L zerolog.Logger
}

func New() *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.Stamp
	zerolog.ErrorFieldName = "err"

	logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = int(zerolog.InfoLevel)
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.Stamp,

		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
		FormatCaller: func(i interface{}) string {
			path := filepath.Base(fmt.Sprintf("%s", i))
			return fmt.Sprintf("%-18s >", path)
		},
	}

	logger := &Logger{
		L: zerolog.New(consoleWriter).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Caller().
			Logger(),
	}

	return logger
}
