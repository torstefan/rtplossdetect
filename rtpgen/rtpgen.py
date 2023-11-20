import socket
import random
import time
import struct

# Adjusted function to create a random payload
def create_random_payload(size):
    return bytes([random.randint(0, 255) for _ in range(size)])

# Adjusted function to create a single RTP packet
def create_rtp_packet(sequence_number):
    # RTP Header Fields
    version = 2
    padding = 0
    extension = 0
    csrc_count = 0
    marker = 0
    payload_type = 0  # For demonstration, can be changed as needed

    # Construct RTP Header
    header = (
        (version << 6) | (padding << 5) | (extension << 4) | csrc_count,
        (marker << 7) | payload_type,
        sequence_number,
        0  # Timestamp, for simplicity set to 0, can be modified as needed
    )
    rtp_header = struct.pack("!BBHl", *header)

    # Create a random payload
    payload = create_random_payload(payload_size)

    return rtp_header + payload

# RTP packet generator parameters
source_port = 11337
destination_port = 21338
destination_ip = "192.168.10.146"  # Replace with actual IP
payload_size = 160
packet_interval = 0.02

# Function to send RTP packets in a loop
def send_rtp_packets():
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as sock:
        sock.bind(("", source_port))
        sequence_number = 0
        while True:
            packet = create_rtp_packet(sequence_number)
            sock.sendto(packet, (destination_ip, destination_port))
            sequence_number = (sequence_number + 1) % 65536
            time.sleep(packet_interval)

# Uncomment the line below to run the packet sender
send_rtp_packets()

