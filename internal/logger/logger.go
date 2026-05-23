package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
)

func Init(level string) {

	Info = log.New(
		os.Stdout,
		"[INFO] ",
		log.Ldate|log.Ltime,
	)

	Error = log.New(
		os.Stderr,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	debugWriter := io.Discard

	if level == "debug" {
		debugWriter = os.Stdout
	}

	Debug = log.New(
		debugWriter,
		"[DEBUG] ",
		log.Ldate|log.Ltime,
	)
}