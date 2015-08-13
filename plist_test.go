package plist

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalEmpty(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"></plist>`

	buffer := new(bytes.Buffer)
	err := Marshal(nil, buffer)

	assert.NoError(t, err)
	assert.Equal(t, expected, buffer.String())
}

func TestDictRoot(t *testing.T) {
	dict := NewEmptyDict()
	dict.Set("string", "string")
	dict.Set("true", true)
	dict.Set("false", false)
	dict.Set("int", int(42))
	dict.Set("uint", uint(42))
	dict.Set("float32", float32(4.2))
	dict.Set("array", []string{"foo", "bar"})

	// These indents MUST be tab characters (\t).
	expected := wrap(`
	<dict>
		<key>string</key>
		<string>string</string>
		<key>true</key>
		<true></true>
		<key>false</key>
		<false></false>
		<key>int</key>
		<integer>42</integer>
		<key>uint</key>
		<integer>42</integer>
		<key>float32</key>
		<real>4.199999809265137</real>
		<key>array</key>
		<array>
			<string>foo</string>
			<string>bar</string>
		</array>
	</dict>`)

	buffer := new(bytes.Buffer)
	err := Marshal(dict, buffer)
	assert.NoError(t, err)
	assert.Equal(t, expected, buffer.String())
}

func TestArrayRoot(t *testing.T) {
	array := &Array{}
	array.Append("string")
	array.Append(true)
	array.Append(false)
	array.Append(int(42))
	array.Append(uint(42))
	array.Append(float32(4.2))
	array.Append(NewDictFromMap(map[string]string{"foo": "bar"}))

	// These indents MUST be tab characters (\t).
	expected := wrap(`
	<array>
		<string>string</string>
		<true></true>
		<false></false>
		<integer>42</integer>
		<integer>42</integer>
		<real>4.199999809265137</real>
		<dict>
			<key>foo</key>
			<string>bar</string>
		</dict>
	</array>`)

	buffer := new(bytes.Buffer)
	err := Marshal(array, buffer)

	assert.NoError(t, err)
	assert.Equal(t, expected, buffer.String())
}

func wrap(contents string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">` +
		contents + "\n" +
		`</plist>`
}
