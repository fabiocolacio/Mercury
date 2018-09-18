# Mercury Server

This is the Mercury's chat server.

## Installing

First, [install golang](https://golang.org/dl/).

Then, you can compile and install the server binary with:

```sh
# Set $GOPATH, where the source and executable will be stored
mkdir ~/go && export GOPATH=~/go

# Download and compile the mercury-server source code into your $GOPATH
go get github.com/fabiocolacio/mercury/mercury-server

# Move the executable into the secure_path, so you can more easily
# run the server with sudo (useful for binding to restricted ports 80 and 443).
sudo mv $GOPATH/bin/mercury-server /usr/local/bin
```

## Configuration

In order to run the server, you must first write a configuration file in the ``toml`` format, specifying the following details:

* HttpAddr
  * The address and port to bind the HTTP server to.
  * HTTP requests are simply redirected to the HTTPS server.
* HttpsPort
  * The address and port to bind the HTTPS server to.
* CertFile
  * The absolute path to your server's certificate.
  * For information about acquiring a certificate, see [the wiki](https://github.com/fabiocolacio/Mercury/wiki/Acquiring-an-SSL-Certificate)
* KeyFile
  * The absolute path to your server's private key.
  * For information about acquiring a key, see [the wiki](https://github.com/fabiocolacio/Mercury/wiki/Acquiring-an-SSL-Certificate).

You should save your configuration file to ``~/.config/mercury/server.toml``.

A sample configuration file, ``sample-server.toml``, can be found in this directory.

## Usage

If you placed your ``server.toml`` file in ``~/.config/mercury``, you can run the server by issuing the command:

```
$ mercury-server
```

To run the server with a different configuration file, run:

```
$ mercury-server ~/path/to/configuration/file/server.toml
```
