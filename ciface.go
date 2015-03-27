package ciface

import (
	"bytes"
	"encoding/csv"
	"errors"
	"math"
	"strconv"
	"strings"
)

// CsvInterface is a structure describing the structure of the parser.
// Structure is in line with that of the standard libraries Csv parser.
type CsvInterface struct {
	Reader    *csv.Reader
	Data      []byte
	Header    []string
	Precision int64
}

// NewParser is used to initialize a new CsvInterface with data
// Other settings should be applied directly to the new object
func NewParser(data []byte) *CsvInterface {
	reader := csv.NewReader(bytes.NewBuffer(data))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	return &CsvInterface{
		Reader:    reader,
		Precision: 4,
	}
}

// Parse takes the CsvInterface struct and parses the []byte data to an []interface{}
func (cif *CsvInterface) Parse() ([]interface{}, error) {
	// Run the Csv reader with the provided settings
	raw, err := cif.Reader.ReadAll()

	// Start processing the raw line data. If no header is configured
	// the first line will be automatically parsed as the header.
	var output []interface{}

	for count, line := range raw {
		if count == 0 && cif.Header == nil {
			cif.Header = line
		} else {
			if lineInterface, lineErr := cif.LineConverter(line); lineErr == nil {
				output = append(output, lineInterface)
			} else {
				err = lineErr
			}
		}
	}

	return output, err
}

// LineConverter processes the csv line data to the proper json types and returns
// the resulting data structure as a map[string]interface{} (golang json structure).
func (cif *CsvInterface) LineConverter(line []string) (interface{}, error) {
	var err error
	doc := make(map[string]interface{})

	if len(line) == len(cif.Header) {

		for count, value := range line {
			if BooleanString(value) {
				bs, _ := strconv.ParseBool(value)
				doc[cif.Header[count]] = bs
			} else if number, err := strconv.ParseFloat(value, 64); err == nil {
				doc[cif.Header[count]] = Round(number, cif.Precision)
			} else if value == "" {
				doc[cif.Header[count]] = nil
			} else {
				doc[cif.Header[count]] = value
			}
		}
	} else {
		err = errors.New("[CSV] - mismatching header and item count")
	}

	return doc, err
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
