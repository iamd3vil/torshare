package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/iamd3vil/torshare/models"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/secretbox"
)

func startSender(dir, filename, relay string) {
	var password string
	fmt.Printf("Enter Password:")
	_, err := fmt.Scan(&password)
	if err != nil {
		log.Fatalf("error while reading password: %v", err)
	}

	if password == "" {
		log.Fatalln("Password cannot be empty")
	}

	// Initialize a tor service
	t, onion, err := newTorService()
	if err != nil {
		log.Fatalf("error while starting tor service: %v", err)
	}

	defer t.Close()
	defer onion.Close()

	rMsg := models.RelayMsg{
		TorAddr:  fmt.Sprintf("http://%s.onion", onion.ID),
		FileName: filename,
	}

	channel, err := sendMsgToRelay(rMsg, relay, password)
	if err != nil {
		log.Fatalf("error while handshake: %v", err)
	}

	log.Printf("Channel Name(Has to be communicated with Receiver): %s", channel)

	// Start service
	log.Fatalln(http.Serve(onion, http.FileServer(http.Dir(dir))))
}

// sendMsgToRelay sends initial handshake message to relay
// This returns a channel name which the receiver has to send to relay
// for getting details.
func sendMsgToRelay(rMsg models.RelayMsg, relay, password string) (string, error) {
	// Marshal data
	j, err := json.Marshal(rMsg)
	if err != nil {
		return "", err
	}

	// Encrypt the marshalled data using the password.
	// We can use secretbox which uses XSalsa20 and Poly1305 to encrypt and authenticate messages
	// For getting key from password, we will hash the password with Blake2b-256
	key := blake2b.Sum256(j)

	// Generate random nonce
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", fmt.Errorf("error while encrypting data: %v", err)
	}

	// Seal using secretbox
	encrypted := secretbox.Seal(nonce[:], []byte("hello world"), &nonce, &key)

	c := http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := c.Post(relay, "application/json", bytes.NewReader(encrypted))
	if err != nil {
		return "", fmt.Errorf("error while sending message to relay: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error in response from relay: %v", resp.StatusCode)
	}

	reply := models.RelayReply{}
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return "", errors.New("error reading JSON response from relay")
	}

	return reply.Channel, nil
}
