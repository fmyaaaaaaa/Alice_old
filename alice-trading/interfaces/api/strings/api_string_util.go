package strings

import (
	"net/url"
)

// string型のUrlをURL型に変換して返却します。
func ParsedUrl(targetUrl string) *url.URL {
	result, err := url.ParseRequestURI(targetUrl)
	if err != nil {
		panic(err)
	}
	return result
}
