#!/bin/bash

#TODO check if binary and config exist
#TODO will have to write removal scripts
#TODO version stuff? probably bake into a makefile

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

DEB_LOC=$DIR/mykeystore
BIN_LOC=$DEB_LOC/usr/local/bin
CONFIG_LOC=$DEB_LOC/usr/local/etc/mykeystore
DEB_NAME=mykeystore-0.1.0.deb

rm -rf $DEB_LOC 2>1
rm $DIR/$DEB_NAME 2>1

mkdir -p $DEB_LOC/DEBIAN

mkdir -p $BIN_LOC
mkdir -p $CONFIG_LOC

cp -a $DIR/DEBIAN $DIR/etc $DEB_LOC

cd $DIR

cp ../../bin/keystore $BIN_LOC
cp ../../config.production.json $CONFIG_LOC/config.json
cp -a ../../assets $CONFIG_LOC/

#get rid of any swap files that may have been copied over
find $DEB_LOC -type f -name *.sw* -exec rm {} \;

dpkg-deb -b $DEB_LOC $DEB_NAME
