package json_diff

import (
	"bytes"
	"fmt"
	"sort"
)

type LineType int

const (
	LINE_TYPE_SAME LineType = 0
	LINE_TYPE_ADD  LineType = 1
	LINE_TYPE_DEL  LineType = 2
)

var (
	printOnlyDiff = false
)

type DiffPair struct {
	oldValue interface{}
	newValue interface{}
	oldExist bool
	newExist bool
	modify   bool
}

func NewDiffPair(oldValue, newValue interface{}) *DiffPair {
	node := &DiffPair{}
	node.oldValue = oldValue
	node.newValue = newValue
	node.oldExist = true
	node.newExist = true
	node.modify = true
	return node
}

func (pair *DiffPair) Modify() bool {
	return pair.modify
}

func (pair *DiffPair) Print(key string, deviation int) (result []*Line) {
	if !pair.Modify() {
		if printOnlyDiff {
			return nil
		}
		return GetLines(key, pair.oldValue, deviation, LINE_TYPE_SAME)
	}
	if pair.oldExist {
		result = append(result, GetLines(key, pair.oldValue, deviation, LINE_TYPE_DEL)...)
	}
	if pair.newExist {
		result = append(result, GetLines(key, pair.newValue, deviation, LINE_TYPE_ADD)...)
	}
	return
}

type Line struct {
	deviation int
	line      string
	lineType  LineType
}

func (this *Line) Pop() {
	this.line = this.line[:len(this.line)-1]
}

func GetLines(name string, value interface{}, deviation int, lineType LineType) (result []*Line) {
	if name != "" {
		name = fmt.Sprintf(`"%s": `, name)
	}
	switch value.(type) {
	case map[string]interface{}:
		result = append(result, &Line{deviation: deviation, lineType: lineType, line: fmt.Sprintf(`%s{ `, name)})
		m := value.(map[string]interface{})
		for _, key := range sortedKeys(m) {
			result = append(result, GetLines(key, m[key], deviation+1, lineType)...)
		}
		result[len(result)-1].Pop()
		result = append(result, &Line{deviation: deviation, lineType: lineType, line: `},`})
	case []interface{}:
		result = append(result, &Line{deviation: deviation, lineType: lineType, line: fmt.Sprintf(`%s[ `, name)})
		m := value.([]interface{})
		for _, val := range m {
			result = append(result, GetLines("", val, deviation+1, lineType)...)
		}
		result[len(result)-1].Pop()
		result = append(result, &Line{deviation: deviation, lineType: lineType, line: `],`})
	case string:
		result = append(result, &Line{deviation: deviation, lineType: lineType, line: fmt.Sprintf(`%s"%+v",`, name, value)})
	case nil:
		result = append(result, &Line{deviation: deviation, lineType: lineType, line: fmt.Sprintf(`%snull,`, name)})
	default:
		result = append(result, &Line{deviation: deviation, lineType: lineType, line: fmt.Sprintf(`%s%+v,`, name, value)})
	}
	return
}

func saveSpace(buffer *bytes.Buffer, number int) {
	for i := 0; i < number; i++ {
		buffer.WriteString("  ")
	}
}

func GetLineString(lines []*Line) string {
	b := bytes.NewBuffer(nil)
	for i, line := range lines {
		if i == len(lines)-1 {
			line.Pop()
		}
		switch line.lineType {
		case LINE_TYPE_ADD:
			b.WriteString("\x1b[30;42m+ ")
			saveSpace(b, line.deviation)
			b.WriteString(line.line)
			b.WriteString("\x1b[0m")
		case LINE_TYPE_DEL:
			b.WriteString("\x1b[30;41m- ")
			saveSpace(b, line.deviation)
			b.WriteString(line.line)
			b.WriteString("\x1b[0m")
		case LINE_TYPE_SAME:
			saveSpace(b, line.deviation+1)
			b.WriteString(line.line)
		default:
			panic(line.lineType)
		}
		b.WriteRune('\n')
	}
	return Byte2String(b.Bytes())
}

func sortedKeys(m map[string]interface{}) (keys []string) {
	keys = make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}
