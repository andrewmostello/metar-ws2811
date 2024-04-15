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
	Format    string
	Level     string
	AddSource bool
}

func GetLog() Log {
	return Log{
		Format:    viper.GetString(cfgKeyLogFormat),
		Level:     viper.GetString(cfgKeyLogLevel),
		AddSource: viper.GetBool(cfgKeyLogAddSource),
	}
}

const (
	cfgKeyLogFormat    = "log.format"
	cfgKeyLogLevel     = "log.level"
	cfgKeyLogAddSource = "log.add_source"
)

func AddLogFlags(cmd *cobra.Command) {
	flag := "log-format"
	cmd.PersistentFlags().String(flag, "text", "Log format. Options are text and json.")
	viper.BindPFlag(cfgKeyLogFormat, cmd.PersistentFlags().Lookup(flag))

	flag = "log-level"
	cmd.PersistentFlags().String(flag, "warn", "Log Level. Options are error, warn, info, and debug.")
	viper.BindPFlag(cfgKeyLogLevel, cmd.PersistentFlags().Lookup(flag))

	flag = "log-add-source"
	cmd.PersistentFlags().Bool(flag, false, "Add source to log output.")
	viper.BindPFlag(cfgKeyLogAddSource, cmd.PersistentFlags().Lookup(flag))
}

func NewLogger() *slog.Logger {

	cfg := GetLog()

	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
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
