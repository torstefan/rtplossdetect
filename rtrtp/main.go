package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	deviceFlag   = flag.String("i", "", "Interface device name e.g., eth2 (required)")
	host         = flag.String("host", "", "Host address to filter by e.g., 10.20.30.1 (required)")
	hostType     = flag.String("t", "src", "Type of host [src|dst]")
	portRange    = flag.String("pr", "10000-20000", "Port range")
	promiscuous  = flag.Bool("promiscuous", false, "Run in promiscuous mode")
	packetCount  int
	streams      map[string]StreamInfo = make(map[string]StreamInfo)
)

type StreamInfo struct {
	lastSeq        uint16
	lastSeen       time.Time
	sequenceErrors int
}

func main() {
	const snapshotLen int32 = 65536

	flag.Parse()

	if *deviceFlag == "" || *host == "" {
		fmt.Println("Missing required flags:")
		flag.Usage()
		os.Exit(1)
	}

	handle, err := pcap.OpenLive(*deviceFlag, snapshotLen, *promiscuous, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Error opening device %s: %v", *deviceFlag, err)
	}
	defer handle.Close()

	filter := fmt.Sprintf("udp and portrange %s and %s host %s", *portRange, *hostType, *host)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatalf("Error setting BPF filter: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		processPacket(packet)
	}
}

func processPacket(packet gopacket.Packet) {
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return
	}
	udp, _ := udpLayer.(*layers.UDP)

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return
	}
	ip, _ := ipLayer.(*layers.IPv4)

	streamID := fmt.Sprintf("%s:%d -> %s:%d", ip.SrcIP, udp.SrcPort, ip.DstIP, udp.DstPort)
	currentTime := time.Now()
	streamInfo, exists := streams[streamID]

	if !exists || currentTime.Sub(streamInfo.lastSeen) > time.Minute*5 {
		fmt.Printf("[%s] New RTP stream detected: %s\n", currentTime.Format(time.RFC3339), streamID)
		streams[streamID] = StreamInfo{
			lastSeq:        getRTPSequenceNumber(udp.Payload),
			lastSeen:       currentTime,
			sequenceErrors: 0,
		}
	} else {
		rtpSeqNum := getRTPSequenceNumber(udp.Payload)
		expectedSeqNum := (streamInfo.lastSeq + 1) & 0xFFFF // Handle sequence number wrapping
		if rtpSeqNum != expectedSeqNum {
			streamInfo.sequenceErrors++
			fmt.Printf("[%s] Packet loss detected in stream %s: Expected %d, got %d (sequence errors so far: %d)\n",
				currentTime.Format(time.RFC3339), streamID, expectedSeqNum, rtpSeqNum, streamInfo.sequenceErrors)
		}
		streamInfo.lastSeq = rtpSeqNum
		streamInfo.lastSeen = currentTime
		streams[streamID] = streamInfo
	}
	packetCount++
}

func getRTPSequenceNumber(udpPayload []byte) uint16 {
	if len(udpPayload) < 4 {
		return 0
	}
	return binary.BigEndian.Uint16(udpPayload[2:4])
}

