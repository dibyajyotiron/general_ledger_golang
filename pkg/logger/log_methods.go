package logger

import "github.com/sirupsen/logrus"

// Debug logs a message at level Debug on the standard logger.
func Debug(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Debugf(format, args...)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Infof(format, args...)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Warnf(format, args...)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		logger.SetReportCaller(true)
		entry := logger.WithFields(logrus.Fields{})
		entry.Errorf(format, args...)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Fatalf(format, args...)
	}
}
