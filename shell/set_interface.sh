#!/bin/bash

# simulates ip link set to bring network interface up/down

# set_interface.sh - Network interface state management


# Function to display usage info
show_usage() {
    echo "Usage: $0 <interface> <state>"
    echo "Example: $0 eth0 up"
    echo "Example: $0 eth1 down"
    echo ""
    echo "Valid interfaces: eth0, eth1, eth2, ethx"
    echo "Valid states: up, down"
    echo ""
    echo "simulates: ip link set <interface> <state>"
}

# Check if both args provided
if [ $# -ne 2 ]; then
    echo "ERROR: Incorrect number of arguments"
    show_usage
    exit 1
fi

# Get arguments
INTERFACE="$1"
STATE="$2"

# Validate interface name
case "$INTERFACE" in
    eth0|eth1|eth2|ethx)
        echo "Valid interface: $INTERFACE"
        ;;
    *)
        echo "ERROR: Invalid interface '$INTERFACE'"
        echo "Valid interfaces: eth0, eth1, eth2, ethx"
        exit 1
        ;;
esac

# Validate state arg, accept up/down or UP/DOWN
case "$STATE" in
    up|UP)
        STATE="up"
        ACTION="bringing up"
        ;;
    down|DOWN)
        STATE="down"
        ACTION="bringing down"
        ;;
    *)
        echo "ERROR: Invalid state '$STATE'"
        echo "Valid states: up, down"
        exit 1
        ;;
esac

# Simulate interface state change
echo "SwitchConfigSim: $ACTION interface $INTERFACE"
echo "Simulating: ip link set $INTERFACE $STATE"

# Simulate output based on the operation
if [ "$STATE" = "up" ]; then
    echo "Interface $INTERFACE: Link is now UP"
    echo "Interface $INTERFACE: Carrier detected"
else
    echo "Interface $INTERFACE: Link is now DOWN"
    echo "Interface $INTERFACE: Carrier lost"
fi

# Simulate checking the result
echo "Interface state change completed successfully"
echo "Interface: $INTERFACE, State: $STATE"

# Return success exit code
exit 0 