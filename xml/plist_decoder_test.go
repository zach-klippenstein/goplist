package xml

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeEmptyArrayPlist(t *testing.T) {
	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<array>
	</array>
</plist>`
	decoder := NewDecoder(bytes.NewReader([]byte(plist)))

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.IsType(t, StartDecodingArray{}, value)

	value, err = decoder.NextValue()
	assert.NoError(t, err)
	assert.IsType(t, EndDecodingContainer{}, value)

	value, err = decoder.NextValue()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, value)
}

func TestDecodeEmptyDictPlist(t *testing.T) {
	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
	</dict>
</plist>`
	decoder := NewDecoder(bytes.NewReader([]byte(plist)))

	value, err := decoder.NextValue()
	assert.NoError(t, err)
	assert.IsType(t, StartDecodingDict{}, value)

	value, err = decoder.NextValue()
	assert.NoError(t, err)
	assert.IsType(t, EndDecodingContainer{}, value)

	value, err = decoder.NextValue()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, value)
}

func ExamplePlistDecoder_array() {
	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<array>
		<string>hello world</string>
		<true></true>
		<integer>42</integer>
		<real>3.14</real>
		<date>2015-08-01T02:03:04Z</date>
		<data>aGVsbG8gd29ybGQ=</data>
	</array>
</plist>`

	decoder := NewDecoder(bytes.NewReader([]byte(plist)))
	for {
		value, err := decoder.NextValue()
		if err != nil {
			log.Fatalln(err)
		}

		switch value.(type) {
		case EndDecodingContainer:
			return
		case StartDecodingArray:
			continue
		default:
			fmt.Println(value)
		}
	}

	// Output:
	// hello world
	// true
	// 42
	// 3.14
	// 2015-08-01 02:03:04 +0000 UTC
	// [104 101 108 108 111 32 119 111 114 108 100]
}

func ExamplePlistDecoder_dict() {
	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>a</key>
		<string>hello world</string>
		<key>b</key>
		<true></true>
		<key>c</key>
		<integer>42</integer>
		<key>d</key>
		<real>3.14</real>
		<key>e</key>
		<date>2015-08-01T02:03:04Z</date>
		<key>f</key>
		<data>aGVsbG8gd29ybGQ=</data>
	</dict>
</plist>`

	decoder := NewDecoder(bytes.NewReader([]byte(plist)))
	for {
		value, err := decoder.NextValue()
		if err != nil {
			log.Fatalln(err)
		}

		if _, ok := value.(EndDecodingContainer); ok {
			return
		}

		if entry, ok := value.(DictEntry); ok {
			fmt.Printf("%s: %v\n", entry.Key, entry.Value)
		}
	}

	// Output:
	// a: hello world
	// b: true
	// c: 42
	// d: 3.14
	// e: 2015-08-01 02:03:04 +0000 UTC
	// f: [104 101 108 108 111 32 119 111 114 108 100]
}
