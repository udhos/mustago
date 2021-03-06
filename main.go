package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	//"github.com/cbroglie/mustache"
	"github.com/gopherjs/gopherjs/compiler"
	"github.com/gopherjs/gopherjs/js"
	"github.com/udhos/mustache"
	"gopkg.in/yaml.v2"
	"honnef.co/go/js/dom"
)

const version = "0.9"

func main() {

	param := docQuery("#parameters")
	input := docQuery("#input")
	logbox := docQuery("#log")
	button := docQuery("#button-output-copy")
	checkEscaping := docQuery("#escaping")

	if box, ok := logbox.(*dom.HTMLTextAreaElement); ok {
		box.Value = "" // clear log
	}

	logf("main: mustago %s, go %s, gopherjs %s", version, runtime.Version(), compiler.Version)

	ver := docQuery("#version")
	if v, ok := ver.(*dom.HTMLSpanElement); ok {
		v.SetInnerHTML("mustago " + version + ` - simple golang <a href="http://mustache.github.io">mustache</a> template evaluation in html`)
	}

	param.AddEventListener("change", false, listenerParam)
	input.AddEventListener("change", false, listenerInput)
	param.AddEventListener("keyup", false, listenerParam)
	input.AddEventListener("keyup", false, listenerInput)
	button.AddEventListener("click", false, buttonOutputCopy)
	checkEscaping.AddEventListener("change", false, toggleEscaping)

	updateOutput() // the first update
}

func toggleEscaping(ev dom.Event) {
	escaping := docQuery("#escaping")
	i, isInput := escaping.(*dom.HTMLInputElement)
	if !isInput {
		return
	}
	logf("toggleEscaping: %v checked=%v", i, i.Checked)
	updateOutput()
}

func buttonOutputCopy(ev dom.Event) {

	output := docQuery("#output")

	o, isTextArea := output.(*dom.HTMLTextAreaElement)
	if !isTextArea {
		return
	}

	logf("buttonOutputCopy: %v", o)

	o.Select()                                            // select text area
	js.Global.Get("document").Call("execCommand", "copy") // document.execCommand('copy');
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

	escaping := docQuery("#escaping")
	e, isInput := escaping.(*dom.HTMLInputElement)
	if !isInput {
		return
	}
	disableHTMLEscape := !e.Checked

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
	result, errRender = mustache.RenderRaw(i.Value, disableHTMLEscape, doc)
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
