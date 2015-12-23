package xml

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"time"
)

const plistFieldTagName = "plist"

type DictEncodingFunc func(*DictEncoder) error
type ArrayEncodingFunc func(*ArrayEncoder) error

type baseEncoder struct {
	// Need to hang on to the underlying writer so we can control formatting.
	writer     io.Writer
	xmlEncoder *xml.Encoder

	// Set by StartArray or StartDict, and automatically closed when this encoder
	// is used again.
	encodingContainer bool

	// Set to true when writeEndTag() is called.
	finished bool
}

func newBaseEncoder(w io.Writer) *baseEncoder {
	e := xml.NewEncoder(w)
	e.Indent("", "\t")

	return &baseEncoder{
		writer:     w,
		xmlEncoder: e,
	}
}

// copy returns a *baseEncoder with the same writer and xml encoder but
// fresh state flags.
func (e *baseEncoder) copy() *baseEncoder {
	return &baseEncoder{
		writer:     e.writer,
		xmlEncoder: e.xmlEncoder,
	}
}

// writeEndTag encodes an end element and marks the encoder as finished.
// Any subsequent operations will panic.
func (e *baseEncoder) writeEndTag(endElement xml.EndElement) error {
	if err := e.xmlEncoder.EncodeToken(endElement); err != nil {
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
func (e *baseEncoder) writeArrayFunc(encode ArrayEncodingFunc) error {
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
func (e *baseEncoder) writeDictFunc(encode DictEncodingFunc) error {
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
		return e.EncodeElement("", boolTrueElement)
	}
	return e.EncodeElement("", boolFalseElement)
}

func writeFloat(e *xml.Encoder, val float64) error {
	return e.EncodeElement(val, realStartElement)
}

func writeBigFloat(e *xml.Encoder, val *big.Float) error {
	return e.EncodeElement(val.String(), realStartElement)
}

func writeInt(e *xml.Encoder, val int64) error {
	return e.EncodeElement(val, integerStartElement)
}

func writeUint(e *xml.Encoder, val uint64) error {
	return e.EncodeElement(val, integerStartElement)
}

func writeBigInt(e *xml.Encoder, val *big.Int) error {
	return e.EncodeElement(val.String(), integerStartElement)
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

func arrayWriter(val reflect.Value) ArrayEncodingFunc {
	return func(e *ArrayEncoder) (err error) {
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if err = e.write(elem); err != nil {
				return err
			}
		}
		return nil
	}
}

func mapWriter(val reflect.Value) DictEncodingFunc {
	return func(e *DictEncoder) (err error) {
		for _, key := range val.MapKeys() {
			elem := val.MapIndex(key)

			if key.Kind() != reflect.String {
				return fmt.Errorf("key must be a string: %v", key)
			}

			if err = e.write(key.String(), elem); err != nil {
				return err
			}
		}
		return nil
	}
}

func structWriter(val reflect.Value) DictEncodingFunc {
	return func(e *DictEncoder) (err error) {
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			fieldVal := val.Field(i)
			key := keyForField(field)

			if err = e.write(key, fieldVal); err != nil {
				return err
			}
		}
		return nil
	}
}

func keyForField(field reflect.StructField) string {
	if tag := field.Tag.Get(plistFieldTagName); tag != "" {
		return tag
	}
	return field.Name
}
