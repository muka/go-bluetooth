package util

import (
	"log"
	"os"
)

// NewLogger return a new instance of a logger
func NewLogger(name string) *log.Logger {
	return log.New(os.Stdout, name+": ", log.Ldate|log.Ltime|log.Lshortfile)
}

//Logger exposes a default logger
var Logger = NewLogger("bluetooth")
