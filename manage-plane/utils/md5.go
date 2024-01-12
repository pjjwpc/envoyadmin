package utils

// 创建md5加密 方法
import (
	"crypto/md5"
	"encoding/hex"
)

func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
