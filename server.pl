#!/bin/perl -w

use strict;

use threads;

use IO::Socket;
use IO::Socket::SSL;

$0 = "chatter-server";
$| = 1;

$SIG{INT} = \&clean_exit;

my $port = shift || 8080;
my $host = "0.0.0.0";
my $proto = "tcp";
my $backlog = 5;

my $max_thread_count = shift || 200;

my $log_file = "$0.log";
my $log_handle;

my $active :shared = 1;

my $server;

sub log_message {
    if (defined($log_handle)) {
        my $message = shift;
        my $timestamp = localtime;
        print $log_handle "[$timestamp] $message\n";
    }
}

sub clean_exit {
    my $tid = threads->tid;
    if ($tid == 0) {
        my $timestamp = localtime;

        log_message "Shutting down server..";
        $active = 0;
        $server->close;

        my @threads = threads->list(threads::all);
        if (scalar @threads > 0) {
            $timestamp = localtime;
            log_message "Fulfilling unresolved requests...";
            foreach my $thread (@threads) {
                $thread->join;
            }
        }

        close $log_handle;

        exit 0;
    }
}

sub handle_client {
    my $client = shift;
    
    my $timestamp = localtime;
    my $peerhost = $client->peerhost;
    
    log_message "$peerhost => connection established";

    my $message = "";
    my $message_size = 0;
    my $chunk_size = 256;
    my $eof = 0;

    while ($message !~ m/\r\n/) {
        $message_size += read($client, $message, $chunk_size, $message_size);
    }

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

    $timestamp = localtime;
    log_message "$peerhost => connection closed";
}

open $log_handle, ">", $log_file
    or warn "Failed to open file '$log_file': $!";

$server = new IO::Socket::INET(
    LocalHost => $host,
    LocalPort => $port,
    Proto     => $proto,
    Listen    => $backlog)
    or die "Failed to bind socket: $!";

log_message "$0 listening on port $port";

for (;;) {
    my $thread_count = threads->list(threads::running);
    if ($active && $thread_count < $max_thread_count) {
        my $client = $server->accept;
        threads->new(\&handle_client, $client);
    }
}

