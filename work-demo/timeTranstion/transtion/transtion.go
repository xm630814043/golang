package transtion

import (
    "fmt"
    "strings"
    "time"
)

var (
	timeFormat = "2006-01-02 15:04:05"
    nowTime = time.Now() // 当日时间
)
//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
    return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
    return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

//获取某一天的0点时间
func GetStartTime(d time.Time) time.Time {
    return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func GetEndTime(d time.Time) time.Time {
    return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
}


// 获取当天时间的开始和结束
func GetDayStartOrEndTime(d time.Time) (string,string) {
    curTime := time.Unix(d.Unix(), 0).Format(timeFormat)
    todayCurTime := strings.Split(curTime, " ")
    startTime := todayCurTime[0] + " 00:00:00" // 当日开始时间
    endTime := todayCurTime[0] + " 59:59:59"   // 当日结束时间
    return startTime,endTime
}

// 获取当月时间的开始和结束
func GetMonthStartOrEndTime(d time.Time) (string,string) {
    timeStart := GetFirstDateOfMonth(d)
    timeEnd := GetLastDateOfMonth(d)
    monthTimeStart := time.Unix(timeStart.Unix(), 0).Format("2006-01-02 15:04:05") // 当月开始时间
    monthTimeEnd := time.Unix(timeEnd.Unix(), 0).Format("2006-01-02 15:04:05")
    splitTimeEnd := strings.Split(monthTimeEnd, " ")
    monthEndTime := splitTimeEnd[0] + " 59:59:59" // 当月结束时间
    return monthTimeStart,monthEndTime
}


func Demo()  {
    t := nowTime.Unix()
    t += 86400 //增加一天
    curTime := time.Unix(t, 0).Format(timeFormat)
    fmt.Println("当天时间加一天：",curTime)
    timeStart,timeEnd := GetDayStartOrEndTime(nowTime)
    fmt.Println("当天开始时间：",timeStart,"当天结束时间：",timeEnd)
    monthTimeStart,monthTimeEnd := GetMonthStartOrEndTime(nowTime)
    fmt.Println("当月开始时间：",monthTimeStart,"当月结束时间：",monthTimeEnd)
    fmt.Println("当天开始时间：",GetStartTime(nowTime),"当天结束时间：",GetEndTime(nowTime))
}