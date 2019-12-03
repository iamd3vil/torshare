package main

import (
	"flag"
	"path"
)

func main() {
	filePath := flag.String("file", "", "File Path")
	relay := flag.String("relay", "https://torshare.sarat.dev", "Relay Address")
	flag.Parse()

	dir, filename := path.Split(*filePath)

	if *filePath != "" {
		startSender(dir, filename, *relay)
	} else {
		startReceiver(*relay)
	}
}
