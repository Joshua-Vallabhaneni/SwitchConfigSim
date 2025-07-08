#!/bin/bash

# get_status.sh - Simulates getting system status info
# simulates checking health of system: systemctl status, uptime, ip link show, and more

# func to display usage info
show_usage() {
    echo "Usage: $0 [component]"
    echo "Example: $0"
    echo "Example: $0 interfaces"
    echo "Example: $0 system"
    echo ""
    echo "Components:"
    echo "  system      - Show system health and uptime"
    echo "  interfaces  - Show all interface states"
    echo "  (no args)   - Show complete status overview"
    echo ""
    echo "simulates: systemctl status, uptime, ip link show, etc."
}

# func to show system info
show_system_status() {
    echo "=== System Status ==="
    echo "Hostname: $(echo $HOSTNAME)"
    echo "Status: Operational"
    echo "Uptime: $(uptime | awk '{print $3, $4}' | sed 's/,//')"
    echo "Load Average: $(uptime | awk '{print $(NF-2), $(NF-1), $NF}')"
    echo "Memory Usage: 2.1GB / 8.0GB (26%)"
    echo "Temperature: 45°C (Normal)"
    echo "Fan Status: Normal"
    echo ""
}

# func to show interface states
show_interface_status() {
    echo "=== Interface Status ==="
    echo "eth0: UP   - 1000Mbps Full Duplex - Carrier: Detected"
    echo "eth1: DOWN - Interface administratively down"
    echo "eth2: UP   - 1000Mbps Full Duplex - Carrier: Detected"
    echo "ethx: UP   - 10000Mbps Full Duplex - Carrier: Detected"
    echo ""
}

# func to show network config
show_network_status() {
    echo "=== Network Configuration ==="
    echo "Management IP: 192.168.1.100/24"
    echo "Default Gateway: 192.168.1.1"
    echo "DNS Servers: 8.8.8.8, 8.8.4.4"
    echo "SNMP Status: Enabled"
    echo "SSH Status: Enabled"
    echo ""
}

# func to show complete overview
show_complete_status() {
    echo "SwitchConfigSim Status Report"
    echo "Generated: $(date)"
    echo "=========================================="
    echo ""
    
    show_system_status
    show_interface_status
    show_network_status
    
    echo "=== Recent Events ==="
    echo "$(date): Interface eth1 administratively down"
    echo "5 minutes ago: System boot completed"
    echo "10 minutes ago: All interfaces initialized"
    echo ""
    
    echo "Overall Status: HEALTHY"
}

# Main script logic
case "${1:-complete}" in
    system)
        show_system_status
        ;;
    interfaces)
        show_interface_status
        ;;
    network)
        show_network_status
        ;;
    complete|"")
        show_complete_status
        ;;
    -h|--help|help)
        show_usage
        ;;
    *)
        echo "ERROR: Unknown component '$1'"
        show_usage
        exit 1
        ;;
esac

# Return success exit code
exit 0 