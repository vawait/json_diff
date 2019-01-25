package json_diff

func SetPrintOnlyDiff(s bool) {
	printOnlyDiff = s
}

func SetDefaultValueZero(s bool) {
	defaultValueZero = s
}

func SetAutoDecodeByte(s bool) {
	autoDecodeByte = s
}

func SetBase64Decode(s bool) {
	base64Decode = s
}

func EnableAllConfig() {
	SetDefaultValueZero(true)
	SetAutoDecodeByte(true)
	SetBase64Decode(true)
	SetPrintOnlyDiff(true)
}
