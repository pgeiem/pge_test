package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/teambition/rrule-go"
)

// RecurrentDate represents an interface for recurrent date operations.
type RecurrentDate interface {
	Next(now time.Time) (time.Time, error)
	Prev(now time.Time) (time.Time, error)
	Between(from, to time.Time) []time.Time
	String() string
}

var functionRegex = regexp.MustCompile(`^(\w+)\((.+)\)$`)

func ParseRecurrentDate(pattern string) (RecurrentDate, error) {
	recurrentDateTypes := map[string]func(string) (RecurrentDate, error){
		// Periodic reccurence
		"periodic": func(arg string) (RecurrentDate, error) {
			r := RecurrentDatePeriodic{}
			err := r.Parse(arg)
			return r, err
		},
		// Duration is an alias for periodic reccurence
		"duration": func(arg string) (RecurrentDate, error) {
			r := RecurrentDatePeriodic{}
			err := r.Parse(arg)
			return r, err
		},
		// Date pattern based reccurence
		"pattern": func(arg string) (RecurrentDate, error) {
			r := RecurrentDatePattern{}
			err := r.ParseFromDatePattern(arg)
			return r, err
		},
		// RRule RFC 5545 based reccurence
		"rrule": func(arg string) (RecurrentDate, error) {
			r := RecurrentDatePattern{}
			err := r.ParseFromRRule(arg)
			return r, err
		},
	}

	// Split function name and function arguments
	matches := functionRegex.FindStringSubmatch(pattern)
	if matches == nil || len(matches) != 3 {
		return nil, fmt.Errorf("error while parsing %s pattern, invalid pattern, expected form 'type(pattern)'", pattern)
	}

	// Create the recurrent date object
	createRecurrentDate, exists := recurrentDateTypes[matches[1]]
	if !exists {
		return nil, fmt.Errorf("error while parsing %s pattern, unknown type %s", pattern, matches[1])
	}
	return createRecurrentDate(matches[2])
}

// RecurrentDatePeriodic represents a periodic recurrent date.
type RecurrentDatePeriodic struct {
	Period Duration
}

func (r *RecurrentDatePeriodic) Parse(pattern string) error {
	duration, err := ParseDuration(pattern)
	if err != nil {
		return err
	}
	r.Period = duration
	return nil
}

// Next returns the next occurrence based on the current time.
func (r RecurrentDatePeriodic) Next(now time.Time) (time.Time, error) {
	return now.Add(r.Period.toDuration()), nil
}

// Prev returns the previous occurrence based on the current time.
func (r RecurrentDatePeriodic) Prev(now time.Time) (time.Time, error) {
	return now.Add(-r.Period.toDuration()), nil
}

// Between returns the time segments between the given time range.
func (r RecurrentDatePeriodic) Between(from, to time.Time) []time.Time {
	segments := []time.Time{}
	for t := from; t.Before(to); t = t.Add(r.Period.toDuration()) {
		segments = append(segments, t)
	}
	return segments
}

// Stringer for RecurrentDatePeriodic, print the period
func (r RecurrentDatePeriodic) String() string {
	return fmt.Sprintf("periodic(%s)", r.Period.toDuration().String())
}

// RecurrentDatePattern represents a pattern based recurrent date.
type RecurrentDatePattern struct {
	rule        *rrule.RRule
	origPattern string
}

func (r *RecurrentDatePattern) ParseFromDatePattern(pattern string) error {
	rule, err := builRRuleFromDatePattern(pattern)
	if err != nil {
		return fmt.Errorf("error while parsing %s rule pattern, %v", pattern, err)
	}
	rule.DTStart(time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)) //TODO: Start date must be before the current date to find the previous occurrence, see if any smarter thing can be done
	r.rule = rule
	r.origPattern = pattern
	return nil
}

func (r *RecurrentDatePattern) ParseFromRRule(pattern string) error {
	rule, err := rrule.StrToRRule(pattern)
	if err != nil {
		return fmt.Errorf("error while parsing %s rule pattern, %v", pattern, err)
	}
	rule.DTStart(time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)) //TODO: Start date must be before the current date to find the previous occurrence, see if any smarter thing can be done
	r.rule = rule
	r.origPattern = pattern
	return nil
}

// Next returns the next occurrence based on the current time.
func (r RecurrentDatePattern) Next(now time.Time) (time.Time, error) {
	//TODO: check if now is not too much in the past, before DTStart constant date
	next := r.rule.After(now, false)
	if next.IsZero() {
		return next, fmt.Errorf("no next occurrence found")
	}
	return next, nil
}

// Prev returns the previous occurrence based on the current time.
func (r RecurrentDatePattern) Prev(now time.Time) (time.Time, error) {
	//TODO: check if now is not too much in the past, before DTStart constant date
	prev := r.rule.Before(now, false)
	if prev.IsZero() {
		return prev, fmt.Errorf("no previous occurrence found")
	}
	return prev, nil
}

// Between returns the time segments between the given time range.
func (r RecurrentDatePattern) Between(from, to time.Time) []time.Time {
	return r.rule.Between(from, to, false)
}

// Stringer for RecurrentDatePattern, print the rule
func (r RecurrentDatePattern) String() string {
	return fmt.Sprintf("rrule(%s)", r.origPattern)
}

// Take a string describing a list or a range or a mix of both and return a list of integers representing the expanded list of values
func expandDateComponentList(pattern string) ([]int, error) {

	if pattern == "*" || pattern == "" {
		return []int{}, nil
	}

	convertToInt := func(component string) (int, error) {
		dayMap := map[string]int{
			"MON": 0, "TUE": 1, "WED": 2, "THU": 3, "FRI": 4, "SAT": 5, "SUN": 6,
		}
		monthMap := map[string]int{
			"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4, "MAY": 5, "JUN": 6,
			"JUL": 7, "AUG": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12,
		}

		day, exists := dayMap[strings.ToUpper(component)]
		if exists {
			return day, nil
		}
		month, exists := monthMap[strings.ToUpper(component)]
		if exists {
			return month, nil
		}
		number, err := strconv.Atoi(component)
		if err != nil {
			return 0, fmt.Errorf("error while parsing %s date component, invalid date, %v", component, err)
		}
		return number, nil
	}

	components := strings.Split(pattern, ",")
	var output []int
	for _, component := range components {
		if strings.Contains(component, "-") {
			limits := strings.Split(component, "-")
			if len(limits) != 2 {
				return nil, fmt.Errorf("error while parsing %s date component, invalid range, expected form 'start-end'", component)
			}
			start, err := convertToInt(limits[0])
			if err != nil {
				return nil, err
			}
			end, err := convertToInt(limits[1])
			if err != nil {
				return nil, err
			}
			if start > end || end-start > 35 {
				return nil, fmt.Errorf("error while parsing %s date component, invalid range, start is greater than end or range is too large", component)
			}
			for i := start; i <= end; i++ {
				output = append(output, i)
			}
		} else {
			day, err := convertToInt(component)
			if err != nil {
				return nil, err
			}
			output = append(output, day)
		}
	}
	return output, nil

}

// Regular expression to parse a date pattern in the form of "<yyyy/>mm/dd <weekdays> hh:mm<:ss> <extra>"
var rrule_regex = regexp.MustCompile(`^(?:([\d\-,*]+)\/)?([\d\-,*]+)\/([\d\-,*]+)\s+(?:([\w\-,*]*)\s+)?([\d\-,*]+):([\d\-,*]+)(?::([\d\-,*]*))?(?: (.*))?$`)

// builRRuleFromDatePattern takes a date pattern in the form of "<yyyy/>mm/dd <weekdays> hh:mm<:ss> <extra>" and returns a RRule object
func builRRuleFromDatePattern(pattern string) (*rrule.RRule, error) {

	//Parse the date pattern
	matches := rrule_regex.FindStringSubmatch(pattern)
	if matches == nil || len(matches) != 9 {
		return nil, fmt.Errorf("error while parsing %s pattern, invalid pattern, expected 'yyyy/mm/dd <weekdays> hh:mm:ss <extra>'", pattern)
	}
	extra_str := matches[8]

	// If year is not provided, set it to default "*"
	if matches[1] == "" {
		matches[1] = "*"
	}
	// Find the frequency looking for the first "*" field
	frequency := rrule.YEARLY
	frequencyList := []rrule.Frequency{rrule.YEARLY, rrule.MONTHLY, rrule.WEEKLY, rrule.DAILY, rrule.HOURLY, rrule.MINUTELY, rrule.SECONDLY}
	matchFrequencyList := []string{matches[1], matches[2], matches[4], matches[3], matches[5], matches[6], matches[7]}
	found := false
	for i, match := range matchFrequencyList {
		if match == "*" {
			frequency = frequencyList[i]
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("error while parsing '%s' pattern, no '*' field found to determine frequency", pattern)
	}

	rropt := &rrule.ROption{}
	if extra_str != "" {
		var err error
		extra_str = "FREQ=" + frequency.String() + ";" + extra_str //Add the frequency to the extra part as this is mendatory accordingly to RFC 5545
		rropt, err = rrule.StrToROption(extra_str)
		if err != nil {
			return nil, fmt.Errorf("error while parsing '%s' pattern, invalid extra part, %v", pattern, err)
		}
	}
	rropt.Freq = frequency
	rropt.Wkst = rrule.MO

	//Decode month
	byMonth, err := expandDateComponentList(matches[2])
	if err != nil {
		return nil, err
	}
	if len(byMonth) > 0 {
		rropt.Bymonth = byMonth
	}
	//Decode day
	byMonthDay, err := expandDateComponentList(matches[3])
	if err != nil {
		return nil, err
	}
	if len(byMonthDay) > 0 {
		rropt.Bymonthday = byMonthDay
	}
	//Decode weekday
	byWeekday, err := expandDateComponentList(matches[4])
	if err != nil {
		return nil, err
	}
	if len(byWeekday) > 0 {
		//Convert list of int to list of rrule.Weekday
		var weekdays []rrule.Weekday
		weekdayList := []rrule.Weekday{rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR, rrule.SA, rrule.SU}
		for _, day := range byWeekday {
			weekdays = append(weekdays, weekdayList[day])
		}
		rropt.Byweekday = weekdays
	}
	//Decode hour
	byHour, err := expandDateComponentList(matches[5])
	if err != nil {
		return nil, err
	}
	if len(byHour) > 0 {
		rropt.Byhour = byHour
	}
	//Decode minute
	byMinute, err := expandDateComponentList(matches[6])
	if err != nil {
		return nil, err
	}
	if len(byMinute) > 0 {
		rropt.Byminute = byMinute
	}
	//Decode second
	bySecond, err := expandDateComponentList(matches[7])
	if err != nil {
		return nil, err
	}
	if len(bySecond) > 0 {
		rropt.Bysecond = bySecond
	}

	//Create the RRule
	rrule, err := rrule.NewRRule(*rropt)
	if err != nil {
		return nil, err
	}

	return rrule, nil
}
