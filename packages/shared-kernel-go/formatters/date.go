package formatters

import "time"

const dateFormat = "02/01/2006"

func FormatDate(t time.Time) string {
	return t.Format(dateFormat)
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse(dateFormat, s)
}
