#!/bin/perl -w

use strict;
use threads;
use IO::Socket;
use IO::Socket::SSL;

#== Global Declarations ==#

$0 = "chatter-server";    # Set the program name
$| = 1;                   # Set flush-on-write on
$SIG{INT} = \&clean_exit; # Exit program on sig-int (Ctrl-C)
$SIG{CHLD} = 'IGNORE';    # Keep the process table clean, ignore dead child processes

my $http_port = 8080;     # Port for http connections
my $https_port = 8090;    # Port for https connections
my $host = "0.0.0.0";     # Bind to all devices
my $proto = "tcp";        # Use TCP protocol
my $backlog = 5;          # Size of queue for pending requests
my $http_server;          # Host HTTP server socket
my $https_server;         # Host HTTPS server socket

my $use_ssl = 0;          # Set to 0 to disable SSL (for debugging purposes only!)
my $ssl_cert_file;        # Path to an SSL certificate file
my $ssl_key_file;         # Path to an SSL key file

my $num_processes = 1;    # Total number of processes
my $max_processes = 200;  # Maximum number of processes to use

my $log_file = "$0.log";  # Name of log file that will be created
my $log_handle;           # File handle for the log file

#== Subroutines ==#

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
        $http_server->close;
        $https_server->close if ($use_ssl);

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

    $timestamp = localtime;
    log_message "$peerhost => connection closed";
}

#== Program Logic ==#

open $log_handle, ">", $log_file
    or warn "Failed to open file '$log_file': $!";

$http_server = new IO::Socket::INET(
    LocalHost => $host,
    LocalPort => $http_port,
    Proto     => $proto,
    Listen    => $backlog)
    or die "Failed to bind socket: $!";

$https_server = new IO::Socket::SSL(
    LocalAddr     => $host,
    LocalPort     => $https_port,
    Listen        => $backlog,
    SSL_cert_file => $ssl_cert_file,
    SSL_key_file  => $ssl_key_file)
    or die "Failed to bind SSL socket: $!"
    if ($use_ssl); 

log_message "$0 listening on port $http_port";

for (;;) {
    if ($num_processes < $max_processes) {
        my $client = $http_server->accept;

        my $pid = fork;
        if (defined $pid) {
            $num_processes += 1;

            if ($pid == 0) {
                handle_client $client;
                close $client;
                exit;
            }
        }
        else {
            log_message "Failed to create child process.";
        }
    }
    else {
        log_message "Unable to accept clients. Max process count reached.";
    }
}

