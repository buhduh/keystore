#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

go-bindata data/
go install
keystore --assets $DIR/assets
