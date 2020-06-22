package util

import (
	"log"
	"strconv"
)

// 文字列をintに変換します。
func ParseInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Println("fail to parse int")
	}
	return i
}

// 文字列をfloat64に変換します。
func ParseFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Print("fail to parse float64")
	}
	return f
}
