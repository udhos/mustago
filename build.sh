#!/bin/sh

[ -z "$GOPATH" ] && export GOPATH=$HOME/go

echo GOPATH=$GOPATH

go get -insecure honnef.co/go/tools/cmd/unused
go get -insecure honnef.co/go/tools/cmd/gosimple
go get -insecure honnef.co/go/tools/cmd/staticcheck
go get github.com/golang/lint/golint

go get github.com/gopherjs/gopherjs
#go get github.com/cbroglie/mustache
go get github.com/udhos/mustache
go get -insecure honnef.co/go/js/dom
go get gopkg.in/yaml.v2

gofmt -s -w main.go
go tool fix main.go
go tool vet .

hash gosimple && gosimple main.go
hash golint && golint main.go
hash staticcheck && staticcheck main.go
hash unused && unused main.go

gopherjs install -v

ls $GOPATH/bin/mustago.js* | while read i; do
    cmd="cp $i ."
    echo $cmd
    $cmd
done
