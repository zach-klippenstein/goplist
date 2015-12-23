package xml

import (
	"bytes"
	"encoding/xml"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriteString(t *testing.T) {
	var buffer bytes.Buffer
	e := xml.NewEncoder(&buffer)
	assert.NoError(t, writeString(e, "hello world"))
	assert.Equal(t, `<string>hello world</string>`, buffer.String())
}

func TestWriteBool(t *testing.T) {
	var buffer bytes.Buffer
	e := xml.NewEncoder(&buffer)
	assert.NoError(t, writeBool(e, false))
	assert.Equal(t, `<false></false>`, buffer.String())
}

func TestWriteInt(t *testing.T) {
	var buffer bytes.Buffer
	e := xml.NewEncoder(&buffer)
	assert.NoError(t, writeInt(e, -42))
	assert.Equal(t, `<integer>-42</integer>`, buffer.String())
}

func TestWriteUint(t *testing.T) {
	var buffer bytes.Buffer
	e := xml.NewEncoder(&buffer)
	assert.NoError(t, writeUint(e, 42))
	assert.Equal(t, `<integer>42</integer>`, buffer.String())
}

func TestWriteBigInt(t *testing.T) {
	var buffer bytes.Buffer
	e := xml.NewEncoder(&buffer)
	assert.NoError(t, writeBigInt(e, big.NewInt(42)))
	assert.Equal(t, `<integer>42</integer>`, buffer.String())
}

func TestWriteFloat(t *testing.T) {
	var buffer bytes.Buffer
	e := xml.NewEncoder(&buffer)
	assert.NoError(t, writeFloat(e, 4.2))
	assert.Equal(t, `<real>4.2</real>`, buffer.String())
}

func TestWriteBigFloat(t *testing.T) {
	var buffer bytes.Buffer
	e := xml.NewEncoder(&buffer)
	assert.NoError(t, writeBigFloat(e, big.NewFloat(4.2)))
	assert.Equal(t, `<real>4.2</real>`, buffer.String())
}

func TestWriteDate(t *testing.T) {
	var buffer bytes.Buffer
	e := xml.NewEncoder(&buffer)
	assert.NoError(t, writeDate(e, time.Date(2015, time.August, 1, 2, 3, 4, 5, time.UTC)))
	assert.Equal(t, `<date>2015-08-01T02:03:04Z</date>`, buffer.String())
}

func TestWriteData(t *testing.T) {
	var buffer bytes.Buffer
	e := xml.NewEncoder(&buffer)
	assert.NoError(t, writeData(e, []byte("hello world")))
	assert.Equal(t, `<data>aGVsbG8gd29ybGQ=</data>`, buffer.String())
}

func TestArrayWriter(t *testing.T) {
	data := []int{1, 2}

	var buffer bytes.Buffer
	e, err := newArrayEncoder(newBaseEncoder(&buffer))
	assert.NoError(t, err)

	w := arrayWriter(reflect.ValueOf(data))
	assert.NoError(t, w(e))

	assert.Equal(t, `<array>
	<integer>1</integer>
	<integer>2</integer>`, buffer.String())
}

func TestArrayWriterEmpty(t *testing.T) {
	data := []int{}

	var buffer bytes.Buffer
	e, err := newArrayEncoder(newBaseEncoder(&buffer))
	assert.NoError(t, err)

	w := arrayWriter(reflect.ValueOf(data))
	assert.NoError(t, w(e))

	assert.Equal(t, ``, buffer.String())
}

func TestMapWriter(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}

	var buffer bytes.Buffer
	e, err := newDictEncoder(newBaseEncoder(&buffer))
	assert.NoError(t, err)

	w := mapWriter(reflect.ValueOf(data))
	assert.NoError(t, w(e))

	// Go randomizes map iteration order, so we have to check all permutations.
	perms := []string{
		`<dict>
	<key>a</key>
	<integer>1</integer>
	<key>b</key>
	<integer>2</integer>`,
		`<dict>
	<key>b</key>
	<integer>2</integer>
	<key>a</key>
	<integer>1</integer>`,
	}

	match := false
	for _, perm := range perms {
		if perm == buffer.String() {
			match = true
		}
	}
	assert.True(t, match, "expected one of %v, got %s", perms, buffer.String())
}

func TestMapWriterEmpty(t *testing.T) {
	data := map[string]int{}

	var buffer bytes.Buffer
	e, err := newDictEncoder(newBaseEncoder(&buffer))
	assert.NoError(t, err)

	w := mapWriter(reflect.ValueOf(data))
	assert.NoError(t, w(e))

	assert.Equal(t, ``, buffer.String())
}

func TestMapWriterNonStringKey(t *testing.T) {
	data := map[int]int{1: 2}

	var buffer bytes.Buffer
	e, err := newDictEncoder(newBaseEncoder(&buffer))
	assert.NoError(t, err)

	w := mapWriter(reflect.ValueOf(data))
	assert.EqualError(t, w(e), "key must be a string: 1")
}

func TestStructWriter(t *testing.T) {
	data := struct {
		field1 string
		field2 string `plist:"field2 name"`
	}{
		"a", "b",
	}

	var buffer bytes.Buffer
	e, err := newDictEncoder(newBaseEncoder(&buffer))
	assert.NoError(t, err)

	w := structWriter(reflect.ValueOf(data))
	assert.NoError(t, w(e))

	assert.Equal(t, `<dict>
	<key>field1</key>
	<string>a</string>
	<key>field2 name</key>
	<string>b</string>`, buffer.String())
}
