#!/bin/sh

[ -z "$GOPATH" ] && export GOPATH=$HOME/go

echo GOPATH=$GOPATH

go get github.com/gopherjs/gopherjs
go get honnef.co/go/js/dom

gofmt -s -w main.go
go tool fix main.go
go tool vet .

[ -x $GOPATH/bin/gosimple ] && $GOPATH/bin/gosimple main.go
[ -x $GOPATH/bin/golint ] && $GOPATH/bin/golint main.go
[ -x $GOPATH/bin/staticcheck ] && $GOPATH/bin/staticcheck main.go

gopherjs install -v

ls $GOPATH/bin/mustago.js* | while read i; do
    cmd="cp $i ."
    echo $cmd
    $cmd
done
