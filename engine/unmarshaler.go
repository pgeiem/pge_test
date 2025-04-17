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
	if duration != nil {
		*duration = d
	}
	return nil
}

func unmarshalRecurrentDate(rec *RecurrentDate, data []byte) error {
	str := strings.Trim(string(data), `"`)
	tmp, err := ParseRecurrentDate(str)
	if err != nil {
		return err
	}
	if rec != nil {
		*rec = tmp
	}
	return nil
}
