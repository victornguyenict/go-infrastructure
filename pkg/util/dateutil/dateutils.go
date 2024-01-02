package dateutil

import "time"

// AddDays adds n days to a given date.
func AddDays(t time.Time, n int) time.Time {
	return t.AddDate(0, 0, n)
}

// DayOfWeek returns the day of the week for a given date.
func DayOfWeek(t time.Time) string {
	return t.Weekday().String()
}

// FormatDate formats a date according to a specified format.
func FormatDate(t time.Time, format string) string {
	return t.Format(format)
}

// ParseDate parses a date string according to a specified format.
func ParseDate(dateStr, format string) (time.Time, error) {
	return time.Parse(format, dateStr)
}

// GetQuarter returns the quarter of the year for a given date.
func GetQuarter(t time.Time) int {
	month := t.Month()
	if month <= 3 {
		return 1
	} else if month <= 6 {
		return 2
	} else if month <= 9 {
		return 3
	} else {
		return 4
	}
}

// ListDays returns a slice of time.Time objects for each day between two dates.
func ListDays(start, end time.Time) []time.Time {
	var days []time.Time
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}
	return days
}

// ToFirstDayOfMonth sets the date to the first day of the current month.
func ToFirstDayOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// ToLastDayOfMonth sets the date to the last day of the current month.
func ToLastDayOfMonth(t time.Time) time.Time {
	return ToFirstDayOfMonth(t).AddDate(0, 1, -1)
}

// IsLeapYear checks if a year is a leap year.
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
