package util

import (
	"log"
	"os"
)

// NewLogger return a new instance of a logger
func NewLogger(name string) *log.Logger {
	return log.New(os.Stdout, name+": ", log.Ldate)
}

//Logger exposes a default logger
var Logger = NewLogger("bluetooth")
