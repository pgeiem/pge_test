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

type Segment struct {
	Start time.Time `yaml:"start" validate:"required"`
	End   time.Time `yaml:"end" validate:"required"`
}

func (s *Segment) Duration() time.Duration {
	return s.End.Sub(s.Start)
}

func (s *Segment) IsWithin(t time.Time) bool {
	return TimeAfterOrEqual(t, s.Start) && t.Before(s.End)
}

func (s *Segment) String() string {
	return s.Start.String() + " -> " + s.End.String()
}

func (s *Segment) ToRelativeTimeSpan(now time.Time) RelativeTimeSpan {
	return RelativeTimeSpan{
		From: s.Start.Sub(now),
		To:   s.End.Sub(now),
	}
}

type RecurrentSegment struct {
	Start RecurrentDate `yaml:"start" validate:"required"`
	End   RecurrentDate `yaml:"end" validate:"required"`
}

// Create a recurrent segment from two recurrent dates pattern
func NewRecurrentSegmentFromPatterns(start, end string) (RecurrentSegment, error) {
	startrule, err := ParseRecurrentDate(start)
	if err != nil {
		return RecurrentSegment{}, err
	}
	endrule, err := ParseRecurrentDate(end)
	if err != nil {
		return RecurrentSegment{}, err
	}
	return RecurrentSegment{
		Start: startrule,
		End:   endrule,
	}, nil
}

func (rs *RecurrentSegment) Next(now time.Time) (Segment, error) {
	var err error
	s := Segment{}
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

func (rs *RecurrentSegment) Prev(now time.Time) (Segment, error) {
	var err error
	s := Segment{}
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

func (rs *RecurrentSegment) Between(from, to time.Time) []Segment {
	var segments []Segment
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

func (rs *RecurrentSegment) BetweenIterator(from, to time.Time, iterator func(Segment) bool) {
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

func (rs *RecurrentSegment) IsWithin(t time.Time) (bool, error) {
	within, _, err := rs.IsWithinWithSegment(t)
	return within, err
}

func (rs *RecurrentSegment) IsWithinWithSegment(t time.Time) (bool, Segment, error) {

	s, err := rs.Prev(t)
	if err != nil {
		return false, Segment{}, err
	}
	if s.IsWithin(t) {
		return true, s, nil
	}

	// Handle case where t is the start of the segment
	s, err = rs.Next(s.Start)
	if err != nil {
		return false, Segment{}, err
	}
	if s.IsWithin(t) {
		return true, s, nil
	}
	return false, Segment{}, err
}

// Stringer for RecurrentSegment display start and end
func (rs RecurrentSegment) String() string {
	return rs.Start.String() + " -> " + rs.End.String()
}
