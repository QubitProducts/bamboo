// +build linux

package lumberjack_test

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/natefinch/lumberjack"
)

// Example of how to rotate in response to SIGHUP.
func ExampleLogger_Rotate() {
	l := &lumberjack.Logger{}
	log.SetOutput(l)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			l.Rotate()
		}
	}()
}
