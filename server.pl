#!/bin/perl -w

use strict;

use IO::Socket;
use IO::Socket::SSL;

my $port = shift || 8080;
my $host = "0.0.0.0";
my $proto = "tcp";
my $backlog = 5;

my $server = new IO::Socket::INET(
    LocalHost => $host,
    LocalPort => $port,
    Proto     => $proto,
    Listen    => $backlog
) or die "Failed to bind socket: $!";

for (;;) {
    my $client = $server->accept;

    my $message = "";
    my $message_size = 0;
    my $chunk_size = 256;
    my $eof = 0;

    while ($message !~ m/\r\n/) {
        $message_size += read($client, $message, $chunk_size, $message_size);
    }

    print $message;

    print($client
        "HTTP/1.1 200 OK\r\n\r\n
        <!doctype html>
        <html>
            <head>
                <title>hi</title>
            </head>
            <body>
                <h1>Hello World!</h1>
            </body>
        </html>");

    close($client);
}

$server->close;

