# Mercury Server

This is the Mercury's chat server.

## Installing

If you have Go installed, navigate your terminal to this directory (``mercury-server``) then run:

```
$ go build
```

The executable ``mercury-server`` will be produced in the current working directory.

You can move the executable anywhere you please. To move it to the ``$PATH`` type:

```
$ sudo mv mercury-server /usr/local/bin
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
