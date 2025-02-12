package parser

import (
	"time"
)

func TimeAfterOrEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.After(t2)
}

func TimeBeforeOrEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2) || t1.Before(t2)
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

type RecurrentSegment struct {
	Start RecurrentDate `yaml:"start" validate:"required"`
	End   RecurrentDate `yaml:"end" validate:"required"`
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

/*
func (rs *RecurrentSegment) IsWithin(t time.Time) bool {

}
*/

// Stringer for RecurrentSegment display start and end
func (rs RecurrentSegment) String() string {
	return rs.Start.String() + " -> " + rs.End.String()
}
