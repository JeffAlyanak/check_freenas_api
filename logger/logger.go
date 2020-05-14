package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Logger struct bundles log.Logger and filename
type Logger struct {
	filename string
	*log.Logger
}

var lg *Logger
var once sync.Once

// Get Returns a Logger
func Get() (lg *Logger, err error) { // TODO: Add user configurable log directory

	once.Do(func() {
		lg, err = makeLogger("check_freenas_api.log")
	})
	if err != nil {
		return nil, err
	}
	return lg, nil
}

func makeLogger(fname string) (*Logger, error) {
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("Could not open log file at: ", fname)
		return nil, err
	}

	return &Logger{
		filename: fname,
		Logger:   log.New(file, "", log.Ldate|log.Ltime),
	}, nil
}
