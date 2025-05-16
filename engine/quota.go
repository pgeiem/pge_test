package engine

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/iem-rd/quote-engine/timeutils"
)

// DurationDetail represents the details of a parking duration
type DurationDetail struct {
	Type     DurationType  `yaml:"type"`     // Type of the duration (Free, Paid, etc.)
	Start    time.Time     `yaml:"start"`    // Start date of the duration (used for advanced quota types)
	Duration time.Duration `yaml:"duration"` // Duration of the parking
}

// AssignedRight represents the parking assigned rights (a ticket)
type AssignedRight struct {
	TariffCode      string           `yaml:"tariffCode"`      // Identifier of the tariff
	Flags           []string         `yaml:"flags"`           // list of flags of parking assigned rights (such as PMR, etc.)
	LayerCode       string           `yaml:"layerCode"`       // Zone code
	LayerCodes      []string         `yaml:"layerCodes"`      // Zone codes
	StartDate       time.Time        `yaml:"startDate"`       // Start date of the parking assigned right
	DurationDetails []DurationDetail `yaml:"durationDetails"` // List of duration details
}

type AssignedRights []AssignedRight

func (ar AssignedRight) MatchLayerCode(pattern string) (bool, error) {
	if len(ar.LayerCodes) == 0 {
		return globMatch(pattern, ar.LayerCode)
	}
	for _, layercode := range ar.LayerCodes {
		match, err := globMatch(pattern, layercode)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func (ar AssignedRight) MatchFlags(pattern string) (bool, error) {
	if len(ar.Flags) == 0 {
		return globMatch(pattern, "")
	}
	for _, flag := range ar.Flags {
		match, err := globMatch(pattern, flag)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func (ar AssignedRight) MatchTariffCode(pattern string) (bool, error) {
	return globMatch(pattern, ar.TariffCode)
}

// LoadAssignedRights loads the assigned rights from a JSON/YAML payload
func NewAssignedRightHistoryFromYAML(data []byte) (AssignedRights, error) {
	var history AssignedRights
	err := yaml.Unmarshal(data, &history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

// MatchingRule represents a rule to match the parking assigned rights to be used in a quota
type MatchingRule struct {
	TariffCodePattern   string `yaml:"tariff"`
	LayerCodePattern    string `yaml:"layer"`
	DurationTypePattern string `yaml:"type"`
	FlagsPattern        string `yaml:"flags"`
}

// Stringer for MatchingRule, print the area and type patterns
func (m MatchingRule) String() string {
	return fmt.Sprintf("(%s, %s, %s, %s)", m.TariffCodePattern, m.LayerCodePattern, m.DurationTypePattern, m.FlagsPattern)
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
	Update(now time.Time, history AssignedRights) error
	IsExausted() bool
	UseDuration(duration time.Duration) time.Duration
	String() string
}

// AbstractQuota is a helper to ease the implementation of different quotas types
type AbstractQuota struct {
	Name                       string                  `yaml:"name"`
	MatchingRules              MatchingRules           `yaml:"matching"`
	PeriodicityRule            timeutils.RecurrentDate `yaml:"periodicity"`
	DefaultTariffCodePattern   string                  `yaml:"-"`
	DefaultLayerCodePattern    string                  `yaml:"-"`
	DefaultDurationTypePattern string                  `yaml:"-"`
	DefaultFlagsPattern        string                  `yaml:"-"`
}

func (q AbstractQuota) GetName() string {
	return q.Name
}

// SelectReferenceTime selects the reference time to be used to filter the assigned rights based on the matching rules
func SelectReferenceTime(rule MatchingRule, detail DurationDetail, right AssignedRight) time.Time {
	reftime := detail.Start
	if reftime.IsZero() {
		reftime = right.StartDate
	}
	return reftime
}

// Helper function to match a glob string pattern, in a case-insensitive way
func globMatch(pattern, name string) (bool, error) {
	return filepath.Match(strings.ToLower(pattern), strings.ToLower(name))
}

// Filter filters the history of assigned rights based on the matching rules and calls the matchHandler for each matching detail
func (q AbstractQuota) Filter(from time.Time, history AssignedRights, matchAssignedRightHandler func(right AssignedRight),
	matchDurationDetailsHandler func(detail DurationDetail)) error {
	rules := q.MatchingRules
	if len(rules) == 0 {
		rules = []MatchingRule{{}}
	}
	// Iterate over all matching rules of the quota
	for _, rule := range rules {

		fmt.Println(" >> Matching rule", rule.String(), "from", from)

		// Iterate over all assigned rights in the history
		for i, right := range history {

			// Check if the assigned right tariffCode matches
			tariffCodePattern := rule.TariffCodePattern
			if tariffCodePattern == "" {
				tariffCodePattern = q.DefaultTariffCodePattern
			}
			matchTariffCode, err := right.MatchTariffCode(tariffCodePattern)
			if err != nil {
				return err
			}

			// Check if the assigned right layerCode matches
			layerCodePattern := rule.LayerCodePattern
			if layerCodePattern == "" {
				layerCodePattern = q.DefaultLayerCodePattern
			}
			matchLayerCode, err := right.MatchLayerCode(layerCodePattern)
			if err != nil {
				return err
			}

			// Check if the assigned right flags matches
			flagsPattern := rule.FlagsPattern
			if flagsPattern == "" {
				flagsPattern = q.DefaultFlagsPattern
			}
			matchFlags, err := right.MatchFlags(flagsPattern)
			if err != nil {
				return err
			}

			// Check if all the matching rules match and if set, call the Assigned Right callback
			match := matchTariffCode && matchLayerCode && matchFlags
			if match {
				fmt.Println("   >> Assigned right", i, "starting on", right.StartDate, "matches", rule.String())
			}
			if match && matchAssignedRightHandler != nil {
				matchAssignedRightHandler(right)
			}
			// If set, check duration details matches and call the Duration Detail callback
			if match && matchDurationDetailsHandler != nil {
				typePattern := rule.DurationTypePattern
				if typePattern == "" {
					typePattern = q.DefaultDurationTypePattern
				}
				for _, detail := range right.DurationDetails {
					match, err := globMatch(typePattern, detail.Type.ShortString())
					if err != nil {
						return err
					}
					if match {
						reftime := SelectReferenceTime(rule, detail, right)
						if !reftime.IsZero() && timeutils.TimeAfterOrEqual(reftime, from) {
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

func NewDurationQuota(name string, allowance time.Duration, period timeutils.RecurrentDate, rules []MatchingRule) *DurationQuota {
	return &DurationQuota{
		AbstractQuota: AbstractQuota{
			Name:                       name,
			MatchingRules:              rules,
			PeriodicityRule:            period,
			DefaultTariffCodePattern:   "*",
			DefaultLayerCodePattern:    "*",
			DefaultFlagsPattern:        "*",
			DefaultDurationTypePattern: FreeDuration.ShortString(),
		},
		Allowance: allowance,
	}
}

// Update updates the quota based on the history of assigned rights
func (q *DurationQuota) Update(now time.Time, history AssignedRights) error {
	fmt.Println("Update duration quota", q.Name)
	total := time.Duration(0)
	// Compute the start period of quota calculation
	start, err := q.PeriodStart(now)
	if err != nil {
		return err
	}
	// Compute the total duration of matching assigned rights
	err = q.Filter(start, history, nil, func(detail DurationDetail) {
		total += detail.Duration
		fmt.Println("   >> Duration detail", detail.Type, "duration", detail.Duration, "from", detail.Start)
	})
	q.used = total
	fmt.Println("Duration quota", q.Name, "used", q.used, "out of", q.Allowance)
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

func (q *DurationQuota) IsExausted() bool {
	return q.Available() <= 0
}

func (q *DurationQuota) UseDuration(duration time.Duration) time.Duration {
	if duration > q.Available() {
		duration = q.Available()
	}
	q.used += duration
	return duration
}

// Stringer for DurationQuota, print the name and the used/allowed values
func (q DurationQuota) String() string {
	return fmt.Sprintf("DurationQuota(%s): Usage %s/%s %v", q.Name, q.used, q.Allowance, q.AbstractQuota)
}

// CounterQuota represents a quota based on the number of parking assigned rights
type CounterQuota struct {
	AbstractQuota `yaml:",inline"`
	Allowance     int `yaml:"allowance"`
	used          int
}

func NewCounterQuota(name string, allowance int, period timeutils.RecurrentDate, rules []MatchingRule) *CounterQuota {
	return &CounterQuota{
		AbstractQuota: AbstractQuota{
			Name:                       name,
			MatchingRules:              rules,
			PeriodicityRule:            period,
			DefaultTariffCodePattern:   "*",
			DefaultLayerCodePattern:    "*",
			DefaultFlagsPattern:        "*",
			DefaultDurationTypePattern: FreeDuration.ShortString(),
		},
		Allowance: allowance,
	}
}

// Update updates the quota based on the history of assigned rights
func (q *CounterQuota) Update(now time.Time, history AssignedRights) error {
	fmt.Println("Update counter quota", q.Name)
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
	fmt.Println("Counter quota", q.Name, "used", q.used, "out of", q.Allowance)
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

func (q *CounterQuota) IsExausted() bool {
	return q.Available() <= 0
}

func (q *CounterQuota) UseDuration(duration time.Duration) time.Duration {
	if q.Available() <= 0 {
		duration = 0
	}
	q.used++
	return duration
}

// Stringer for CounterQuota, print the name and the used/allowed values
func (q CounterQuota) String() string {
	return fmt.Sprintf("CounterQuota(%s): Usage %d/%d %v", q.Name, q.used, q.Allowance, q.AbstractQuota)
}

type QuotaInventory map[string]Quota

func (qi QuotaInventory) Update(now time.Time, history AssignedRights) error {

	// Iterate over all quotas and update them
	for _, quota := range qi {
		err := quota.Update(now, history)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stringer for QuotaInventory, iterate over all quotas and print some details
func (qi QuotaInventory) String() string {
	var str strings.Builder
	str.WriteString("Quotas:\n")
	for _, quota := range qi {
		str.WriteString(" - ")
		str.WriteString(quota.String())
		str.WriteString("\n")
	}
	str.WriteString("\n")
	return str.String()
}

func (qi *QuotaInventory) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {
	temp := []struct {
		DurationQuota *DurationQuota `yaml:"duration"`
		CounterQuota  *CounterQuota  `yaml:"counter"`
	}{}

	// Unmarshal the quota into a temp struct
	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	*qi = make(QuotaInventory)
	for _, t := range temp {
		quota := Quota(nil)
		// TODO return an error if both DurationQuota and CounterQuota are set
		if t.DurationQuota != nil {
			quota = NewDurationQuota(t.DurationQuota.Name, t.DurationQuota.Allowance, t.DurationQuota.PeriodicityRule, t.DurationQuota.MatchingRules)
		} else if t.CounterQuota != nil {
			quota = NewCounterQuota(t.CounterQuota.Name, t.CounterQuota.Allowance, t.CounterQuota.PeriodicityRule, t.CounterQuota.MatchingRules)
		}
		if quota != nil {
			if quota.GetName() == "" {
				return fmt.Errorf("missing quota name")
			}
			(*qi)[quota.GetName()] = quota
		}
	}
	return nil
}

/*
type ParsableQuota struct {
	Quota Quota
}

func (q *ParsableQuota) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {
	temp := struct {
		QuotaName string `yaml:"quota"`
	}{}

	// Unmarshal the quota into a temp struct
	// dont check the error, as we we are interested only by the quota name
	unmarshal(&temp)

	// Search the coresponding quota
	fmt.Println(">>>>>> Quota name:", temp.QuotaName)
	// Get the quota from the context
	quota, exists := ContextGetQuotaByName(ctx, temp.QuotaName)
	if !exists {
		return fmt.Errorf("unknown quota: %s", temp.QuotaName)
	}
	q.Quota = quota

	return nil
}
*/
