# CSV Utility Functions Documentation

This document lists all the functions available in the CSV Utility package along with their signatures.

## Function Signatures

### 1. ReadCsvFile
```go
func ReadCsvFile(filePath string) ([][]string, error)
```
Reads a CSV file and returns the data as a slice of slices of strings.

### 2. WriteCsvFile
```go
func WriteCsvFile(filePath string, records [][]string) error
```
Writes a slice of slices of strings to a CSV file.

### 3. PrintCsvData
```go
func PrintCsvData(data [][]string)
```
Utility function to print CSV data.

### 4. FilterCsvData
```go
func FilterCsvData(records [][]string, filterFunc func([]string) bool) [][]string
```
Filters rows based on a condition.

### 5. AppendToCsvFile
```go
func AppendToCsvFile(filePath string, newRecords [][]string) error
```
Appends new rows to an existing CSV file.

### 6. ReadCsvWithHeader
```go
func ReadCsvWithHeader(filePath string) ([]map[string]string, error)
```
Reads a CSV file and returns a map for easier access to columns by header names.

### 7. WriteCsvWithHeader
```go
func WriteCsvWithHeader(filePath string, data []map[string]string) error
```
Writes data to a CSV file using a map, assuming the map keys as headers.



