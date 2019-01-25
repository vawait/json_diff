# json_diff

用法
===
调用函数CheckJsonDiff
输入：需要对比的两个值(可以为结构体，会自动转成json)，以及需要忽略diff的key
返回：是否有变化、错误

参数
===
SetDefaultValueZero：      会把所有类型默认值当成空值等价处理
SetAutoDecodeByte：        会自动把string或者[]byte解析json
SetBase64Decode：          需要开启AutoDecodeByte，先进行base64解码再解析json
SetPrintOnlyDiff：         只输出有diff
EnableAllConfig：          开启上面所有参数

样例
===
    package main
    
    import (
    	"github.com/vawait/json_diff"
    )
    
    func main() {
    	json_diff.EnableAllConfig()
    	json_diff.SetPrintOnlyDiff(false)
    	a := map[string]interface{}{
    		"a": 1,
    		"b": nil,
    		"c": map[string]interface{}{"a": "kk", "b": "22"},
    		"d": []interface{}{1, 2.4, 3, 423232323234},
    	}
    	b := map[string]interface{}{
    		"a": 2,
    		"b": []int{},
    		"c": map[string]interface{}{"a": "kk", "b": "221"},
    		"d": []interface{}{1, 2.4, 1, 4, 423232323234, 5},
    	}
    	json_diff.CheckJsonDiff(a, b, "a")
    }
