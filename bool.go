package plist

import "encoding/xml"

type plistBool bool

var (
	boolTrueElement  = xmlElement("true")
	boolFalseElement = xmlElement("false")
)

func (b plistBool) value() interface{} {
	return b
}

func (b plistBool) marshal(e *xml.Encoder) error {
	// TODO this probably won't work, since it won't auto-close.
	if b {
		return e.EncodeElement("", boolTrueElement)
	} else {
		return e.EncodeElement("", boolFalseElement)
	}
}
