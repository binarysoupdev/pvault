package logger

import (
	"log"
	"os"

	"github.com/binarysoupdev/go-extensions/errors"
)

var logger *log.Logger

func Open(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Chain(err, "error opening/creating log file")
	}

	logger = log.New(file, "", log.Ldate|log.Ltime)
	return file, nil
}
