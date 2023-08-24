package utils

import (
	"github.com/wonderivan/logger"
	"regexp"
	"time"
)

//获取当前时间戳
func GetNowTimestamp() int64 {
	return time.Now().Unix() * 1000
}

//获取当前时间
func GetNowTime() time.Time {
	return time.Now()
}

//获取当前时间字符串
func GetNowTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//获取当前时间字符串
func GetNowTimeString2() string {
	return time.Now().Format("20060102150405")
}

func RemoveDuplicateElement(languages []uint) []uint {
	result := make([]uint, 0, len(languages))
	temp := map[uint]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func CheckEn(myString string) bool {
	b, err := regexp.MatchString("^([a-zA-Z0-9_]+)$", myString)
	if err != nil {
		logger.Error(err.Error())
	}
	return b
}
