package timeParser

import "time"

func ParseToDate(date string) (time.Time, error) {
	const layout = "2006-01-02"
	return time.Parse(layout, date)
}
