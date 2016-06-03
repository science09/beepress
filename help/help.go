package help

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

func StrToInt32(str string) int32 {
	i, _ := strconv.ParseInt(str, 10, 0)
	return int32(i)
}

func MD5(str string) string {
	data := []byte(str)
	result := md5.Sum(data)
	md5Digest := hex.EncodeToString(result[:])
	return md5Digest
}
