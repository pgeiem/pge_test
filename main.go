package main

import (
	"fmt"

	parser "github.com/iem-rd/tarif2/parser"
)

func main() {

	// 	sampleyaml := `
	// quotas:
	//   - duration:
	//       name: "Plop"
	//       allowance: 2h10m
	//       periodicity: duration(4h)
	//       matching:
	//         - area: z*
	//           type: paying
	//         - area: t*
	//           type: nonpaying
	//   - counter:
	//       name: "Plip"
	//       allowance: 42
	//       periodicity: pattern(*/* 12:00)
	// `

	sampleyaml := `
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
  plop:
      type: duration
      name: "Plop"
      allowance: 2h10m
      periodicity: duration(4h)
      matching: 
        - area: z*
          type: paying
        - area: t*
          type: nonpaying
  plip:
      type: counter
      name: "Plip"
      allowance: 42
      periodicity: pattern(*/* 12:00)
`

	x, err := parser.ParseTariffDefinitionString(sampleyaml)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", x)
}
