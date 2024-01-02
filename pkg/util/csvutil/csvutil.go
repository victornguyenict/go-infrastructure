package csvutil

import (
	"encoding/csv"
	"errors"
	"os"
)

// ReadCsvFile reads a CSV file and returns the records as a slice of slices of strings.
func ReadCsvFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// WriteCsvFile writes the given records to a CSV file.
func WriteCsvFile(filePath string, records [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// PrintCsvData is a utility function to print CSV data.
func PrintCsvData(data [][]string) {
	for _, row := range data {
		for _, record := range row {
			print(record + " ")
		}
		println()
	}
}

// FilterCsvData filters the records based on a provided filter function.
func FilterCsvData(records [][]string, filterFunc func([]string) bool) [][]string {
	var filteredRecords [][]string
	for _, record := range records {
		if filterFunc(record) {
			filteredRecords = append(filteredRecords, record)
		}
	}
	return filteredRecords
}

// AppendToCsvFile appends new records to an existing CSV file.
func AppendToCsvFile(filePath string, newRecords [][]string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range newRecords {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// ReadCsvWithHeader reads a CSV file and returns a slice of maps for easy column access.
func ReadCsvWithHeader(filePath string) ([]map[string]string, error) {
	records, err := ReadCsvFile(filePath)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.New("CSV file is empty")
	}

	headers := records[0]
	var result []map[string]string
	for _, record := range records[1:] {
		if len(record) != len(headers) {
			continue // Skip malformed rows
		}
		rowMap := make(map[string]string)
		for i, header := range headers {
			rowMap[header] = record[i]
		}
		result = append(result, rowMap)
	}

	return result, nil
}

// WriteCsvWithHeader writes data to a CSV file using a map with headers.
func WriteCsvWithHeader(filePath string, data []map[string]string) error {
	if len(data) == 0 {
		return errors.New("no data to write")
	}

	// Extract headers from the first row
	headers := make([]string, 0, len(data[0]))
	for header := range data[0] {
		headers = append(headers, header)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, rowMap := range data {
		row := make([]string, len(headers))
		for i, header := range headers {
			row[i] = rowMap[header]
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
