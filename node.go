package json_diff

import (
	"fmt"
	"sort"
)

type INode interface {
	Modify() bool
	Similarity() int64
	Print(key string, deviation int) []*Line
}

type ValueNode struct {
	*DiffPair
}

func NewValueNode(oldValue, newValue interface{}, modify ...bool) *ValueNode {
	node := &ValueNode{DiffPair: NewDiffPair(oldValue, newValue)}
	if len(modify) > 0 {
		node.modify = modify[0]
	}
	return node
}

func (node *ValueNode) Similarity() int64 {
	if node.Modify() {
		return 1
	}
	return 0
}

type MapNode struct {
	*DiffPair
	value map[string]INode
}

func NewMapNode(oldValue, newValue interface{}) *MapNode {
	return &MapNode{value: map[string]INode{}, DiffPair: NewDiffPair(oldValue, newValue)}
}

func (node *MapNode) Similarity() (result int64) {
	for _, pair := range node.value {
		if !pair.Modify() {
			result++
		}
	}
	return
}

func (node *MapNode) Modify() bool {
	return node.Similarity() != int64(len(node.value))
}

func (node *MapNode) Print(key string, deviation int) (result []*Line) {
	if printOnlyDiff && !node.Modify() {
		return nil
	}

	if key != "" {
		key = fmt.Sprintf(`"%s": `, key)
	}
	result = append(result, &Line{deviation: deviation, line: fmt.Sprintf(`%s{ `, key)})
	keys := make([]string, 0, len(node.value))
	for key, _ := range node.value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		result = append(result, node.value[key].Print(key, deviation+1)...)
	}
	result[len(result)-1].Pop()
	result = append(result, &Line{deviation: deviation, line: `},`})
	return
}

type ListNode struct {
	*DiffPair
	similarity int64
	value      []INode
	keys       []string
}

func NewListNode(oldValue, newValue interface{}) *ListNode {
	return &ListNode{DiffPair: NewDiffPair(oldValue, newValue)}
}

func (node *ListNode) Similarity() int64 {
	return node.similarity
}

func (node *ListNode) Modify() bool {
	for _, val := range node.value {
		if val.Modify() {
			return true
		}
	}
	return false
}

func (node *ListNode) AddValue(item INode, key string) {
	node.value = append(node.value, item)
	node.keys = append(node.keys, key)
}

func (node *ListNode) Print(key string, deviation int) (result []*Line) {
	if printOnlyDiff && !node.Modify() {
		return nil
	}

	if key != "" {
		key = fmt.Sprintf(`"%s": `, key)
	}
	result = append(result, &Line{deviation: deviation, line: fmt.Sprintf(`%s[ `, key)})
	for i, pair := range node.value {
		result = append(result, pair.Print(node.keys[i], deviation+1)...)
	}
	result[len(result)-1].Pop()
	result = append(result, &Line{deviation: deviation, line: `],`})
	return
}

var _ INode = &ValueNode{}
var _ INode = &MapNode{}
var _ INode = &ListNode{}
