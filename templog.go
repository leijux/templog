package templog

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LevelWriter struct {
	loggers   map[zerolog.Level]io.WriteCloser
	allWriter io.WriteCloser // 所有日志的汇总文件
}

// NewLevelWriter 创建分级日志写入器
func NewLevelWriter() *LevelWriter {
	var (
		basePath = "logs"

		maxSize    = 100 // 每个日志文件的最大大小，单位为 MB
		maxBackups = 3   // 保留的最大备份文件数
		maxAge     = 30  // 保留的最大天数
		compress   = true

		lw = &LevelWriter{
			loggers: make(map[zerolog.Level]io.WriteCloser),
		}

		levels = []zerolog.Level{
			zerolog.DebugLevel,
			zerolog.InfoLevel,
			zerolog.WarnLevel,
			zerolog.ErrorLevel,
		}
	)
	for _, level := range levels {
		lw.loggers[level] = &lumberjack.Logger{
			Filename:   filepath.Join(basePath, level.String(), level.String()+".log"),
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}
	}

	lw.allWriter = &lumberjack.Logger{
		Filename:   filepath.Join(basePath, "all.log"),
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	return lw
}

func (lw *LevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	// 写入对应级别的日志文件
	if writer, ok := lw.loggers[level]; ok {
		n, err = writer.Write(p)
		if err != nil {
			return n, err
		}
	}

	return lw.allWriter.Write(p)
}

func (lw *LevelWriter) Write(p []byte) (n int, err error) {
	return lw.allWriter.Write(p)
}

func (lw *LevelWriter) Close() error {
	var errs error

	for _, logger := range lw.loggers {
		if err := logger.Close(); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if err := lw.allWriter.Close(); err != nil {
		errs = errors.Join(errs, err)
	}

	if errs != nil {
		return fmt.Errorf("failed to close LevelWriter: %w", errs)
	}

	return nil
}

func init() {
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	levelWriter := NewLevelWriter()
	zerologLogger := zerolog.New(io.MultiWriter(levelWriter, zerolog.NewConsoleWriter())).With().Caller().Logger()
	slog.SetDefault(slog.New(slogzerolog.Option{Level: slog.LevelDebug, Logger: &zerologLogger}.NewZerologHandler()))

	runtime.SetFinalizer(levelWriter, func(lw *LevelWriter) { lw.Close() })
}
