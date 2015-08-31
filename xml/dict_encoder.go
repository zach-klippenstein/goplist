package xml

import "time"

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
	return e.writeArray(encode)
}

func (e *DictEncoder) WriteDict(key string, encode DictEncodingFunc) error {
	e.assertReady()
	if err := e.writeKey(key); err != nil {
		return err
	}
	return e.writeDict(encode)
}

func (e *DictEncoder) writeEndTag() error {
	e.assertReady()
	return e.baseEncoder.writeEndTag(dictStartElement.End())
}

func (e *DictEncoder) writeKey(key string) error {
	return e.xmlEncoder.EncodeElement(key, dictKeyElement)
}
