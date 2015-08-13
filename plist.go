package plist

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

var indentation = strings.Repeat(" ", 4)

var procInst = xml.ProcInst{
	Target: "xml",
	Inst:   []byte(`version="1.0" encoding="UTF-8"`),
}

var doctype = xml.Directive([]byte(
	`DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"`))

var plistStartElement = xml.StartElement{
	Name: xml.Name{"", "plist"},
	Attr: []xml.Attr{xml.Attr{xml.Name{"", "version"}, "1.0"}},
}

type plistElement interface {
	// value returns the value of the element as an unwrapped Go type.
	value() interface{}
	marshal(e *xml.Encoder) error
}

func Marshal(v interface{}, writer io.Writer) error {
	root := getOrCreateRootElement(v)

	e := xml.NewEncoder(writer)
	if err := encodeHeader(writer, e); err != nil {
		return err
	}
	if root != nil {
		if err := root.marshal(e); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(plistStartElement.End()); err != nil {
		return err
	}
	return e.Flush()
}

func encodeHeader(w io.Writer, e *xml.Encoder) error {
	e.Indent("", "\t")
	if err := e.EncodeToken(procInst); err != nil {
		return err
	}
	// Flush before writing to the Writer ourselves.
	if err := e.Flush(); err != nil {
		return err
	}
	// Encoder won't add a newline after the processing instruction, so we have
	// to do it manually.
	fmt.Fprintln(w)
	if err := e.EncodeToken(doctype); err != nil {
		return err
	}
	// Flush before writing to the Writer ourselves.
	if err := e.Flush(); err != nil {
		return err
	}
	// Encoder won't add a newline after the processing instruction, so we have
	// to do it manually.
	fmt.Fprintln(w)
	if err := e.EncodeToken(plistStartElement); err != nil {
		return err
	}
	return nil
}

func getOrCreateRootElement(v interface{}) plistElement {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case *Dict:
		return t
	case *Array:
		return t
	case plistElement:
		panic("root element must be a *Dict or *Array")
	default:
		return getOrCreateRootElement(createPlistElement(v))
	}
}

func xmlElement(name string) xml.StartElement {
	return xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: name,
		},
	}
}

func createPlistElement(value interface{}) plistElement {
	if value == nil {
		return nil
	}
	switch t := value.(type) {
	case plistElement:
		return t
	case time.Time, *time.Time:
		// TODO
		panic("date plist values not supported")
	case []byte, *bytes.Buffer:
		// TODO
		panic("data plist values not supported")
	case string:
		return plistString(t)
	case bool:
		return plistBool(t)
	default:
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			return plistInt{signed: v.Int()}
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			return plistInt{unsigned: v.Uint()}
		case reflect.Float32, reflect.Float64:
			return plistReal(v.Float())
		case reflect.Map:
			return createDict(reflect.ValueOf(value))
		case reflect.Slice, reflect.Array:
			return createArray(reflect.ValueOf(value))
		}
	}

	panic(fmt.Sprintf("%s plist values not supported", reflect.TypeOf(value)))
}
