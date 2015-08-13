package plist

import "encoding/xml"

type plistReal float64

var realStartElement = xmlElement("real")

func (r plistReal) value() interface{} {
	return r
}

func (r plistReal) marshal(e *xml.Encoder) error {
	return e.EncodeElement(r, realStartElement)
}
