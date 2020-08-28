package logger

import (
	"dubhe-ci/config"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/orandin/lumberjackrus"
	"github.com/sirupsen/logrus"
)

const (
	ActionKey = "action"
	RoleKey   = "role"
)

func InitLogger(log config.Log) error {
	logrus.SetLevel(logrus.Level(log.Level))
	formatter := getFormatter(log.Format, true)
	logrus.SetFormatter(formatter)

	if log.InfoPath != "" {
		var opts *lumberjackrus.LogFileOpts = nil
		if log.ErrorPath != "" {
			opts = &lumberjackrus.LogFileOpts{
				logrus.ErrorLevel: &lumberjackrus.LogFile{
					Filename:   log.ErrorPath,
					MaxSize:    100,
					MaxAge:     1,
					MaxBackups: 1,
					LocalTime:  true,
					Compress:   false,
				},
			}
		}
		hook, err := lumberjackrus.NewHook(
			&lumberjackrus.LogFile{
				Filename:   log.InfoPath,
				MaxSize:    100,
				MaxAge:     1,
				MaxBackups: 1,
				LocalTime:  true,
				Compress:   false,
			},
			logrus.GetLevel(),
			getFormatter(log.Format, true),
			opts,
		)
		if err != nil {
			return err
		}

		logrus.AddHook(hook)
	}

	return nil
}

func getFormatter(format string, isFile bool) logrus.Formatter {
	var formatter logrus.Formatter
	switch format {
	case "json":
		formatter = new(logrus.JSONFormatter)
	case "nested":
		formatter = &nested.Formatter{
			NoColors:        isFile,
			TimestampFormat: "2006-01-02 15:04:05",
		}
	default:
		formatter = new(logrus.TextFormatter)
	}

	return formatter
}

func WithAction(action string) *logrus.Entry {
	return logrus.WithField(ActionKey, action)
}
func WithRole(role string) *logrus.Entry {
	return logrus.WithField(RoleKey, role)
}
