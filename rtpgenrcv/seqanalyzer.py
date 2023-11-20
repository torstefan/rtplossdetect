import sys

def main():
    print("Starting RTP sequence number check...")
    prev_seq = None

    try:
        for line in sys.stdin:
            try:
                seq = int(line.strip())
                if prev_seq is not None and seq != prev_seq + 1:
                    print(f"Missing packet(s) between sequence numbers {prev_seq} and {seq} - Lost: {seq-prev_seq}")
                prev_seq = seq
            except ValueError:
                print("Invalid input. Please ensure input is only RTP sequence numbers.")
    except KeyboardInterrupt:
        print("\nSequence number check stopped.")

if __name__ == "__main__":
    main()
