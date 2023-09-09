#!/bin/bash
#
# Build and deploy the sensor application to a TCG4.
#
# This script requires a working SSH connection to the TCG4. The S50sensor init.d script must be
# set up correctly to remotely stop and start the service. Go, ssh, and sshpass need to be
# installed on this system to run this script.
#
# Warning: keep your passwords safe. Don't store sensitive (production) passwords in this file. Use
# a temporary development password if needed.

# Define the remote server details
remote_user="user"
remote_host="192.168.222.80"
remote_password="useruser" 

echo "Deploying to $remote_user@$remote_host..."

# Go to project root
cd "$(dirname $(realpath $0))/../"
mkdir -p build

# Set environment variables
export GOOS=linux
export GOARCH=arm

# Build for TCG4
go build -o build/sensor -ldflags "-X main.AppVersion=dev+$(date +%Y%m%d%H%M%S)" ./cmd/sensor

if [ $? != 0 ]
then
    echo "Error: build failed"
    exit 1
fi

# Deploy
sshpass -p "$remote_password" scp -r -o StrictHostKeyChecking=no build/sensor user@192.168.222.80:/home/user &> /dev/null

if [ $? != 0 ]
then
    echo "Error: copy failed"
    exit 1
fi

echo "$remote_password" | sshpass -p "$remote_password"  ssh -o StrictHostKeyChecking=no -tt "$remote_user@$remote_host" "sudo /etc/init.d/S50sensor stop && sudo mv /home/user/sensor /usr/local/bin && sudo chmod +x /usr/local/bin/sensor && sudo chown root:root /usr/local/bin/sensor && sudo /etc/init.d/S50sensor start" &> /dev/null

if [ $? != 0 ]
then
    echo "Error: configuration failed"
    exit 1
fi

echo "Build and deploy OK!"
exit 0