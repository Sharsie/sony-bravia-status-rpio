package logger

import (
	"log"

	"github.com/Sharsie/tv-status-rpio/cmd/is-on/config"
)

type Log struct{}

func (l *Log) Debug(message string, args ...interface{}) {
	if config.Debug {
		log.Printf(message+"\n", args...)
	}
}
