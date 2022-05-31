package tools

import "time"

func GetNextDayStartTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.Local)
}

func GetNextDayEndTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 23, 59, 59, 59, time.Local)
}

func GetLastDayStartTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 0, 0, 0, 0, time.Local)
}

func GetLastDayEndTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 23, 59, 59, 59, time.Local)
}

func GetNextWeekStartTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+7, 0, 0, 0, 0, time.Local)
}

func GetNextWeekEndTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+7, 23, 59, 59, 59, time.Local)
}

func GetLastWeekStartTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-7, 0, 0, 0, 0, time.Local)
}

func GetLastWeekEndTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-7, 23, 59, 59, 59, time.Local)
}

func GetNextMonthStartTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.Local)
}

func GetNextMonthEndTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month()+2, 1, 0, 0, 0, 0, time.Local).Add(-1 * time.Second)
}

func GetLastMonthStartTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month()-1, 1, 0, 0, 0, 0, time.Local)
}

func GetLastMonthEndTime() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local).Add(-1 * time.Second)
}
