package wenex

import (
	"log"
	"os"
	"path"
	"sync"
)

func newLogger(wnx *Wenex, name string) (func(string) *log.Logger, error) {
	filePrefix, err := wnx.Config.String("log.filePrefix")
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

		var file *os.File

		switch filePrefix {
		case "+stdout":
			file = os.Stdout
		case "+stderr":
			file = os.Stderr
		default:
			if err = os.MkdirAll(path.Dir(filePrefix+name), 0755); err == nil {
				if file, err = os.OpenFile(filePrefix+name+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644); err != nil {
					file = os.Stdout
				}
			} else {
				file = os.Stdout
			}
		}

		if loggers[""] != nil {
			loggers[name] = log.New(file, loggers[""].Prefix(), loggers[""].Flags())
		} else {
			if file.Name() == "/dev/stdout" || file.Name() == "/dev/stderr" {
				filePrefix = name + ": "
			} else {
				filePrefix = ""
			}

			loggers[name] = log.New(file, "[!] "+filePrefix, log.LstdFlags)
		}

		return loggers[name]
	}

	loggers[""] = f(name)
	return f, err
}
