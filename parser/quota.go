package parser

import (
	"path/filepath"
	"strings"
	"time"
)

// DurationType represents the different type of parking duration
type DurationType string

const (
	FreeDuration      DurationType = "free"
	NonPayingDuration DurationType = "nonpaying"
	PayingDuration    DurationType = "paying"
	BannedDuration    DurationType = "banned"
)

// DurationDetail represents the details of a parking duration
type DurationDetail struct {
	Type     DurationType
	Start    time.Time
	Duration time.Duration
}

// AssignedRight represents the parking assigned rights (a ticket)
type AssignedRight struct {
	Area    string
	Start   time.Time
	End     time.Time
	Details []DurationDetail
}

// MatchingRule represents a rule to match the parking assigned rights to be used in a quota
type MatchingRule struct {
	AreaPattern string
	TypePattern string
}

// Quota represents a quota to be used to limit the parking assigned rights
type Quota interface {
	Update(now time.Time, history []AssignedRight) error
}

// AbstractQuota is a helper to ease the implementation of different quotas types
type AbstractQuota struct {
	MatchingRules   []MatchingRule
	PeriodicityRule RecurrentDate
}

// SelectReferenceTime selects the reference time to be used to filter the assigned rights based on the matching rules
func SelectReferenceTime(rule MatchingRule, detail DurationDetail, right AssignedRight) time.Time {
	reftime := detail.Start
	if reftime.IsZero() {
		reftime = right.Start
	}
	return reftime
}

const DefaultAreaPattern = "*"
const DefaultTypePattern = string(FreeDuration)

// Helper function to match a glob string pattern, in a case-insensitive way
func globMatch(pattern, name string) (bool, error) {
	return filepath.Match(strings.ToLower(pattern), strings.ToLower(name))
}

// Filter filters the history of assigned rights based on the matching rules and calls the matchHandler for each matching detail
func (q AbstractQuota) Filter(from time.Time, history []AssignedRight, matchHandler func(detail DurationDetail)) error {
	rules := q.MatchingRules
	if len(rules) == 0 {
		rules = []MatchingRule{{}}
	}
	for _, rule := range rules {
		for _, right := range history {
			areaPattern := rule.AreaPattern
			if areaPattern == "" {
				areaPattern = DefaultAreaPattern
			}
			match, err := globMatch(areaPattern, right.Area)
			if err != nil {
				return err
			}
			if match {
				for _, detail := range right.Details {
					typePattern := rule.TypePattern
					if typePattern == "" {
						typePattern = DefaultTypePattern
					}
					match, err := globMatch(typePattern, string(detail.Type))
					if err != nil {
						return err
					}
					if match {
						reftime := SelectReferenceTime(rule, detail, right)
						if !reftime.IsZero() && TimeAfterOrEqual(reftime, from) {
							matchHandler(detail)
						}
					}
				}
			}
		}
	}
	return nil
}

func (q AbstractQuota) PeriodStart(now time.Time) (time.Time, error) {
	return q.PeriodicityRule.Prev(now)
}

// DurationQuota represents a quota based on the duration of the parking assigned rights
type DurationQuota struct {
	AbstractQuota
	Duration time.Duration
}

// Update updates the quota based on the history of assigned rights
func (q *DurationQuota) Update(now time.Time, history []AssignedRight) error {
	total := time.Duration(0)
	start, err := q.PeriodStart(now)
	if err != nil {
		return err
	}
	err = q.Filter(start, history, func(detail DurationDetail) {
		total += detail.Duration
	})
	q.Duration = total
	return err
}

// CounterQuota represents a quota based on the number of parking assigned rights
type CounterQuota struct {
	AbstractQuota
	Counter int
}

// Update updates the quota based on the history of assigned rights
func (q *CounterQuota) Update(now time.Time, history []AssignedRight) error {
	counter := 0
	start, err := q.PeriodStart(now)
	if err != nil {
		return err
	}
	err = q.Filter(start, history, func(detail DurationDetail) {
		counter++
	})
	q.Counter = counter
	return err
}
