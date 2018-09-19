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
CONF_FILE=$DATA_DIR/config.toml

go build
mv mercury-server $BIN_DIR
mkdir $DATA_DIR
cp res/sample-config.toml $DATA_DIR

echo ""
echo "You can install mercury-server as a service."
echo "Doing so will launch the server on system-startup, and keep it alive if it fails."
echo "These systems also provide useful command-line tools to manage the service."
echo "Would you like to install mercury-server as a service?"
echo "1) Create a launchd service. (macOS)"
echo "2) Create a systemd service. (linux)"
echo "3) don't create a service"
read service
echo ""

case $service in
    "1")
        cp res/com.github.fabiocolacio.mercury-server.plist /Library/LaunchDaemons
        echo "launchd service created"
        echo "To start the service manually, use the following command:"
        echo "'sudo launchctl load /Library/LaunchDaemons/com.github.fabiocolacio.mercury-server.plist'\n"
        ;;
    "2")
        cp res/mercury-server.service /lib/systemd/system
        echo "systemd service created"
        echo "To start the service manually, use the following command:"
        echo "'sudo systemctl start mercury-server'\n"
        ;;
    *)
        echo "No service files were installed."
        echo "Use the command 'mercury-server' to start the server.\n"
        ;;
esac

echo "Before you can run the server, you need to create a configuration file for it."
echo "Do you want to interactively create one now? (y/n)"
read createconfig
echo ""

if [ "$createconfig" = "y" ]
then
    echo "What address and port should Mercury listen to HTTP connections on?"
    echo "To bind to all interfaces (recommended), use the address '0.0.0.0'"
    echo "Enter it in the format 'address:port'"
    read httpaddr
    echo "HttpAddr = \"$httpaddr\"" > $CONF_FILE
    echo ""

    echo "What address and port should Mercury listen to HTTPS connections on?"
    echo "To bind to all interfaces (recommended), use the address '0.0.0.0'"
    echo "Enter it in the format 'address:port'"
    read httpsaddr
    echo "HttpsAddr = \"$httpsaddr\"" >> $CONF_FILE
    echo ""

    echo "What is the absolute path to your SSL Certificate?"
    read certfile
    echo "CertFile = \"$certfile\"" >> $CONF_FILE
    echo ""

    echo "What is the absolute path to your Private Key?"
    read keyfile
    echo "KeyFile = \"$keyfile\"" >> $CONF_FILE
    echo ""

    echo "Configuration file written to '$CONF_FILE'"
    echo ""
else
    echo "Please place a config in the following location:"
    echo "$DATA_DIR\n"

    echo "Alternatively, you can pass a config file as a cli argument"
    echo ""
fi

echo "Installation Complete!"
