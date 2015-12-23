package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
)

func Encode(w io.Writer, v interface{}) error {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		return EncodeArrayPlist(w, arrayWriter(val))
	case reflect.Map:
		return EncodeDictPlist(w, mapWriter(val))
	case reflect.Struct:
		return EncodeDictPlist(w, structWriter(val))
	default:
		return fmt.Errorf("invalid type for encoding: %v (%v)", val, val.Kind())
	}
}

func EncodeArrayPlist(w io.Writer, encode ArrayEncodingFunc) error {
	encoder, err := startPlist(w)
	if err != nil {
		return err
	}
	if err = encoder.writeArrayFunc(encode); err != nil {
		return err
	}
	return writePlistEndTag(encoder)
}

func EncodeDictPlist(w io.Writer, encode DictEncodingFunc) error {
	encoder, err := startPlist(w)
	if err != nil {
		return err
	}
	if err = encoder.writeDictFunc(encode); err != nil {
		return err
	}
	return writePlistEndTag(encoder)
}

func startPlist(w io.Writer) (*baseEncoder, error) {
	base := newBaseEncoder(w)
	if err := writePlistHeader(base); err != nil {
		return nil, fmt.Errorf("error writing plist header: %s", err)
	}

	if err := base.xmlEncoder.EncodeToken(plistStartElement); err != nil {
		return nil, err
	}
	return base, nil
}

func writePlistHeader(e *baseEncoder) error {
	if err := e.xmlEncoder.EncodeToken(procInst); err != nil {
		return err
	}
	// Encoder won't add a newline after the processing instruction, so we have
	// to do it manually.
	if err := writeNewline(e.writer, e.xmlEncoder); err != nil {
		return err
	}

	if err := e.xmlEncoder.EncodeToken(doctype); err != nil {
		return err
	}
	if err := writeNewline(e.writer, e.xmlEncoder); err != nil {
		return err
	}
	return nil
}

func writePlistEndTag(e *baseEncoder) error {
	e.assertReady()
	if err := e.writeEndTag(plistStartElement.End()); err != nil {
		return err
	}
	return e.xmlEncoder.Flush()
}

// writeNewline is used to force conventional XML plist formatting when
// writing the file header.
func writeNewline(w io.Writer, e *xml.Encoder) error {
	// Flush before writing to the Writer ourselves.
	if err := e.Flush(); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return nil
}
