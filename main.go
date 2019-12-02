package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
)

func main() {
	filePath := flag.String("file", "", "File Path")
	relay := flag.String("relay", "https://torshare.sarat.dev", "Relay Address")
	flag.Parse()

	dir, filename := path.Split(*filePath)

	done := make(chan struct{})

	c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-c:
			done <- struct{}{}
		}
	}()

	if *filePath != "" {
		startSender(dir, filename, *relay)
	}

	<-done
	log.Println("Exiting.....")
}
