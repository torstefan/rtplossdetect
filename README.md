# RtpLossDetect

RtpLossDetect is a small series of scripts / programs that can help you detect packetloss in your network using RTP as baseline protocoll.

## Installation

Programs need Golang with pcap and Python3

```bash
# rtpgen needs golang and golangpcap for sniffing packets
sudo apt-get install  golang 
sudo apt-get install  golang-github-akrennmair-gopcap-dev/stable
cd rtpgen
go build

# rtpgen / rtpgenrcv needs python3
sudo apt-get install python3

# loss.sh needs tc - traffic control part of iproute2
sudo apt-get install iproute2
```

## Usage
Please, before use read the differenet scripts. Get to know them. 

Modify rtpgen.py with the right vars

### Setup - Send side
```python
# On the send side
# RTP packet generator parameters
source_port = 11337
destination_port = 21338
destination_ip = "192.168.10.146"  # Replace with actual IP
payload_size = 160
packet_interval = 0.02
```

### Setup - RCV Side
```bash
# On the rcv side - edit rtpgenrcv/start.sh - set the right vars
PORT=21338
NIC=wlan0
```

### Usage

```bash
# Start generating packets
python rtpgen/rtpgen.py

# Simulate loss
bash rtpgen/loss.sh 1%
```

On the rcv side
```bash
# After setting the right vars in start.sh
bash rtpgenrcv/start.sh
```

start.sh will generate output if it detects packetloss from the rtpgen. E.g. if sequence number in rtp stream does not add up.

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
