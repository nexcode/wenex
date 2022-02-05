package wenex

import (
	"github.com/nexcode/joneva"
	"io"
	"log"
	"os"
	"path"
	"sync"
)

// LogWriter used to connect custom loggers to wenex
type LogWriter interface {
	GetWriter(string) (io.Writer, error)
}

func NewLogger(config *joneva.Joneva, logWriter LogWriter) (func(string) *log.Logger, error) {
	defaultName, err := config.String("logger.defaultName")
	if err != nil {
		return nil, err
	}

	if defaultName == "" {
		return nil, ErrDefaultLogEmpty
	}

	namePrefix, err := config.String("logger.namePrefix")
	if err != nil {
		return nil, err
	}

	usePrefix, err := config.String("logger.usePrefix")
	if err != nil {
		return nil, err
	}

	useFlag, err := config.Float64("logger.useFlag")
	if err != nil {
		return nil, err
	}

	var mutex sync.Mutex
	loggers := make(map[string]*log.Logger)

	f := func(name string) *log.Logger {
		mutex.Lock()
		defer mutex.Unlock()

		if logger, ok := loggers[name]; ok {
			return logger
		}

		var writer io.Writer

		if logWriter == nil {
			if err = os.MkdirAll(path.Dir(namePrefix+name), 0755); err != nil {
				return loggers[""]
			}

			if writer, err = os.OpenFile(namePrefix+name+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644); err != nil {
				return loggers[""]
			}
		} else if writer, err = logWriter.GetWriter(namePrefix + name); err != nil {
			return loggers[""]
		}

		loggers[name] = log.New(writer, usePrefix, int(useFlag))
		return loggers[name]
	}

	loggers[""] = f(defaultName)
	return f, err
}
