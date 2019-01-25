package json_diff

import (
	"fmt"
)

func CheckJsonDiff(a, b interface{}, ignore ...string) (bool, error) {
	a, err := NewJsonItem(a)
	if err != nil {
		return false, err
	}
	b, err = NewJsonItem(b)
	if err != nil {
		return false, err
	}
	if autoDecodeByte {
		a = deepDecode(a)
		b = deepDecode(b)
	}

	ignoreMap := map[string]bool{}
	for _, key := range ignore {
		ignoreMap[key] = true
	}

	node := compare(a, b, ignoreMap)
	lines := node.Print("", 0)
	output := GetLineString(lines)
	fmt.Println(output)

	return node.Modify(), nil
}

func compare(a, b interface{}, ignoreMap map[string]bool) INode {
	switch a.(type) {
	case map[string]interface{}:
		if bMap, ok := b.(map[string]interface{}); ok {
			return compareMap(a.(map[string]interface{}), bMap, ignoreMap)
		}
	case []interface{}:
		if bSlice, ok := b.([]interface{}); ok {
			return compareSlice(a.([]interface{}), bSlice, ignoreMap)
		}
	}
	node := NewValueNode(a, b)
	if DeepEqual(a, b) {
		node.modify = false
	}
	return node
}

func compareMap(a, b map[string]interface{}, ignoreMap map[string]bool) INode {
	mapNode := NewMapNode(a, b)
	for key, _ := range a {
		if ignoreMap[key] {
			mapNode.value[key] = NewValueNode(a[key], b[key], false)
			continue
		}
		if _, ok := b[key]; !ok {
			node := NewValueNode(a[key], b[key])
			node.newExist = false
			mapNode.value[key] = node
		} else {
			mapNode.value[key] = compare(a[key], b[key], ignoreMap)
		}
	}

	for key, val := range b {
		if _, ok := a[key]; ok {
			continue
		}
		if ignoreMap[key] {
			mapNode.value[key] = NewValueNode(a[key], b[key], false)
			continue
		}
		node := NewValueNode(nil, val)
		node.oldExist = false
		mapNode.value[key] = node
	}
	return mapNode
}

func compareSlice(a, b []interface{}, ignoreMap map[string]bool) INode {
	listNode := NewListNode(a, b)
	length1 := len(a)
	length2 := len(b)
	f := New2DArray(length1+1, length2+1)
	for i := length1 - 1; i >= 0; i-- {
		for j := length2 - 1; j >= 0; j-- {
			if DeepEqual(a[i], b[j], ignoreMap) {
				f[i][j] = f[i+1][j+1] + 1
			} else if f[i+1][j] > f[i][j+1] {
				f[i][j] = f[i+1][j]
			} else {
				f[i][j] = f[i][j+1]
			}
		}
	}

	fp1 := make([]int, length1)
	fp2 := make([]int, length2+1)
	i := 0
	j := 0
	for i < length1 && j < length2 {
		if DeepEqual(a[i], b[j], ignoreMap) {
			fp1[i] = j + 1
			fp2[j] = i + 1
			i++
			j++
		} else if f[i+1][j] == f[i][j] {
			i++
		} else {
			j++
		}
	}

	j = 0
	for i := 0; i < length1; i++ {
		for j < length2 && fp1[i] > 0 && fp1[i] != j+1 {
			node := NewValueNode(nil, b[j])
			node.oldExist = false
			listNode.AddValue(node, "")
			j++
		}
		if j < length2 && (fp1[i] == j+1 || fp2[j] == 0 && fp1[i] == 0) {
			listNode.AddValue(compare(a[i], b[j], ignoreMap), fmt.Sprintf("%d", i))
			j++
		} else {
			node := NewValueNode(a[i], nil)
			node.newExist = false
			listNode.AddValue(node, fmt.Sprintf("%d", i))
		}
	}
	for j < length2 {
		node := NewValueNode(nil, b[j])
		node.oldExist = false
		listNode.AddValue(node, "")
		j++
	}

	return listNode
}

func New2DArray(n, m int) [][]int {
	f := make([][]int, n)
	for i, _ := range f {
		f[i] = make([]int, m)
	}
	return f
}

func getSlice(data []interface{}, index int) interface{} {
	if index < 0 || index >= len(data) {
		return nil
	}
	return data[index]
}
