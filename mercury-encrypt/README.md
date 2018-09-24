# Mercury Encrypter

This is a cli program that uses Mercury's encryption algorithm to encrypt and decrypt messages.

## Installation

First, [install golang](https://golang.org/dl/).

Then, you can compile and install the binary with:

```sh
# Set $GOPATH, where the source and executable will be stored.
# This can be substituted with any directory of your choosing.
export GOPATH=~/go && mkdir $GOPATH && export PATH=$PATH:$GOPATH/bin

# Download and compile the mercury-encrypt source code into your $GOPATH.
# The binary can be found in $GOPATH/bin, and can be moved as you see fit.
go get github.com/fabiocolacio/mercury/mercury-encrypt
```

## Usage

```sh
# Encrypt a message (outputs to stdout)
mercury-encrypt -m "message to encrypt" -k "/path/to/public/key.pem"

# Decrypt a message (outputs to stdout) with the -d flag
mercury-encrypt -m "json-encoded encrypted message" -k "/path/to/private/key.pem" -d
```
