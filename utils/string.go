package utils

import (
	"code.google.com/p/mahonia"
)

//把gbk字符串转为utf8
func StrGbk2Utf8(str string) string {
	return StrConvert("gb2312", str)
}

//转换字符串为指定的编码
func StrConvert(encode_type, str string) string {
	dec := mahonia.NewDecoder(encode_type)
	return dec.ConvertString(str)
}
