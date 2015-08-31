package xml

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"time"
)

type baseDecoder struct {
	parent     containerDecoder
	xmlDecoder *xml.Decoder
}

func (d *baseDecoder) NextValue() (interface{}, error) {
	token, err := nextStartOrEndElement(d.xmlDecoder)
	if err != nil {
		return nil, err
	}

	switch token := token.(type) {
	case xml.StartElement:
		switch token.Name {
		case stringStartElement.Name:
			return finishReadingString(d.xmlDecoder)
		case boolTrueElement.Name:
			return finishReadingBool(d.xmlDecoder, true, boolTrueElement.End())
		case boolFalseElement.Name:
			return finishReadingBool(d.xmlDecoder, false, boolFalseElement.End())
		case integerStartElement.Name:
			return finishReadingInteger(d.xmlDecoder)
		case realStartElement.Name:
			return finishReadingReal(d.xmlDecoder)
		case dateStartElement.Name:
			return finishReadingDate(d.xmlDecoder)
		case dataStartElement.Name:
			return finishReadingData(d.xmlDecoder)
		case arrayStartElement.Name:
			return newArrayDecoder(d, d.xmlDecoder), nil
		case dictStartElement.Name:
			return newDictDecoder(d, d.xmlDecoder), nil
		}

	case xml.EndElement:
		return EndDecodingContainer{}, nil
	}

	return nil, fmt.Errorf("Invalid element: %+v", token)
}

func (d *baseDecoder) ParentDecoder() containerDecoder {
	return d.parent
}

func nextStartOrEndElement(xmlDecoder *xml.Decoder) (xml.Token, error) {
	for {
		token, err := nextInterestingToken(xmlDecoder)
		if err != nil {
			return nil, err
		}

		switch token.(type) {
		case xml.StartElement, xml.EndElement:
			return token, nil
		}
	}
}

// finishReadingString consumes tokens, concatenating all CharDatas,
// until a </string> is encountered.
func finishReadingString(xmlDecoder *xml.Decoder) (interface{}, error) {
	str, err := readCharDataUntilEnd(xmlDecoder, stringStartElement.End())
	if err != nil {
		return nil, err
	}
	return str, nil
}

func finishReadingBool(xmlDecoder *xml.Decoder, value bool, end xml.EndElement) (interface{}, error) {
	for {
		token, err := xmlDecoder.Token()
		if err != nil {
			return nil, err
		}

		if token == end {
			return value, nil
		}
	}
}

// finishReadingInteger tries parsing an int64, then a uint64, and finally a big.Int,
// as required by the size of the value.
func finishReadingInteger(xmlDecoder *xml.Decoder) (interface{}, error) {
	raw, err := readCharDataUntilEnd(xmlDecoder, integerStartElement.End())
	if err != nil {
		return nil, err
	}

	var value interface{}
	value, err = strconv.ParseInt(raw, 10, 64)
	if err == nil {
		return value, nil
	} else if !isErrOutOfRange(err) {
		return nil, err
	}

	value, err = strconv.ParseUint(raw, 10, 64)
	if err == nil {
		return value, nil
	} else if !isErrOutOfRange(err) {
		return nil, err
	}

	var bigValue big.Int
	if _, ok := bigValue.SetString(raw, 10); ok {
		return bigValue, nil
	}

	return nil, fmt.Errorf("Could not parse '%s' as an integer.", raw)
}

// finishReadingReal tries parsing a float64, then a big.Float, as required by the
// size of the value.
func finishReadingReal(xmlDecoder *xml.Decoder) (interface{}, error) {
	raw, err := readCharDataUntilEnd(xmlDecoder, realStartElement.End())
	if err != nil {
		return nil, err
	}

	var value interface{}
	value, err = strconv.ParseFloat(raw, 64)
	if err == nil {
		return value, nil
	} else if !isErrOutOfRange(err) {
		return nil, err
	}

	var bigValue big.Float
	if _, _, err := bigValue.Parse(raw, 10); err != nil {
		return nil, err
	}
	return bigValue, nil
}

func isErrOutOfRange(err error) bool {
	if err, ok := err.(*strconv.NumError); ok && err.Err == strconv.ErrRange {
		return true
	}
	return false
}

func finishReadingDate(xmlDecoder *xml.Decoder) (interface{}, error) {
	raw, err := readCharDataUntilEnd(xmlDecoder, dateStartElement.End())
	if err != nil {
		return nil, err
	}

	date, err := time.Parse(dateFormat, raw)
	if err != nil {
		return nil, err
	}
	return date, nil
}

func finishReadingData(xmlDecoder *xml.Decoder) (interface{}, error) {
	raw, err := readCharDataUntilEnd(xmlDecoder, dataStartElement.End())
	if err != nil {
		return nil, err
	}

	encoded := bytes.NewReader([]byte(raw))
	decoder := base64.NewDecoder(base64.StdEncoding, encoded)
	data, err := ioutil.ReadAll(decoder)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func readCharDataUntilEnd(xmlDecoder *xml.Decoder, end xml.EndElement) (string, error) {
	var str bytes.Buffer

	for {
		token, err := xmlDecoder.Token()
		if err != nil {
			return "", err
		}

		switch token := token.(type) {
		case xml.CharData:
			str.Write(token)
		case xml.EndElement:
			if token == end {
				return str.String(), nil
			}
			return "", fmt.Errorf("Expected %#v, found %#v", end, token)
		}
	}
}
