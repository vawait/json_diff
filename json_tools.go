package json_diff

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
)

var (
	defaultValueZero = false
	autoDecodeByte   = false
	base64Decode     = false
)

func IsZero(data interface{}) bool {
	if data == nil {
		return true
	}

	if !defaultValueZero {
		return false
	}

	switch data.(type) {
	case float64:
		return int64(data.(float64)) == 0
	case bool:
		return data.(bool) == false
	case string:
		return data.(string) == ""
	case int:
		return data.(int) == 0
	case int32:
		return data.(int32) == 0
	case int64:
		return data.(int64) == 0
	case json.Number:
		return data.(json.Number) == "0"
	case []interface{}:
		return len(data.([]interface{})) == 0
	case map[string]interface{}:
		return len(data.(map[string]interface{})) == 0
	}
	return false
}

func NewJsonItem(item interface{}) (interface{}, error) {
	var err error
	byteData, ok := GetByte(item)
	if !ok {
		byteData, err = json.Marshal(item)
		if err != nil {
			return item, err
		}
	}
	var newItem interface{}
	d := json.NewDecoder(bytes.NewReader(byteData))
	d.UseNumber()
	if err := d.Decode(&newItem); err != nil {
		return item, err
	}
	return newItem, nil
}

func DeepEqual(a, b interface{}, ignoreMap ...map[string]bool) bool {
	if IsZero(a) && IsZero(b) {
		return true
	}
	if IsZero(a) || IsZero(b) {
		return false
	}
	switch a.(type) {
	case []interface{}:
		newA := a.([]interface{})
		newB, ok := b.([]interface{})
		if !ok || len(newB) != len(newA) {
			return false
		}
		for i, _ := range newA {
			if !DeepEqual(newA[i], newB[i], ignoreMap...) {
				return false
			}
		}
	case map[string]interface{}:
		newA := a.(map[string]interface{})
		newB, ok := b.(map[string]interface{})
		if !ok || len(newB) != len(newA) {
			return false
		}
		for key, val1 := range newA {
			if len(ignoreMap) > 0 && ignoreMap[0][key] {
				continue
			}
			val2, ok := newB[key]
			if !ok || !DeepEqual(val1, val2, ignoreMap...) {
				return false
			}
		}
	default:
		return reflect.DeepEqual(a, b)
	}
	return true
}

func NewDecodeValue(item interface{}) (interface{}, bool) {
	byteData, ok := GetByte(item)
	if !ok {
		return item, false
	}
	if base64Decode {
		if newByte, err := base64.StdEncoding.DecodeString(Byte2String(byteData)); err == nil {
			byteData = newByte
		}
	}
	newVal := map[string]interface{}{}
	d := json.NewDecoder(bytes.NewReader(byteData))
	d.UseNumber()
	if d.Decode(&newVal) == nil {
		return newVal, true
	}
	return item, false
}

func deepDecode(item interface{}) interface{} {
	switch item.(type) {
	case map[string]interface{}:
		m := item.(map[string]interface{})
		for key, val := range m {
			m[key] = deepDecode(val)
		}
		return m
	case []interface{}:
		m := item.([]interface{})
		for index, val := range m {
			m[index] = deepDecode(val)
		}
		return m
	default:
		var ok bool
		if item, ok = NewDecodeValue(item); ok {
			item = deepDecode(item)
		}
	}
	return item
}
