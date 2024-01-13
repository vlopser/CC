package utils

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func WaitForSigkill(wg *sync.WaitGroup) {

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGINT)

	<-c

	wg.Done()

	os.Exit(0)
}
