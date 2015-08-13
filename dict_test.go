package plist

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeysFromMapSortedAlpha(t *testing.T) {
	dict := NewDictFromMap(map[string]int{
		"c": 1,
		"a": 2,
		"b": 3,
	})
	assert.ObjectsAreEqual([]string{"a", "b", "c"}, dict.Keys())

	// These indents MUST be tab characters (\t).
	expected := wrap(`
	<dict>
		<key>a</key>
		<integer>2</integer>
		<key>b</key>
		<integer>3</integer>
		<key>c</key>
		<integer>1</integer>
	</dict>`)

	buffer := new(bytes.Buffer)
	err := Marshal(dict, buffer)
	assert.NoError(t, err)
	assert.Equal(t, expected, buffer.String())
}

func TestKeysAddedSorted(t *testing.T) {
	dict := NewEmptyDict()
	dict.Set("c", 1)
	dict.Set("b", 2)
	dict.Set("a", 3)
	assert.ObjectsAreEqual([]string{"c", "b", "a"}, dict.Keys())

	// These indents MUST be tab characters (\t).
	expected := wrap(`
	<dict>
		<key>c</key>
		<integer>1</integer>
		<key>b</key>
		<integer>2</integer>
		<key>a</key>
		<integer>3</integer>
	</dict>`)

	buffer := new(bytes.Buffer)
	err := Marshal(dict, buffer)
	assert.NoError(t, err)
	assert.Equal(t, expected, buffer.String())
}
