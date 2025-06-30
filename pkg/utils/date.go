package utils

import (
	"time"
)

func StartOfDay(t time.Time, timezone string) time.Time {
	location, _ := time.LoadLocation(timezone)
	year, month, day := t.In(location).Date()

	return time.Date(year, month, day, 0, 0, 0, 0, location)
}

func EndOfDay(t time.Time, timezone string) time.Time {
	location, _ := time.LoadLocation(timezone)
	year, month, day := t.In(location).Date()

	return time.Date(year, month, day, 23, 59, 59, 0, location)
}

func DiffInDays(start time.Time, end time.Time) int {
	return int(end.Sub(start).Hours() / 24)
	// days := end.Sub(start).Hours() / 24
	// return int(math.Round(days))
}

func IsSameDay(first time.Time, second time.Time) bool {
	return first.YearDay() == second.YearDay() && first.Year() == second.Year()
}

func IsLeapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}
