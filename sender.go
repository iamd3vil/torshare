package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/howeyc/gopass"
	"github.com/iamd3vil/torshare/models"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/secretbox"
)

type hub struct {
	shutdownChan chan int
	srv          *http.Server
}

func (h *hub) Serve(l net.Listener) error {
	go func() {
		<-h.shutdownChan

		log.Println("Transfer complete. Shutting down in 5 secs...")
		time.Sleep(5 * time.Second)

		h.srv.Shutdown(context.Background())
	}()
	return h.srv.Serve(l)
}

func (h *hub) handleTransferComplete(w http.ResponseWriter, r *http.Request) {
	// Shutdown sender since transfer is complete
	close(h.shutdownChan)
	w.Write([]byte("OK"))
}

func startSender(dir, filename, relay string) {
	fmt.Printf("Enter Password: ")

	password, err := gopass.GetPasswdMasked()
	if err != nil {
		log.Fatalf("error while reading password: %v", err)
	}

	if string(password) == "" {
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

	channel, err := sendMsgToRelay(rMsg, relay, string(password))
	if err != nil {
		log.Printf("error while handshake: %v", err)
		return
	}

	log.Printf("Channel Name(Has to be communicated with Receiver): %s", channel)

	h := hub{
		shutdownChan: make(chan int),
	}

	r := mux.NewRouter()
	r.PathPrefix("/file/").Handler(http.StripPrefix("/file/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/v1/complete", h.handleTransferComplete)

	srv := &http.Server{
		Handler: r,
	}

	h.srv = srv

	// Start service
	if err := h.Serve(onion); err != http.ErrServerClosed {
		log.Fatalln(err)
	}
	log.Println("Shutting down...")
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

	resp, err := c.Post(fmt.Sprintf("%s/v1/relay", relay), "application/json", bytes.NewReader(encrypted))
	if err != nil {
		return "", fmt.Errorf("error while sending message to relay: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error in response from relay: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response from relay: %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	reply := models.RelayReply{}
	err = json.Unmarshal(body, &reply)
	if err != nil {
		return "", errors.New("error reading JSON response from relay")
	}

	return reply.Data.Channel, nil
}
