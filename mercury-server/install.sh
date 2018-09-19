#!/bin/bash

if [ "$EUID" -ne 0 ]; then
    echo "Please run this script as root!"
    exit
fi

if ! [ -x "$(command -v go)" ]; then
    echo "You must install golang to compile this program!"
    echo "You can get it at https://golang.org/dl/"
    exit
fi

BIN_DIR=/usr/local/bin
DATA_DIR=/usr/local/share/com.github.fabiocolacio.mercury-server

go build
mv mercury-server $BIN_DIR
mkdir $DATA_DIR
cp res/sample-config.toml $DATA_DIR

echo "Install a service file?"
echo "1) launchd (macOS)"
echo "2) systemd (linux)"
echo "3) don't create a service"
read service

case $service
    1)
        cp res/com.github.fabiocolacio.mercury-server.plist /Library/LaunchDaemons
        echo "launchd service created"
        echo "Use the command 'sudo launchctl load /Library/LaunchDaemons/com.github.fabiocolacio.mercury-server.plist' to start the service."
        ;;
    2)
        cp res/mercury-server.service /lib/systemd/system
        echo "systemd service created"
        echo "Use the command 'sudo systemctl start mercury-server' to start the service."
        ;;
    *)
        echo "No service files were installed."
        echo "Use to command 'mercury-server' to start the server."
        ;;
esac

echo "Installation Complete!"
