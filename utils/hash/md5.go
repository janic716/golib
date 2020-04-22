package utils

import (
	"crypto/md5"
	"fmt"
)

func Md5(value interface{}) string {
	if v, ok := value.(string); ok {
		return fmt.Sprintf("%x", md5.Sum([]byte(v)))
	} else {
		return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(value))))
	}
}

func Md5File(filePath string) string {
	//todo
	return ""
}
