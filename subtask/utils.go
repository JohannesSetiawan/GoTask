package subtask

import (
	"time"

	"go-tasker/database"

)

func validateDatetime(dateStr string) error {
	if _, err := time.Parse(time.RFC3339, dateStr); err != nil {
		if _, err := time.Parse(time.RFC3339Nano, dateStr); err != nil {
			return err
		}
	}
	return nil
}

func validateTaskStatus(status database.TaskStatus) bool {
	switch status {
	case database.StatusCreated, database.StatusInProgress, database.StatusDone:
		return true
	default:
		return false
	}
}