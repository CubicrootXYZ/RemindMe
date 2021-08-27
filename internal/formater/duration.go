package formater

import (
	"fmt"
	"time"
)

// ToNiceDuration formats a time.Duration into a nice string
func ToNiceDuration(d time.Duration) string {
	pre := ""
	if d < 0 {
		d *= -1
		pre = "-"
	}

	if d < time.Minute {
		return fmt.Sprintf("%s%.0f seconds", pre, float64(d/time.Second))
	} else if d < time.Hour {
		return fmt.Sprintf("%s%.0f minutes", pre, float64(d/time.Second))
	} else if d < 48*time.Hour {
		return fmt.Sprintf("%s%.0f hours", pre, float64(d/time.Hour))
	}
	return fmt.Sprintf("%s%.0f days", pre, float64(d/(24*time.Hour)))
}
