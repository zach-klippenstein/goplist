package xml

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"io"
	"time"
)

type DictEncodingFunc func(*DictEncoder) error
type ArrayEncodingFunc func(*ArrayEncoder) error

type baseEncoder struct {
	// Need to hang on to the underlying writer so we can control formatting.
	w io.Writer
	e *xml.Encoder

	// Set by StartArray or StartDict, and automatically closed when this encoder
	// is used again.
	encodingContainer bool

	// Set to true when writeEndTag() is called.
	finished bool
}

func newBaseEncoder(w io.Writer) *baseEncoder {
	e := xml.NewEncoder(w)
	e.Indent("", "\t")

	return &baseEncoder{w: w, e: e}
}

// copy returns a *baseEncoder with the same writer and xml encoder but
// fresh state flags.
func (e *baseEncoder) copy() *baseEncoder {
	return &baseEncoder{w: e.w, e: e.e}
}

// writeEndTag encodes an end element and marks the encoder as finished.
// Any subsequent operations will panic.
func (e *baseEncoder) writeEndTag(endElement xml.EndElement) error {
	if err := e.e.EncodeToken(endElement); err != nil {
		return err
	}
	e.finished = true
	return nil
}

/*
assertReady panics if the end tag has already been written or a call to
writeArray/writeDict has not returned.

It should be called before every exported operation. baseEncoder doesn't call
itself since DictEncoder needs to also check before writing the key.
*/
func (e *baseEncoder) assertReady() {
	if e.finished {
		panic("cannot write to encoder, end tag has already been written")
	}
	if e.encodingContainer {
		panic("cannot encode to parent container before closing child container")
	}
}

// writeArray locks this encoder and calls encode with an encoder that
// can be used to write array entries.
// The return value of encode is returned.
func (e *baseEncoder) writeArray(encode ArrayEncodingFunc) error {
	e.startContainer()
	defer e.endContainer()

	subEncoder, err := newArrayEncoder(e.copy())
	if err != nil {
		return err
	}
	if err = encode(subEncoder); err != nil {
		// Doesn't make sense to finish writing the container if
		// there was an error.
		return err
	}
	return subEncoder.writeEndTag()
}

// writeDict locks this encoder and calls encode with an encoder that
// can be used to write dictionary entries.
// The return value of encode is returned.
func (e *baseEncoder) writeDict(encode DictEncodingFunc) error {
	e.startContainer()
	defer e.endContainer()

	subEncoder, err := newDictEncoder(e.copy())
	if err != nil {
		return err
	}
	if err = encode(subEncoder); err != nil {
		// Doesn't make sense to finish writing the container if
		// there was an error.
		return err
	}
	return subEncoder.writeEndTag()
}

func (e *baseEncoder) startContainer() {
	e.encodingContainer = true
}

func (e *baseEncoder) endContainer() {
	e.encodingContainer = false
}
func writeString(e *xml.Encoder, val string) error {
	return e.EncodeElement(val, stringStartElement)
}

func writeBool(e *xml.Encoder, val bool) error {
	if val {
		// return encodeClosedElement(w, boolTrueElementName)
		return e.EncodeElement("", boolTrueElement)
	}
	// return encodeClosedElement(w, boolFalseElementName)
	return e.EncodeElement("", boolFalseElement)
}

func writeFloat(e *xml.Encoder, val float64) error {
	return e.EncodeElement(val, realStartElement)
}

func writeInt(e *xml.Encoder, val int64) error {
	return e.EncodeElement(val, integerStartElement)
}

func writeUint(e *xml.Encoder, val uint64) error {
	return e.EncodeElement(val, integerStartElement)
}

// writeDate encodes val as an ISO 8601/RFC 3339 date string.
func writeDate(e *xml.Encoder, val time.Time) error {
	encodedDate := val.Format(dateFormat)
	return e.EncodeElement(encodedDate, dateStartElement)
}

// writeData base64-encodes val.
func writeData(e *xml.Encoder, val []byte) error {
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	// Just writing to a buffer, can't fail.
	encoder.Write(val)
	encoder.Close()

	return e.EncodeElement(encoded.String(), dataStartElement)
}
