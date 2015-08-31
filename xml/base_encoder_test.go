package xml

import (
	"bytes"
	"encoding/xml"
	"math/big"
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
