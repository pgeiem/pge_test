package main

import (
	"fmt"
	"time"

	dps "github.com/markusmobius/go-dateparser"
)

func main() {
	// Current time is 2015-07-31 12:00:00 UTC
	cfg := &dps.Configuration{
		CurrentTime: time.Date(2015, 7, 31, 12, 0, 0, 0, time.UTC),
	}

	fmt.Println(dps.Parse(cfg, "December 2015"))
	// time: 2015-12-31 00:00:00 UTC (day from current time)

	fmt.Println(dps.Parse(cfg, "February 2020"))
	// time: 2020-02-29 00:00:00 UTC (day from current time, corrected for leap year)

	fmt.Println(dps.Parse(cfg, "December"))
	// time: 2015-12-31 00:00:00 UTC (year and day from current time)

	fmt.Println(dps.Parse(cfg, "2015"))
	// time: 2015-07-31 00:00:00 UTC (day and month from current time)

	fmt.Println(dps.Parse(cfg, "Sunday"))
	// time: 2015-07-26 00:00:00 UTC (the closest Sunday from current time)

	fmt.Println(dps.Parse(cfg, "12:23:34"))

}
