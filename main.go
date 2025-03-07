package main

import (
	"fmt"
	"time"

	"github.com/iem-rd/quoteengine/parser"
)

func main() {

	sampleyaml := `
version: "0.1"

nonpaying:
  - name: "christmas"
    start: pattern(12/25 00:00)
    end: pattern(12/25 23:59:59)
  - name: "new year"
    start: pattern(01/01 00:00)
    end: pattern(01/01 23:59:59)
  - name: "sunday"
    start: pattern(*/* SUN 00:00)
    end: pattern(*/* SUN 23:59:59)
  - name: "lunch"
    start: pattern(*/* * 12:00)
    end: pattern(*/* * 13:00)
  - name: "night"
    start: pattern(*/* * 22:00)
    end: pattern(*/* * 06:00)

quotas:
  - duration:
      name: "plop"
      allowance: 2h10m
      periodicity: duration(4h)
      matching:
        - area: z*
          type: paying
        - area: t*
          type: nonpaying

  - counter:
      name: "plip"
      allowance: 42
      periodicity: pattern(*/* 12:00)

sequences:
  - name: "weekend"
    start: pattern(*/* SAT 00:00)	
    end: pattern(*/* MON 00:00)
    quota: plop
  - name: "weekdays"
    start: pattern(*/* SAT 00:00)	
    end: pattern(*/* MON 00:00)
    quota: plip
    nonpaying:
      - name: "saturday"
        start: pattern(*/* SUN 00:00)
        end: pattern(*/* SUN 23:59:59)
    rules: 
      - linear:
          name: "A"
          duration: 1h
          hourlyrate: 1.0
      - flatrate:
          name: "B"
          duration: 1h
          amount: 1.0 
`

	t, err := parser.ParseTariffDefinition([]byte(sampleyaml))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", t)

	t.Sequences.Solve(time.Now(), time.Hour*24*7)

}

// plop:
//       type: duration
//       allowance: 2h10m
//       periodicity: duration(4h)
//       matching:
//         - area: z*
//           type: paying
//         - area: t*
//           type: nonpaying
//   plip:
//       type: counter
//       allowance: 42
//       periodicity: pattern(*/* 12:00)
