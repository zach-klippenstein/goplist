package xml

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteDict(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>foo</key>
		<string>bar</string>
		<key>true</key>
		<true></true>
		<key>the answer to everything</key>
		<integer>42</integer>
	</dict>
</plist>`
	var buffer bytes.Buffer
	assert.NoError(t, EncodeDictPlist(&buffer, func(e *DictEncoder) error {
		assert.NoError(t, e.WriteString("foo", "bar"))
		assert.NoError(t, e.WriteBool("true", true))
		assert.NoError(t, e.WriteInt("the answer to everything", 42))
		return nil
	}))
	assert.Equal(t, expected, buffer.String())
}

func TestWriteRecursiveDict(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>rabbit hole</key>
		<dict>
			<key>foo</key>
			<string>bar</string>
		</dict>
	</dict>
</plist>`
	var buffer bytes.Buffer
	assert.NoError(t, EncodeDictPlist(&buffer, func(e *DictEncoder) error {
		assert.NoError(t, e.WriteDict("rabbit hole", func(e *DictEncoder) error {
			assert.NoError(t, e.WriteString("foo", "bar"))
			return nil
		}))
		return nil
	}))
	assert.Equal(t, expected, buffer.String())
}

func TestWriteDictInterface(t *testing.T) {
	// TODO
}
