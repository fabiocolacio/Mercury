#!/bin/perl -w

use strict;
use Socket;

my $port = 8080;
my $proto = getprotobyname "tcp";

socket(SERVER, PF_INET, SOCK_STREAM, $proto)
    or die "socket: $!";

setsockopt(SERVER, SOL_SOCKET, SO_REUSEADDR, pack("l", 1))
    or die "setsockopt: $!";

bind(SERVER, sockaddr_in($port, INADDR_ANY))
    or die "bind: $!";

listen(SERVER, SOMAXCONN);

sub handle_request {
    
}

for (;;) {
    my $paddr = accept(CLIENT, SERVER);

    my $message = "";
    my $message_size = 0;
    my $chunk_size = 256;
    my $eof = 0;

    while ($message !~ m/\r\n/) {
        $message_size += read(CLIENT, $message, $chunk_size, $message_size);
    }

    print $message;

    print(CLIENT
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

    close(CLIENT);
}

