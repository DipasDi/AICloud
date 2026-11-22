package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ruMonthPrefix = map[string]time.Month{
	"янв": time.January,
	"фев": time.February,
	"мар": time.March,
	"апр": time.April,
	"май": time.May,
	"мая": time.May,
	"июн": time.June,
	"июл": time.July,
	"авг": time.August,
	"сен": time.September,
	"окт": time.October,
	"ноя": time.November,
	"дек": time.December,
}

var (
	reClock     = regexp.MustCompile(`(?:^|\s|в\s*)(\d{1,2}):(\d{2})(?:\s*|$|в)`)
	reClockHour = regexp.MustCompile(`(?:^|\s)в\s*(\d{1,2})(?:\s*|$|ч)`)
	reRuMonth   = regexp.MustCompile(`(\d{1,2})\s*(янв[а-я]*|фев[а-я]*|мар[а-я]*|апр[а-я]*|ма[яй][а-я]*|июн[а-я]*|июл[а-я]*|авг[а-я]*|сен[а-я]*|окт[а-я]*|ноя[а-я]*|дек[а-я]*)(?:\s*(\d{4}))?`)
	reDotDate   = regexp.MustCompile(`(\d{1,2})\.(\d{1,2})(?:\.(\d{4}))?`)
	reDayNumber = regexp.MustCompile(`(\d{1,2})\s*числ[а-я]*`)
)

func ParseRemindTime(input string) (time.Duration, time.Time, error) {
	s := strings.ToLower(strings.TrimSpace(input))
	now := time.Now()

	// 1. Относительные интервалы: 10m, 2h, 1d, через 15 минут, 2 часа
	cleanRel := strings.TrimPrefix(s, "через")
	cleanRel = strings.TrimSpace(cleanRel)
	cleanRel = strings.ReplaceAll(cleanRel, "минут", "m")
	cleanRel = strings.ReplaceAll(cleanRel, "мин", "m")
	cleanRel = strings.ReplaceAll(cleanRel, "часов", "h")
	cleanRel = strings.ReplaceAll(cleanRel, "часа", "h")
	cleanRel = strings.ReplaceAll(cleanRel, "час", "h")
	cleanRel = strings.ReplaceAll(cleanRel, "дней", "d")
	cleanRel = strings.ReplaceAll(cleanRel, "дня", "d")
	cleanRel = strings.ReplaceAll(cleanRel, "день", "d")
	cleanRel = strings.ReplaceAll(cleanRel, " ", "")

	if strings.HasSuffix(cleanRel, "d") {
		numStr := strings.TrimSuffix(cleanRel, "d")
		days, err := strconv.Atoi(numStr)
		if err == nil && days > 0 {
			dur := time.Duration(days) * 24 * time.Hour
			return dur, now.Add(dur), nil
		}
	}

	dur, err := time.ParseDuration(cleanRel)
	if err == nil && dur > 0 {
		return dur, now.Add(dur), nil
	}

	// 2. Извлекаем время (часы и минуты)
	hour := -1
	min := -1

	if m := reClock.FindStringSubmatch(s); len(m) == 3 {
		h, _ := strconv.Atoi(m[1])
		minute, _ := strconv.Atoi(m[2])
		if h >= 0 && h < 24 && minute >= 0 && minute < 60 {
			hour = h
			min = minute
		}
	} else if m := reClockHour.FindStringSubmatch(s); len(m) == 2 {
		h, _ := strconv.Atoi(m[1])
		if h >= 0 && h < 24 {
			hour = h
			min = 0
		}
	}

	defHour, defMin := 9, 0

	// 3. "завтра" / "послезавтра" / "сегодня"
	isPosleZavtra := strings.Contains(s, "послезавтра") || strings.Contains(s, "после завтра")
	isZavtra := !isPosleZavtra && strings.Contains(s, "завтра")
	isSegodnya := strings.Contains(s, "сегодня")

	if isPosleZavtra {
		h, m := defHour, defMin
		if hour != -1 {
			h, m = hour, min
		}
		target := time.Date(now.Year(), now.Month(), now.Day()+2, h, m, 0, 0, now.Location())
		return target.Sub(now), target, nil
	}

	if isZavtra {
		h, m := defHour, defMin
		if hour != -1 {
			h, m = hour, min
		}
		target := time.Date(now.Year(), now.Month(), now.Day()+1, h, m, 0, 0, now.Location())
		return target.Sub(now), target, nil
	}

	if isSegodnya {
		if hour == -1 {
			hour, min = 18, 0
		}
		target := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
		if target.Before(now) {
			return 0, time.Time{}, fmt.Errorf("время %02d:%02d сегодня уже прошло", hour, min)
		}
		return target.Sub(now), target, nil
	}

	// 4. "23 марта" / "23 марта 2026" / "23 марта в 19:00"
	if m := reRuMonth.FindStringSubmatch(s); len(m) >= 3 {
		day, _ := strconv.Atoi(m[1])
		monthPrefix := string([]rune(m[2]))
		if len([]rune(monthPrefix)) >= 3 {
			monthPrefix = string([]rune(monthPrefix)[:3])
		}
		month, ok := ruMonthPrefix[monthPrefix]
		if ok && day >= 1 && day <= 31 {
			year := now.Year()
			if len(m) == 4 && m[3] != "" {
				if y, err := strconv.Atoi(m[3]); err == nil {
					year = y
				}
			}

			h, minute := defHour, defMin
			if hour != -1 {
				h, minute = hour, min
			}

			target := time.Date(year, month, day, h, minute, 0, 0, now.Location())
			if len(m) < 4 || m[3] == "" {
				if target.Before(now) {
					target = time.Date(year+1, month, day, h, minute, 0, 0, now.Location())
				}
			}
			return target.Sub(now), target, nil
		}
	}

	// 5. "23.10" / "23.10.2026" / "23.10 в 18:00"
	if m := reDotDate.FindStringSubmatch(s); len(m) >= 3 {
		day, _ := strconv.Atoi(m[1])
		monthNum, _ := strconv.Atoi(m[2])
		if day >= 1 && day <= 31 && monthNum >= 1 && monthNum <= 12 {
			year := now.Year()
			if len(m) == 4 && m[3] != "" {
				if y, err := strconv.Atoi(m[3]); err == nil {
					year = y
				}
			}

			h, minute := defHour, defMin
			if hour != -1 {
				h, minute = hour, min
			}

			target := time.Date(year, time.Month(monthNum), day, h, minute, 0, 0, now.Location())
			if len(m) < 4 || m[3] == "" {
				if target.Before(now) {
					target = time.Date(year+1, time.Month(monthNum), day, h, minute, 0, 0, now.Location())
				}
			}
			return target.Sub(now), target, nil
		}
	}

	// 6. "20 числа" / "20 числа в 13:00"
	if m := reDayNumber.FindStringSubmatch(s); len(m) >= 2 {
		day, _ := strconv.Atoi(m[1])
		if day >= 1 && day <= 31 {
			h, minute := defHour, defMin
			if hour != -1 {
				h, minute = hour, min
			}

			target := time.Date(now.Year(), now.Month(), day, h, minute, 0, 0, now.Location())
			if target.Before(now) {
				// Переносим на следующий месяц
				nextMonth := now.AddDate(0, 1, 0)
				target = time.Date(nextMonth.Year(), nextMonth.Month(), day, h, minute, 0, 0, now.Location())
			}
			return target.Sub(now), target, nil
		}
	}

	// 7. Просто время "19:00" / "в 19:00" / "в 19"
	if hour != -1 {
		target := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
		if target.Before(now) {
			target = target.Add(24 * time.Hour)
		}
		return target.Sub(now), target, nil
	}

	return 0, time.Time{}, fmt.Errorf("не удалось распознать дату и время: %s", input)
}
