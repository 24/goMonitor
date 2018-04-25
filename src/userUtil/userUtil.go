package userUtil

import (
	"github.com/axgle/mahonia"
	"time"
)

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func Datetime2timeStamp(datetime string) int64 {
	//获取本地location
	timeLayout := "2006-01-02 15:04:05"                           //转化所需模板
	loc, _ := time.LoadLocation("Local")                          //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, datetime, loc) //使用模板在对应时区转化为time.time类型
	return theTime.Unix()                                         //转化为时间戳 类型是int64
}

func Timestame2Datetime(timestamp int64) string {
	timeLayout := "2006-01-02 15:04:05"
	return time.Unix(timestamp, 0).Format(timeLayout)
}
