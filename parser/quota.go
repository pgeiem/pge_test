package parser

import (
	"fmt"
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
	ParkingArea []string
	Start       time.Time
	//End         time.Time
	Details []DurationDetail
}

func (ar AssignedRight) MatchParkingArea(pattern string) (bool, error) {
	for _, area := range ar.ParkingArea {
		match, err := globMatch(pattern, area)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

// MatchingRule represents a rule to match the parking assigned rights to be used in a quota
type MatchingRule struct {
	ParkingAreaPattern  string `yaml:"area"`
	DurationTypePattern string `yaml:"type"`
}

// Stringer for MatchingRule, print the area and type patterns
func (m MatchingRule) String() string {
	return fmt.Sprintf("(%s, %s)", m.ParkingAreaPattern, m.DurationTypePattern)
}

// MatchingRules is a list of MatchingRule
type MatchingRules []MatchingRule

// Stringer for MatchingRules, iterate over all rules and print them
func (m MatchingRules) String() string {
	var str strings.Builder
	str.WriteString("[")
	for _, rule := range m {
		str.WriteString(rule.String())
		str.WriteString(" ")
	}
	str.WriteString("]")
	return str.String()
}

// Quota represents a quota to be used to limit the parking assigned rights
type Quota interface {
	GetName() string
	Update(now time.Time, history []AssignedRight) error
}

// AbstractQuota is a helper to ease the implementation of different quotas types
type AbstractQuota struct {
	Name               string        `yaml:"name" validate:"required"`
	MatchingRules      MatchingRules `yaml:"matching"`
	PeriodicityRule    RecurrentDate `yaml:"periodicity" validate:"required"`
	DefaultAreaPattern string        `yaml:"-"`
	DefaultTypePattern string        `yaml:"-"`
}

func (q AbstractQuota) GetName() string {
	return q.Name
}

// SelectReferenceTime selects the reference time to be used to filter the assigned rights based on the matching rules
func SelectReferenceTime(rule MatchingRule, detail DurationDetail, right AssignedRight) time.Time {
	reftime := detail.Start
	if reftime.IsZero() {
		reftime = right.Start
	}
	return reftime
}

// Helper function to match a glob string pattern, in a case-insensitive way
func globMatch(pattern, name string) (bool, error) {
	return filepath.Match(strings.ToLower(pattern), strings.ToLower(name))
}

// Filter filters the history of assigned rights based on the matching rules and calls the matchHandler for each matching detail
func (q AbstractQuota) Filter(from time.Time, history []AssignedRight, matchAssignedRightHandler func(right AssignedRight),
	matchDurationDetailsHandler func(detail DurationDetail)) error {
	rules := q.MatchingRules
	if len(rules) == 0 {
		rules = []MatchingRule{{}}
	}
	// Iterate over all matching rules of the quota
	for _, rule := range rules {
		// Iterate over all assigned rights in the history
		for _, right := range history {
			areaPattern := rule.ParkingAreaPattern
			if areaPattern == "" {
				areaPattern = q.DefaultAreaPattern
			}
			match, err := right.MatchParkingArea(areaPattern)
			if err != nil {
				return err
			}
			// If set, call the Assigned Right callback
			if match && matchAssignedRightHandler != nil {
				matchAssignedRightHandler(right)
			}
			// If set, check duration details matches and call the Duration Detail callback
			if match && matchDurationDetailsHandler != nil {
				typePattern := rule.DurationTypePattern
				if typePattern == "" {
					typePattern = q.DefaultTypePattern
				}
				for _, detail := range right.Details {
					match, err := globMatch(typePattern, string(detail.Type))
					if err != nil {
						return err
					}
					if match {
						reftime := SelectReferenceTime(rule, detail, right)
						if !reftime.IsZero() && TimeAfterOrEqual(reftime, from) {
							matchDurationDetailsHandler(detail)
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

// Stringer for AbstractQuota, print the matching rule and periodicity rule
func (q AbstractQuota) String() string {
	return fmt.Sprintf("PeriodicityRule: %v, MatchingRules: %v", q.PeriodicityRule, q.MatchingRules)
}

// DurationQuota represents a quota based on the duration of the parking assigned rights
type DurationQuota struct {
	AbstractQuota `yaml:",inline"`
	Allowance     time.Duration `yaml:"allowance"`
	used          time.Duration
}

func NewDurationQuota(allowance time.Duration, period RecurrentDate, rules []MatchingRule) *DurationQuota {
	return &DurationQuota{
		AbstractQuota: AbstractQuota{
			MatchingRules:      rules,
			PeriodicityRule:    period,
			DefaultAreaPattern: "*",
			DefaultTypePattern: string(FreeDuration),
		},
		Allowance: allowance,
	}
}

// Update updates the quota based on the history of assigned rights
func (q *DurationQuota) Update(now time.Time, history []AssignedRight) error {
	total := time.Duration(0)
	// Compute the start period of quota calculation
	start, err := q.PeriodStart(now)
	if err != nil {
		return err
	}
	// Compute the total duration of matching assigned rights
	err = q.Filter(start, history, nil, func(detail DurationDetail) {
		total += detail.Duration
	})
	q.used = total
	return err
}

func (q *DurationQuota) Available() time.Duration {
	available := time.Duration(0)
	if q.Allowance > q.used {
		available = q.Allowance - q.used
	}
	return available
}

func (q *DurationQuota) Used() time.Duration {
	return q.used
}

// Stringer for DurationQuota, print the name and the used/allowed values
func (q DurationQuota) String() string {
	return fmt.Sprintf("DurationQuota(%s): Usage  %s/%s %v", q.Name, q.used, q.Allowance, q.AbstractQuota)
}

// CounterQuota represents a quota based on the number of parking assigned rights
type CounterQuota struct {
	AbstractQuota `yaml:",inline"`
	Allowance     int `yaml:"allowance"`
	used          int
}

func NewCounterQuota(allowance int, period RecurrentDate, rules []MatchingRule) *CounterQuota {
	return &CounterQuota{
		AbstractQuota: AbstractQuota{
			MatchingRules:      rules,
			PeriodicityRule:    period,
			DefaultAreaPattern: "*",
		},
		Allowance: allowance,
	}
}

// Update updates the quota based on the history of assigned rights
func (q *CounterQuota) Update(now time.Time, history []AssignedRight) error {
	counter := 0
	// Compute the start period of quota calculation
	start, err := q.PeriodStart(now)
	if err != nil {
		return err
	}
	// Compute the number of matching assigned rights
	err = q.Filter(start, history, func(detail AssignedRight) {
		counter++
	}, nil)
	q.used = counter
	return err
}

func (q *CounterQuota) Available() int {
	var available int
	if q.Allowance > q.used {
		available = q.Allowance - q.used
	}
	return available
}

func (q *CounterQuota) Used() int {
	return q.used
}

// Stringer for CounterQuota, print the name and the used/allowed values
func (q CounterQuota) String() string {
	return fmt.Sprintf("CounterQuota(%s): Usage %d/%d %v", q.Name, q.used, q.Allowance, q.AbstractQuota)
}

type QuotaInventory map[string]Quota

/*struct {
	*orderedmap.OrderedMap[string, Quota]
}*/

/*func (qi QuotaInventory) Update(now time.Time, history []AssignedRight) error {
	for quota := range qi.Values() {
		err := quota.Update(now, history)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qi QuotaInventory) GetDurationQuota(name string) (*DurationQuota, bool) {
	quota, ok := qi.Get(name)
	if !ok {
		return nil, false
	}
	dq, ok := quota.(*DurationQuota)
	return dq, ok
}

func (qi QuotaInventory) GetCounterQuota(name string) (*CounterQuota, bool) {
	quota, ok := qi.Get(name)
	if !ok {
		return nil, false
	}
	cq, ok := quota.(*CounterQuota)
	return cq, ok
}

// Stringer for QuotaInventory, iterate over all quotas and print some details
func (qi QuotaInventory) String() string {
	var str strings.Builder
	for key, quota := range qi.AllFromFront() {
		str.WriteString(fmt.Sprintf("%s: %s\n", key, quota))
	}
	return str.String()
}*/

/*
func (qi *QuotaInventory) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Temporary struct to unmarshal the different types of quotas
	temp := []struct {
		DurationQuota *DurationQuota `yaml:"duration"`
		CounterQuota  *CounterQuota  `yaml:"counter"`
	}{}

	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	// Convert from the temporary struct to the QuotaInventory
	qi.OrderedMap = orderedmap.NewOrderedMapWithCapacity[string, Quota](len(temp))
	var quota Quota
	for _, t := range temp {
		quota = nil
		if t.DurationQuota != nil {
			quota = t.DurationQuota
		} else if t.CounterQuota != nil {
			quota = t.CounterQuota
		}
		if quota != nil {
			if quota.GetName() == "" {
				return fmt.Errorf("missing quota name")
			}
			qi.OrderedMap.Set(quota.GetName(), quota)
		}
	}
	return nil
}
*/
/*
func (qi *QuotaInventory) UnmarshalYAML(data []byte) error {
	temp := map[string]struct {
		Type string `yaml:"type"`
	}{}

	err := yaml.Unmarshal(data, &temp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal quota type: %w", err)
	}

	fmt.Println("unmarshalQuotaInventory", temp, string(data))

	*qi = make(QuotaInventory)
	for name, t := range temp {
		//var q map[string]Quota
		switch t.Type {
		case "duration":
			var q DurationQuota
			if err := yaml.UnmarshalWithOptions(data, &q, yaml.Strict()); err != nil {
				return fmt.Errorf("failed to decode duration quota %q: %w", name, err)
			}

			(*qi)[name] = &q
		case "counter":
			var q CounterQuota
			if err := yaml.UnmarshalWithOptions(data, &q, yaml.Strict()); err != nil {
				return fmt.Errorf("failed to decode counter quota %q: %w", name, err)
			}

			(*qi)[name] = &q
		default:
			return fmt.Errorf("unknown quota type: %s", t.Type)
		}

	}
	/*var q Quota
	switch temp.Type {
	case "duration":
		q = &DurationQuota{}
	case "counter":
		q = &CounterQuota{}
	default:
		return fmt.Errorf("unknown quota type: %s", temp.Type)
	}

	if err := yaml.UnmarshalWithOptions(data, q, yaml.Strict()); err != nil {
		return err
	}

	*quota = q
	return nil
}
*/
