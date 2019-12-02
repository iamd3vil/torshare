package models

// RelayMsg is the message to be sent to and received from the relay.
type RelayMsg struct {
	TorAddr  string `json:"tor_addr"`
	FileName string `json:"filename"`
}

// RelayReply is the reply from the relay
type RelayReply struct {
	Channel string `json:"channel"`
}
