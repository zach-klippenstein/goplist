package xml

import (
	"fmt"
	"math/big"
	"reflect"
	"time"
)

type DictEncoder struct {
	*baseEncoder
}

var dictStartElement = xmlElement("dict")
var dictKeyElement = xmlElement("key")

func newDictEncoder(base *baseEncoder) (*DictEncoder, error) {
	if err := base.xmlEncoder.EncodeToken(dictStartElement); err != nil {
		return nil, err
	}
	return &DictEncoder{base}, nil
}

func (e *DictEncoder) Write(key string, val interface{}) error {
	e.assertReady()
	return e.write(key, reflect.ValueOf(val))
}

func (e *DictEncoder) write(key string, val reflect.Value) error {
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		return e.WriteArray(key, arrayWriter(val))
	case reflect.Map:
		return e.WriteDict(key, mapWriter(val))
	case reflect.Struct:
		return e.WriteDict(key, structWriter(val))
	case reflect.Bool:
		return e.WriteBool(key, val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.WriteInt(key, val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.WriteUint(key, val.Uint())
	case reflect.String:
		return e.WriteString(key, val.String())
	case reflect.Ptr, reflect.Interface:
		return e.write(key, val.Elem())
	default:
		return fmt.Errorf("invalid type for encoding: %v (%v)", val, val.Kind())
	}
}

func (e *DictEncoder) WriteString(key string, val string) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return writeString(e.xmlEncoder, val)
}

func (e *DictEncoder) WriteBool(key string, val bool) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return writeBool(e.xmlEncoder, val)
}

func (e *DictEncoder) WriteFloat(key string, val float64) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return writeFloat(e.xmlEncoder, val)
}

func (e *DictEncoder) WriteBigFloat(key string, val *big.Float) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return writeBigFloat(e.xmlEncoder, val)
}

func (e *DictEncoder) WriteInt(key string, val int64) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return writeInt(e.xmlEncoder, val)
}

func (e *DictEncoder) WriteUint(key string, val uint64) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return writeUint(e.xmlEncoder, val)
}

func (e *DictEncoder) WriteBigInt(key string, val *big.Int) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return writeBigInt(e.xmlEncoder, val)
}

func (e *DictEncoder) WriteDate(key string, val time.Time) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return writeDate(e.xmlEncoder, val)
}

func (e *DictEncoder) WriteData(key string, val []byte) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return writeData(e.xmlEncoder, val)
}

func (e *DictEncoder) WriteArray(key string, encode ArrayEncodingFunc) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return e.writeArrayFunc(encode)
}

func (e *DictEncoder) WriteDict(key string, encode DictEncodingFunc) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return e.writeDictFunc(encode)
}

func (e *DictEncoder) writeEndTag() error {
	e.assertReady()
	return e.baseEncoder.writeEndTag(dictStartElement.End())
}

func (e *DictEncoder) writeKey(key string) error {
	return e.xmlEncoder.EncodeElement(key, dictKeyElement)
}
