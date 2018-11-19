# Mercury Server

This is the Mercury chat server.

## Installing

After installing golang and MySQL, run the installer with this command:

```
curl https://raw.githubusercontent.com/fabiocolacio/Mercury/master/install.sh -sSf | sh
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
* LogFile (Optional)
  * The absolute path to a log file to maintain
  * If this is not set, the server logs to stdout, which can be redirected as you please
  * The file will be appended to if it already exists
* SQLUser
  * The user which will perform operations on the SQL database
* SQLPass
  * The password for SQLUser
* SQLDb
  * The SQL database that mercury will use for its tables

You should save your configuration file to ``/usr/local/share/com.github.fabiocolacio.mercury-server/config.toml``.

A sample configuration file, ``sample-server.toml``, can be found in the ``res`` directory.

## Usage


Run mercury with the following command:

```
$ mercury -c ~/path/to/configuration/file/config.toml
```
