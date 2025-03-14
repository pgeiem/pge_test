package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type OutputSegment struct {
	SegName      string       `json:"n,omitempty"`
	Trace        []string     `json:"dbg,omitempty"`
	Duration     int          `json:"d"`
	Amount       Amount       `json:"a"`
	Islinear     bool         `json:"l"`
	DurationType DurationType `json:"dt,omitempty"`
	Meta         MetaData     `json:"m,omitempty"`
}

func (seg OutputSegment) String() string {
	at := time.Duration(seg.Duration) * time.Second
	return fmt.Sprintf(" - %s: %s (isLinear %t)", at, seg.Amount, seg.Islinear)
}

type OutputSegments []OutputSegment

func (segs OutputSegments) String() string {
	out := "OutputSegments:\n"
	for i := range segs {
		out += segs[i].String() + "\n"
	}
	return out
}

type Output struct {
	Now   time.Time      `json:"now"`
	Table OutputSegments `json:"table"`
}

func (segs Output) ToJson() ([]byte, error) {
	return json.Marshal(segs)
}

func (s *SolverRules) GenerateOutput(now time.Time, detailed bool) Output {
	var out Output
	var previous SolverRule

	out.Now = now

	for _, rule := range *s {
		fmt.Println("Rule", rule)
		// If there is a gap between the previous rule and the current one this is the end of the output
		if previous.To != rule.From {
			break
		}
		seg := OutputSegment{
			Duration:     int(math.Round(rule.To.Seconds() - previous.To.Seconds())),
			Amount:       rule.EndAmount,
			Islinear:     !rule.IsFlatRate(),
			DurationType: rule.DurationType,
			Meta:         rule.Meta,
		}
		if detailed {
			seg.SegName = rule.Name()
			seg.Trace = rule.Trace
		}
		out.Table = append(out.Table, seg)
		previous = rule
	}
	return out
}
