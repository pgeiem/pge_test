package engine

import (
	"encoding/json"
	"fmt"
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
	Now        time.Time      `json:"now"`
	ExpiryDate time.Time      `json:"expiry"`
	Table      OutputSegments `json:"table"`
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
		// If the segement is linear and is longer than the target duration, we need to calculate the amount for the remaining duration
		if seg.Islinear && targetDuration < totDuration+segDuration {
			fmt.Println("   >> Segment (partial linear)", seg.SegName, "Total amount", totAmount, "Total duration", totDuration)
			return Amount(float64(seg.Amount)*float64(targetDuration-totDuration)/float64(segDuration)) + totAmount
		}
		// If the segment is not linear and is longer or egual to the target duration, include it in the total
		if !seg.Islinear && targetDuration <= totDuration+segDuration {
			fmt.Println("   >> Segment (fixed)", seg.SegName, "Total amount", totAmount, "Total duration", totDuration)
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
