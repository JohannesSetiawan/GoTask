package task

import "time"

// validateDatetime validates if the provided string is a valid datetime in RFC3339 format
func validateDatetime(dateStr string) error {
	// Try to parse as RFC3339 (standard format: 2006-01-02T15:04:05Z07:00)
	_, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// If RFC3339 fails, try RFC3339Nano for more precision
		_, err = time.Parse(time.RFC3339Nano, dateStr)
	}
	return err
}
