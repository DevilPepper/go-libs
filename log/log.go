package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/DevilPepper/go-libs/environment"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func GetLogger() *log.Logger {
	return log.Default()
}

func GetLogLevel() log.Level {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "warn"
	}
	level, _ := log.ParseLevel(logLevel)
	return level
}

func GetLogOptions() log.Options {
	logLevel := GetLogLevel()
	timeFormat := time.RFC3339
	if environment.IS_DEV {
		timeFormat = time.TimeOnly
	}
	return log.Options{
		TimeFormat:      timeFormat,
		Level:           logLevel,
		ReportTimestamp: true,
		ReportCaller:    logLevel == log.DebugLevel,
		Formatter:       log.TextFormatter,
	}
}

func GetLogStyles() *log.Styles {
	styles := log.DefaultStyles()
	logLevels := []log.Level{
		log.DebugLevel,
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
	}

	for _, l := range logLevels {
		styles.Levels[l] = styles.Levels[l].
			MaxWidth(5).
			PaddingRight(5)
	}

	styles.Timestamp = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(0))
	styles.Caller = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(0))
	styles.Message = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(7))
	styles.Value = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(4))

	return styles
}

func InitLogger() {
	log.SetDefault(log.NewWithOptions(os.Stderr, GetLogOptions()))
	// TODO: idk why this was a problem
	if environment.IS_DEV {
		GetLogger().SetStyles(GetLogStyles())
	}
	slog.SetDefault(slog.New(GetLogger()))
}
