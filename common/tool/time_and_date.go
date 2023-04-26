package tool

import (
	"fmt"
	"math"
	"td_report/app/bean"
	"td_report/vars"
	"time"
)

// GetWeek 获取当前时间是一年中第几周
func GetWeek(datetime string) (y, w int) {
	loc, _ := time.LoadLocation(vars.Timezon)
	tmp, _ := time.ParseInLocation(vars.TimeFormatTpl, datetime, loc)
	return tmp.ISOWeek()
}

// GetIoc 获取时区对应的ioc
func GetIoc() *time.Location {
	loc, _ := time.LoadLocation(vars.Timezon)
	return loc
}

// GetMonth 获取当前时间是一年中第几周
func GetMonth(datetime string) (y, w int) {
	loc, _ := time.LoadLocation(vars.Timezon)
	tmp, _ := time.ParseInLocation(vars.TimeFormatTpl, datetime, loc)
	return tmp.Year(), int(tmp.Month())
}

func GetLastDays(daysCount int, layout string) (dayList []string) {
	today := time.Now()
	for i := 0; i < daysCount; i++ {
		v := today.Add(-1 * time.Duration(i) * time.Hour * 24)
		dayList = append(dayList, v.Format(layout))
	}

	return
}

func ParaseWithLoc(datetime string, layout string) (time.Time, error) {
	loc, _ := time.LoadLocation(vars.Timezon)
	return time.ParseInLocation(layout, datetime, loc)
}

func ParaseWithChinaLoc(datetime string, layout string) (time.Time, error) {
	chinaZon, err := GetChinaZon()
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(layout, datetime, chinaZon)
}

func GetChinaZon() (*time.Location, error) {
	return time.LoadLocation("Asia/Shanghai")
}

// GetBetweenmonth 针对新用户获取调度的月的时间，即一个月有中挑选6天来调度
func GetBetweenmonth(sdate, edate string, layout string) []bean.Mydate {
	var mydate bean.Mydate
	d := make([]bean.Mydate, 0)
	date, err := ParaseWithLoc(sdate, layout)
	if err != nil {
		return d
	}
	date2, err := ParaseWithLoc(edate, layout)
	if err != nil {
		return d
	}
	if date2.Before(date) {
		return d
	}

	loc, _ := time.LoadLocation(vars.Timezon)
	s := date.Month()
	monthStart := time.Date(date.Year(), s, 1, 0, 0, 0, 0, loc)
	monthEnd := monthStart.AddDate(0, 1, -1)

	mydate.StartDate = monthStart.Format(vars.TimeFormatTpl)
	mydate.EndDate = monthEnd.Format(vars.TimeFormatTpl)
	if date2.Before(monthEnd) {
		mydate.EndDate = date2.Format(vars.TimeFormatTpl)
		d = append(d, mydate)
		return d
	}

	for {

		date = monthStart.AddDate(0, 1, 0)
		s = date.Month()
		monthStart = time.Date(date.Year(), s, 1, 0, 0, 0, 0, loc)
		monthEnd = monthStart.AddDate(0, 1, -1)
		mydate.StartDate = monthStart.Format(vars.TimeFormatTpl)
		mydate.EndDate = monthEnd.Format(vars.TimeFormatTpl)

		if date2.Before(monthEnd) {
			mydate.EndDate = date2.Format(vars.TimeFormatTpl)
			d = append(d, mydate)
			break
		} else {
			d = append(d, mydate)
		}
	}

	return d
}

// GetBetweenWeek  获取时间段内周一，三、六的时间周期，针对新用户调度
func GetBetweenWeek(sdate, edate string, layout string) []bean.Mydate {
	var mydate bean.Mydate
	d := make([]bean.Mydate, 0)
	date, err := ParaseWithLoc(sdate, layout)
	if err != nil {
		return d
	}
	date2, err := ParaseWithLoc(edate, layout)
	if err != nil {
		return d
	}

	if date2.Before(date) {
		return d
	}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	s := date.Weekday()
	cha := int(s - time.Sunday)
	g, _ := time.ParseDuration(fmt.Sprintf("-%dh", 24*cha))
	d1 := date.Add(g)
	month := date.Month()

	weekStart := time.Date(date.Year(), month, d1.Day()+1, 0, 0, 0, 0, loc)
	offset := int(time.Monday - weekStart.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekEnd := weekStart.AddDate(0, 0, offset+6)
	//weekEnd := weekStart.AddDate(0, 0, 6)
	mydate.StartDate = weekStart.Format(vars.TimeFormatTpl)
	mydate.EndDate = weekEnd.Format(vars.TimeFormatTpl)
	if date2.Before(weekEnd) {
		mydate.EndDate = date2.Format(vars.TimeFormatTpl)
		d = append(d, mydate)
		return d
	} else {
		d = append(d, mydate)
	}

	for {

		date = weekEnd.AddDate(0, 0, 1)
		month = date.Month()
		da := date.Day()
		weekStart = time.Date(date.Year(), month, da, 0, 0, 0, 0, loc)
		weekEnd = date.AddDate(0, 0, 6)

		mydate.StartDate = weekStart.Format(vars.TimeFormatTpl)
		mydate.EndDate = weekEnd.Format(vars.TimeFormatTpl)
		if date2.Before(weekEnd) {
			mydate.EndDate = date2.Format(vars.TimeFormatTpl)
			d = append(d, mydate)
			break
		} else {
			d = append(d, mydate)
		}
	}

	return d
}

// GetDaysWithTime 按照顺序排列日期
func GetDaysWithTime(start time.Time, end time.Time, layout string) ([]string, error) {
	cha := int(math.Ceil(end.Sub(start).Hours() / 24))
	daylist := make([]string, 0)
	for i := 0; i <= cha; i++ {
		day := start.Add(24 * time.Hour * time.Duration(i))
		daylist = append(daylist, day.Format(layout))
	}

	return daylist, nil
}

// GetDaysWithTimeRever 按照倒叙排列日期
func GetDaysWithTimeRever(start time.Time, end time.Time, layout string) ([]string, error) {
	cha := int(math.Ceil(end.Sub(start).Hours() / 24))
	daylist := make([]string, 0)
	for i := 0; i <= cha; i++ {
		day := end.Add(24 * time.Hour * time.Duration(-1*i))
		daylist = append(daylist, day.Format(layout))
	}

	return daylist, nil
}

// GetDaysDsp 专门准对dsp的时间格式化
func GetDaysDsp(start, end, layout string) ([]string, error) {
	opstart, err := ParaseWithLoc(start, layout)
	if err != nil {
		return nil, err
	}

	opend, err := ParaseWithLoc(end, layout)
	if err != nil {
		return nil, err
	}

	cha := int(math.Ceil(opend.Sub(opstart).Hours() / 24))
	daylist := make([]string, 0)
	for i := 0; i <= cha; i++ {
		date := opstart.Add(24 * time.Hour * time.Duration(i))
		daylist = append(daylist, date.Format(vars.TIMEDSPFILE))
	}

	return daylist, nil
}

// GetDays 根据开始日期和结束日期获取对应的天数列表
func GetDays(start, end, layout string) ([]string, error) {
	opstart, err := ParaseWithLoc(start, layout)
	if err != nil {
		return nil, err
	}

	opend, err := ParaseWithLoc(end, layout)
	if err != nil {
		return nil, err
	}

	cha := int(math.Ceil(opend.Sub(opstart).Hours() / 24))
	daylist := make([]string, 0)
	for i := 0; i <= cha; i++ {
		day := opstart.Add(24 * time.Hour * time.Duration(i))
		daylist = append(daylist, day.Format(layout))
	}

	return daylist, nil
}

// GetDaysPeriod 根据开始日期和结束日期获取对应的天数列表,period是间隔天数
func GetDaysPeriod(start, end, layout string, period int) ([]*bean.Mydate, error) {
	opstart, err := ParaseWithLoc(start, layout)
	if err != nil {
		return nil, err
	}

	opend, err := ParaseWithLoc(end, layout)
	if err != nil {
		return nil, err
	}

	cha := int(math.Ceil(opend.Sub(opstart).Hours() / float64(24*period)))
	daylist := make([]*bean.Mydate, 0)
	var startday, endday time.Time
	startday = opstart
	endday = opstart.Add(24 * time.Hour * time.Duration(1*period))
	opgstart := &bean.Mydate{
		StartDate: startday.Format(layout),
		EndDate:   endday.Format(layout),
	}

	daylist = append(daylist, opgstart)
	for i := 0; i <= cha; i++ {
		startday = endday.Add(24 * time.Hour * time.Duration(1))
		endday = startday.Add(24 * time.Hour * time.Duration(period))
		var mystart *bean.Mydate
		if startday.After(opend) {
			break
		} else {

			if endday.After(opend) {
				mystart = &bean.Mydate{
					StartDate: startday.Format(layout),
					EndDate:   opend.Format(layout),
				}
				daylist = append(daylist, mystart)
				break
			} else {
				mystart = &bean.Mydate{
					StartDate: startday.Format(layout),
					EndDate:   endday.Format(layout),
				}
			}
		}

		daylist = append(daylist, mystart)
	}

	return daylist, nil
}

// GetCurrentMonthGay 获取当前月份的天数
func GetCurrentMonthGay() int {
	current := time.Now()
	month := current.Format("1")
	year := current.Year()
	if month == "2" {
		if (year%4 == 0 && year%100 != 0) || (year%400 == 0) {
			return 29
		} else {
			return 28
		}
	} else {
		if month == "1" || month == "3" || month == "5" || month == "7" || month == "8" || month == "10" || month == "12" {
			return 31
		} else {
			return 30
		}
	}
}
