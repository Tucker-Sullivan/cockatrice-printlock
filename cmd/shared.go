package cmd

import "strings"

func parseCSV(input string) []string {
	// Optional: parse --sets "ABC,DEF,GHI"
	var values []string
	if strings.TrimSpace(input) != "" {
		parts := strings.Split(input, ",")
		values = make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				values = append(values, p)
			}
		}
	}
	return values
}
