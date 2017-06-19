package main

import (
	//"fmt"
	"log"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

const version = "0.0"

func main() {

	logf("mustago version: %s", version)

	document := js.Global.Get("document")

	param := docQuery("#parameters")
	input := docQuery("#input")
	output := docQuery("#output")
	logbox := docQuery("#log")

	logf("main: document=%v param=%v input=%v output=%v log=%v", document, param, input, output, logbox)
}

func docQuery(query string) dom.Element {
	return dom.GetWindow().Document().QuerySelector(query)
}

func logf(format string, a ...interface{}) {
	log.Printf(format, a...)
}
