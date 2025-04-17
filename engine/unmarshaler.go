package engine

import (
	"strings"
	"time"

	"github.com/iem-rd/quote-engine/timeutils"
)

func unmarshalTimeDuration(duration *time.Duration, data []byte) error {
	str := strings.Trim(string(data), `"`)
	d, err := timeutils.ParseDuration(str)
	if err != nil {
		return err
	}
	if duration != nil {
		*duration = d
	}
	return nil
}

func unmarshalRecurrentDate(rec *timeutils.RecurrentDate, data []byte) error {
	str := strings.Trim(string(data), `"`)
	tmp, err := timeutils.ParseRecurrentDate(str)
	if err != nil {
		return err
	}
	if rec != nil {
		*rec = tmp
	}
	return nil
}
