package main

import (
	"context"
	"fmt"

	"github.com/cretz/bine/tor"
)

func newTorService() (*tor.Tor, *tor.OnionService, error) {
	t, err := tor.Start(context.Background(), &tor.StartConf{
		DataDir: "/tmp/datadir",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error while starting Tor service: %v", err)
	}

	onion, err := t.Listen(context.Background(), &tor.ListenConf{
		RemotePorts: []int{80},
		Version3:    true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error while starting Onion service: %v", err)
	}

	return t, onion, nil
}

// func startTorService(dir string) error {
// 	t, err := tor.Start(context.Background(), &tor.StartConf{
// 		DataDir: "/tmp/datadir",
// 	})
// 	if err != nil {
// 		log.Fatalf("error while starting Tor service: %v", err)
// 	}

// 	defer t.Close()

// 	onion, err := t.Listen(context.Background(), &tor.ListenConf{
// 		RemotePorts: []int{80},
// 		Version3:    true,
// 	})
// 	if err != nil {
// 		log.Fatalf("error while starting Onion service: %v", err)
// 	}
// 	defer onion.Close()

// 	log.Printf("Started onion service on: %v", onion.ID)

// 	return http.Serve(onion, http.FileServer(http.Dir(dir)))
// }
