package helpers

import (
	"fmt"
	"strings"
	"time"
)

func DateFormat(t time.Time) string {
	month := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	day := t.Day()
	monthName := month[t.Month()-1]
	year := t.Year()

	return fmt.Sprintf("%d %s %d", day, monthName, year)
}

// di gunakan di service employee (create default employee)
func NormalizeName(name string, email string) string {
	if name != "" {
		return name
	}

	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return strings.Title(parts[0])
	}

	return "User"
}