package engine

import (
	"fmt"
	"time"
)

func TimeAfterOrEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.After(t2)
}

func TimeBeforeOrEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.Before(t2)
}

type RelativeTimeSpan struct {
	From time.Duration
	To   time.Duration
}

func (r RelativeTimeSpan) String() string {
	return fmt.Sprintf("[%s, %s]", r.From, r.To)
}

func (r RelativeTimeSpan) Duration() time.Duration {
	return r.To - r.From
}

func (r RelativeTimeSpan) IsValid() bool {
	return r.From <= r.To
}

type AbsTimeSpan struct {
	Start time.Time `yaml:"start"`
	End   time.Time `yaml:"end"`
}

func (s *AbsTimeSpan) Duration() time.Duration {
	return s.End.Sub(s.Start)
}

func (s *AbsTimeSpan) IsWithin(t time.Time) bool {
	return TimeAfterOrEqual(t, s.Start) && t.Before(s.End)
}

func (s *AbsTimeSpan) String() string {
	return s.Start.String() + " -> " + s.End.String()
}

func (s *AbsTimeSpan) ToRelativeTimeSpan(now time.Time) RelativeTimeSpan {
	return RelativeTimeSpan{
		From: s.Start.Sub(now),
		To:   s.End.Sub(now),
	}
}

type RecurrentTimeSpan struct {
	Start RecurrentDate `yaml:"start"`
	End   RecurrentDate `yaml:"end"`
}

// Create a recurrent timespan from two recurrent dates pattern
func NewRecurrentTimeSpanFromPatterns(start, end string) (RecurrentTimeSpan, error) {
	startrule, err := ParseRecurrentDate(start)
	if err != nil {
		return RecurrentTimeSpan{}, err
	}
	endrule, err := ParseRecurrentDate(end)
	if err != nil {
		return RecurrentTimeSpan{}, err
	}
	return RecurrentTimeSpan{
		Start: startrule,
		End:   endrule,
	}, nil
}

func (rs *RecurrentTimeSpan) Next(now time.Time) (AbsTimeSpan, error) {
	var err error
	s := AbsTimeSpan{}
	s.Start, err = rs.Start.Next(now)
	if err != nil {
		return s, err
	}
	s.End, err = rs.End.Next(s.Start)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (rs *RecurrentTimeSpan) Prev(now time.Time) (AbsTimeSpan, error) {
	var err error
	s := AbsTimeSpan{}
	s.Start, err = rs.Start.Prev(now)
	if err != nil {
		return s, err
	}
	s.End, err = rs.End.Next(s.Start)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (rs *RecurrentTimeSpan) Between(from, to time.Time) []AbsTimeSpan {
	var segments []AbsTimeSpan
	now := from
	for now.Before(to) {
		segment, err := rs.Next(now)
		if err != nil {
			break
		}
		if segment.Start.After(to) {
			break
		}
		segments = append(segments, segment)
		now = segment.Start
	}
	return segments
}

func (rs *RecurrentTimeSpan) BetweenIterator(from, to time.Time, iterator func(AbsTimeSpan) bool) {
	now := from
	for now.Before(to) {
		segment, err := rs.Next(now)
		if err != nil {
			break
		}
		if segment.Start.After(to) {
			break
		}
		if !iterator(segment) {
			break
		}
		now = segment.Start
	}
}

func (rs *RecurrentTimeSpan) IsWithin(t time.Time) (bool, AbsTimeSpan, error) {

	s, err := rs.Prev(t)
	if err != nil {
		return false, AbsTimeSpan{}, err
	}
	if s.IsWithin(t) {
		return true, s, nil
	}

	// Handle case where t is the start of the segment
	s, err = rs.Next(s.Start)
	if err != nil {
		return false, AbsTimeSpan{}, err
	}
	if s.IsWithin(t) {
		return true, s, nil
	}
	return false, AbsTimeSpan{}, err
}

// Stringer for RecurrentSegment display start and end
func (rs RecurrentTimeSpan) String() string {
	if rs.Start != nil && rs.End != nil {
		return rs.Start.String() + " -> " + rs.End.String()
	}
	return "<nil>"
}
