package parser

import (
	"fmt"
	"time"

	"github.com/teambition/rrule-go"
)

func printTimeSlice(ts []time.Time) {
	for _, t := range ts {
		fmt.Println(t)
	}
}

/*
func ParseDatePattern(pattern string, now time.Time) (*rrule.RRule, error) {
	re := regexp.MustCompile(`^(\d+|\*)\/(\d+|\*)\/(\d+|\*)(?:\s(.+)\s|\s)(\d+|\*):(\d+|\*):(\d+|\*)$`)
	matches := re.FindStringSubmatch(pattern)
	if matches == nil || len(matches) != 8 {
		return nil, fmt.Errorf("Invalid date pattern: %s", pattern)
	}

	year := matches[1]
	month := matches[2]
	day := matches[3]
	weekdays := matches[4]
	hour := matches[5]
	minute := matches[6]
	second := matches[7]

	option := rrule.ROption{}

	if year != "*" {
		option.By
	}
	if month != "*" {
		option.
	}

	time.Date(1996, 05, 5, 9, 0, 0, 0, time.UTC)

	for i, match := range matches {
		fmt.Println(i, " => ", match)
	}

	r, _ := rrule.NewRRule(rrule.ROption{
		Freq:       rrule.YEARLY,
		Interval:   4,
		Count:      3,
		Bymonth:    []int{11},
		Byweekday:  []rrule.Weekday{rrule.TU},
		Bymonthday: []int{2, 3, 4, 5, 6, 7, 8},
		Dtstart:    time.Date(1996, 05, 5, 9, 0, 0, 0, time.UTC),
	})

}*/

func Toto() {

	r, _ := rrule.NewRRule(rrule.ROption{
		Freq:       rrule.YEARLY,
		Interval:   4,
		Count:      3,
		Bymonth:    []int{11},
		Byweekday:  []rrule.Weekday{rrule.TU},
		Bymonthday: []int{2, 3, 4, 5, 6, 7, 8},
		Dtstart:    time.Date(1996, 05, 5, 9, 0, 0, 0, time.UTC),
	})

	fmt.Println(r.Before(time.Now(), true))
	printTimeSlice(r.All())
}
