package xml

import "encoding/xml"

type arrayDecoder struct {
	baseDecoder
}

var _ containerDecoder = &arrayDecoder{}

func newArrayDecoder(parent containerDecoder, xmlDecoder *xml.Decoder) *arrayDecoder {
	return &arrayDecoder{baseDecoder{parent, xmlDecoder}}
}
