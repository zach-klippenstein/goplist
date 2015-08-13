package plist

import "encoding/xml"

type plistString string

var stringStartElement = xmlElement("string")

func (s plistString) value() interface{} {
	return s
}

func (s plistString) marshal(e *xml.Encoder) error {
	return e.EncodeElement(s, stringStartElement)
}
