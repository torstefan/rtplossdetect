NIC=wlxa8637d2da500
LOSS=$1

sudo tc qdisc del dev $NIC root
sudo tc qdisc add dev $NIC root netem loss $LOSS
