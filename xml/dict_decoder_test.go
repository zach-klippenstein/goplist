package xml

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeDictEntries(t *testing.T) {
	data := `<dict>
		<key>foo</key>
		<string>bar</string>
	</dict>`
	rootDecoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := rootDecoder.NextValue()
	assert.NoError(t, err)
	assert.IsType(t, &dictDecoder{}, value)

	decoder := value.(*dictDecoder)
	value, err = decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, DictEntry{"foo", "bar"}, value)

	value, err = decoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, EndDecodingContainer{}, value)
}

func TestDecodeDictArray(t *testing.T) {
	data := `<dict>
		<key>foo</key>
		<array>
			<string>foobar</string>
		</array>
	</dict>`
	rootDecoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := rootDecoder.NextValue()
	assert.NoError(t, err)
	assert.IsType(t, &dictDecoder{}, value)

	dictDecoder := value.(*dictDecoder)
	value, err = dictDecoder.NextValue()
	assert.NoError(t, err)
	entry := value.(DictEntry)
	assert.IsType(t, &arrayDecoder{}, entry.Value)

	arrayDecoder := entry.Value.(*arrayDecoder)
	value, err = arrayDecoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, "foobar", value)

	value, err = arrayDecoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, EndDecodingContainer{}, value)

	value, err = dictDecoder.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, EndDecodingContainer{}, value)
}

func TestDecodeDictDict(t *testing.T) {
	data := `<dict>
		<key>foo</key>
		<dict>
			<key>foo</key>
			<string>bar</string>
		</dict>
	</dict>`
	rootDecoder := baseDecoder{nil, xml.NewDecoder(bytes.NewReader([]byte(data)))}

	value, err := rootDecoder.NextValue()
	assert.NoError(t, err)
	assert.IsType(t, &dictDecoder{}, value)

	dictDecoder1 := value.(*dictDecoder)
	value, err = dictDecoder1.NextValue()
	assert.NoError(t, err)
	entry := value.(DictEntry)
	assert.IsType(t, &dictDecoder{}, entry.Value)

	dictDecoder2 := entry.Value.(*dictDecoder)
	value, err = dictDecoder2.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, DictEntry{"foo", "bar"}, value)

	value, err = dictDecoder2.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, EndDecodingContainer{}, value)

	value, err = dictDecoder1.NextValue()
	assert.NoError(t, err)
	assert.Equal(t, EndDecodingContainer{}, value)
}
