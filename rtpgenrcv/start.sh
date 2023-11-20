#!/bin/bash
PORT=21338
NIC=wlan0
export TMPDIR=$(pwd)
tcpdump -i $NIC -w - -U udp port $PORT | stdbuf -oL tshark -i - -T fields -e rtp.seq -d udp.port==$PORT,rtp | python3 seqanalyzer.py

