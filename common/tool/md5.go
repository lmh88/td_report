package tool

import (
	"crypto/md5"
	"fmt"
)

func GetMd5(str string) string {
	hs := md5.New()
	hs.Write([]byte(str))
	res := hs.Sum(nil)
	return fmt.Sprintf("%x", res)
}
