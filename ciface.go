package ciface

import (
	"bytes"
	"encoding/csv"
	"log"
	"math"
	"strconv"
	"strings"
)

// CsvInterface is a structure describing the structure of the parser.
// Structure is in line with that of the standard libraries Csv parser.
type CsvInterface struct {
	Data      []byte
	Header    []string
	Delimiter rune
	Comment   rune
	Precision int64
}

// NewParser is used to initialize a new CsvInterface with data
// Other settings should be applied directly to the new object
func NewParser(data []byte) *CsvInterface {
	return &CsvInterface{
		Data:      data,
		Precision: 4,
	}
}

// Parse takes the CsvInterface struct and parses the []byte data to an []interface{}
func (cif *CsvInterface) Parse() []interface{} {
	//cif.Delimiter = ','
	bbuf := bytes.NewBuffer(cif.Data)

	// Configure the CSV reader
	reader := csv.NewReader(bbuf)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	if cif.Delimiter != 0 {
		reader.Comma = cif.Delimiter
	}

	if cif.Comment != 0 {
		reader.Comment = cif.Comment
	}

	// Run the Csv reader with the provided settings
	raw, err := reader.ReadAll()

	if err != nil {
		log.Fatal("Error while attempting to parse csv data", err)
	}

	// Start processing the raw line data. If no header is configured
	// the first line will be automatically parsed as the header.
	var output []interface{}

	for count, line := range raw {
		if count == 0 && cif.Header == nil {
			cif.Header = line
		} else {
			output = append(output, cif.LineConverter(line))
		}
	}

	return output
}

// LineConverter processes the csv line data to the proper json types and returns
// the resulting data structure as a map[string]interface{} (golang json structure).
func (cif *CsvInterface) LineConverter(line []string) interface{} {
	doc := make(map[string]interface{})

	for count, value := range line {
		if BooleanString(value) {
			bs, _ := strconv.ParseBool(value)
			doc[cif.Header[count]] = bs
		} else if number, err := strconv.ParseFloat(value, 64); err == nil {
			doc[cif.Header[count]] = Round(number, cif.Precision)
		} else {
			doc[cif.Header[count]] = value
		}
	}

	return doc
}

// BooleanString is used to detect if a string contains a boolean value
func BooleanString(s string) bool {
	low := strings.ToLower(s)
	return low == "true" || low == "false"
}

// Round is used to round numbers to a certain precision
func Round(number float64, p int64) float64 {
	precision := math.Pow(10, float64(p))
	corrected := float64(int(number*precision)) / precision

	// Make sure to round up the last digit
	if (number-corrected)*precision > 0.5 {
		corrected = corrected + (1 / precision)
	}

	return corrected
}
