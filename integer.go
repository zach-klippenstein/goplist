package plist

import "encoding/xml"

type plistInt struct {
	signed   int64
	unsigned uint64
}

var integerStartElem = xmlElement("integer")

func (i plistInt) value() interface{} {
	if i.signed != 0 {
		return i.signed
	}
	return i.unsigned
}

func (i plistInt) marshal(e *xml.Encoder) error {
	if i.signed != 0 {
		return e.EncodeElement(i.signed, integerStartElem)
	}
	return e.EncodeElement(i.unsigned, integerStartElem)
}
