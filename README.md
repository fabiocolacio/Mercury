# Mercury

An end-to-end encrypted chat service.

## I'm lost... help!

* Fret not fellow homosapien!
* You can test the client-side encryption/decryption algorithm with the encryption-tester [here](https://github.com/fabiocolacio/Mercury/tree/master/mercury-encrypt#mercury-encrypter).
* Instructions for setting up the server are [here](https://github.com/fabiocolacio/Mercury/tree/master/mercury-server#mercury-server).
* The client is coming soon.
* Other tidbits can be found at the [wiki](https://github.com/fabiocolacio/Mercury/wiki).

## How to Navigate This Repository

```
root-folder: The Mercury Library
|
--> mercury-encrypter: The code for the encryption-tester executable
|
--> mercury-server: The code for the mercury-server executable
|
--> mercury-client: The code for the mercury-client executable
```

Each of the sub-directories have their own ``README.md`` files with information about the setup and usage of that particular program.

If you have [golang](https://golang.org/dl/) installed, you can install the individual binaries as you please to ``$GOPATH/bin`` with ``go get``:

```sh
# Clone the entire repository to $GOPATH/src
go get github.com/fabiocolacio/mercury

# OR Install the server binary to $GOPATH/bin
go get github.com/fabiocolacio/mercury/mercury-server

# OR Install the client binary to $GOPATH/bin
go get github.com/fabiocolacio/mercury/mercury-client

# OR Install the encryption-tester binary to $GOPATH/bin
go get github.com/fabiocolacio/mercury/mercury-encrypt
```

## Whats with the name?

This program is named after the [Roman God](https://en.wikipedia.org/wiki/Mercury_(mythology)) of messages and communication!

Evidently, he is also the patron god of thieves, so perhaps it's not the best name for a secure, end-to-end encrypted messaging platform...
