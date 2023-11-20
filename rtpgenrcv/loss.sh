NIC=wlan0
LOSS=$1

sudo tc qdisc del dev $NIC root
sudo tc qdisc add dev $NIC root netem loss $LOSS
