package templog

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"

	"github.com/leijux/templog/pkg/lwriter"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

var Close func() error

func init() {
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	name := "templog"
	info, ok := debug.ReadBuildInfo()
	if ok {
		name = info.Main.Path
	}

	_, ok = os.LookupEnv("TEMPLOG_DISABLE_FILE")

	levelWriter := lwriter.New(name, !ok, zerolog.NewConsoleWriter())
	zerologLogger := zerolog.New(levelWriter).With().Caller().Logger()
	slog.SetDefault(slog.New(slogzerolog.Option{Level: slog.LevelDebug, Logger: &zerologLogger}.NewZerologHandler()))

	Close = levelWriter.Close
}
