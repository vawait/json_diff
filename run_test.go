package json_diff

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	m := map[string]interface{}{
		"a": 1,
		"b": nil,
		"c": map[string]interface{}{"a": "kk", "b": "22"},
		"d": []interface{}{1, 2.4, 3, 423232323234},
	}
	lines := GetLines("", m, 3, LINE_TYPE_SAME)
	line := GetLineString(lines)
	fmt.Println(line)
}

func TestDiff(t *testing.T) {
	EnableAllConfig()
	SetPrintOnlyDiff(false)
	a := map[string]interface{}{
		"a": 1,
		"b": nil,
		"c": map[string]interface{}{"a": "kk", "b": 22},
		"d": []interface{}{1, 2.4, 3, 423232323234},
	}
	b := map[string]interface{}{
		"a": 2,
		"b": []int{},
		"c": map[string]interface{}{"a": "kk", "b": "22"},
		"d": []interface{}{1, 2.4, 1, 4, 423232323234, 5},
	}
	CheckJsonDiff(a, b, "a")
}
