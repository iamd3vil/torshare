package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cretz/bine/tor"
)

func newTorService() (*tor.Tor, *tor.OnionService, error) {
	os.RemoveAll("/tmp/datadirsender")
	t, err := tor.Start(context.Background(), &tor.StartConf{
		DataDir: "/tmp/datadirsender",
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
