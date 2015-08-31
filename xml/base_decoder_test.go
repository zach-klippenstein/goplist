package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecodeNothing(t *testing.T) {
	data := ""
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, value)
}

func TestDecodeString(t *testing.T) {
	data := "<string>foo</string>"
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, "foo", value)

	value, err = decoder.NextValue()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, value)
}

func TestDecodeTrueSingleTag(t *testing.T) {
	data := "<true/>"
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, true, value)

	value, err = decoder.NextValue()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, value)
}

func TestDecodeTrueContainerTag(t *testing.T) {
	data := "<true></true>"
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, true, value)

	value, err = decoder.NextValue()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, value)
}

func TestDecodeFalse(t *testing.T) {
	data := "<false/>"
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, false, value)
}

func TestDecodePositiveInt(t *testing.T) {
	data := "<integer>42</integer>"
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), value)
}

func TestDecodeIntInvalid(t *testing.T) {
	data := "<integer>foo</integer>"
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.EqualError(t, err, `strconv.ParseInt: parsing "foo": invalid syntax`)
	assert.Nil(t, value)
}

func TestDecodeNegativeInt(t *testing.T) {
	data := "<integer>-42</integer>"
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, int64(-42), value)
}

func TestDecodeUint64(t *testing.T) {
	var tooBigForAnInt uint64 = uint64(math.MaxInt64) + 1
	data := fmt.Sprintf("<integer>%d</integer>", tooBigForAnInt)
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, tooBigForAnInt, value)
}

func TestDecodeBigInt(t *testing.T) {
	var hugeInt big.Int
	hugeInt.SetString("9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)
	data := fmt.Sprintf("<integer>%s</integer>", hugeInt.String())
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, hugeInt, value)
}

func TestDecodeReal(t *testing.T) {
	data := "<real>3.14</real>"
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, float64(3.14), value)
}

func TestDecodeBigReal(t *testing.T) {
	var hugeFloat big.Float
	hugeFloat.SetString("3.14e+99999")
	data := fmt.Sprintf("<real>%s</real>", hugeFloat.String())
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, hugeFloat, value)
}

func TestDecodeDate(t *testing.T) {
	data := "<date>2015-08-01T02:03:04Z</date>"
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2015, time.August, 1, 2, 3, 4, 0, time.UTC), value)
}

func TestDecodeData(t *testing.T) {
	data := "<data>aGVsbG8gd29ybGQ=</data>"
	decoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello world"), value)
}

func TestDecodeArray(t *testing.T) {
	data := "<array></array>"
	rootDecoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := rootDecoder.NextValue()
	assert.NoError(t, err)
	assert.IsType(t, &arrayDecoder{}, value)

	decoder := value.(*arrayDecoder)
	value, err = decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, EndDecodingContainer{}, value)

	value, err = rootDecoder.NextValue()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, value)
}

func TestDecodeDict(t *testing.T) {
	data := "<dict></dict>"
	rootDecoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := rootDecoder.NextValue()
	assert.NoError(t, err)
	assert.IsType(t, &dictDecoder{}, value)

	decoder := value.(*dictDecoder)
	value, err = decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, EndDecodingContainer{}, value)

	value, err = rootDecoder.NextValue()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, value)
}
