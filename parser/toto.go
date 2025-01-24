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
}

var functionRegex = regexp.MustCompile(`^(\w+)\((.+)\)$`)

func ParseRecurrentDate(pattern string) (RecurrentDate, error) {
	recurrentDateTypes := map[string]func(string) (RecurrentDate, error){
		"periodic": func(arg string) (RecurrentDate, error) {
			r := RecurrentDatePeriodic{}
			err := r.Parse(arg)
			return r, err
		},
		"pattern": func(arg string) (RecurrentDate, error) {
			r := RecurrentDatePattern{}
			err := r.ParseFromDatePattern(arg)
			return r, err
		},
		"rrule": func(arg string) (RecurrentDate, error) {
			r := RecurrentDatePattern{}
			err := r.ParseFromRRule(arg)
			return r, err
		},
	}

	// Split function name anf function arguments
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

// RecurrentDatePattern represents a pattern based recurrent date.
type RecurrentDatePattern struct {
	rule *rrule.RRule
}

func (r *RecurrentDatePattern) ParseFromDatePattern(pattern string) error {
	rule, err := BuilRRuleFromDatePattern(pattern)
	if err != nil {
		return fmt.Errorf("error while parsing %s rule pattern, %v", pattern, err)
	}
	rule.DTStart(time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)) //TODO: Start date must be before the current date to find the previous occurrence, see if any smarter thing can be done
	r.rule = rule
	return nil
}

func (r *RecurrentDatePattern) ParseFromRRule(pattern string) error {
	rule, err := rrule.StrToRRule(pattern)
	if err != nil {
		return fmt.Errorf("error while parsing %s rule pattern, %v", pattern, err)
	}
	rule.DTStart(time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)) //TODO: Start date must be before the current date to find the previous occurrence, see if any smarter thing can be done
	r.rule = rule
	return nil
}

// Next returns the next occurrence based on the current time.
func (r RecurrentDatePattern) Next(now time.Time) (time.Time, error) {
	//TODO: check if now is not too much in the past, before DTStart constant date
	fmt.Println("Next: ", r.rule.String(), now)
	next := r.rule.After(now, false)
	if next.IsZero() {
		return next, fmt.Errorf("no next occurrence found")
	}
	return next, nil
}

// Prev returns the previous occurrence based on the current time.
func (r RecurrentDatePattern) Prev(now time.Time) (time.Time, error) {
	//TODO: check if now is not too much in the past, before DTStart constant date
	fmt.Println("Prev: ", r.rule.String(), now)
	prev := r.rule.Before(now, false)
	if prev.IsZero() {
		return prev, fmt.Errorf("no previous occurrence found")
	}
	return prev, nil
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
		fmt.Println("handling component", component)

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

// BuilRRuleFromDatePattern takes a date pattern in the form of "<yyyy/>mm/dd <weekdays> hh:mm<:ss> <extra>" and returns a RRule object
func BuilRRuleFromDatePattern(pattern string) (*rrule.RRule, error) {

	//Parse the date pattern
	matches := rrule_regex.FindStringSubmatch(pattern)
	if matches == nil || len(matches) != 9 {
		return nil, fmt.Errorf("error while parsing %s pattern, invalid pattern, expected 'yyyy/mm/dd <weekdays> hh:mm:ss <extra>'", pattern)
	}
	extra_str := matches[8]

	fmt.Println("\nHandling pattern: ", pattern)
	fmt.Printf("%#v\n", matches)

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

	fmt.Println(rrule.String())

	return rrule, nil
}

/*
// RecurrentDateList represents a list of RecurrentDate.
type RecurrentDateList []RecurrentDate

func (rlist *RecurrentDateList) Add(pattern string) error {
	r, err := ParseRecurrentDate(pattern)
	if err != ninl {
		return err
	}
	liste = rlist.append(r)
}

// Next returns the next occurrence based on the current time.
func (rlist *RecurrentDateList) Next(now time.Time) time.Time {
	next := time.Date())
	for _, r := range rlist {
		n := r.Next(now)
		if n < next {
			next = n
		}
	}
	return nil
}

// Prev returns the previous occurrence based on the current time.
func (r *RecurrentDateList) Prev(now time.Time) time.Time {
}
*/
/*
// RecurentSegment represents a segment of recurrent dates.
type RecurentSegment struct {
	Recurences RecurrentDateList
	Duration   time.Duration
}

// Between returns the time segments between the given time range.
func (rs *RecurentSegment) Between(from, to time.Time) TimeSegments {
	// unroll the recuruences list and generate the segments
	// do not forget to include one prev if needed
}

// Check if t is included in the recurent segment, this means that t is between the start and the end of one of the segment generated by the recurent date
// This is used for example to enable conditionally some rules depending on date included in a rule
func (rs *RecurentSegment) IsWithin(t time.Time) bool {
}
*/

// Duration represents a duration with second resolution.
type Duration time.Duration

// toDuration converts Duration to time.Duration.
func (d Duration) toDuration() time.Duration {
	return time.Duration(d)
}

// ParseDuration parses a duration string with units: seconds (s), minutes (m), hours (h), days (d), weeks (w).
func ParseDuration(s string) (Duration, error) {
	var totalSeconds int64
	var multipliers = map[byte]int64{
		'w': 7 * 24 * 60 * 60,
		'd': 24 * 60 * 60,
		'h': 60 * 60,
		'm': 60,
		's': 1,
	}

	matches := regexp.MustCompile(`^(\d+w)?(\d+d)?(\d+h)?(\d+m)?(\d+s)?$`).FindStringSubmatch(s)
	if matches == nil || s == "" {
		return 0, fmt.Errorf("error while parsing %s duration, invalid pattern", s)
	}

	for _, match := range matches[1:] {
		if match == "" {
			continue
		}

		num, err := strconv.ParseInt(match[:len(match)-1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("error while parsing %s duration, invalid number, %s", s, err)
		}

		unit := match[len(match)-1]
		multiplier, exists := multipliers[unit]
		if !exists {
			return 0, fmt.Errorf("error while parsing %s duration, invalid unit (%c)", s, unit)
		}

		totalSeconds += num * multiplier
	}

	return Duration(totalSeconds * int64(time.Second)), nil
}
