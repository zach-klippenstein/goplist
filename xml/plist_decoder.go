package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// StartDecodingArray is returned by NextValue when an array
// start element is read.
type StartDecodingArray struct{}

// StartDecodingDict is returned by NextValue when an dict
// start element is read.
type StartDecodingDict struct{}

// EndDecodingContainer is returned by NextValue when an array
// or dict end element is read.
type EndDecodingContainer struct{}

// PlistDecoder parses XML plist data.
type PlistDecoder struct {
	xmlDecoder     *xml.Decoder
	currentDecoder containerDecoder
}

type containerDecoder interface {
	// NextValue reads the next value out of the plist.
	// Returns one of the plist scalar type (int, uint, float64, string, []byte, DictEntry),
	// or another containerDecoder. When the current container end tag is read,
	// returns (EndDecodingContainer{}, nil).
	NextValue() (interface{}, error)
	ParentDecoder() containerDecoder
}

// NewDecoder creates a decoder that reads a plist file from r.
func NewDecoder(r io.Reader) *PlistDecoder {
	return &PlistDecoder{
		xmlDecoder: xml.NewDecoder(r),
	}
}

/*
NextValue decodes the next value out of the plist.
Returns one of the plist scalar types (int, uint, float64, string, []byte, DictEntry),
or one of the container sentry types: StartDecodingArray, StartDecodingDict,
or EndDecodingContainer.

If StartDecodingArray is returned, the decoder encountered an <array> opening tag,
and will return the values of the array until a matching EndDecodingContainer is returned.

StartDecodingDict means the same thing for dictionaries. Dictionary entries are returned
as DictEntry values.
*/
func (d *PlistDecoder) NextValue() (interface{}, error) {
	var value interface{}

	// The first time NextValue() is called, we need to skip past
	// all the XML header stuff.
	decoder, err := d.consumeHeader()
	if err != nil {
		return nil, err
	}

	if decoder != nil {
		// This is the first time NextValue has been called.
		d.currentDecoder = decoder
		value = decoder
	} else {
		value, err = d.currentDecoder.NextValue()
		if err != nil {
			return nil, err
		}
	}

	// Handle container starts/ends.
	switch value := value.(type) {
	case *arrayDecoder:
		// Push the current decoder.
		d.currentDecoder = value
		return StartDecodingArray{}, nil
	case *dictDecoder:
		d.currentDecoder = value
		return StartDecodingDict{}, nil
	case EndDecodingContainer:
		// Pop the current decoder.
		d.currentDecoder = d.currentDecoder.ParentDecoder()
		return EndDecodingContainer{}, nil
	}

	return value, nil
}

// consumeHeader reads tokens until we find a <plist> or an error.
// It then reads the next container start tag, and returns the appropriate decoder for
// that type.
// Returns (nil, nil) if the header has already been read.
func (d *PlistDecoder) consumeHeader() (containerDecoder, error) {
	if d.currentDecoder != nil {
		return nil, nil
	}

	token, err := nextStartElement(d.xmlDecoder)
	if err != nil {
		return nil, err
	}

	if token.Name == plistStartElement.Name {
		return d.startFirstContainer()
	}

	return nil, fmt.Errorf("Expected %s, found %#v", plistStartElement, token)

}

// startFirstContainer reads until the next <array> or <dict> and returns
// the appopriate containerDecoder.
func (d *PlistDecoder) startFirstContainer() (containerDecoder, error) {
	token, err := nextStartElement(d.xmlDecoder)
	if err != nil {
		return nil, err
	}

	switch token.Name.Local {
	case arrayStartElement.Name.Local:
		return newArrayDecoder(nil, d.xmlDecoder), nil
	case dictStartElement.Name.Local:
		return newDictDecoder(nil, d.xmlDecoder), nil
	}

	return nil, fmt.Errorf("Expected container start element, found %#v", token)
}

func nextStartElement(xmlDecoder *xml.Decoder) (xml.StartElement, error) {
	for {
		token, err := nextInterestingToken(xmlDecoder)
		if err != nil {
			return xml.StartElement{}, err
		}

		if token, ok := token.(xml.StartElement); ok {
			return token, nil
		}
	}
}

// nextInterestingToken returns the next token from xmlDecoder that isn't
// a comment or blank char data.
func nextInterestingToken(xmlDecoder *xml.Decoder) (xml.Token, error) {
	for {
		token, err := xmlDecoder.Token()
		if err != nil {
			return nil, err
		}

		switch token := token.(type) {
		case xml.Comment:
			continue
		case xml.CharData:
			if len(strings.TrimSpace(string(token))) == 0 {
				continue
			}
		}

		return token, nil
	}
}
