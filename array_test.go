package plist

import "testing"

func TestAppendDictToArray(t *testing.T) {
	array := &Array{}
	array.Append(NewDictFromMap(map[string]string{"foo": "bar"}))
}

func TestAppendMapToArray(t *testing.T) {
	array := &Array{}
	array.Append(map[string]string{"foo": "bar"})
}
