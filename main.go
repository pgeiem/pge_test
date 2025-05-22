package main

import (
	"fmt"
	"time"

	"github.com/iem-rd/quote-engine/engine"
)

func main() {

	plop :=
		`
version: "0.1"
config:
  window: 5d
quotas:
  - duration:
      name: "plop"
      periodicity: duration(4h)
      allowance: 2h10m
      matching:
        - layer: z*
          type: paying
        - layer: t*	
          type: nonpaying

sequences:
- name: "BC01"
  rules:
  - absflatrate:
      name: "daily"
      start: pattern(*/* 00:00)
      end: pattern(*/* 00:00)
      amount: 3.50
  - fixedrate:
      name: "minamount"
      duration: 42m
      amount: 0.70
      quota: "plop"
  - linear:
      name: "hourlyrate"
      hourlyrate: 1.0
      duration: 10h
      quota: "plop"

`

	t, err := engine.ParseTariffDefinition([]byte(plop))
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
