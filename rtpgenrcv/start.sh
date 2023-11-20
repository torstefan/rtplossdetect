#!/bin/bash
PORT=$1
NIC=A1

export TMPDIR=$(pwd)

tcpdump -i $NIC -w - -U udp port $PORT | stdbuf -oL tshark -i - -w hei.pcap -b filesize:1024 -b files:1 -T fields -e rtp.seq -d udp.port==$PORT,rtp | python2 seqanalyzer.py
