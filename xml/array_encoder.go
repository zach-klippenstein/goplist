package xml

import "time"

type ArrayEncoder struct {
	*baseEncoder
}

var arrayStartElement = xmlElement("array")

func newArrayEncoder(base *baseEncoder) (*ArrayEncoder, error) {
	if err := base.e.EncodeToken(arrayStartElement); err != nil {
		return nil, err
	}
	return &ArrayEncoder{base}, nil
}

func (e *ArrayEncoder) WriteString(val string) error {
	e.assertReady()
	return writeString(e.e, val)
}

func (e *ArrayEncoder) WriteBool(val bool) error {
	e.assertReady()
	return writeBool(e.e, val)
}

func (e *ArrayEncoder) WriteFloat(val float64) error {
	e.assertReady()
	return writeFloat(e.e, val)
}

func (e *ArrayEncoder) WriteInt(val int64) error {
	e.assertReady()
	return writeInt(e.e, val)
}

func (e *ArrayEncoder) WriteUint(val uint64) error {
	e.assertReady()
	return writeUint(e.e, val)
}

func (e *ArrayEncoder) WriteDate(val time.Time) error {
	e.assertReady()
	return writeDate(e.e, val)
}

func (e *ArrayEncoder) WriteData(val []byte) error {
	e.assertReady()
	return writeData(e.e, val)
}

func (e *ArrayEncoder) WriteArray(encode ArrayEncodingFunc) error {
	e.assertReady()
	return e.writeArray(encode)
}

func (e *ArrayEncoder) WriteDict(encode DictEncodingFunc) error {
	e.assertReady()
	return e.writeDict(encode)
}

func (e *ArrayEncoder) writeEndTag() error {
	e.assertReady()
	return e.baseEncoder.writeEndTag(arrayStartElement.End())
}
