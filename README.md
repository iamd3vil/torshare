# TorShare :rocket: :zap:

TorShare enables file sharing between two parties over Tor.

### Design

TorShare needs a relay server which just relays information between two parties. The relay server can't see the data being relayed since the data is encrypted using a shared key agreed between two parties.

The sender starts a Tor hidden service for the file which needs to be transmitted and the address is sent to the receiver using the relay server. The hidden node address and the file name to be received is encrypted using a password agreed between both sender and receiver. 

When the sender starts the hidden node and sends the info the relay, the relay will return a channel name of 6 word phrases. The receiver needs to know both the password and the channel name to receive the file.

The relay shouldn't be able to figure out anything related to the transfer. Since the file transfer is also done over the Tor network, the transfer is encrypted as well.

Once the transfer is completed, the receiver will send a message directly to sender and sender will teardown the Tor hidden service.

### Installation

Latest version of Tor has to be installed before running the binary.

For ubuntu, instead of installing ubuntu's repos, see [here](https://2019.www.torproject.org/docs/debian.html.en) and use Option two in that page to get the latest version of Tor.

Precompiled binaries can be downloaded from [here](https://github.com/iamd3vil/torshare/releases).

### Usage

```bash
$ torshare -h
Usage of ./dist/torshare.bin:
  -file string
    	File Path
  -relay string
    	Relay Address (default "https://torshare.sarat.dev")

```

#### Sending a file

```bash
$ torshare -file sample.jpeg
Enter Password: **********
2019/12/07 23:18:22 Channel Name(Has to be communicated with Receiver): luckily-slogan-twitter-haunt-joystick-earwig
```

Channel name and password has to be shared with the receiver

#### Receiving a file

```bash
$ torshare
Enter channel name: luckily-slogan-twitter-haunt-joystick-earwig
Enter Password: **********
2019/12/08 14:19:07 Starting download for Burton Stein - The New Cambridge History of India_ Vijayanagara-Cambridge University Press (1990).pdf
757.84 KiB / 9.04 MiB [===============>-----------------------------------]   8.19% 86.77 KiB/s 01m37s
```
