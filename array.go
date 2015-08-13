package plist

import (
	"encoding/xml"
	"reflect"
)

type Array struct {
	array []plistElement
}

var arrayStartElement = xmlElement("array")

func createArray(v reflect.Value) plistElement {
	length := v.Len()
	array := &Array{make([]plistElement, length)}
	for i := 0; i < length; i++ {
		array.array[i] = createPlistElement(v.Index(i).Interface())
	}
	return array
}

func (a *Array) Get(index int) interface{} {
	return a.array[index].value()
}

func (a *Array) Append(value interface{}) {
	a.array = append(a.array, createPlistElement(value))
}

func (a *Array) Set(index int, value interface{}) {
	a.array[index] = createPlistElement(value)
}

func (a *Array) value() interface{} {
	array := make([]interface{}, len(a.array))
	for i, v := range a.array {
		array[i] = v
	}
	return array
}

func (a *Array) marshal(e *xml.Encoder) error {
	if err := e.EncodeToken(arrayStartElement); err != nil {
		return err
	}
	for _, v := range a.array {
		if err := v.marshal(e); err != nil {
			return err
		}
	}
	return e.EncodeToken(arrayStartElement.End())
}
