package rtfprinter

/*
Copyright © 2021 Pontus Ullgren

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"encoding/csv"
	"io"
	"log"
	"strings"

	"github.com/kataras/tablewriter"
)

type RTFPrinter struct {
	format string
	out    io.Writer
}

/*
New creates a new RTF Minion printer
*/
func New(w io.Writer, format string) (RTFPrinter, error) {
	// TODO: Validate format
	return RTFPrinter{
		format: format,
		out:    w,
	}, nil
}

/*
Print the given data
*/
func (p *RTFPrinter) Print(data TableFormater) {

	switch strings.ToLower(p.format) {
	case "pretty":
		tablewriter := tablewriter.NewWriter(p.out)

		tabledata := data.GetTableData()
		tablewriter.SetHeader(tabledata.Headers)

		for _, row := range tabledata.Rows {
			tablewriter.Append(row)
		}
		tablewriter.Render() // Send output
	case "csv":
		w := csv.NewWriter(p.out)
		defer w.Flush()

		tabledata := data.GetTableData()
		if err := w.Write(tabledata.Headers); err != nil {
			log.Fatalln("error writing headers to file", err)
		}
		for _, record := range tabledata.Rows {
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to file", err)
			}
		}
	}

}

type TableFormater interface {
	GetTableData() TableData
}

type TableData struct {
	Headers []string
	Rows    [][]string
}
