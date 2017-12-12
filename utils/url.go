package utils

import (
	"net/url"
)

//url编码
func UrlEncode(url_str string) string {
	return url.PathEscape(url_str)
}
