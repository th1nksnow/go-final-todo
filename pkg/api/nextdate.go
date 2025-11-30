package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Валидация параметров
	if date == "" {
		http.Error(w, "date parameter is required", http.StatusBadRequest)
		return
	}

	if repeat == "" {
		http.Error(w, "repeat parameter is required", http.StatusBadRequest)
		return
	}

	var now time.Time
	var err error

	if nowStr != "" {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid now parameter: %v", err), http.StatusBadRequest)
			return
		}
	} else {
		now = time.Now()
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule cannot be empty")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("failed to parse date: %v", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", errors.New("invalid repeat format")
	}

	ruleType := parts[0]

	switch ruleType {
	case "d":
		return nextDateDaily(now, date, parts)
	case "y":
		return nextDateYearly(now, date)
	case "w":
		return nextDateWeekly(now, date, parts)
	case "m":
		return nextDateMonthly(now, date, parts)
	default:
		return "", fmt.Errorf("unsupported repeat format: %s", ruleType)
	}
}

// nextDateDaily обрабатывает правило "d <дни>"
func nextDateDaily(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("invalid daily format: expected 'd <days>'")
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid days value: %v", err)
	}

	if days <= 0 || days > 400 {
		return "", errors.New("days must be between 1 and 400")
	}

	result := date
	for {
		result = result.AddDate(0, 0, days)
		if afterNow(result, now) {
			break
		}
	}

	return result.Format(DateFormat), nil
}

// nextDateYearly обрабатывает правило "y"
func nextDateYearly(now, date time.Time) (string, error) {
	result := date
	for {
		result = result.AddDate(1, 0, 0)
		if afterNow(result, now) {
			break
		}
	}
	return result.Format(DateFormat), nil
}

// nextDateWeekly обрабатывает правило "w <дни_недели>"
func nextDateWeekly(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("invalid weekly format: expected 'w <days>'")
	}

	dayStrs := strings.Split(parts[1], ",")
	weekdays := make(map[int]bool)

	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			return "", fmt.Errorf("invalid weekday value: %v", err)
		}
		if day < 1 || day > 7 {
			return "", errors.New("weekday must be between 1 and 7")
		}
		weekdays[day] = true
	}

	if len(weekdays) == 0 {
		return "", errors.New("no valid weekdays specified")
	}

	result := date
	for {
		result = result.AddDate(0, 0, 1)
		weekday := int(result.Weekday()) // A Weekday specifies a day of the week (Sunday = 0, ...)
		if weekday == 0 {
			weekday = 7 // Воскресенье = 7
		}

		if weekdays[weekday] && afterNow(result, now) {
			break
		}
	}

	return result.Format(DateFormat), nil
}

// nextDateMonthly обрабатывает правило "m <дни> [месяцы]"
func nextDateMonthly(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("invalid monthly format: expected 'm <days> [months]'")
	}

	dayStrs := strings.Split(parts[1], ",")
	days := make([]int, 0, len(dayStrs))

	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			return "", fmt.Errorf("invalid day value: %v", err)
		}
		if day < -2 || day == 0 || day > 31 {
			return "", errors.New("day must be between 1-31, -1 or -2")
		}
		days = append(days, day)
	}

	// Парсим месяцы (optional)
	months := make(map[int]bool)
	if len(parts) > 2 {
		monthStrs := strings.Split(parts[2], ",")
		for _, monthStr := range monthStrs {
			month, err := strconv.Atoi(strings.TrimSpace(monthStr))
			if err != nil {
				return "", fmt.Errorf("invalid month value: %v", err)
			}
			if month < 1 || month > 12 {
				return "", errors.New("month must be between 1 and 12")
			}
			months[month] = true
		}
	}

	result := date
	for {
		result = result.AddDate(0, 0, 1)

		// Проверяем месяц
		if len(months) > 0 {
			currentMonth := int(result.Month())
			if !months[currentMonth] {
				continue
			}
		}

		if isValidDayInMonth(result, days) && afterNow(result, now) {
			break
		}
	}

	return result.Format(DateFormat), nil
}

// isValidDayInMonth проверяет, подходит ли день месяца
func isValidDayInMonth(date time.Time, days []int) bool {
	for _, day := range days {
		switch day {
		case -1:
			lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
			if date.Day() == lastDay {
				return true
			}
		case -2:
			lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
			if date.Day() == lastDay-1 {
				return true
			}
		default:
			if date.Day() == day {
				return true
			}
		}
	}
	return false
}

// afterNow возвращает true, если date больше now
func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}
