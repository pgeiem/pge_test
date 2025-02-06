package main

import (
	"fmt"

	parser "github.com/iem-rd/tarif2/parser"
)

func main() {

	sampleyaml := `
quotas:
  - duration: 
      name: "Plop"
      allowance: 2h10m
      periodicity: duration(4h)
      matching: 
        - area: z*
          type: paying
        - area: t*
          type: nonpaying
  - counter:
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
