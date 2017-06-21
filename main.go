package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/cbroglie/mustache"
	"gopkg.in/yaml.v2"
	"honnef.co/go/js/dom"
)

const version = "0.3"

func main() {

	param := docQuery("#parameters")
	input := docQuery("#input")
	logbox := docQuery("#log")

	if box, ok := logbox.(*dom.HTMLTextAreaElement); ok {
		box.Value = "" // clear log
	}

	logf("main: mustago version: %s", version)

	ver := docQuery("#version")
	if v, ok := ver.(*dom.HTMLSpanElement); ok {
		v.SetInnerHTML("mustache " + version)
	}

	param.AddEventListener("change", false, listenerParam)
	input.AddEventListener("change", false, listenerInput)
	param.AddEventListener("keyup", false, listenerParam)
	input.AddEventListener("keyup", false, listenerInput)
}

func listenerParam(ev dom.Event) {
	t := ev.Target()
	p, ok := t.(*dom.HTMLTextAreaElement)
	if !ok {
		return
	}
	logf("listenerParam: event=%v param=%q", ev, p.Value)
	updateOutput()
}

func listenerInput(ev dom.Event) {
	t := ev.Target()
	i, ok := t.(*dom.HTMLTextAreaElement)
	if !ok {
		return
	}
	logf("listenerInput: event=%v input=%q", ev, i.Value)
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

	var result string

	buf := []byte(p.Value)
	var doc interface{}
	errParse := yaml.Unmarshal(buf, &doc)
	if errParse != nil {
		msg := fmt.Errorf("updateOutput: YAML parse error; %v", errParse)
		logf("%s", msg)
		result = msg.Error()
		setOutput(o, p.Value, i.Value, result)
		return
	}

	var errRender error
	result, errRender = mustache.Render(i.Value, doc)
	if errRender != nil {
		msg := fmt.Errorf("updateOutput: mustache render error; %v", errRender)
		logf("%s", msg)
		result = msg.Error()
	}

	setOutput(o, p.Value, i.Value, result)
}

func setOutput(output *dom.HTMLTextAreaElement, param, input, result string) {
	//output.Value = fmt.Sprintf("updated: param=%q input=%q\noutput:\n%s", param, input, result)
	output.Value = result
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
	old := box.Value
	if l := len(old); l > 0 && old[l-1] != '\n' {
		old += "\n"
	}
	box.Value = old + msg + "\n"

	max := 10
	lines := strings.Split(box.Value, "\n")
	if len(lines) <= max {
		return
	}
	tail := lines[len(lines)-max-1:]
	box.Value = strings.Join(tail, "\n")
}
