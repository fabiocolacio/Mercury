# Mercury Chat

Mercury is my end-to-end encrypted chat protocol.
This repository is currently the home of the server-side code.
The prototype client can currently be found [here](https://github.com/fabiocolacio/quicksilver).

## Protocol

The current protocol is relatively simple.
It was inspired by the *Double Ratchet Algorithm* implemented by Open Whisper Systems for their app *Signal*.

Each Alice would like to initiate a conversation with a new peer, they must both create elliptic-curve Diffie-Hellman parameters, and exchange their public parameters.
Every time Alice sends a message to a peer (such as Bob), she creates a shared secret using her private parameters and Bob's public parameters.
She encrypts the message using this shared key, and sends it to Bob.
Each time she sends a message, she also creates new Diffie-Hellman parameters, and sends the public component alongside the message.

No two messages are ever encrypted with the same secret, because each time a message is received, the Diffie-Hellman parameters are updated and the shared-secret changes.

This is the structure of a message:

```
Sid: The id of the sender's DH parameters used for this message
Rid: The id of the receiver's DH parameters used for this message
 IV: The initialization vector used for AES-CBC encryption of encrypted fields
Nxt: The sender's newly created public DH parameter for use on the next message (encrypted)
Msg: The encrypted message
Key: The encrypted HMAC key
Tag: HMAC integrity tag - HMAC(Nxt || Msg, Decrypted HMAC key))
```

## Future Plans

* Group chat
* Voice/video support
* Blockchain-based trustless keyserver integration (removes need to manually coordinate initial key-exchange)
* Docker image

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

When running mercury for the first time, use the ``--init`` flag.
This will create all necessary tables in the database.
Use this flag with caution!
If the tables already exist, *it will delete and reset them*.

```
mercury --init
```

mercury looks for a configuration file in ``~/.config/mercury/config.toml`` by default.
You may specify another config file with the ``--config`` flag.

```
mercury --config ~/path/to/configuration/file/config.toml
```
