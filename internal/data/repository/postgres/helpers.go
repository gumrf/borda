package postgres

import (
	"fmt"
)

// FormatLimitOffset returns a SQL string for a given limit & offset.
// Clauses are only added if limit and/or offset are greater than zero.
func formatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)

	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)

	} else if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}

	return ""
}

// FormatError is a helper func for dealing with errors.
func formatError(err error) (int, error) {
	return -1, fmt.Errorf("TaskRepository.Create: %v", err)
}
