package main

import (
	"fmt"
	"log"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

const version = "0.1"

func main() {

	logf("mustago version: %s", version)

	document := js.Global.Get("document")

	param := docQuery("#parameters")
	input := docQuery("#input")
	output := docQuery("#output")
	logbox := docQuery("#log")

	logf("main: document=%v param=%v input=%v output=%v log=%v", document, param, input, output, logbox)
	logf("main: param=%q input=%q output=%q log=%q", param.TextContent(), input.TextContent(), output.TextContent(), logbox.TextContent())

	param.AddEventListener("change", false, listenerParam)
	input.AddEventListener("change", false, listenerInput)
}

func listenerParam(ev dom.Event) {
	t := ev.Target()
	p, ok := t.(*dom.HTMLTextAreaElement)
	if !ok {
		return
	}
	logf("listenerParam: param=%q", p.Value)
	updateOutput()
}

func listenerInput(ev dom.Event) {
	t := ev.Target()
	i, ok := t.(*dom.HTMLTextAreaElement)
	if !ok {
		return
	}
	logf("listenerInput: input=%q", i.Value)
	updateOutput()
}

func updateOutput() {
	param := docQuery("#parameters")
	input := docQuery("#input")
	output := docQuery("#output")

	p, isParam := param.(*dom.HTMLTextAreaElement)
	if !isParam {
		return
	}
	i, isInput := input.(*dom.HTMLTextAreaElement)
	if !isInput {
		return
	}
	o, isOutput := output.(*dom.HTMLTextAreaElement)
	if !isOutput {
		return
	}

	o.Value = fmt.Sprintf("updated: param=%q input=%q", p.Value, i.Value)
}

func docQuery(query string) dom.Element {
	return dom.GetWindow().Document().QuerySelector(query)
}

func logf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	log.Print(msg) // write to browser console

	logbox := docQuery("#log")
	box, ok := logbox.(*dom.HTMLTextAreaElement)
	if !ok {
		return // dom element not found
	}

	// write to dom element
	old := box.TextContent()
	if l := len(old); l > 0 && old[l-1] != '\n' {
		old += "\n"
	}
	box.SetTextContent(old + msg + "\n")
}
