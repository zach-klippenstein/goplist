/*
Package xml implements Apple's XML plist format.

A plist file is encoded as a top-level array or dictionary, to which
primitives and other arrays/dictionaries can be written.

Eventually I might get around to adding support for encoding Go container types:

	xml.Write(os.Stdout, map[string]interface{}{
		"name": "Bilbo Baggins",
		"age": 111,
		"acquaintances": []string{
			"Gandalf the Grey",
			"Frodo Baggins",
			"Samwise Gamgee",
		}})

and

	xml.Write(os.Stdout, struct {
			name          string
			age           uint
			acquaintances []string
		} {
			name:          "Bilbo Baggins",
			age:           111,
			acquaintances: []string{
				"Gandalf the Grey",
				"Frodo Baggins",
				"Samwise Gamgee",
			}})

More Information

https://developer.apple.com/library/mac/documentation/Darwin/Reference/ManPages/man5/plist.5.html#//apple_ref/doc/man/5/plist

https://developer.apple.com/library/mac/documentation/CoreFoundation/Conceptual/CFPropertyLists/Articles/XMLTags.html#//apple_ref/doc/uid/20001172-CJBEJBHH

https://en.wikipedia.org/wiki/Property_list
*/
package xml

import (
	"encoding/xml"
	"strings"
	"time"
)

var indentation = strings.Repeat(" ", 4)

var procInst = xml.ProcInst{
	Target: "xml",
	Inst:   []byte(`version="1.0" encoding="UTF-8"`),
}

var doctype = xml.Directive([]byte(
	`DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"`))

var plistStartElement = xml.StartElement{
	Name: xml.Name{"", "plist"},
	Attr: []xml.Attr{xml.Attr{xml.Name{"", "version"}, "1.0"}},
}

var stringStartElement = xmlElement("string")
var realStartElement = xmlElement("real")
var boolTrueElement = xmlElement("true")   // var boolTrueElementName = "true"
var boolFalseElement = xmlElement("false") // var boolFalseElementName = "false"
var integerStartElement = xmlElement("integer")
var dateStartElement = xmlElement("date")
var dataStartElement = xmlElement("data")

const dateFormat = time.RFC3339

func xmlElement(name string) xml.StartElement {
	return xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: name,
		},
	}
}
