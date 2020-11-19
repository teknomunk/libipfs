#!/bin/bash

go build -o libipfs.so -buildmode=c-shared *.go
cp libipfs.so /usr/local/lib/
