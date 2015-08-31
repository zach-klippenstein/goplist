package xml

import (
	"encoding/xml"
	"fmt"
)

// DictEntry is returned from NextValue() when parsing a dictionary.
type DictEntry struct {
	Key   string
	Value interface{}
}

type dictDecoder struct {
	baseDecoder
}

var _ containerDecoder = &dictDecoder{}

func newDictDecoder(parent containerDecoder, xmlDecoder *xml.Decoder) *dictDecoder {
	return &dictDecoder{baseDecoder{parent, xmlDecoder}}
}

func (d *dictDecoder) NextValue() (interface{}, error) {
	for {
		token, err := nextInterestingToken(d.xmlDecoder)
		if err != nil {
			return nil, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name == dictKeyElement.Name {
				return d.finishReadingEntry()
			}
		case xml.EndElement:
			if token.Name == dictStartElement.End().Name {
				return EndDecodingContainer{}, nil
			}
		}

		return nil, fmt.Errorf("Expected %s or %s, found %#v",
			dictKeyElement.Name.Local, dictStartElement.End().Name.Local, token)
	}
}

func (d *dictDecoder) finishReadingEntry() (interface{}, error) {
	key, err := finishReadingKey(d.xmlDecoder)
	if err != nil {
		return nil, err
	}

	value, err := d.baseDecoder.NextValue()
	if err != nil {
		return nil, err
	}

	return DictEntry{key, value}, nil
}

func finishReadingKey(xmlDecoder *xml.Decoder) (string, error) {
	value, err := readCharDataUntilEnd(xmlDecoder, dictKeyElement.End())
	if err != nil {
		return "", err
	}
	return value, nil
}
