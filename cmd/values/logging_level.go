package values

import (
	log "github.com/sirupsen/logrus"
)

// From github.com/spf13/pflag
//type Value interface {
//	String() string
//	Set(string) error
//	Type() string
//}

type LoggingLevelValue struct {
	Level log.Level
}

func (loggingLevel *LoggingLevelValue) String() string {
	return loggingLevel.Level.String()
}

func (loggingLevel *LoggingLevelValue) Set(lvl string) error {
	l, err := log.ParseLevel(lvl)
	if err != nil {
		return err
	}

	loggingLevel.Level = l
	return nil
}

func (loggingLevel *LoggingLevelValue) Type() string {
	return "loggingLevel"
}

func (loggingLevel *LoggingLevelValue) Get() interface{} {
	return loggingLevel.Level
}
