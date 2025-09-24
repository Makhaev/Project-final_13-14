package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	const maxDayInterval = 400
	const maxSearchDays = 365 * 5

	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат dstart: %v", err)
	}

	if repeat == "" {
		return "", fmt.Errorf("repeat не указан — задача одноразовая")
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("некорректное правило d")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n <= 0 || n > maxDayInterval {
			return "", fmt.Errorf("некорректный интервал: должен быть 1..%d", maxDayInterval)
		}

		start := now.AddDate(0, 0, 1)
		if !start.After(startDate) {
			start = startDate.AddDate(0, 0, 1)
		}
		daysSinceStart := int(start.Sub(startDate).Hours() / 24)
		steps := (daysSinceStart + n - 1) / n
		nextDate := startDate.AddDate(0, 0, steps*n)
		return nextDate.Format(DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("некорректное правило y")
		}
		start := now.AddDate(0, 0, 1)
		if !start.After(startDate) {
			start = startDate.AddDate(0, 0, 1)
		}
		date := startDate
		for date.Before(start) || date.Equal(start) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(DateFormat), nil

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("некорректное правило m")
		}

		dayParts := strings.Split(parts[1], ",")
		var days []int
		for _, dp := range dayParts {
			day, err := strconv.Atoi(dp)
			if err != nil || day == 0 || day < -31 || day > 31 {
				return "", fmt.Errorf("некорректный день месяца: %s", dp)
			}
			days = append(days, day)
		}

		var months map[time.Month]bool
		if len(parts) == 3 {
			months = make(map[time.Month]bool)
			for _, mp := range strings.Split(parts[2], ",") {
				m, err := strconv.Atoi(mp)
				if err != nil || m < 1 || m > 12 {
					return "", fmt.Errorf("некорректный месяц: %s", mp)
				}
				months[time.Month(m)] = true
			}
		}

		start := now.AddDate(0, 0, 1)
		if !start.After(startDate) {
			start = startDate.AddDate(0, 0, 1)
		}

		for i, d := 0, start; i < maxSearchDays; d, i = d.AddDate(0, 0, 1), i+1 {
			if months != nil && !months[d.Month()] {
				continue
			}
			for _, day := range days {
				checkDay := day
				if day < 0 {
					lastDay := time.Date(d.Year(), d.Month()+1, 0, 0, 0, 0, 0, d.Location()).Day()
					checkDay = lastDay + day + 1
				}
				if checkDay < 1 || checkDay > 31 {
					continue // недопустимая дата
				}
				if d.Day() == checkDay {
					return d.Format(DateFormat), nil
				}
			}
		}
		return "", nil

	case "w":
		if len(parts) < 2 {
			return "", fmt.Errorf("некорректное правило w")
		}
		var weekdays []time.Weekday
		for _, part := range strings.Split(parts[1], ",") {
			n, err := strconv.Atoi(part)
			if err != nil || n < 1 || n > 7 {
				return "", fmt.Errorf("некорректный день недели: %s", part)
			}
			weekdays = append(weekdays, time.Weekday(n%7))
		}

		start := now.AddDate(0, 0, 1)
		if !start.After(startDate) {
			start = startDate.AddDate(0, 0, 1)
		}

		for i, d := 0, start; i < maxSearchDays; d, i = d.AddDate(0, 0, 1), i+1 {
			for _, w := range weekdays {
				if d.Weekday() == w {
					return d.Format(DateFormat), nil
				}
			}
		}
		return "", nil

	default:
		return "", fmt.Errorf("неподдерживаемое правило: %s", rule)
	}
}
