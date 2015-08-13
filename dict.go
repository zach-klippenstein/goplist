package plist

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"sort"
)

type Dict struct {
	dict     map[string]plistElement
	ordering []string
}

var (
	dictStartElement = xmlElement("dict")
	dictKeyElement   = xmlElement("key")
)

func NewEmptyDict() *Dict {
	return &Dict{
		dict: make(map[string]plistElement),
	}
}

// NewDict returns a *Dict with the contents of stringMap, which must be
// a map of string to some plist-compatible type.
func NewDictFromMap(stringMap interface{}) *Dict {
	return createDict(reflect.ValueOf(stringMap))
}

// createDict creates a Dict from a map[] with keys sorted alphabetically.
func createDict(v reflect.Value) *Dict {
	if v.Kind() != reflect.Map {
		panic(fmt.Sprint("cannot create dict from", v.Type()))
	}
	if v.Type().Key().Kind() != reflect.String {
		panic("plist dict keys must be strings")
	}

	// Get all the string keys so we can sort them.
	keyValues := v.MapKeys()
	keys := make([]string, len(keyValues))
	for i, k := range keyValues {
		keys[i] = k.String()
	}

	sort.Strings(keys)
	dict := NewEmptyDict()

	for _, key := range keys {
		value := v.MapIndex(reflect.ValueOf(key))
		dict.Set(key, createPlistElement(value.Interface()))
	}
	return dict
}

// Set associates key with value, and if it doesn't exist, adds it to the end
// of the dict.
// Entries will be written out in the order they're added.
// Returns an error if value is not a valid plist type.
func (d *Dict) Set(key string, value interface{}) error {
	// Validate the value first.
	elem := createPlistElement(value)
	if _, exists := d.dict[key]; !exists {
		d.ordering = append(d.ordering, key)
	}
	d.dict[key] = elem
	return nil
}

func (d *Dict) Get(key string) (value interface{}, exists bool) {
	if elem, exists := d.dict[key]; exists {
		return elem.value(), true
	}
	return nil, false
}

func (d *Dict) Keys() []string {
	keys := make([]string, len(d.ordering))
	copy(keys, d.ordering)
	return keys
}

func (d *Dict) value() interface{} {
	var dict map[string]interface{}
	for k, v := range d.dict {
		dict[k] = v.value()
	}
	return dict
}

func (d *Dict) marshal(e *xml.Encoder) error {
	if err := e.EncodeToken(dictStartElement); err != nil {
		return err
	}

	for _, k := range d.ordering {
		// Encode key first.
		if err := e.EncodeElement(k, dictKeyElement); err != nil {
			return err
		}

		// Then encode value.
		v := d.dict[k]
		if err := v.marshal(e); err != nil {
			return err
		}
	}

	return e.EncodeToken(dictStartElement.End())
}
