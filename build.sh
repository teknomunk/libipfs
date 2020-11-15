#!/bin/bash

go build -o libipfs.so -buildmode=c-shared libipfs.go
cp libipfs.so /usr/local/lib/
