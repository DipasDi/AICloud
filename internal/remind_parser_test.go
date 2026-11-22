package internal

import (
	"testing"
	"time"
)

func TestParseRemindTime(t *testing.T) {
	testCases := []struct {
		input       string
		expectError bool
		check       func(t *testing.T, dur time.Duration, target time.Time)
	}{
		{
			input: "23 марта в 19:00",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if target.Month() != time.March || target.Day() != 23 || target.Hour() != 19 || target.Minute() != 0 {
					t.Errorf("expected March 23 19:00, got %v", target)
				}
			},
		},
		{
			input: "23 марта",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if target.Month() != time.March || target.Day() != 23 || target.Hour() != 9 || target.Minute() != 0 {
					t.Errorf("expected March 23 09:00, got %v", target)
				}
			},
		},
		{
			input: "23.10.2026",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if target.Year() != 2026 || target.Month() != time.October || target.Day() != 23 || target.Hour() != 9 {
					t.Errorf("expected 2026-10-23 09:00, got %v", target)
				}
			},
		},
		{
			input: "23.10",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if target.Month() != time.October || target.Day() != 23 || target.Hour() != 9 {
					t.Errorf("expected October 23 09:00, got %v", target)
				}
			},
		},
		{
			input: "23.10 в 18:30",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if target.Month() != time.October || target.Day() != 23 || target.Hour() != 18 || target.Minute() != 30 {
					t.Errorf("expected October 23 18:30, got %v", target)
				}
			},
		},
		{
			input: "Завтра",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				now := time.Now()
				tomorrow := now.AddDate(0, 0, 1)
				if target.Day() != tomorrow.Day() || target.Hour() != 9 {
					t.Errorf("expected tomorrow 09:00, got %v", target)
				}
			},
		},
		{
			input: "После завтра",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				now := time.Now()
				afterTomorrow := now.AddDate(0, 0, 2)
				if target.Day() != afterTomorrow.Day() || target.Hour() != 9 {
					t.Errorf("expected day after tomorrow 09:00, got %v", target)
				}
			},
		},
		{
			input: "послезавтра в 14:00",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				now := time.Now()
				afterTomorrow := now.AddDate(0, 0, 2)
				if target.Day() != afterTomorrow.Day() || target.Hour() != 14 || target.Minute() != 0 {
					t.Errorf("expected day after tomorrow 14:00, got %v", target)
				}
			},
		},
		{
			input: "Завтра в 15:30",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				now := time.Now()
				tomorrow := now.AddDate(0, 0, 1)
				if target.Day() != tomorrow.Day() || target.Hour() != 15 || target.Minute() != 30 {
					t.Errorf("expected tomorrow 15:30, got %v", target)
				}
			},
		},
		{
			input: "20 числа в 13:00",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if target.Day() != 20 || target.Hour() != 13 || target.Minute() != 0 {
					t.Errorf("expected 20th at 13:00, got %v", target)
				}
			},
		},
		{
			input: "20 числа",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if target.Day() != 20 || target.Hour() != 9 || target.Minute() != 0 {
					t.Errorf("expected 20th at 09:00, got %v", target)
				}
			},
		},
		{
			input: "в 19:00",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if target.Hour() != 19 || target.Minute() != 0 {
					t.Errorf("expected 19:00, got %v", target)
				}
			},
		},
		{
			input: "в 19",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if target.Hour() != 19 || target.Minute() != 0 {
					t.Errorf("expected 19:00, got %v", target)
				}
			},
		},
		{
			input: "15m",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if dur < 14*time.Minute || dur > 16*time.Minute {
					t.Errorf("expected ~15m, got %v", dur)
				}
			},
		},
		{
			input: "через 2 часа",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if dur < 119*time.Minute || dur > 121*time.Minute {
					t.Errorf("expected ~2h, got %v", dur)
				}
			},
		},
		{
			input: "через 3 дня",
			check: func(t *testing.T, dur time.Duration, target time.Time) {
				if dur < 71*time.Hour || dur > 73*time.Hour {
					t.Errorf("expected ~72h, got %v", dur)
				}
			},
		},
	}

	for _, tc := range testCases {
		dur, target, err := ParseRemindTime(tc.input)
		if tc.expectError {
			if err == nil {
				t.Errorf("input '%s': expected error, got nil", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("input '%s': unexpected error: %v", tc.input, err)
			continue
		}
		if tc.check != nil {
			tc.check(t, dur, target)
		}
	}
}
