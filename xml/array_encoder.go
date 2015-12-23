package xml

import (
	"fmt"
	"math/big"
	"reflect"
	"time"
)

type ArrayEncoder struct {
	*baseEncoder
}

var arrayStartElement = xmlElement("array")

func newArrayEncoder(base *baseEncoder) (*ArrayEncoder, error) {
	if err := base.xmlEncoder.EncodeToken(arrayStartElement); err != nil {
		return nil, err
	}
	return &ArrayEncoder{base}, nil
}

func (e *ArrayEncoder) Write(val interface{}) error {
	e.assertReady()
	return e.write(reflect.ValueOf(val))
}

func (e *ArrayEncoder) write(val reflect.Value) error {
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		return e.WriteArray(arrayWriter(val))
	case reflect.Map:
		return e.WriteDict(mapWriter(val))
	case reflect.Struct:
		return e.WriteDict(structWriter(val))
	case reflect.Bool:
		return e.WriteBool(val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.WriteInt(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.WriteUint(val.Uint())
	case reflect.String:
		return e.WriteString(val.String())
	case reflect.Ptr, reflect.Interface:
		return e.write(val.Elem())
	default:
		return fmt.Errorf("invalid type for encoding: %v (%v)", val, val.Kind())
	}
}

func (e *ArrayEncoder) WriteString(val string) error {
	e.assertReady()
	return writeString(e.xmlEncoder, val)
}

func (e *ArrayEncoder) WriteBool(val bool) error {
	e.assertReady()
	return writeBool(e.xmlEncoder, val)
}

func (e *ArrayEncoder) WriteFloat(val float64) error {
	e.assertReady()
	return writeFloat(e.xmlEncoder, val)
}

func (e *ArrayEncoder) WriteBigFloat(val *big.Float) error {
	e.assertReady()
	return writeBigFloat(e.xmlEncoder, val)
}

func (e *ArrayEncoder) WriteInt(val int64) error {
	e.assertReady()
	return writeInt(e.xmlEncoder, val)
}

func (e *ArrayEncoder) WriteUint(val uint64) error {
	e.assertReady()
	return writeUint(e.xmlEncoder, val)
}

func (e *ArrayEncoder) WriteBigInt(val *big.Int) error {
	e.assertReady()
	return writeBigInt(e.xmlEncoder, val)
}

func (e *ArrayEncoder) WriteDate(val time.Time) error {
	e.assertReady()
	return writeDate(e.xmlEncoder, val)
}

func (e *ArrayEncoder) WriteData(val []byte) error {
	e.assertReady()
	return writeData(e.xmlEncoder, val)
}

func (e *ArrayEncoder) WriteArray(encode ArrayEncodingFunc) error {
	e.assertReady()
	return e.writeArrayFunc(encode)
}

func (e *ArrayEncoder) WriteDict(encode DictEncodingFunc) error {
	e.assertReady()
	return e.writeDictFunc(encode)
}

func (e *ArrayEncoder) writeEndTag() error {
	e.assertReady()
	return e.baseEncoder.writeEndTag(arrayStartElement.End())
}
