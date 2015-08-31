package xml

import (
	"encoding/xml"
	"fmt"
	"io"
)

func WriteArrayPlist(w io.Writer, encode ArrayEncodingFunc) error {
	encoder, err := startPlist(w)
	if err != nil {
		return err
	}
	if err = encoder.writeArray(encode); err != nil {
		return err
	}
	return writePlistEndTag(encoder)
}

func WriteDictPlist(w io.Writer, encode DictEncodingFunc) error {
	encoder, err := startPlist(w)
	if err != nil {
		return err
	}
	if err = encoder.writeDict(encode); err != nil {
		return err
	}
	return writePlistEndTag(encoder)
}

func startPlist(w io.Writer) (*baseEncoder, error) {
	base := newBaseEncoder(w)
	if err := writePlistHeader(base); err != nil {
		return nil, fmt.Errorf("error writing plist header: %s", err)
	}

	if err := base.e.EncodeToken(plistStartElement); err != nil {
		return nil, err
	}
	return base, nil
}

func writePlistHeader(e *baseEncoder) error {
	if err := e.e.EncodeToken(procInst); err != nil {
		return err
	}
	// Encoder won't add a newline after the processing instruction, so we have
	// to do it manually.
	if err := writeNewline(e.w, e.e); err != nil {
		return err
	}

	if err := e.e.EncodeToken(doctype); err != nil {
		return err
	}
	if err := writeNewline(e.w, e.e); err != nil {
		return err
	}
	return nil
}

func writePlistEndTag(e *baseEncoder) error {
	e.assertReady()
	if err := e.writeEndTag(plistStartElement.End()); err != nil {
		return err
	}
	return e.e.Flush()
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
