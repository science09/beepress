package help

import "strconv"

func StrToInt32(str string) int32 {
	i, _ := strconv.ParseInt(str, 10, 0)
	return int32(i)
}