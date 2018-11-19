#!/bin/bash

# This script installs Mercury with it's default
# configuration. The default configuration doesn't
# require admin priveleges to install or run, as it
# will install to ~/.go/bin, and bind to ports 8080
# and 9090.

#== DEFAULT CONFIGURATION ==#
CONF_DIR=$HOME/.config/mercury
CONF_FILE=$CONF_DIR/config.toml
LOG_FILE=$CONF_DIR/log.txt
KEY_FILE=$CONF_DIR/privkey.pem
CERT_FILE=$CONF_DIR/cert.pem
CERT_EXP=365
SQL_USER="root"
SQL_PASS=""
SQL_DB="mercury"
HTTP_ADDR=":8080"
HTTPS_ADDR=":9090"

#== INSTALL BINARY ==#
if [ -z "$(command -v go)" ]; then
    echo "go tools are not installed!"
    echo "Please install golang at https://golang.org/dl and try again."
    exit 1
fi

if [ -z "$GOPATH" ]; then
    echo 'export GOPATH=$HOME/.go'
    echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
    source ~/.bashrc
fi

echo "Installing mercury to $GOPATH/bin"
go get github.com/fabiocolacio/mercury

#== INSTALL CONFIG FILES ==#
if [ -z "$(command -v openssl)" ]; then
    echo "openssl tools are not installed!"
    echo "Please install openssl and try again."
    exit 1
fi

mkdir -p $CONF_DIR

if [ ! -e $KEY_FILE ] || [ ! -e $CERT_FILE ]; then
    echo "Creating private key as $KEY_FILE"
    openssl genrsa -out $KEY_FILE 2048

    echo "Creating self-signed certificate as $CERT_FILE"
    openssl req -new -x509 -sha256 -key $KEY_FILE -out $CERT_FILE -days $CERT_EXP
fi

if [ ! -e $LOG_FILE ]; then
    echo "Creating log file as $LOG_FILE"
    touch $LOG_FILE
fi

if [ ! -e $CONF_FILE ]; then
    echo "Installing config as $CONF_FILE"
    echo "HttpAddr = \"$HTTP_ADDR\"" >> $CONF_FILE
    echo "HttpsAddr = \"$HTTPS_ADDR\"" >> $CONF_FILE
    echo "CertFile = \"$CERT_FILE\"" >> $CONF_FILE
    echo "KeyFile = \"$KEY_FILE\"" >> $CONF_FILE
    echo "LogFile = \"$LOG_FILE\"" >> $CONF_FILE
    echo "SQLUser = \"$SQL_USER\"" >> $CONF_FILE
    echo "SQLPass = \"$SQL_PASS\"" >> $CONF_FILE
    echo "SQLDb = \"$SQL_DB\"" >> $CONF_FILE
fi

echo "Installation complete."

