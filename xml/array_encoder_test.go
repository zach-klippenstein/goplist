package xml

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteArray(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<array>
		<string>foo</string>
		<true></true>
		<real>4.2</real>
	</array>
</plist>`
	var buffer bytes.Buffer
	assert.NoError(t, WriteArrayPlist(&buffer, func(e *ArrayEncoder) error {
		assert.NoError(t, e.WriteString("foo"))
		assert.NoError(t, e.WriteBool(true))
		assert.NoError(t, e.WriteFloat(4.2))
		return nil
	}))
	assert.Equal(t, expected, buffer.String())
}

func TestWriteRecursiveArray(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<array>
		<array>
			<string>hi</string>
		</array>
	</array>
</plist>`
	var buffer bytes.Buffer
	assert.NoError(t, WriteArrayPlist(&buffer, func(e *ArrayEncoder) error {
		assert.NoError(t, e.WriteArray(func(e *ArrayEncoder) error {
			assert.NoError(t, e.WriteString("hi"))
			return nil
		}))
		return nil
	}))
	assert.Equal(t, expected, buffer.String())
}
