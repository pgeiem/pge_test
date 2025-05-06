module github.com/iem-rd/quote-engine/engine

go 1.23.4

require (
	github.com/goccy/go-yaml v1.15.17
	github.com/google/btree v1.1.3
	github.com/iem-rd/quote-engine/table v0.0.1
	github.com/iem-rd/quote-engine/timeutils v0.1.2
)

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/teambition/rrule-go v1.8.2 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

//replace github.com/iem-rd/quote-engine/table => ../table

//replace github.com/iem-rd/quote-engine/timeutils => ../timeutils
