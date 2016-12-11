Keystore:
=============
A password store intended to be run on a Rasberry PI.
Uses MFA with google authenticator on a smart phone for
verification.

Dependencies:
-------------
1. Mysql
2. go-bindata (go get -u github.com/jteeuwen/go-bindata/...)

Installation:
-------------
1.  Clone repo to $GOPATH/src/keystore.  $git clone git@github.com:buhduh/keystore.git $GOPATH/src/keystore
2.  cd $GOPATH/src/keystore && go get -u
3.  go-bindata (go get -u github.com/jteeuwen/go-bindata/...)
4.  Install mysql and make an empty database called "keystore" on the root account
with no password.  $cd database && mysql -u root keystore < create_database.sql
5.  Install go-bindata from https://github.com/jteeuwen/go-bindata
6.  run build.go

If everything is set up, the keystore can be navigated to by going to http://localhost:8080 in the browser

Plans:
-------------
1.  More intelligent mysql account creation, detection so root with no password isn't used.
2.  makefile for builds, tests, and deployments
