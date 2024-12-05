package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

// logEntry represents a log entry structure
type logEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// regexPattern defines the regex to parse log lines
var regexPattern = regexp.MustCompile(`^(?P<timestamp>\S+) \[(?P<level>[A-Z]+)\] (?P<message>.+)$`)

// extractFields parses a single log line into a logEntry if valid
func extractFields(log string) (logEntry, bool) {
	matches := regexPattern.FindStringSubmatch(log)
	if matches == nil {
		return logEntry{}, false
	}
	return logEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
	}, true
}

// processBatch handles a batch of log lines, sending valid entries to the results channel
func processBatch(batch []string, output chan logEntry, seen *sync.Map) {
	for _, logLine := range batch {
		if entry, valid := extractFields(logLine); valid {
			// Check for duplicates using a concurrent map
			entryKey := fmt.Sprintf("%s|%s|%s", entry.Timestamp, entry.Level, entry.Message)
			if _, exists := seen.LoadOrStore(entryKey, struct{}{}); !exists {
				output <- entry
			}
		}
	}
}

// splitFileIntoChunks reads a file in chunks, splitting by the chunk size
func splitFileIntoChunks(filepath string, chunkSize int) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunks [][]string
	scanner := bufio.NewScanner(file)
	var chunk []string
	for scanner.Scan() {
		chunk = append(chunk, scanner.Text())
		if len(chunk) >= chunkSize {
			chunks = append(chunks, chunk)
			chunk = nil
		}
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks, scanner.Err()
}

func main() {
	// Path for the log file and the output JSON file
	filePath := "C:/Users/majjuri.kanak/Documents/Assignment2/log_file_500mb.log"
	outputFile := strings.Replace(filePath, ".log", ".json", 1)
	chunkSize := 1000 // Number of lines to read per chunk

	// Read file in chunks
	chunks, err := splitFileIntoChunks(filePath, chunkSize)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Channels and sync.Map for handling results and tracking duplicates
	results := make(chan logEntry, chunkSize)
	var wg sync.WaitGroup
	var seen sync.Map

	// Process chunks concurrently
	for _, chunk := range chunks {
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			processBatch(batch, results, &seen)
		}(chunk)
	}

	// Close results channel after all processing is done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results from the channel
	var logEntries []logEntry
	for entry := range results {
		logEntries = append(logEntries, entry)
	}

	// Save results to JSON file
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	// Encode log entries into the output file in JSON format
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(logEntries); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}
}
