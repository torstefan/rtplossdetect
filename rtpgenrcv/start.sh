#!/bin/bash
PORT=21338
export TMPDIR=$(pwd)
tcpdump -i wlan0 -w - -U udp port $PORT | stdbuf -oL tshark -i - -T fields -e rtp.seq -d udp.port==$PORT,rtp | python3 seqanalyzer.py

