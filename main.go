package main

import (
	"fmt"

	parser "github.com/iem-rd/tarif2/parser"
)

func main() {

	sampleyaml := `
quotas:
  - duration: 
      name: "duration"
      allowance: 2h10m
  - counter:
      name: "counter"
      allowance: 42
  - counter:
      allowance: 42
      name: plop
`

	x, err := parser.ParseTariffDefinitionString(sampleyaml)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", x)
}
