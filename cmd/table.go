package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func newTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	t.Style().Format.Header = text.FormatDefault

	return t
}

func sprintfTransformer(format string) text.Transformer {
	return func(val interface{}) string {
		return fmt.Sprintf(format, val)
	}
}

func scoreTransformer(val interface{}) string {
	if number, ok := val.(float32); ok {
		str := fmt.Sprintf("%.2f", number)

		if number >= 30.0 {
			return text.Colors{text.FgHiGreen}.Sprint(str)
		} else if number >= 10.0 {
			return text.Colors{text.FgHiYellow}.Sprint(str)
		}

		return text.Colors{text.FgHiRed}.Sprint(str)
	}

	return ""
}
