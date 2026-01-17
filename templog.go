package templog

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LevelWriter struct {
	loggers   map[zerolog.Level]io.WriteCloser
	allWriter []io.WriteCloser // 所有日志的汇总文件
}

// NewLevelWriter 创建分级日志写入器
func NewLevelWriter(writers ...io.WriteCloser) *LevelWriter {
	lw := &LevelWriter{}

	_, ok := os.LookupEnv("TEMPLOG_DISABLE_FILE")
	if ok {
		lw.allWriter = writers
		return lw
	}

	var (
		basePath = "logs"

		maxSize    = 100 // 每个日志文件的最大大小，单位为 MB
		maxBackups = 3   // 保留的最大备份文件数
		maxAge     = 30  // 保留的最大天数
		compress   = true

		levels = []zerolog.Level{
			zerolog.DebugLevel,
			zerolog.InfoLevel,
			zerolog.WarnLevel,
			zerolog.ErrorLevel,
		}
	)

	lw.loggers = make(map[zerolog.Level]io.WriteCloser, 4)
	for _, level := range levels {
		lw.loggers[level] = &lumberjack.Logger{
			Filename:   filepath.Join(basePath, level.String(), level.String()+".log"),
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}
	}

	lw.allWriter = append(append(make([]io.WriteCloser, 0, len(writers)+1), &lumberjack.Logger{
		Filename:   filepath.Join(basePath, "all.log"),
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}), writers...)

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

	for _, writer := range lw.allWriter {
		n, err = writer.Write(p)
		if err != nil {
			return n, err
		}
	}

	return len(p), nil
}

func (lw *LevelWriter) Write(p []byte) (n int, err error) {
	for _, writer := range lw.allWriter {
		n, err = writer.Write(p)
		if err != nil {
			return n, err
		}
	}

	return len(p), nil
}

func (lw *LevelWriter) Close() error {
	var errs []error

	for _, logger := range lw.loggers {
		if err := logger.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	for _, writer := range lw.allWriter {
		if err := writer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if errs != nil {
		return fmt.Errorf("failed to close LevelWriter: %w", errors.Join(errs...))
	}

	return nil
}

var Close func() error

func init() {
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	levelWriter := NewLevelWriter(zerolog.NewConsoleWriter())
	zerologLogger := zerolog.New(levelWriter).With().Caller().Logger()
	slog.SetDefault(slog.New(slogzerolog.Option{Level: slog.LevelDebug, Logger: &zerologLogger}.NewZerologHandler()))

	Close = levelWriter.Close
}
