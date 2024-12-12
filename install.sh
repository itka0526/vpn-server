#!/bin/sh

wget https://go.dev/dl/go1.21.2.freebsd-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" > $HOME/.profile
go version