package parser

import (
	"time"

	"github.com/iem-rd/quoteengine/engine"
)

func unmarshalTimeDuration(duration *time.Duration, data []byte) error {
	d, err := engine.ParseDuration(string(data))
	if err != nil {
		return err
	}
	*duration = d.ToDuration()
	return nil
}

func unmarshalDuration(duration *engine.Duration, data []byte) error {
	d, err := engine.ParseDuration(string(data))
	if err != nil {
		return err
	}
	*duration = d
	return nil
}

func unmarshalRecurrentDate(rec *engine.RecurrentDate, data []byte) error {
	tmp, err := engine.ParseRecurrentDate(string(data))
	if err != nil {
		return err
	}
	*rec = tmp
	return nil
}
