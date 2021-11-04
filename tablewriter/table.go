package tablewriter

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

type Table struct {
	outBuf    bytes.Buffer
	tabWriter *tabwriter.Writer
}

func New() *Table {
	f := new(Table)
	f.tabWriter = tabwriter.NewWriter(&f.outBuf, 0, 0, 5, ' ', 0)

	return f
}

func (f *Table) WriteRow(cols ...interface{}) {
	format := strings.Repeat("%v\t", len(cols))
	_, _ = fmt.Fprintf(f.tabWriter, "\n"+format, cols...)
}

func (f *Table) Writef(format string, message ...interface{}) {
	_, _ = fmt.Fprintf(f.tabWriter, format, message...)
}

func (f *Table) String() string {
	if err := f.tabWriter.Flush(); err != nil {
		return err.Error()
	}

	return f.outBuf.String()
}
