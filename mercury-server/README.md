# Mercury Server

This is the Mercury's chat server.

## Installing

First follow the instructions [in the wiki](https://github.com/fabiocolacio/Mercury/wiki/Setting-Up-The-Build-Environment) to setup your build environment.

Next, navigate to the ``mercury-server`` directory, and install the server:

```
$ cd $GOPATH/src/github.com/fabiocolacio/mercury/mercury-server
$ go install
```

The source will be compiled, and the executable will be placed in ``$GOPATH/bin``. If you followed the instructions in the wiki, the binary will already be in your ``$PATH``. You are free to move the executable as you please.

If you plan to run this command with ``sudo``, you may want to move the executable to the system's ``secure_path`` (you will have to do this if you plan to run the server on restricted ports 80 and 443):

```
$ sudo mv $GOPATH/bin/mercury-server /usr/local/bin
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
