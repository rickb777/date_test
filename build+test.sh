#!/bin/bash -e
cd $(dirname $0)
PATH=$HOME/go/bin:$PATH
unset GOPATH

if ! type -p goveralls; then
  echo go get github.com/mattn/goveralls
  go get github.com/mattn/goveralls
  echo go install github.com/mattn/goveralls
  go install github.com/mattn/goveralls
fi

echo Sqlite test...
go test -v -covermode=count -coverprofile=date.out .
go tool cover -func=date.out
[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=date.out -service=travis-ci -repotoken $COVERALLS_TOKEN

if [ "$1" = "mysql" ]; then
  shift
  echo
  echo Mysql test...
  export GO_DRIVER='mysql'
  export GO_DSN='testuser:TestPasswd.9.9.9@/test'
  go test -v .
fi

if [ "$1" = "postgres" ]; then
  shift
  echo
  echo Postgres test...
  export GO_DRIVER='postgres'
  export GO_DSN='postgres://testuser:TestPasswd.9.9.9@/test'
  go test -v .
fi
