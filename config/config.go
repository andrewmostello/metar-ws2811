package config

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Log struct {
	Format string
	Level  string
}

func GetLog() Log {
	return Log{
		Format: viper.GetString(cfgKeyLogFormat),
		Level:  viper.GetString(cfgKeyLogLevel),
	}
}

const (
	cfgKeyLogFormat = "log.format"
	cfgKeyLogLevel  = "log.level"
)

func AddLogFlags(cmd *cobra.Command) {
	viper.SetDefault(cfgKeyLogFormat, "text")
	viper.SetDefault(cfgKeyLogLevel, "warn")

	cmd.PersistentFlags().String("log-format", "text", "Log format. Options are text and json. Default is text.")
	viper.BindPFlag(cfgKeyLogFormat, cmd.PersistentFlags().Lookup("log-format"))

	cmd.PersistentFlags().String("log-level", "warn", "Log Level. Options are error, warn, info, and debug. Default is info.")
	viper.BindPFlag(cfgKeyLogLevel, cmd.PersistentFlags().Lookup("log-level"))
}

func NewLogger() *slog.Logger {

	cfg := GetLog()

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	switch strings.ToLower(cfg.Level) {
	case "error":
		opts.Level = slog.LevelError
	case "warn":
		opts.Level = slog.LevelWarn
	case "info":
		opts.Level = slog.LevelInfo
	case "debug":
		opts.Level = slog.LevelDebug
	}

	var out io.Writer = os.Stdout
	var h slog.Handler

	switch cfg.Format {
	case "json":
		h = slog.NewJSONHandler(out, nil)
	case "text":
		fallthrough
	default:
		h = slog.NewTextHandler(out, opts)
	}

	return slog.New(h)
}
