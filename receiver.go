package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/howeyc/gopass"
	"github.com/iamd3vil/torshare/models"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/secretbox"
)

func startReceiver(relay string) {
	var channel string
	fmt.Printf("Enter Channel Name: ")
	fmt.Scan(&channel)

	// Ask relay for metadata
	rMsg, err := askRelayForMeta(relay, channel)
	if err != nil {
		log.Fatalln(err)
	}

	os.RemoveAll("/tmp/datadirreceiver")
	t, err := tor.Start(context.Background(), &tor.StartConf{
		DataDir: "/tmp/datadirreceiver",
	})
	defer t.Close()

	log.Printf("Starting download for %s", rMsg.FileName)
	err = DownloadFile(t, fmt.Sprintf("%s/file/%s", rMsg.TorAddr, rMsg.FileName), rMsg.FileName)
	if err != nil {
		log.Printf("error while downloading file: %v", err)
		return
	}

	err = sendCompleteSignal(t, rMsg)
	if err != nil {
		log.Println(err)
		return
	}
}

func askRelayForMeta(relay, channel string) (*models.RelayMsg, error) {
	client := http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s/v1/relay?channel=%s", relay, channel))
	if err != nil {
		return nil, fmt.Errorf("error while sending message to relay: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while readin response from relay: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp := &struct {
			Message string `json:"message"`
		}{}
		json.Unmarshal(body, errResp)
		return nil, fmt.Errorf("error from relay: %s, status code: %d", errResp.Message, resp.StatusCode)
	}

	fmt.Print("Enter password: ")
	password, err := gopass.GetPasswdMasked()
	if err != nil {
		log.Fatalf("error while reading password: %v", err)
	}

	if string(password) == "" {
		return nil, errors.New("Password cannot be empty")
	}

	key := blake2b.Sum256([]byte(password))

	var decryptNonce [24]byte
	copy(decryptNonce[:], body[:24])

	decrypted, ok := secretbox.Open(nil, body[24:], &decryptNonce, &key)
	if !ok {
		return nil, errors.New("invalid password")
	}

	rMsg := &models.RelayMsg{}

	err = json.Unmarshal(decrypted, rMsg)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling relay response: %v", err)
	}

	return rMsg, nil
}

func sendCompleteSignal(t *tor.Tor, rMsg *models.RelayMsg) error {
	// Make connection
	dialer, err := t.Dialer(context.Background(), nil)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
		},
	}

	_, err = httpClient.Get(fmt.Sprintf("%s/v1/complete", rMsg.TorAddr))
	if err != nil {
		return fmt.Errorf("error while sending transfer complete signal: %v", err)
	}
	return nil
}
