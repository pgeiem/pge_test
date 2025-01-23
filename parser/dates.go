package parser

/*
import (
	"fmt"
	"time"
	"github.com/adhocore/gronx"

// ParserRecurenceRule is a struct that represents a recurence rule in the parser
type ParserRecurenceRule struct {
	Start string
	Duration string
}

// ParserRecurenceRules reprensents a list of ParserRecurenceRule in the parser
type ParserRecurenceRules []ParserRecurenceRule

// Represents a TimeSegment after parsing
type TimeSegment struct {
	Start time.Time
	Duration time.Duration
}

type TimeSegments []TimeSegment

// Parse parses and develop a ParserRecurenceRule into a TimeSegments
func (rr ParserRecurenceRule) Parse() (TimeSegments, error) {

	gronx

	start, err := time.Parse("15:04", rr.Start)
	if err != nil {
		return nil, fmt.Errorf("Error parsing start time: %s", err)
	}

	duration, err := time.ParseDuration(rr.Duration)
	if err != nil {
		return nil, fmt.Errorf("Error parsing duration: %s", err)
	}

	return TimeSegments{TimeSegment{start, duration}}, nil
}
*/
