package table

import "github.com/fatih/color"

func SetDefaultTheme(tbl *Table) {
	color.NoColor = false // Force enabled color output, as color capability is not detected correctly in VSCode
	headerFmt := color.New(color.FgHiMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgHiBlue).SprintfFunc()
	(*tbl).WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
}

func TitleTheme() *color.Color {
	return color.New(color.FgHiYellow).Add(color.Bold)
}
