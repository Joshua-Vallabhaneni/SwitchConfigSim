#!/bin/bash

# simulates hostnamectl set-hostname command

# set_hostname.sh - Sets system hostname
# Function to display usage information
show_usage() {
    echo "Usage: $0 <new_hostname>"
    echo "Example: $0 core-switch-01"
    echo ""
    echo "This script sets the hostname on an NVIDIA switch"
    echo "simulates: hostnamectl set-hostname <hostname>"
}

# Check if hostname arg is provided
if [ $# -eq 0 ]; then
    echo "ERROR: No hostname provided"
    show_usage
    exit 1
fi

# Get the new hostname from CLI
NEW_HOSTNAME="$1"

# Validate hostname format
if [[ ! "$NEW_HOSTNAME" =~ ^[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]$ ]] && [[ ${#NEW_HOSTNAME} -gt 1 ]]; then
    # Allow single character hostnames but not ones that start/end with hyphen
    if [[ ! "$NEW_HOSTNAME" =~ ^[a-zA-Z0-9]$ ]]; then
        echo "ERROR: Invalid hostname format. Hostname must contain only letters, numbers, and hyphens"
        echo "ERROR: Cannot start or end with hyphen"
        exit 1
    fi
fi

# Check hostname length (linx limit -> 64 chars)
if [ ${#NEW_HOSTNAME} -gt 64 ]; then
    echo "ERROR: Hostname too long. Maximum length is 64 characters"
    exit 1
fi

# Hostname change operation
echo "SwitchConfigSim: Setting hostname to '$NEW_HOSTNAME'"
echo "Simulating: hostnamectl set-hostname $NEW_HOSTNAME"

# Successful completion
echo "Hostname change completed successfully"
echo "New hostname: $NEW_HOSTNAME"

# Return success exit code
exit 0 