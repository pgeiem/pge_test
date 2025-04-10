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

func (segs Output) AmountForDuration(targetDuration time.Duration) Amount {
	totAmount := Amount(0)
	totDuration := time.Duration(0)
	fmt.Println("AmountForDuration, Target duration", targetDuration, "nb rules", len(segs.Table))
	for _, seg := range segs.Table {
		segDuration := time.Duration(seg.Duration) * time.Second
		if targetDuration < totDuration+segDuration {
			if seg.Islinear {
				fmt.Println("   >> Segment (partially)", seg.SegName, "Total amount", totAmount, "Total duration", totDuration)
				return Amount(float64(seg.Amount)*float64(targetDuration-totDuration)/float64(segDuration)) + totAmount
			}
			fmt.Println("   >> Segment", seg.SegName, "Total amount", totAmount, "Total duration", totDuration)
			return seg.Amount + totAmount
		}
		totAmount += seg.Amount
		totDuration += segDuration
		fmt.Println("   >> Segment", seg.SegName, "Total amount", totAmount, "Total duration", totDuration)
	}
	if targetDuration > totDuration {
		fmt.Println("WARNING: Duration is greater than the total duration of the output")
		return Amount(0)
	}
	return totAmount
}

func (s *SolverRules) GenerateOutput(now time.Time, detailed bool) Output {
	var out Output
	var previous SolverRule

	out.Now = now

	fmt.Println("Generating output for", len(*s), "rules")
	for _, rule := range *s {
		fmt.Println("   Rule", rule)
		// If there is a gap between the previous rule and the current one this is the end of the output
		if previous.To != rule.From {
			fmt.Println("   >> Gap detected, end of output", previous, rule)
			break
		}
		seg := OutputSegment{
			Duration:     int(math.Round(rule.To.Seconds() - previous.To.Seconds())),
			Amount:       rule.EndAmount.Simplify(),
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
