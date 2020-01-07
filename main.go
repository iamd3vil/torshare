package main

import (
	"flag"
	"path"
)

const (
	// DefaultRelay is the default relay
	DefaultRelay = "http://torwg7vip577vofy2rqtervn7r42f6rh5ptsijmyhkmepr4e6vrddsid.onion"
)

func main() {
	filePath := flag.String("file", "", "File Path")
	relay := flag.String("relay", DefaultRelay, "Relay Address")
	flag.Parse()

	dir, filename := path.Split(*filePath)

	if *filePath != "" {
		startSender(dir, filename, *relay)
	} else {
		startReceiver(*relay)
	}
}
