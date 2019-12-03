package models

// RelayMsg is the message to be sent to and received from the relay.
type RelayMsg struct {
	TorAddr  string `json:"tor_addr"`
	FileName string `json:"filename"`
}

type Channel struct {
	Channel string `json:"channel"`
}

// RelayReply is the reply from the relay
type RelayReply struct {
	Message string  `json:"message"`
	Data    Channel `json:"data"`
}

// RelayReplyMeta returns relay's reply for metadata about a channel
type RelayReplyMeta struct {
	Message string `json:"message"`
	Data    []byte `json:"metadata"`
}
