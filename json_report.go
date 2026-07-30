package main

import (
	"encoding/json"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	// Collect the keys.
	keys := make([]string, 0, len(pages))
	for key := range pages {
		keys = append(keys, key)
	}

	// Sort them.
	sort.Strings(keys)

	// Build the slice in sorted order.
	sorted := make([]PageData, 0, len(keys))
	for _, key := range keys {
		sorted = append(sorted, pages[key])
	}

	// Convert to formatted JSON.
	data, err := json.MarshalIndent(sorted, "", "  ")
	if err != nil {
		return err
	}

	// Write the file.
	return os.WriteFile(filename, data, 0644)
}