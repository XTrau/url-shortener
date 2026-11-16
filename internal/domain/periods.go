package domain

import (
	"errors"
	"slices"
	"strconv"
	"time"
	"unicode/utf8"
)

type Interval struct {
	Start time.Time
	End   time.Time
}

func (i Interval) Label() string {
	return i.Start.Format(time.RFC3339)
}

type Period interface {
	Intervals() []Interval
}

type HourPeriod struct {
	Count int
}

func (p HourPeriod) Intervals() []Interval {
	result := make([]Interval, 0, p.Count)
	now := time.Now().UTC()
	d := now.Truncate(time.Hour)
	step := time.Hour
	for range p.Count {
		result = append(result, Interval{
			Start: d,
			End:   d.Add(step),
		})
		d = d.Add(-step)
	}
	return result
}

type DayPeriod struct {
	Count int
}

func (p DayPeriod) Intervals() []Interval {
	result := make([]Interval, 0, p.Count)
	now := time.Now().UTC()
	year, month, day := now.Date()
	d := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	step := time.Hour * 24
	for range p.Count {
		result = append(result, Interval{
			Start: d,
			End:   d.Add(step),
		})
		d = d.Add(-step)
	}
	return result
}

type MonthPeriod struct {
	Count int
}

func (p MonthPeriod) Intervals() []Interval {
	result := make([]Interval, 0, p.Count)
	now := time.Now().UTC()
	year, month, _ := now.Date()

	cur := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	prev := now

	for range p.Count {
		result = append(result, Interval{
			Start: cur,
			End:   prev,
		})

		month--
		if month <= 0 {
			month = 12
			year--
		}

		prev = cur
		cur = time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	}
	return result
}

func ParsePeriod(s string) (Period, error) {
	len := utf8.RuneCountInString(s)
	if len < 2 || len > 3 {
		return nil, errors.New("Incorrect period format")
	}

	goodSymbols := []string{"h", "d", "m"}
	lastSymbol := s[len-1 : len]

	if !slices.Contains(goodSymbols, lastSymbol) {
		return nil, errors.New("Incorrect period format")
	}

	periodNumStr := s[:len-1]
	periodNum, err := strconv.Atoi(periodNumStr)
	if err != nil {
		return nil, errors.New("Incorrect period format")
	}

	switch lastSymbol {
	case "h":
		if periodNum < 0 || periodNum > 24 {
			return nil, errors.New("Incorrect period format: for hours - min 1, max 24")
		}

		return HourPeriod{Count: periodNum}, nil
	case "d":
		if periodNum < 0 || periodNum > 30 {
			return nil, errors.New("Incorrect period format: for days - min 1, max 30")
		}

		return DayPeriod{Count: periodNum}, nil
	case "m":
		if periodNum < 0 || periodNum > 30 {
			return nil, errors.New("Incorrect period format: for month - min 1, max 12")
		}

		return MonthPeriod{Count: periodNum}, nil
	}

	return nil, errors.New("Incorrect period format")
}
