package main

import (
	"fmt"
	"time"

	"github.com/iem-rd/quote-engine/engine"
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
    rules:
    - linear:
        name: "A"
        duration: 15h
        hourlyrate: 3.0
  - name: "weekdays"
    quota: plip
    rules: 
      - nonpaying:
          name: "Sunday"
          start: pattern(*/* SUN 00:00)
          end: pattern(*/* SUN 23:59:59)
      - nonpaying:
          name: "night"
          start: pattern(*/* 20:00)
          end: pattern(*/* 8:00)
      - absflatrate:
          name: "Forfait du samedi"
          start: pattern(*/* SAT 08:00)
          end: pattern(*/* SAT 16:00)
      - abslinear:
          name: "Prix spécial du lundi"
          start: pattern(*/* MON 08:00)
          end: pattern(*/* MON 16:00)
          hourlyrate: 1.5
      - linear:
          name: "A"
          duration: 1h
          hourlyrate: 1.0
          duration: 20h        
`

	t, err := engine.ParseTariffDefinition([]byte(sampleyaml2))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", t)

	//now := time.Now()
	now, _ := time.ParseInLocation("2006-01-02T15:04:05", "2025-03-14T15:54:30", time.Local)
	out := t.Compute(now, []engine.AssignedRight{})

	json, _ := out.ToJson()
	fmt.Println(string(json))

	//t.Sequences.Solve(time.Now(), time.Hour*24*7)

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
