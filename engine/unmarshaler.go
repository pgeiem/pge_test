package engine

import (
	"strings"
	"time"
)

func unmarshalTimeDuration(duration *time.Duration, data []byte) error {
	str := strings.Trim(string(data), `"`)
	d, err := ParseDuration(str)
	if err != nil {
		return err
	}
	*duration = d.ToDuration()
	return nil
}

func unmarshalDuration(duration *Duration, data []byte) error {
	str := strings.Trim(string(data), `"`)
	d, err := ParseDuration(str)
	if err != nil {
		return err
	}
	*duration = d
	return nil
}

func unmarshalRecurrentDate(rec *RecurrentDate, data []byte) error {
	str := strings.Trim(string(data), `"`)
	tmp, err := ParseRecurrentDate(str)
	if err != nil {
		return err
	}
	*rec = tmp
	return nil
}
