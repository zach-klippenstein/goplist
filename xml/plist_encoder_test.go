package xml

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteEmptyPlist(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"></plist>`
	var buffer bytes.Buffer
	encoder, err := startPlist(&buffer)
	assert.NoError(t, err)
	assert.NoError(t, writePlistEndTag(encoder))
	assert.Equal(t, expected, buffer.String())
}

func Example_encoding() {
	err := EncodeDictPlist(os.Stdout, func(e *DictEncoder) error {
		if err := e.WriteString("name", "Bilbo Baggins"); err != nil {
			return err
		}
		if err := e.WriteUint("age", 111); err != nil {
			return err
		}
		return e.WriteArray("acquaintances", func(e *ArrayEncoder) error {
			for _, name := range []string{
				"Gandalf the Grey",
				"Frodo Baggins",
				"Samwise Gamgee",
			} {
				if err := e.WriteString(name); err != nil {
					return err
				}
			}
			return nil
		})
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	// <plist version="1.0">
	// 	<dict>
	// 		<key>name</key>
	// 		<string>Bilbo Baggins</string>
	// 		<key>age</key>
	// 		<integer>111</integer>
	// 		<key>acquaintances</key>
	// 		<array>
	// 			<string>Gandalf the Grey</string>
	// 			<string>Frodo Baggins</string>
	// 			<string>Samwise Gamgee</string>
	// 		</array>
	// 	</dict>
	// </plist>
}
