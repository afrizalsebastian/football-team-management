package helper

import "time"

func MatchDateTimeInPast(dateStr string, timeStr string) bool {
	parsedDate, dateErr := time.Parse("02-01-2006", dateStr)
	parsedTime, timeErr := time.Parse("15:04", timeStr)
	if dateErr != nil || timeErr != nil {
		return true
	}

	combinedDateTime := time.Date(
		parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
		parsedTime.Hour(), parsedTime.Minute(), 0, 0,
		time.Local,
	)

	if combinedDateTime.Before(time.Now()) {
		return true
	}

	return false
}
