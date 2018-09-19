# Mercury Server

This is the Mercury's chat server.

## Installing

First, [install golang](https://golang.org/dl/).

Then, you can compile and install the server binary with:

```sh
# Set $GOPATH, where the source and executable will be stored.
# This can be substituted with any directory of your choosing.
export GOPATH=~/go && mkdir $GOPATH

# Download and compile the mercury-server source code into your $GOPATH.
# The binary can be found in $GOPATH/bin, and can be moved as you see fit.
go get github.com/fabiocolacio/mercury/mercury-server

# Optionally, you can run the install script, which affords you a few benefits:
# - Installs the binary to the secure_path in /usr/local/bin
# - Optionally creates a systemd/launchd service
# - Interactively creates a config file for you
cd $GOPATH/src/github.com/fabiocolacio/mercury/mercury-server
sudo sh install.sh
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

You should save your configuration file to ``/usr/local/share/com.github.fabiocolacio.mercury-server/config.toml``.

A sample configuration file, ``sample-server.toml``, can be found in this directory.

## Usage

If you placed your ``config.toml`` file in ``/usr/local/share/com.github.fabiocolacio.mercury-server/``, you can run the server by issuing the command:

```
$ mercury-server
```

To run the server with a different configuration file, run:

```
$ mercury-server ~/path/to/configuration/file/config.toml
```
