package engine

import (
	"time"
)

func unmarshalTimeDuration(duration *time.Duration, data []byte) error {
	d, err := ParseDuration(string(data))
	if err != nil {
		return err
	}
	*duration = d.ToDuration()
	return nil
}

func unmarshalDuration(duration *Duration, data []byte) error {
	d, err := ParseDuration(string(data))
	if err != nil {
		return err
	}
	*duration = d
	return nil
}

func unmarshalRecurrentDate(rec *RecurrentDate, data []byte) error {
	tmp, err := ParseRecurrentDate(string(data))
	if err != nil {
		return err
	}
	*rec = tmp
	return nil
}
